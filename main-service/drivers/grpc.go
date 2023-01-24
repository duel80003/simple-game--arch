package drivers

import (
	tools "github.com/duel80003/my-tools"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"os"
	"sync"
	"time"
)

var (
	gameGrpcConn     *grpc.ClientConn
	gameGrpcConnOnce sync.Once
)

func InitGameGrpcConn() {
	gameGrpcConnOnce.Do(func() {
		addr := os.Getenv("GAME_SERVER_ADDR")
		tools.Logger.Infof("game service service addr: %s", addr)
		conn, err := grpc.Dial(addr, grpc.WithBlock(), grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             time.Second * 3,
			PermitWithoutStream: true,
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			tools.Logger.Panicf("game service grpc connection error: %s", addr)
		}
		gameGrpcConn = conn
	})
}

func CloseGameGrpcConn() {
	err := gameGrpcConn.Close()
	if err != nil {
		tools.Logger.Errorf("closing state grpc connection error: %s", err)
		return
	}
	tools.Logger.Infof("state grpc disconnected")
}

func GetGsGRpcConn() *grpc.ClientConn {
	return gameGrpcConn
}
