package app

import (
	"io"
	"net/http"
)

type Service interface {
	ZipURL(url string) (string, error)
	UnzipURL(key string) (string, error)
}

type ZipURLHandler struct {
	service Service
}

func NewZipURLHandler() *ZipURLHandler {
	var h ZipURLHandler
	h.service = NewZipService()
	return &h
}

func (z *ZipURLHandler) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		z.getMethodHandler(w, r)
	case http.MethodPost:
		z.postMethodHandler(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported request type"))
	}

}

func (z *ZipURLHandler) postMethodHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if string(b) == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
		return
	} else {
		res, _ := z.service.ZipURL(string(b))
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(res))
		return
	}
}

func (z *ZipURLHandler) getMethodHandler(w http.ResponseWriter, r *http.Request) {
	key := r.RequestURI[1:]
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		res, err := z.service.UnzipURL(key)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Header().Set("Location", res)
		return
	}
}
