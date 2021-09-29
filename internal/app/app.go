package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	config2 "github.com/portnyagin/practicum_go/internal/app/config"
	"github.com/portnyagin/practicum_go/internal/app/custommiddleware"
	"github.com/portnyagin/practicum_go/internal/app/handler"
	"github.com/portnyagin/practicum_go/internal/app/infrastructure"
	"github.com/portnyagin/practicum_go/internal/app/repository"
	service2 "github.com/portnyagin/practicum_go/internal/app/service"
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
	fileRepository, err := repository.NewFileRepository(config.FileStorage)

	if err != nil {
		fmt.Println("can't init file repository", err)
		return
	}
	postgresHandler, err := infrastructure.NewPostgresqlHandler(context.Background(), config.DatabaseDSN)
	if err != nil {
		fmt.Println("can't init postgres handler", err)
		return
	}
	if config.Reinit {
		err = repository.ClearDatabase(postgresHandler)
		if err != nil {
			fmt.Println("can't clear database structure", err)
			return
		}
	}
	err = repository.InitDatabase(postgresHandler)
	if err != nil {
		fmt.Println("can't init database structure", err)
		return
	}

	postgresRepository, err := repository.NewPostgresRepository(postgresHandler)
	if err != nil {
		fmt.Println("can't init postgres repository", err)
		return
	}
	service := service2.NewZipService(fileRepository, config.BaseURL)
	cs, _ := service2.NewCryptoService()
	userService := service2.NewUserService(postgresRepository, config.BaseURL)
	//zip, _ := service2.NewCompressService()
	//h := handler.NewZipURLHandler(service)
	uh := handler.NewUserHandler(userService, service, cs)
	router := chi.NewRouter()
	router.Use(middleware.CleanPath)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(custommiddleware.Compress)
	router.Route("/", func(r chi.Router) {
		r.Get("/", uh.HelloHandler)
		r.Get("/{id}", uh.GetMethodHandler)
		r.Get("/user/urls", uh.GetUserURLsHandler)
		r.Get("/ping", uh.PingHandler)
		r.Post("/api/shorten", uh.PostAPIShortenHandler)
		r.Post("/api/shorten/batch", uh.PostShortenBatchHandler)
		r.Post("/", uh.PostMethodHandler)
		r.Put("/", uh.DefaultHandler)
		r.Patch("/", uh.DefaultHandler)
		r.Delete("/", uh.DefaultHandler)
	})

	err = http.ListenAndServe(config.ServerAddress, router)
	if err != nil {
		fmt.Println("can't start service")
		fmt.Println(err)
	}
}
