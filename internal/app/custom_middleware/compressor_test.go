package custom_middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func stubHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func TestCompressGET(t *testing.T) {
	type args struct {
		acceptEncoding string
		method         string
		pattern        string
	}
	type wants struct {
		responseCode    int
		contentType     string
		contentEncoding string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{name: "GET compress test #1 ",
			args:  args{acceptEncoding: "gzip", method: "GET", pattern: "/shortURL"},
			wants: wants{responseCode: http.StatusOK, contentType: "application/json", contentEncoding: "gzip"},
		},
		{name: "GET compress test #2 ",
			args:  args{acceptEncoding: "deflate", method: "GET", pattern: "/shortURL"},
			wants: wants{responseCode: http.StatusOK, contentType: "application/json", contentEncoding: ""},
		},
		{name: "GET compress test #3 ",
			args:  args{acceptEncoding: "", method: "GET", pattern: "/shortURL"},
			wants: wants{responseCode: http.StatusOK, contentType: "application/json", contentEncoding: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, tt.args.pattern, nil)
			request.Header.Set("Accept-Encoding", tt.args.acceptEncoding)
			h := Compress(http.HandlerFunc(stubHandler))
			w := httptest.NewRecorder()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
			assert.Equal(t, tt.wants.contentType, res.Header.Get("Content-Type"), "Expected Content-Type %d, got %d", tt.wants.responseCode, res.StatusCode)
			assert.Equal(t, tt.wants.contentEncoding, res.Header.Get("Content-Encoding"), "Expected Content-Encoding %d, got %d", tt.wants.responseCode, res.StatusCode)
		})
	}
}

// Пусть этот тест будет интеграционным
func TestCompressPOST(t *testing.T) {
	type args struct {
		acceptEncoding  string
		contentEncoding string
		method          string
		pattern         string
		body            string
	}
	type wants struct {
		responseCode    int
		contentType     string
		contentEncoding string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{name: "POST compress test #1 ",
			args:  args{acceptEncoding: "", contentEncoding: "", method: "POST", pattern: "/api/shorten", body: "{\"url\":\"full_URL\"}"},
			wants: wants{responseCode: http.StatusOK, contentType: "application/json", contentEncoding: ""},
		},
		{name: "POST compress test #2 ",
			args:  args{acceptEncoding: "gzip", contentEncoding: "gzip", method: "POST", pattern: "/api/shorten", body: "{\"url\":\"full_URL\"}"},
			wants: wants{responseCode: http.StatusOK, contentType: "application/json", contentEncoding: "gzip"},
		},
		{name: "POST compress test #3 ",
			args:  args{acceptEncoding: "", contentEncoding: "gzip", method: "POST", pattern: "/api/shorten", body: "{\"url\":\"full_URL\"}"},
			wants: wants{responseCode: http.StatusOK, contentType: "application/json", contentEncoding: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var zp bytes.Buffer
			compressor := gzip.NewWriter(&zp)
			_, err := compressor.Write([]byte(tt.args.body))
			if err != nil {
				panic(err)
			}
			err = compressor.Close()
			if err != nil {
				panic(err)
			}
			request := httptest.NewRequest(tt.args.method, tt.args.pattern, &zp)
			if tt.args.acceptEncoding != "" {
				request.Header.Set("Accept-Encoding", tt.args.acceptEncoding)
			}
			if tt.args.contentEncoding != "" {
				request.Header.Set("Content-Encoding", tt.args.acceptEncoding)
			}
			h := Compress(http.HandlerFunc(stubHandler))
			w := httptest.NewRecorder()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
			assert.Equal(t, tt.wants.contentType, res.Header.Get("Content-Type"), "Expected Content-Type %d, got %d", tt.wants.responseCode, res.StatusCode)
			assert.Equal(t, tt.wants.contentEncoding, res.Header.Get("Content-Encoding"), "Expected Content-Encoding %d, got %d", tt.wants.responseCode, res.StatusCode)
		})
	}
}
