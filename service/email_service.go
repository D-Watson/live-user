package service

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/D-Watson/live-safety/log"
	"gopkg.in/gomail.v2"
	"live-user/consts"
	"live-user/dbs/redis"
	"live-user/entity"
)

func SendEmail(ctx context.Context, req *entity.SendCodeReq) (res *entity.SendCodeResp) {
	res = &entity.SendCodeResp{
		SendSucc: false,
	}
	lock := redis.NewEmailSendLock(req.Email)
	defer func(lock *redis.RedisLock, ctx context.Context) {
		err := lock.Release(ctx)
		if err != nil {
			log.Errorf(ctx, "[redisLock] release error,err=", err)
		}
	}(lock, ctx)
	if ok, _ := lock.Acquire(ctx); !ok {
		return
	}
	err := sendEmail(req.Email)
	if err != nil {
		log.Errorf(ctx, "[Email] send err=", err)
		return
	}
	res.SendSucc = true
	return
}

func generateCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func sendEmail(to string) error {
	code := generateCode()
	// 创建邮件对象
	m := gomail.NewMessage()
	m.SetHeader("From", "3833340167@qq.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "您的验证码")
	m.SetBody("text/plain", fmt.Sprintf("验证码：%s，5分钟内有效", code))

	d := gomail.NewDialer("smtp.qq.com", 465, "3833340167@qq.com", os.Getenv(consts.SMTP_PASSWORD))
	return d.DialAndSend(m)
}
