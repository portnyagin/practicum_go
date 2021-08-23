package app

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//type Service interface {
//	ZipURL(url string) (string, error)
//	UnzipURL(key string) (string, error)
//}

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

func TestZipURLHandler_postMethodHandler(t *testing.T) {
	service := new(ServiceMock)
	service.On("ZipURL", "full_URL").Return("short_URL", nil)
	service.On("ZipURL", "").Return("", errors.New("URL is empty"))
	service.On("UnzipURL", "short_URL").Return("full_URL", nil)
	service.On("UnzipURL", "xxx").Return("", errors.New("key not found"))
	handler := NewZipURLHandler(service)

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
		{name: "Post test #1 (Negotive). Empty body",
			args:  args{requestBody: ""},
			wants: wants{responseCode: http.StatusBadRequest, resultResponse: ""},
		},
		{name: "Post test #2 (Positive)",
			args:  args{requestBody: "full_URL"},
			wants: wants{responseCode: http.StatusCreated, resultResponse: "short_URL"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/", strings.NewReader(tt.args.requestBody))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(handler.postMethodHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			fmt.Println(res)

			if res.StatusCode != tt.wants.responseCode {
				t.Errorf("Expected status %d, got %d", tt.wants.responseCode, res.StatusCode)
			}
			if res.StatusCode == http.StatusCreated {
				responseBody, err := io.ReadAll(res.Body)
				defer res.Body.Close()
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
