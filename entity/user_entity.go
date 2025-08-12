package entity

import (
	"time"
)

type LoginReq struct {
}

type RegisterReq struct {
}

// Users 用户基本信息表
type Users struct {
	Id            int64     `json:"id" gorm:"id"`                         // 用户ID，主键
	Username      string    `json:"username" gorm:"username"`             // 用户名，用于登录
	PasswordHash  string    `json:"passwordHash" gorm:"password_hash"`    // 加盐值加密后的密码
	Email         string    `json:"email" gorm:"email"`                   // 电子邮箱
	Phone         string    `json:"phone" gorm:"phone"`                   // 手机号码
	Status        int8      `json:"status" gorm:"status"`                 // 状态：0-禁用，1-正常
	Avatar        string    `json:"avatar" gorm:"avatar"`                 // 头像URL
	LastLoginTime time.Time `json:"lastLoginTime" gorm:"last_login_time"` // 最后登录时间
	LastLoginIp   string    `json:"lastLoginIp" gorm:"last_login_ip"`     // 最后登录IP
	CreatedTime   int64     `json:"createdTime" gorm:"created_time"`      // 创建时间
	UpdatedTime   int64     `json:"updatedTime" gorm:"updated_time"`      // 更新时间
}

// TableName 表名称
func (*Users) TableName() string {
	return "users"
}
