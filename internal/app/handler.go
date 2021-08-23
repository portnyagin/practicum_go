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

func NewZipURLHandler(service Service) *ZipURLHandler {
	var h ZipURLHandler
	h.service = service
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
		_, err := w.Write([]byte("Unsupported request type"))
		if err != nil {
			panic("Can't write response")
		}
	}

}

func (z *ZipURLHandler) postMethodHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
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
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				panic("Can't write response")
			}
			return
		}
		w.Header().Set("Location", res)
		w.WriteHeader(http.StatusTemporaryRedirect)
		//w.Write([]byte(res))
		return
	}
}
