package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"strconv"
	"strings"
	"time"
)

type CryptoServiceImpl struct {
	key    []byte
	aesgcm cipher.AEAD
	nonce  []byte
}

func NewCryptoService() (*CryptoServiceImpl, error) {
	var cs CryptoServiceImpl
	// Здесь должен быть криптостойкий ключ, но пока хочется детерминированного поведения между запусками сервиса.
	cs.key = []byte(strings.Repeat("a", aes.BlockSize))

	aesblock, err := aes.NewCipher(cs.key)
	if err != nil {
		return nil, err
	}
	cs.aesgcm, err = cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	cs.nonce = make([]byte, cs.aesgcm.NonceSize())
	_, err = rand.Read(cs.nonce)
	if err != nil {
		return nil, err
	}

	return &cs, nil
}

func (s *CryptoServiceImpl) generateUserID() (string, error) {
	uid := strconv.FormatInt(time.Now().Unix(), 10)
	return uid, nil
}

// function for user_id and encrypted token generation
// returned values:  user_id, token, error
func (s *CryptoServiceImpl) GetNewUserToken() (string, string, error) {
	user, err := s.generateUserID()
	if err != nil {
		return "", "", nil
	}
	token, err := s.encrypt([]byte(user))
	if err != nil {
		return "", "", nil
	}
	return user, string(token), nil

}

// Function try to encrypt given token
// Return true and decrypted user_id. Else - false
func (s *CryptoServiceImpl) Validate(token string) (bool, string) {
	res, err := s.decrypt([]byte(token))
	if err != nil {
		return false, ""
	}
	return true, res
}

func (s *CryptoServiceImpl) decrypt(src []byte) (string, error) {
	res, err := s.aesgcm.Open(nil, s.nonce, src, nil) // расшифровываем
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (s *CryptoServiceImpl) encrypt(userID []byte) ([]byte, error) {
	dst := s.aesgcm.Seal(nil, s.nonce, userID, nil) // зашифровываем
	return dst, nil
}
