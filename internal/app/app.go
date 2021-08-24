package app

import (
	"fmt"
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
	repo := NewBaseRepository()
	service := NewZipService(repo)
	h := NewZipURLHandler(service)
	// маршрутизация запросов обработчику
	http.HandleFunc("/", h.Handler)
	// запуск сервера с адресом localhost, порт 8080

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("can't start service")
		fmt.Println(err)
	}

}
