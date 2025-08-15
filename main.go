package main

import (
	"context"

	cf "github.com/D-Watson/live-safety/conf"
	"live-user/controller"
	"live-user/dbs"
	"live-user/rpc"
)

func main() {
	ctx := context.Background()
	err := cf.ParseConfig(ctx)
	if err != nil {
		return
	}
	dbs.InitDBS(ctx)
	go rpc.InitRpc(ctx)
	controller.InitRouter()
}
