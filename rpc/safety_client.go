package rpc

import (
	"context"
	"time"

	"github.com/D-Watson/live-safety/log"
	"live-user/proto"
)

func DecryptData(ctx context.Context, value string, role int) (decryptData string, err error) {
	c := proto.NewTransferSafeClient(RpcClient)
	decryptData = ""
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	res, err := c.SecureDecrypt(ctx, &proto.Data{
		Role:      int32(role),
		TransData: value,
	})
	if err != nil {
		log.Errorf(ctx, "[RPC] request error , err=", err)
		return
	}
	if len(res.DecryptData) == 0 {
		log.Errorf(ctx, "[RPC] request error, decryptData is nil")
		return
	}
	return
}
