package handler

import (
	"errors"
	"github.com/portnyagin/practicum_go/internal/app/dto"
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
	service.On("ZipURL", "original_URL").Return("short_URL", nil)
	service.On("ZipURL", "bad_URL").Return("short_URL", nil)
	service.On("ZipURL", "").Return("", errors.New("URL is empty"))
	service.On("ZipURLv2", "full_URL").Return("short_URL", nil)
	service.On("ZipURLv2", "").Return("short_URL", errors.New("URL is empty"))
	service.On("UnzipURL", "short_URL").Return("full_URL", nil)
	service.On("UnzipURL", "xxx").Return("", errors.New("key not found"))

	userService = new(UserServiceMock)
	userService.On("GetURLsByUser", "user_id").Return("url-for-user-1", nil)

	var d []dto.UserBatchDTO
	d = append(d, dto.UserBatchDTO{CorrelationID: "correlation1", OriginalURL: "original_url_1"})
	userService.On("SaveBatch", "user_id", d).Return("correlation1", "short_URL_1", nil)

	userService.On("Save", "user_id", "original_URL", "short_URL").Return(nil)
	userService.On("Save", "user_id", "bad_URL", "short_URL").Return(dto.ErrDuplicateKey)

	handler = NewZipURLHandler(service)

	cryptoService := new(CryptoServiceMock)
	cryptoService.On("Validate", "user_id").Return(true, "user_id")

	cryptoService.On("GetNewUserToken").Return("user_id", "valid_user_Token", nil)

	userHandler = NewUserHandler(userService, service, cryptoService)
	os.Exit(m.Run())
}
