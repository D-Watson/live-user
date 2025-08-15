package consts

import (
	"fmt"
	"time"
)

const (
	USER_TOKEN    = "userid:%s:deviceId:%s"
	DEVICE_KEY    = "userid:%s"
	EMAIL_KEY     = "email:%s"
	JWT_SECRET    = "JWT_SECRET"
	SMTP_PASSWORD = "SMTP_PASSWORD"
	JWT_DURATION  = 30 * 24 * time.Hour
	//限制登录设备最多5台
	LIMIT_DEVICE_NUM = 5
	//邮箱验证码过期时间
	EMAIL_TOKEN_EX   = 5 * time.Minute
	LIMIT_SEND_EMAIL = 1 * time.Minute
)

func BuildTokenKey(userId, deviceId string) string {
	return fmt.Sprintf(USER_TOKEN, userId, deviceId)
}
func BuildDeviceKey(userId string) string {
	return fmt.Sprintf(DEVICE_KEY, userId)
}
func BuildEmailKey(email string) string {
	return fmt.Sprintf(EMAIL_KEY, email)
}

const (
	ENCRYPT_ROLE = 2
	DECRYPT_ROLE = 1
)

const (
	USER_PASSWD_ERROR    = 4001
	NETWORK_RPC_ERROR    = 4002
	GENERATE_TOKEN_ERROR = 4003
	DBS_ERROR            = 5001
	DEVICE_OVERLIMIT     = 5002
	EMAIL_SEND_ERROR     = 5003
)
