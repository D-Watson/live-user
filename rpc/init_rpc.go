package rpc

import (
	"context"

	cf "github.com/D-Watson/live-safety/conf"
	"google.golang.org/grpc"
)

var RpcClient *grpc.ClientConn

func InitRpc(ctx context.Context) {
	conn, err := grpc.NewClient(cf.GlobalConfig.Server.Rpc.ServerHost)
	if err != nil {
		return
	}
	RpcClient = conn
}
