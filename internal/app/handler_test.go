package app

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type ServiceMock struct {
	mock.Mock
}

func (s *ServiceMock) ZipURL(url string) (string, error) {
	args := s.Called(url)
	return args.String(0), args.Error(1)
}

func (s *ServiceMock) UnzipURL(key string) (string, error) {
	args := s.Called(key)
	return args.String(0), args.Error(1)
}

var service *ServiceMock
var handler *ZipURLHandler

func TestMain(m *testing.M) {
	service = new(ServiceMock)
	service.On("ZipURL", "full_URL").Return("short_URL", nil)
	service.On("ZipURL", "").Return("", errors.New("URL is empty"))
	service.On("UnzipURL", "short_URL").Return("full_URL", nil)
	service.On("UnzipURL", "xxx").Return("", errors.New("key not found"))

	handler = NewZipURLHandler(service)
	os.Exit(m.Run())
}

func TestZipURLHandler_postMethodHandler(t *testing.T) {
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
			args:  args{requestBody: "full_URL"},
			wants: wants{responseCode: http.StatusCreated, resultResponse: "short_URL"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/", strings.NewReader(tt.args.requestBody))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.Handler)
			h.ServeHTTP(w, request)
			res := w.Result()
			fmt.Println(res)

			if res.StatusCode != tt.wants.responseCode {
				t.Errorf("Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
			}
			if res.StatusCode == http.StatusCreated {
				responseBody, err := io.ReadAll(res.Body)
				defer func() {
					err := res.Body.Close()
					if err != nil {
						log.Fatal(err)
					}
				}()
				if err != nil {
					t.Errorf("Can't read response body, %e", err)
				}
				if string(responseBody) != tt.wants.resultResponse {
					t.Errorf("Expected body is %s, got %s", tt.wants.resultResponse, string(responseBody))
				}
			}
		})
	}
}

func TestZipURLHandler_getMethodHandler(t *testing.T) {
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
		{name: "GET test #2 (Negative).",
			args:  args{shortURLKey: ""},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", fmt.Sprintf("/%s", tt.args.shortURLKey), nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.Handler)
			h.ServeHTTP(w, request)
			res := w.Result()
			fmt.Println(res)

			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)

			if res.StatusCode == tt.wants.responseCode {
				assert.Equal(t, tt.wants.resultResponse, res.Header.Get("Location"))
			}
		})
	}
}

func TestZipURLHandler_other(t *testing.T) {
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
		{name: "Other http method test #1 (Positive).",
			args:  args{method: "PUT"},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: "Unsupported request type"},
		},
		{name: "Other http method test #1 (Positive).",
			args:  args{method: "PATCH"},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: "Unsupported request type"},
		},
		{name: "Other http method test #1 (Positive).",
			args:  args{method: "DELETE"},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: "Unsupported request type"},
		},
		{name: "Other http method test #1 (Positive).",
			args:  args{method: "HEAD"},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: "Unsupported request type"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.args.method, "/", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.Handler)
			h.ServeHTTP(w, request)
			res := w.Result()
			fmt.Println(res)

			assert.Equal(t, tt.wants.responseCode, res.StatusCode, "Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)

			responseBody, err := io.ReadAll(res.Body)
			defer func() {
				err := res.Body.Close()
				if err != nil {
					log.Fatal(err)
				}
			}()
			if err != nil {
				t.Errorf("Can't read response body, %e", err)
			}
			assert.Equal(t, "Unsupported request type", string(responseBody), "Expected body is %s, got %s", tt.wants.resultResponse, string(responseBody))
		})
	}
}
