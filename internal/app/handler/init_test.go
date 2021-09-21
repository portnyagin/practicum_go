package handler

import (
	"errors"
	service2 "github.com/portnyagin/practicum_go/internal/app/service"
	"os"
	"testing"
)

var service *ServiceMock
var userService *UserServiceMock
var cryptoService *CryptoServiceMock
var handler *ZipURLHandler
var userHandler *UserHandler

func TestMain(m *testing.M) {
	service = new(ServiceMock)
	service.On("ZipURL", "full_URL").Return("short_URL", nil)
	service.On("ZipURL", "").Return("", errors.New("URL is empty"))
	service.On("ZipURLv2", "full_URL").Return("short_URL", nil)
	service.On("ZipURLv2", "").Return("short_URL", errors.New("URL is empty"))
	service.On("UnzipURL", "short_URL").Return("full_URL", nil)
	service.On("UnzipURL", "xxx").Return("", errors.New("key not found"))

	userService = new(UserServiceMock)
	userService.On("GetURLsByUser", "user1").Return("url-for-user-1", nil)

	//GetURLsByUser (userID string) ([]string, error)
	handler = NewZipURLHandler(service)

	// TODO: заменить на mock
	cryptoService, _ := service2.NewCryptoService()
	userHandler = NewUserHandler(userService, cryptoService)
	os.Exit(m.Run())
}
