package handler

import (
	"errors"
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

	cryptoService := new(CryptoServiceMock)
	cryptoService.On("Validate").Return(true, "user_id")

	cryptoService.On("GetNewUserToken").Return("user_id", "valid_user_Token", nil)

	userHandler = NewUserHandler(userService, cryptoService)
	os.Exit(m.Run())
}
