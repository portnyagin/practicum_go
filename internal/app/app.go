package app

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	config2 "github.com/portnyagin/practicum_go/internal/app/config"
	"github.com/portnyagin/practicum_go/internal/app/custommiddleware"
	"log"
	"net/http"
)

/*
Напишите сервис для сокращения длинных URL. Требования:
Сервер должен быть доступен по адресу: http://localhost:8080.
Сервер должен предоставлять два эндпоинта: POST / и GET /{id}.
Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает в ответ правильный сокращённый URL.
Эндпоинт GET /{id} принимает в качестве URL параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400.
*/

func Start() {
	config := config2.NewConfig()
	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	repo, err := NewBaseRepository(config.FileStorage)
	if err != nil {
		fmt.Println("can't init repository", err)
		return
	}
	service := NewZipService(repo, config.BaseURL)
	h := NewZipURLHandler(service)
	router := chi.NewRouter()
	router.Use(middleware.CleanPath)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(custommiddleware.Compress)
	router.Route("/", func(r chi.Router) {
		r.Get("/", h.HelloHandler)
		r.Get("/{id}", h.GetMethodHandler)
		r.Post("/api/shorten", h.PostAPIShortenHandler)
		r.Post("/", h.PostMethodHandler)
		r.Put("/", h.DefaultHandler)
		r.Patch("/", h.DefaultHandler)
		r.Delete("/", h.DefaultHandler)
	})

	err = http.ListenAndServe(config.ServerAddress, router)
	if err != nil {
		fmt.Println("can't start service")
		fmt.Println(err)
	}
}
