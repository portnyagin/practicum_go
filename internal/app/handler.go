package app

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Service interface {
	ZipURL(url string) (string, error)
	UnzipURL(key string) (string, error)
	ZipURLv2(url string) (*ShortenResponseDTO, error)
}

type ZipURLHandler struct {
	service Service
}

func NewZipURLHandler(service Service) *ZipURLHandler {
	var h ZipURLHandler
	h.service = service
	return &h
}

func (z *ZipURLHandler) DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte("Unsupported request type"))
	if err != nil {
		panic("Can't write response")
	}
}

func (z *ZipURLHandler) PostMethodHandler(w http.ResponseWriter, r *http.Request) {

	b, err := io.ReadAll(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if string(b) == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("Bad request"))
		if err != nil {
			panic("Can't write response")
		}
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
	b, err := io.ReadAll(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if string(b) == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("Bad request"))
		if err != nil {
			panic("Can't write response")
		}
		return
	} else {
		var req ShortenRequestDTO
		if err := json.Unmarshal(b, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte("Bad request"))
			if err != nil {
				panic("Can't write response")
			}
			return
		}
		resultDTO, _ := z.service.ZipURLv2(req.URL)
		responseBody, err := json.Marshal(resultDTO)
		if err != nil {
			panic("Can't serialize response")
		}
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write(responseBody)
		if err != nil {
			panic("Can't write response")
		}
		return
	}
}

func (z *ZipURLHandler) GetMethodHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "" || r.RequestURI[1:] == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		key := r.RequestURI[1:]
		res, err := z.service.UnzipURL(key)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				panic("Can't write response")
			}
			return
		}
		w.Header().Set("Location", res)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
}
