package handler

import (
	"errors"
	"net/http"
)

type CryptoService interface {
	Validate(token string) bool

	// function for user_id and encrypted token generation
	// returned values:  user_id, token, error
	GetNewUserToken() (string, string, error)
}

type UserService interface {
	GetURLsByUser(userID string) ([]string, error)
}

type UserHandler struct {
	service       UserService
	cryptoService CryptoService
}

func NewUserHandler(service UserService, cs CryptoService) *UserHandler {
	var h UserHandler
	h.service = service
	h.cryptoService = cs
	return &h
}

func (z *UserHandler) bakeCookie() (*http.Cookie, error) {
	var c http.Cookie
	// TODO. user
	_, token, err := z.cryptoService.GetNewUserToken()
	if err != nil {
		return nil, err
	}
	c.Name = "token"
	c.Value = token
	return &c, nil
}

func (z *UserHandler) GetUserURLsHandler(w http.ResponseWriter, r *http.Request) {
	// получить куку
	token, err := r.Cookie("token")
	if errors.Is(err, http.ErrNoCookie) {
		newToken, err := z.bakeCookie()
		if err != nil {

		}
		http.SetCookie(w, newToken)
	}
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if z.cryptoService.Validate(token.Value) {
		// TODO.
		user_id := "key"
		res, err := z.service.GetURLsByUser(user_id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	// проверить
	// Получить урлы
	// Если массив пустой, то отдаем 204

	return
}
