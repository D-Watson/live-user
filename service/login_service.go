package service

import (
	"context"
	"os"
	"strconv"

	consts2 "github.com/D-Watson/live-safety/consts"
	"live-user/consts"
	"live-user/dbs/mysql"
	"live-user/dbs/redis"
	"live-user/entity"
)

func Login(ctx context.Context, req *entity.LoginReq) (resp *consts2.BaseResp) {
	resp = &consts2.BaseResp{
		ErrCode: consts2.HTTP_OK,
	}
	if !req.PasswdVerify && !verifyCode(ctx, req.Email, req.Code) {
		resp.ErrCode = consts.INVALID_CODE
		return
	}
	en := &entity.Users{
		Username: req.UserName,
		Email:    req.Email,
	}
	if req.PasswdVerify {
		passwd := getPasswd(ctx, req.PasswordEncrypted)
		en.PasswordHash = passwd
	}
	user, err := mysql.QueryUser(ctx, en)
	if err != nil {
		resp.ErrCode = consts.USER_PASSWD_ERROR
		return
	}
	userId := strconv.FormatInt(user.Id, 10)
	//设备限制, redis error直接放行
	num, _ := redis.LenUserLoginDevice(ctx, userId)
	if num >= consts.LIMIT_DEVICE_NUM {
		resp.ErrCode = consts.DEVICE_OVERLIMIT
		return
	}
	data := user.BuildToLoginResp()
	//第一次登录返回token
	generateToken, err := GenerateToken(userId, "", []byte(os.Getenv(consts.JWT_SECRET)), consts.JWT_DURATION)
	if err != nil {
		resp.ErrCode = consts.GENERATE_TOKEN_ERROR
		return
	}
	data.Token = generateToken
	resp.Data = data
	err = redis.InsertUserDeviceToken(ctx, userId, req.DevieceId, generateToken)
	if err != nil {
		resp.ErrCode = consts.DBS_ERROR
		return
	}
	return
}

func getPasswd(ctx context.Context, passwd string) string {
	//验证密码
	//passwd, err := rpc.DecryptData(ctx, passwd, consts.DECRYPT_ROLE)
	//if err != nil {
	//	return ""
	//}
	passwd, _ = EncryptPasswd(passwd)
	return passwd
}

func verifyCode(ctx context.Context, email, code string) bool {
	//1. 校验邮箱验证码
	verifiedCode, err := redis.QueryEmailToken(ctx, email)
	if err != nil {
		return false
	}
	if code != verifiedCode {
		return false
	}
	return true
}

//退出登录
