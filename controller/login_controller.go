package controller

import (
	"context"
	"net/http"

	consts2 "github.com/D-Watson/live-safety/consts"
	"github.com/D-Watson/live-safety/log"
	"github.com/gin-gonic/gin"
	"live-user/consts"
	"live-user/entity"
	"live-user/service"
)

// Login 登录
func Login(c *gin.Context) {
	ctx := context.Background()
	req := &entity.LoginReq{}
	resp := &consts2.BaseResp{}
	if err := c.ShouldBindJSON(req); err != nil {
		log.Errorf(context.Background(), "[Param] analyse error, err=", err)
		c.AsciiJSON(http.StatusBadRequest, resp)
		return
	}
	resp = service.LoginService(ctx, req)
	c.AsciiJSON(http.StatusOK, resp)
}

// Register 注册
func Register(c *gin.Context) {
	//ctx := context.Background()
	req := &entity.RegisterReq{}
	resp := &consts2.BaseResp{}
	if err := c.ShouldBindJSON(req); err != nil {
		log.Errorf(context.Background(), "[Param] analyse error, err=", err)
		c.AsciiJSON(http.StatusBadRequest, resp)
		return
	}
}

// SendEmailCode 发送邮箱验证码
func SendEmailCode(c *gin.Context) {
	ctx := context.Background()
	req := &entity.SendCodeReq{}
	res := &entity.SendCodeResp{}
	resp := &consts2.BaseResp{
		ErrCode: consts.EMAIL_SEND_ERROR,
		Data:    res,
	}
	if err := c.ShouldBindJSON(req); err != nil {
		log.Errorf(context.Background(), "[Param] analyse error, err=", err)
		c.AsciiJSON(http.StatusBadRequest, resp)
		return
	}
	res = service.SendEmail(ctx, req)
	if res.SendSucc {
		resp.ErrCode = consts2.HTTP_OK
		resp.Data = res
	}
	c.AsciiJSON(http.StatusOK, resp)
	return
}

// VerifyEmailCode 核对验证码
func VerifyEmailCode(c *gin.Context) {

}
