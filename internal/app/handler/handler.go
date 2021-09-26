package handler

import (
	"encoding/json"
	"github.com/portnyagin/practicum_go/internal/app/dto"
	"net/http"
)

type Service interface {
	ZipURL(url string) (string, error)
	UnzipURL(key string) (string, error)
}

type ZipURLHandler struct {
	service       Service
	cryptoService CryptoService
}

func NewZipURLHandler(service Service) *ZipURLHandler {
	var h ZipURLHandler
	h.service = service
	return &h
}

func (z *ZipURLHandler) HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	_, err := w.Write([]byte("Hello"))
	if err != nil {
		panic("Can't write response")
	}
}

func (z *ZipURLHandler) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte("Unsupported request type"))
	if err != nil {
		panic("Can't write response")
	}
}

func (z *ZipURLHandler) PostMethodHandler(w http.ResponseWriter, r *http.Request) {
	b, err := getRequestBody(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if string(b) == "" {
		writeBadRequest(w)
		return
	} else {
		res, _ := z.service.ZipURL(string(b))
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(res))
		if err != nil {
			panic("Can't write response")
		}
		return
	}
}

func (z *ZipURLHandler) PostAPIShortenHandler(w http.ResponseWriter, r *http.Request) {
	b, err := getRequestBody(r)
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
		res, err := z.service.ZipURL(req.URL)
		if err != nil {
			writeBadRequest(w)
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

// По хорошему надо отключить, так как выхов не авторизованный. Но непонятно что будет  с тестами
func (z *ZipURLHandler) GetMethodHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "" || r.RequestURI[1:] == "" {
		writeBadRequest(w)
		return
	} else {
		key := r.RequestURI[1:]
		res, err := z.service.UnzipURL(key)
		if err != nil {
			writeBadRequest(w)
			return
		}
		w.Header().Set("Location", res)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
}
