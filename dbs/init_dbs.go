package dbs

import (
	"context"
	"os"

	cf "github.com/D-Watson/live-safety/conf"
	"github.com/D-Watson/live-safety/dbs"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	MysqlEngine *gorm.DB
	RedisEngine *redis.Client
)

func initMysql() error {
	var err error
	MysqlEngine, err = dbs.GetDBClient(cf.GlobalConfig.DB.Mysql)
	if err != nil {
		return err
	}
	return nil
}
func initRedis(ctx context.Context) error {
	var err error
	RedisEngine, err = dbs.InitRedisCli(ctx)
	if err != nil {
		return err
	}
	return nil
}

func InitDBS(ctx context.Context) {
	err := os.Setenv("JWT_SECRET", "de5ee29aa3561ab23901d8dad66b67e442ad4006")
	if err != nil {
		return
	}
	os.Setenv("SMTP_PASSWORD", "zgtgkjovcwpaceie")
	err = initRedis(ctx)
	if err != nil {
		return
	}
	err = initMysql()
	if err != nil {
		return
	}
}
