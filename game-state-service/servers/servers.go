package servers

import (
	proto "game-state-service/proto/gen/v1"
	"google.golang.org/grpc"
)

func GRpcServers() *grpc.Server {
	s := grpc.NewServer()
	proto.RegisterGameStateServiceServer(s, &StateService{})
	return s
}
