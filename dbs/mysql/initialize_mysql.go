package mysql

import (
	"context"
	"os"
	"time"

	lg "log"

	"github.com/go-sql-driver/mysql"
	gd "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"live-user/configs"
	"live-user/utils/log"
)

var (
	UserEngineDB *gorm.DB
)

func InitDB(ctx context.Context) {
	var err error
	UserEngineDB, err = GetDBClient(configs.GlobalConf.DBS.Mysql)
	if err != nil {
		log.Fatalf(ctx, "[DB] init failed, err=%v", err)
		panic(err)
	}
}

func GetDBClient(conf *configs.Mysql) (*gorm.DB, error) {
	myconf := mysql.NewConfig()
	myconf.User = conf.UserName
	myconf.Passwd = conf.Password
	myconf.Addr = conf.Address
	myconf.DBName = conf.DBName
	myconf.Net = "tcp"
	myconf.Loc = time.Now().Location()
	myconf.Params = map[string]string{
		"parseTime": "true",
	}
	if conf.Options.Timeout != 0 {
		myconf.Timeout = time.Duration(conf.Options.Timeout) * time.Millisecond
	}

	if conf.Options.ReadTimeout != 0 {
		myconf.ReadTimeout = time.Duration(conf.Options.ReadTimeout) * time.Millisecond
	}

	if conf.Options.WriteTimeout != 0 {
		myconf.WriteTimeout = time.Duration(conf.Options.WriteTimeout) * time.Millisecond
	}
	newLogger := logger.New(
		lg.New(os.Stdout, "\r\n", lg.LstdFlags), // 输出到控制台
		logger.Config{
			SlowThreshold: time.Second, // 慢查询阈值
			LogLevel:      logger.Info, // 日志级别：Info 会打印所有 SQL 和连接信息
			Colorful:      true,        // 彩色输出
		},
	)
	dsn := myconf.FormatDSN()
	cli, err := gorm.Open(gd.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, err
	}
	return cli, nil
}

var globalGormLogger = &gormLogger{}

var _ logger.Interface = new(gormLogger)

type gormLogger struct {
}

func (gl *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return gl
}
func (gl *gormLogger) Info(ctx context.Context, template string, args ...interface{}) {
	log.Infof(ctx, template, args)
}

func (gl *gormLogger) Warn(ctx context.Context, template string, args ...interface{}) {
	log.Warnf(ctx, template, args)
}

func (gl *gormLogger) Error(ctx context.Context, template string, args ...interface{}) {
	log.Errorf(ctx, template, args)
}

func (gl *gormLogger) Trace(ctx context.Context, begin time.Time,
	fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	cost := time.Since(begin).String()
	log.Infof(ctx, "sql[%v], rowsAffected=%v, cost=%v, err=%v", sql, rows, cost, err)
}
