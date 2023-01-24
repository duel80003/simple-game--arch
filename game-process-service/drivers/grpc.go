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
	stateGrpcConn     *grpc.ClientConn
	stateGrpcConnOnce sync.Once
)

func InitStateGrpcConn() {
	stateGrpcConnOnce.Do(func() {
		addr := os.Getenv("STATE_MACHINE_ADDR")
		tools.Logger.Infof("state machine service addr: %s", addr)
		conn, err := grpc.Dial(addr, grpc.WithBlock(), grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second,
			Timeout:             time.Second * 3,
			PermitWithoutStream: true,
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			tools.Logger.Panicf("state grpc connection error: %s", addr)
		}
		stateGrpcConn = conn
	})
}

func CloseStateGrpcConn() {
	err := stateGrpcConn.Close()
	if err != nil {
		tools.Logger.Errorf("closing state grpc connection error: %s", err)
		return
	}
	tools.Logger.Infof("state grpc disconnected")
}

func GetStateGRpcConn() *grpc.ClientConn {
	return stateGrpcConn
}
