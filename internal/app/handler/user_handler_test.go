package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/portnyagin/practicum_go/internal/app/dto"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
			wants: wants{responseCode: http.StatusOK, resultResponse: "full_URL"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/user/urls", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(userHandler.GetUserURLsHandler)
			request.AddCookie(&http.Cookie{Name: "token", Value: "user_id"})

			h.ServeHTTP(w, request)
			res := w.Result()

			defer res.Body.Close()
			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
		})
	}
}

func TestUserHandler_getTokenCookie(t *testing.T) {
	type args struct {
		w        http.ResponseWriter
		r        *http.Request
		userName string
	}
	tests := []struct {
		name string
		//fields  fields
		args args
		//want    *http.Cookie
		wantErr bool
	}{
		{
			name:    "Test #1 (getTokenCookie)",
			args:    args{w: httptest.NewRecorder(), r: httptest.NewRequest("GET", "/user/urls", nil), userName: "user_id"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userHandler.getTokenCookie(tt.args.w, tt.args.r)

			if (err != nil) != tt.wantErr {
				t.Errorf("getTokenCookie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// check returned cookie (name and value)
			assert.Equal(t, tt.args.userName, got, "Expected cookie.name is %s, got %s", got, tt.args.userName)
		})
	}
}

func TestUserHandler_getTokenCookieHeader(t *testing.T) {
	type args struct {
		w          http.ResponseWriter
		r          *http.Request
		cookieName string
		cookieVal  string
	}
	tests := []struct {
		name string
		//fields  fields
		args args
		//want    *http.Cookie
		wantErr bool
	}{
		{
			name:    "Test #2 (getTokenCookie. Check Header)",
			args:    args{w: httptest.NewRecorder(), r: httptest.NewRequest("GET", "/user/urls", nil), cookieName: "token", cookieVal: "valid_user_Token"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userHandler.getTokenCookie(tt.args.w, tt.args.r)

			if (err != nil) != tt.wantErr {
				t.Errorf("getTokenCookie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
			http.SetCookie(tt.args.w, &http.Cookie{Name: "some Namr", Value: "Somne value"})
			http.SetCookie(tt.args.w, &http.Cookie{Name: "some Namr", Value: "Somne value2"})
			http.SetCookie(tt.args.w, &http.Cookie{Name: "some Namr", Value: "Somne value3"})

			// check w for cookie
			tmp := tt.args.w.Header().Get("Set-Cookie")
			assert.NotEmpty(t, tmp, "Can't got cookie from response")
			parsedCookie := strings.Split(tmp, "=")
			assert.ElementsMatch(t, parsedCookie, []string{tt.args.cookieName, tt.args.cookieVal}, "ParsedCookie does not match expected")
		})
	}
}

func TestUserHandler_PostShortenBatchHandler(t *testing.T) {
	type args struct {
		requestBody string
	}
	type wants struct {
		responseCode int
		contentType  string
		responseBody string
	}
	tests := []struct {
		name  string
		wants wants
		args  args
	}{
		{name: "ShortenBatchHandler test#1",
			wants: wants{
				responseCode: http.StatusCreated,
				contentType:  "application/json",
				responseBody: "",
			},
			args: args{requestBody: "[{\"correlation_id\": \"correlation1\",\"original_URL\": \"original_URL_1\"}]"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody := []byte(tt.args.requestBody)

			request := httptest.NewRequest("POST", "/api/shorten/batch", bytes.NewReader(requestBody))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(userHandler.PostShortenBatchHandler)

			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)

		})
	}
}

func TestUserHandler_PostMethodHandler(t *testing.T) {
	type args struct {
		requestBody string
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
		{name: "POST test #1 (Negative). Empty body",
			args:  args{requestBody: ""},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: ""},
		},
		{name: "POST test #2 (Positive)",
			args:  args{requestBody: "original_URL"},
			wants: wants{responseCode: http.StatusCreated, resultResponse: "short_URL"},
		},
		{name: "POST test #4 (Negative)",
			args:  args{requestBody: "bad_URL"},
			wants: wants{responseCode: http.StatusConflict, resultResponse: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/", strings.NewReader(tt.args.requestBody))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(userHandler.PostMethodHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			fmt.Println(res)

			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)

			if res.StatusCode == http.StatusCreated {
				responseBody, err := io.ReadAll(res.Body)

				defer res.Body.Close()
				if err != nil {
					t.Errorf("Can't read response body, %e", err)
				}
				assert.Equal(t, tt.wants.resultResponse, string(responseBody), "Expected body is %s, got %s", tt.wants.resultResponse, string(responseBody))

			}
		})
	}
}

/****/
func TestUserHandler_postApiShortenHandler(t *testing.T) {
	type args struct {
		request *dto.ShortenRequestDTO
	}
	type wants struct {
		responseCode int
		response     string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{name: "POST test #1 (Positive)",
			args: args{request: &dto.ShortenRequestDTO{URL: "original_URL"}},
			wants: wants{responseCode: http.StatusCreated,
				response: "short_URL"},
		},
		{name: "POST test #2 (Empty body)",
			args: args{request: nil},
			wants: wants{responseCode: http.StatusBadRequest,
				response: ""},
		},
		{name: "POST test #3 (Object with empty url)",
			args: args{request: &dto.ShortenRequestDTO{URL: ""}},
			wants: wants{responseCode: http.StatusBadRequest,
				response: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requestBody []byte
			if tt.args.request != nil {
				requestBody, _ = json.Marshal(tt.args.request)
			}
			request := httptest.NewRequest("POST", "/api/shorten", bytes.NewReader(requestBody))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(userHandler.PostAPIShortenHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			fmt.Println(res)

			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)

			if res.StatusCode == http.StatusCreated {
				responseBody, err := io.ReadAll(res.Body)

				defer res.Body.Close()
				if err != nil {
					t.Errorf("Can't read response body, %e", err)
				}
				var resultDTO dto.ShortenResponseDTO
				if err := json.Unmarshal(responseBody, &resultDTO); err != nil {
					t.Error("Can't unmarshal dto", err)
				}
				assert.Equal(t, tt.wants.response, resultDTO.Result, "Expected body is %s, got %s", tt.wants.response, resultDTO.Result)
			}
		})
	}
}

func TestUserHandler_postApiShortenHandler2(t *testing.T) {
	type args struct {
		requestBody string
	}
	type wants struct {
		responseCode int
		response     string
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{name: "POST test #4 (Empty object)",
			args: args{requestBody: "{}"},
			wants: wants{responseCode: http.StatusBadRequest,
				response: ""},
		},
		{name: "POST test #4 (Empty object)",
			args: args{requestBody: ""},
			wants: wants{responseCode: http.StatusBadRequest,
				response: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/api/shorten", strings.NewReader(tt.args.requestBody))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(userHandler.PostAPIShortenHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			fmt.Println(res)

			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)

			if res.StatusCode == http.StatusCreated {
				responseBody, err := io.ReadAll(res.Body)

				defer res.Body.Close()
				if err != nil {
					t.Errorf("Can't read response body, %e", err)
				}
				var resultDTO dto.ShortenResponseDTO
				if err := json.Unmarshal(responseBody, &resultDTO); err != nil {
					t.Error("Can't unmarshal dto", err)
				}
				assert.Equal(t, tt.wants.response, resultDTO.Result, "Expected body is %s, got %s", tt.wants.response, resultDTO.Result)
			}
		})
	}
}

func TestUserHandler_getMethodHandler(t *testing.T) {
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
			wants: wants{responseCode: http.StatusTemporaryRedirect, resultResponse: "original_URL"},
		},
		{name: "GET test #2 (Negative).",
			args:  args{shortURLKey: ""},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", fmt.Sprintf("/%s", tt.args.shortURLKey), nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(userHandler.GetMethodHandler)
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

func TestUserHandler_DefaultHandler(t *testing.T) {
	type args struct {
		method string
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
		{name: "Other http method test #1.",
			args:  args{method: "PUT"},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: "Unsupported request type"},
		},
		{name: "Other http method test #2.",
			args:  args{method: "PATCH"},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: "Unsupported request type"},
		},
		{name: "Other http method test #3.",
			args:  args{method: "DELETE"},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: "Unsupported request type"},
		},
		{name: "Other http method test #4.",
			args:  args{method: "HEAD"},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: "Unsupported request type"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, "/", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(userHandler.DefaultHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
			responseBody, err := io.ReadAll(res.Body)

			if err != nil {
				t.Errorf("Can't read response body, %e", err)
			}
			assert.Equal(t, "Unsupported request type", string(responseBody), "Expected body is %s, got %s", tt.wants.resultResponse, string(responseBody))
		})
	}
}
