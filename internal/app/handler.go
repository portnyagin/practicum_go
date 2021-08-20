package app

import (
	"io"
	"net/http"
)

type Service interface {
	ZipUrl(url string) (string, error)
	UnzipUrl(key string) (string, error)
}

type ZipUrlHandler struct {
	service Service
}

func NewZipUrlHandler() *ZipUrlHandler {
	var h ZipUrlHandler
	h.service = NewZipService()
	return &h
}

func (z *ZipUrlHandler) Handler(w http.ResponseWriter, r *http.Request) {
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

func (z *ZipUrlHandler) postMethodHandler(w http.ResponseWriter, r *http.Request) {
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
		res, _ := z.service.ZipUrl(string(b))
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(res))
		return
	}
}

func (z *ZipUrlHandler) getMethodHandler(w http.ResponseWriter, r *http.Request) {
	key := r.RequestURI[1:]
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		res, err := z.service.UnzipUrl(key)
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
