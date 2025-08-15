package redis

import (
	"context"

	"github.com/D-Watson/live-safety/log"
	"live-user/consts"
	"live-user/dbs"
)

// InsertUserDeviceToken 存储设备token
func InsertUserDeviceToken(ctx context.Context, userId, deviceId, token string) error {
	key := consts.BuildTokenKey(userId, deviceId)
	deviceKey := consts.BuildDeviceKey(userId)
	err := dbs.RedisEngine.LPush(ctx, deviceKey, deviceId, consts.JWT_DURATION).Err()
	if err != nil {
		log.Errorf(ctx, "[Redis] set error, err=", err)
		return err
	}
	err = dbs.RedisEngine.Set(ctx, key, token, consts.JWT_DURATION).Err()
	if err != nil {
		log.Errorf(ctx, "[Redis] set error, err=", err)
		return err
	}
	return nil
}

// LenUserLoginDevice 判断是否超出登录设备数量限制
func LenUserLoginDevice(ctx context.Context, userId string) (int, error) {
	key := consts.BuildDeviceKey(userId)
	nums, err := dbs.RedisEngine.LLen(ctx, key).Result()
	if err != nil {
		log.Errorf(ctx, "[Redis] len error, err=", err)
		return int(nums), err
	}
	return int(nums), nil
}

// InsertEmailToken 存储邮箱验证码
func InsertEmailToken(ctx context.Context, email, code string) error {
	key := consts.BuildEmailKey(email)
	err := dbs.RedisEngine.Set(ctx, key, code, consts.EMAIL_TOKEN_EX).Err()
	if err != nil {
		log.Errorf(ctx, "[redis] set key=%s, err=", key, err)
		return err
	}
	return nil
}

// QueryEmailToken 查询验证码
func QueryEmailToken(ctx context.Context, email string) (string, error) {
	key := consts.BuildEmailKey(email)
	res, err := dbs.RedisEngine.Get(ctx, key).Result()
	if err != nil {
		log.Errorf(ctx, "[redis] get key=%s, err=", key, err)
		return res, err
	}
	return res, nil
}
