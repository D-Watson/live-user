package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// 加密

func EncryptPasswd(passwd string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(passwd), 10)
	if err != nil {

		return "", err
	}
	return string(hashedBytes), nil
}

func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
