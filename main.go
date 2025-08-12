package main

import (
	"context"

	"live-user/configs"
	"live-user/dbs/mysql"
)

func main() {
	ctx := context.Background()
	configs.InitConf(ctx)
	mysql.InitDB(ctx)
}
