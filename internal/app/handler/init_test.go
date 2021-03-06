package handler

import (
	"errors"
	"github.com/portnyagin/practicum_go/internal/app/dto"
	"os"
	"testing"
)

var userService *UserServiceMock
var cryptoService *CryptoServiceMock
var userHandler *UserHandler
var deleteService *DeleteServiceMock

func TestMain(m *testing.M) {
	userService = new(UserServiceMock)
	userService.On("ZipURL", "full_URL").Return("short_URL", "short_URL", nil)
	userService.On("ZipURL", "original_URL").Return("short_URL", "short_URL", nil)
	userService.On("ZipURL", "bad_URL").Return("short_URL", "short_URL", nil)
	userService.On("ZipURL", "").Return("", "", errors.New("URL is empty"))

	userService.On("GetURLsByUser", "user_id").Return("url-for-user-1", nil)

	var d []dto.UserBatchDTO
	d = append(d, dto.UserBatchDTO{CorrelationID: "correlation1", OriginalURL: "original_URL_1"})
	userService.On("SaveBatch", "user_id", d).Return("correlation1", "short_URL_1", nil)

	userService.On("SaveUserURL", "user_id", "original_URL", "short_URL").Return(nil)
	userService.On("SaveUserURL", "user_id", "bad_URL", "short_URL").Return(dto.ErrDuplicateKey)
	userService.On("GetURLByShort", "user_id", "short_URL").Return("original_URL", nil)
	userService.On("GetURLByShort", "", "short_URL").Return("original_URL", nil)
	userService.On("GetURLByShort", "user_id", "badURL").Return("", dto.ErrNotFound)
	userService.On("GetURLByShort", "", "badURL").Return("", dto.ErrNotFound)

	cryptoService := new(CryptoServiceMock)
	cryptoService.On("Validate", "user_id").Return(true, "user_id")

	cryptoService.On("GetNewUserToken").Return("user_id", "valid_user_Token", nil)

	deleteService = new(DeleteServiceMock)

	userHandler = NewUserHandler(userService, cryptoService, deleteService)
	os.Exit(m.Run())
}
