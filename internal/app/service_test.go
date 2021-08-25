package app

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/url"
	"testing"
)

/*


func TestZipService_UnZipUrl(t *testing.T) {
	service := NewZipService()
	key, _ := service.ZipURL("ya.ru")
	u,_ := url.Parse(key)

	res, _ := service.UnzipURL(u.Path[1:])
	if res != "ya.ru" {
		t.Errorf("UnZipUrl test failed.Expected %s, got %s","ya.ru", res )
	}

	res, _ = service.UnzipURL("fake key")
	if res != "" {
		t.Errorf("UnZipUrl test failed.Expected empty string, got %s", res )
	}
}
*/

type RepositoryMock struct {
	mock.Mock
}

func (r *RepositoryMock) Find(key string) (string, error) {
	args := r.Called(key)
	return args.String(0), nil
}

func (r *RepositoryMock) Save(key string, value string) error {
	return nil
}

func mockEncode(str string) string {
	return str
}

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
			s := NewZipService(repo)
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
		{
			name:    "UnzipURL Test 2",
			key:     "",
			want:    "short_URL",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewZipService(repo)

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
