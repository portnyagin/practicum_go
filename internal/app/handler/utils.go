package handler

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func unZip(data []byte) ([]byte, error) {
	var res bytes.Buffer // тут будет результат
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	_, err = res.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	err = r.Close()
	if err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}

func getRequestBody(r *http.Request) ([]byte, error) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}
	if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
		unzipBody, err := unZip(b)
		if err != nil {
			panic(err)
		}
		return unzipBody, nil
	}
	return b, nil
}

func writeBadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte("Bad request"))
	if err != nil {
		panic("Can't write response")
	}
}
