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
	SaveBatch(userID string, srcDTO []dto.UserBatchDTO) ([]dto.UserBatchResultDTO, error)
	GetURLByShort(shortURL string) (string, error)
	ZipURL(url string) (string, error)
	Ping() bool
}

type UserHandler struct {
	userService   UserService
	cryptoService CryptoService
}

func NewUserHandler(userService UserService, cs CryptoService) *UserHandler {
	var h UserHandler
	h.userService = userService
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
	var (
		userID string
		ok     bool
	)
	token, err := r.Cookie("token")
	if err == nil {
		ok, userID = z.cryptoService.Validate(token.Value)
	}
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
	res, err := z.userService.GetURLsByUser(userID)
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
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(responseBody)
		if err != nil {
			panic("Can't write response")
		}
	}
}

func (z *UserHandler) PingHandler(w http.ResponseWriter, r *http.Request) {
	if !z.userService.Ping() {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (z *UserHandler) PostMethodHandler(w http.ResponseWriter, r *http.Request) {
	b, err := getRequestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userID, err := z.getTokenCookie(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if string(b) == "" {
		writeBadRequest(w)
		return
	} else {
		res, _ := z.userService.ZipURL(string(b))
		err = z.userService.Save(userID, string(b), res)
		if errors.Is(err, dto.ErrDuplicateKey) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(res))
		if err != nil {
			panic("Can't write response")
		}
		return
	}
}

func (z *UserHandler) PostAPIShortenHandler(w http.ResponseWriter, r *http.Request) {
	b, err := getRequestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID, err := z.getTokenCookie(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if string(b) == "" {
		writeBadRequest(w)
		return
	} else {
		var req dto.ShortenRequestDTO
		if err := json.Unmarshal(b, &req); err != nil {
			writeBadRequest(w)
			return
		}
		res, err := z.userService.ZipURL(req.URL)
		if err != nil {
			writeBadRequest(w)
			return
		}

		err = z.userService.Save(userID, req.URL, res)
		if errors.Is(err, dto.ErrDuplicateKey) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resultDTO := dto.ShortenResponseDTO{Result: res}

		responseBody, err := json.Marshal(resultDTO)
		if err != nil {
			panic("Can't serialize response")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(responseBody)
		if err != nil {
			panic("Can't write response")
		}
		return
	}
}

func (z *UserHandler) PostShortenBatchHandler(w http.ResponseWriter, r *http.Request) {
	b, err := getRequestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userID, err := z.getTokenCookie(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var req []dto.UserBatchDTO
	if err := json.Unmarshal(b, &req); err != nil {
		writeBadRequest(w)
		return
	}
	resultDTO, err := z.userService.SaveBatch(userID, req)
	if errors.Is(err, dto.ErrDuplicateKey) {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(resultDTO)
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

func (z *UserHandler) GetMethodHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "" || r.RequestURI[1:] == "" {
		writeBadRequest(w)
		return
	} else {
		key := r.RequestURI[1:]
		res, err := z.userService.GetURLByShort(key)
		if err != nil {
			writeBadRequest(w)
			return
		}
		w.Header().Set("Location", res)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
}

func (z *UserHandler) HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	_, err := w.Write([]byte("Hello"))
	if err != nil {
		panic("Can't write response")
	}
}

func (z *UserHandler) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte("Unsupported request type"))
	if err != nil {
		panic("Can't write response")
	}
}
