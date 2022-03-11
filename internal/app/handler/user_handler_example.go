package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Пример вызова с токеном.
func Example() {
	// Создаем  http  клиент.
	client := http.Client{}
	// Делаем первый запрос (не авторизщованный). Сохраняем сокращенный урл, получаем  cookie.
	req, err := http.NewRequest("POST", "http://localhost:8080", strings.NewReader("original_url"))
	req.Header.Add("Content-Type", "text/plain")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	cookies := resp.Cookies()
	var cookieToken *http.Cookie
	for _, c := range cookies {
		if c.Name == "token" {
			cookieToken = c
			break
		}
		fmt.Println(c)
	}

	// Второй запрос с авторизацией. Получаем все  url пользователя.
	req2, err := http.NewRequest("GET", "http://localhost:8080/api/user/urls", nil)
	req2.AddCookie(cookieToken)
	respPOST, err := client.Do(req2)
	if err != nil {
		fmt.Println(errors.Unwrap(err))
	}
	defer respPOST.Body.Close()

	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}
