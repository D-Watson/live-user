package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(userID, username string, secretKey []byte, expiresIn time.Duration) (string, error) {
	// 设置Claims
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "my-app",                                      // 签发者
			Subject:   "user-auth",                                   // 主题
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                // 签发时间
			ID:        uuid.NewString(),                              // JWT ID
		},
	}

	// 创建Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获取完整Token
	return token.SignedString(secretKey)
}

// 加密

func EncryptPasswd(passwd string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(passwd), 256)
	if err != nil {

		return "", err
	}
	return string(hashedBytes), nil
}

func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// ParseToken 解析JWT令牌
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

type Claims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}
