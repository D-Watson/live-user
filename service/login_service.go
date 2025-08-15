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
	"live-user/rpc"
)

func LoginService(ctx context.Context, req *entity.LoginReq) (resp *consts2.BaseResp) {
	resp = &consts2.BaseResp{
		ErrCode: consts2.HTTP_OK,
	}
	//验证密码
	passwd, err := rpc.DecryptData(ctx, req.PasswordEncrypted, consts.DECRYPT_ROLE)
	if err != nil {
		resp.ErrCode = consts.NETWORK_RPC_ERROR
		return
	}
	passwd, _ = EncryptPasswd(passwd)
	en := &entity.Users{
		Username:     req.UserName,
		Email:        req.Email,
		PasswordHash: passwd,
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

//退出登录
