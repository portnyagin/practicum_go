package handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
			wants: wants{responseCode: http.StatusTemporaryRedirect, resultResponse: "full_URL"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/user/urls", nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(userHandler.GetUserURLsHandler)
			request.AddCookie(&http.Cookie{Name: "token", Value: "dfdoskfojfpskf"})

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
			args:    args{w: httptest.NewRecorder(), r: httptest.NewRequest("GET", "/user/urls", nil), userName: "userID1"},
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
