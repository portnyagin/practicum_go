package handler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestZipURLHandler_GetUserURLsHandler(t *testing.T) {
	type args struct {
		shortURLKey string
	}
	type wants struct {
		responseCode   int
		resultResponse string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{name: "GET test #1 (Positive).",
			args:  args{shortURLKey: "short_URL"},
			wants: wants{responseCode: http.StatusTemporaryRedirect, resultResponse: "full_URL"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/user/urls", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(userHandler.GetUserURLsHandler)
			request.AddCookie(&http.Cookie{Name: "token", Value: "dfdoskfojfpskf"})
			//http.SetCookie(w,&http.Cookie{Name:"user_id", Value: "dfdoskfojfpskf"})

			h.ServeHTTP(w, request)
			res := w.Result()

			defer res.Body.Close()
			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)

			if res.StatusCode == tt.wants.responseCode {
				assert.Equal(t, tt.wants.resultResponse, res.Header.Get("Location"))
			}
		})
	}
}
