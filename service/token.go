package service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"live-user/consts"
)

func GenerateToken(userID, deviceId string, secretKey []byte, expiresIn time.Duration) (string, error) {
	// 设置Claims
	now := time.Now()
	ext := time.Now().Add(expiresIn)
	claims := Claims{
		UserID:   userID,
		DeviceID: deviceId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "live-user",             // 签发者
			Subject:   "user-auth",             // 主题
			ExpiresAt: jwt.NewNumericDate(ext), // 过期时间
			IssuedAt:  jwt.NewNumericDate(now), // 签发时间
			ID:        uuid.NewString(),        // JWT ID
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
		return []byte(os.Getenv(consts.JWT_SECRET)), nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

//多设备管理

type Claims struct {
	UserID   string `json:"userId"`
	DeviceID string `json:"deviceId"`
	jwt.RegisteredClaims
}
