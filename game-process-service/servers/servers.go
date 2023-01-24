package servers

import (
	proto "game-process-service/proto/gen/v1"
	"google.golang.org/grpc"
)

func GRpcServers() *grpc.Server {
	s := grpc.NewServer()
	proto.RegisterGameProcessServiceServer(s, &GameService{})
	return s
}
