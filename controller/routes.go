package controller

import (
	"github.com/D-Watson/live-safety/conf"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	c := gin.Default()
	c.Use(CORSMiddleware())
	c.POST("/user/login", Login)
	c.POST("/user/email/send", SendEmailCode)
	err := c.Run(conf.GlobalConfig.Server.Http.Host)
	if err != nil {
		return
	}
}
