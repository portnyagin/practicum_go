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
	"net/http/pprof"
	_ "net/http/pprof" // подключаем пакет pprof
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

	var (
		postgresHandler    *infrastructure.PostgresqlHandler
		postgresRepository *repository.PostgresRepository
		deleteRepository   *repository.DeleteRepository
	)

	if config.DatabaseDSN == "" {
		postgresHandler = nil
		postgresRepository = nil
		deleteRepository = nil
	} else {
		postgresHandler, err = infrastructure.NewPostgresqlHandler(context.Background(), config.DatabaseDSN)
		if err != nil {
			fmt.Println("can't init postgres handler", err)
			return
		}
		if config.Reinit {
			err = repository.ClearDatabase(context.Background(), postgresHandler)
			if err != nil {
				fmt.Println("can't clear database structure", err)
				return
			}
		}
		err = repository.InitDatabase(context.Background(), postgresHandler)
		if err != nil {
			fmt.Println("can't init database structure", err)
			return
		}
		postgresRepository, err = repository.NewPostgresRepository(postgresHandler)
		if err != nil {
			fmt.Println("can't init postgres repository", err)
			return
		}
		deleteRepository, err = repository.NewDeleteRepository(postgresHandler)
		if err != nil {
			fmt.Println("can't init delete repository", err)
			return
		}
	}

	cs, _ := service2.NewCryptoService()
	us := service2.NewUserService(postgresRepository, fileRepository, config.BaseURL)

	ds := service2.NewDeleteService(deleteRepository, config.DeletePoolSize, config.DeleteTaskSize)
	uh := handler.NewUserHandler(us, cs, ds)
	router := chi.NewRouter()
	router.Use(middleware.CleanPath)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(custommiddleware.Compress)
	router.Route("/", func(r chi.Router) {
		r.Get("/", uh.HelloHandler)
		r.Get("/{id}", uh.GetMethodHandler)
		r.Get("/api/user/urls", uh.GetUserURLsHandler)
		r.Get("/ping", uh.PingHandler)
		r.Post("/api/shorten", uh.PostAPIShortenHandler)
		r.Post("/api/shorten/batch", uh.PostShortenBatchHandler)
		r.Post("/", uh.PostMethodHandler)
		r.Put("/", uh.DefaultHandler)
		r.Patch("/", uh.DefaultHandler)
		r.Delete("/", uh.DefaultHandler)
		r.Delete("/api/user/urls", uh.AsyncDeleteHandler)
	})

	debugMux := http.NewServeMux()
	debugMux.HandleFunc("/debug/pprof/", pprof.Index)
	debugMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	debugMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	debugMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	debugMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", debugMux))
	}()

	err = http.ListenAndServe(config.ServerAddress, router)
	if err != nil {
		fmt.Println("can't start service")
		fmt.Println(err)
	}
}
