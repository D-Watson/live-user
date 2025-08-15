package entity

import (
	"time"
)

type TokenEntity struct {
	UserID   string    `json:"userId"`
	DeviceID string    `json:"deviceId"`
	ExpireAt time.Time `json:"expireAt"`
	Iat      time.Time `json:"iat"` //签发时间
}
