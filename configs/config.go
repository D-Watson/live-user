package configs

import (
	"context"
	"fmt"
	"net/url"

	"live-user/utils/log"

	"os"

	"gopkg.in/yaml.v3"
)

var GlobalConf *Config

type Config struct {
	DBS *Databases `yaml:"databases"`
}

type Databases struct {
	Mysql *Mysql `yaml:"mysql"`
	Redis *Redis `Yaml:"redis"`
}

type Mysql struct {
	UserName string  `yaml:"username"`
	Password string  `yaml:"password"`
	Address  string  `yaml:"address"`
	DBName   string  `yaml:"dbname"`
	Options  Options `yaml:"options"`
}
type Options struct {
	MaxIdleConns int `yaml:"max_idle_conns"`
	MaxOpenConns int `yaml:"max_open_conns"`
	Timeout      int `yaml:"timeout"`
	ReadTimeout  int `yaml:"readtimeout"`
	WriteTimeout int `yaml:"writetimeout"`
}
type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

func ParseConfig(ctx context.Context) (*Config, error) {
	data, err := os.ReadFile("./configs/config.yaml")
	if err != nil {
		log.Errorf(ctx, "[configs] read file error", err)
		return nil, err
	}
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Error(ctx, "解析配置失败: %v\n", err)
		return nil, err
	}
	return &config, nil
}

func InitConf(ctx context.Context) {
	var err error
	GlobalConf, err = ParseConfig(ctx)
	ews := url.QueryEscape(GlobalConf.DBS.Mysql.Password)
	fmt.Println(ews)
	if err != nil {
		log.Errorf(ctx, "[configs] parse error, ", err)
		panic(err)
	}
}
