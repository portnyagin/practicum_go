package service

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestZipService_ZipURL(t *testing.T) {
	repo := new(RepositoryMock)
	repo.On("Find", "short_URL").Return("full_URL")
	repo.On("Find", "").Return("full_URL")
	repo.On("Save", "full_URL").Return("short_URL")

	type want struct {
		path   string
		scheme string
		host   string
	}
	tests := []struct {
		name    string
		url     string
		want    want
		wantErr bool
	}{
		{
			name: "ZipURL Test 1",
			url:  "full_URL",
			want: want{
				path:   "/full_URL",
				scheme: "http",
				host:   "localhost:8080",
			},
			wantErr: false,
		},
		{
			name: "ZipURL Test 2",
			url:  "",
			want: want{
				path:   "/short_URL",
				scheme: "http",
				host:   "localhost:8080",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewZipService(repo, "http://localhost:8080/")
			s.encode = mockEncode
			res, err := s.ZipURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZipURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			u, err := url.Parse(res)
			if err != nil {
				t.Errorf("ZipURL() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				assert.Equal(t, tt.want.scheme, u.Scheme)
				assert.Equal(t, tt.want.host, u.Host)
				assert.Equal(t, tt.want.path, u.Path)
			}
		})
	}
}

func TestZipService_UnzipURL(t *testing.T) {
	repo := new(RepositoryMock)
	repo.On("Find", "short_URL").Return("full_URL", nil)
	repo.On("Save", "full_URL").Return("short_URL", nil)

	tests := []struct {
		name    string
		key     string
		want    string
		wantErr bool
	}{
		{
			name:    "UnzipURL Test 1",
			key:     "short_URL",
			want:    "full_URL",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewZipService(repo, "http://localhost:8080/")

			got, err := s.UnzipURL(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZipURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ZipURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
