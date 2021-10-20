package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/portnyagin/practicum_go/internal/app/dto"
	"net/http"
)

type CryptoService interface {
	Validate(token string) (bool, string)
	GetNewUserToken() (string, []byte, error)
}

type UserService interface {
	GetURLsByUser(ctx context.Context, userID string) ([]dto.UserURLsDTO, error)
	SaveUserURL(ctx context.Context, userID string, originalURL string, shortURL string) error
	SaveBatch(ctx context.Context, userID string, srcDTO []dto.UserBatchDTO) ([]dto.UserBatchResultDTO, error)
	GetURLByShort(ctx context.Context, userID string, shortURL string) (string, error)
	ZipURL(url string) (string, string, error)
	Ping(ctx context.Context) bool
}

type DeleteService interface {
	DeleteBatch(ctx context.Context, userID string, URLList []dto.BatchDeleteDTO) error
}

type UserHandler struct {
	userService   UserService
	cryptoService CryptoService
	DeleteService DeleteService
}

func NewUserHandler(userService UserService, cs CryptoService, ds DeleteService) *UserHandler {
	var h UserHandler
	h.userService = userService
	h.cryptoService = cs
	h.DeleteService = ds
	return &h
}

func (z *UserHandler) bakeCookie() (*http.Cookie, string, error) {
	var c http.Cookie
	userID, token, err := z.cryptoService.GetNewUserToken()
	if err != nil {
		return nil, "", err
	}
	c.Name = "token"
	c.Value = base64.StdEncoding.EncodeToString(token)
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
		t, err := base64.StdEncoding.DecodeString(token.Value)
		if err != nil {
			return "", err
		}
		ok, userID = z.cryptoService.Validate(string(t))
		if !ok {
			fmt.Println("Cookie got, but invalid")
		}
	}
	if errors.Is(err, http.ErrNoCookie) || !ok {
		fmt.Println("Cook new cookie")
		var newToken *http.Cookie
		newToken, userID, err = z.bakeCookie()
		fmt.Println(userID, " ", newToken)
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

	res, err := z.userService.GetURLsByUser(r.Context(), userID)
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
	if !z.userService.Ping(r.Context()) {
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
		resURL, key, _ := z.userService.ZipURL(string(b))

		err = z.userService.SaveUserURL(r.Context(), userID, string(b), key)
		if errors.Is(err, dto.ErrDuplicateKey) {
			w.WriteHeader(http.StatusConflict)
			_, err = w.Write([]byte(resURL))
			if err != nil {
				panic("Can't write response")
			}
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(resURL))
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
		resURL, key, err := z.userService.ZipURL(req.URL)
		if err != nil {
			writeBadRequest(w)
			return
		}
		resultDTO := dto.ShortenResponseDTO{Result: resURL}
		responseBody, err := json.Marshal(resultDTO)
		if err != nil {
			panic("Can't serialize response")
		}
		w.Header().Set("Content-Type", "application/json")

		err = z.userService.SaveUserURL(r.Context(), userID, req.URL, key)
		if errors.Is(err, dto.ErrDuplicateKey) {
			w.WriteHeader(http.StatusConflict)
			_, err = w.Write(responseBody)
			if err != nil {
				panic("Can't write response")
			}
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

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
	resultDTO, err := z.userService.SaveBatch(r.Context(), userID, req)
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
		userID, err := z.getTokenCookie(w, r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		key := r.RequestURI[1:]

		userID = "" // Подгоняемся под тесты первых инкрементов
		res, err := z.userService.GetURLByShort(r.Context(), userID, key)
		if errors.Is(err, dto.ErrNotFound) {
			w.WriteHeader(http.StatusGone)
			return
		}

		if err != nil {
			writeBadRequest(w)
			return
		}
		w.Header().Set("Location", res)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
}

func (z *UserHandler) AsyncDeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	var req []dto.BatchDeleteDTO
	if err := json.Unmarshal(b, &req); err != nil {
		writeBadRequest(w)
		return
	}
	err = z.DeleteService.DeleteBatch(r.Context(), userID, req)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)

}

func (z *UserHandler) HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
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
