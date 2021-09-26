package handler

import (
	"encoding/json"
	"errors"
	"github.com/portnyagin/practicum_go/internal/app/dto"
	"net/http"
)

type CryptoService interface {
	Validate(token string) (bool, string)
	GetNewUserToken() (string, string, error)
}

type UserService interface {
	GetURLsByUser(userID string) ([]dto.UserURLsDTO, error)
	Save(userID string, originalURL string, shortURL string) error
	Ping() bool
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

func (z *UserHandler) bakeCookie() (*http.Cookie, string, error) {
	var c http.Cookie
	userID, token, err := z.cryptoService.GetNewUserToken()
	if err != nil {
		return nil, "", err
	}
	c.Name = "token"
	c.Value = token
	return &c, userID, nil
}

// function try to read cookie from request. If cookie isn't set or cookie isn't valid, then generate and set new cookie
// return  userID
func (z *UserHandler) getTokenCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	token, err := r.Cookie("token")
	ok, userID := z.cryptoService.Validate(token.Value)
	if errors.Is(err, http.ErrNoCookie) || !ok {
		var newToken *http.Cookie
		newToken, userID, err = z.bakeCookie()
		if err != nil {
			return "", err
		}
		http.SetCookie(w, newToken)
		token = newToken
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "", err
	}
	return userID, nil
}

func (z *UserHandler) GetUserURLsHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := z.getTokenCookie(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := z.service.GetURLsByUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(res) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	} else {
		// записать ответ
		// res -> DTO
		responseBody, err := json.Marshal(res)
		if err != nil {
			panic("Can't serialize response")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(responseBody)
		if err != nil {
			panic("Can't write response")
		}
	}
}

func (z *UserHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	if !z.service.Ping() {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
