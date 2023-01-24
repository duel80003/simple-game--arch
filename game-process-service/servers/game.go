package servers

import (
	"context"
	"game-process-service/game"
	proto "game-process-service/proto/gen/v1"
)

type GameService struct{}

func (g GameService) Join(_ context.Context, request *proto.JoinRequest) (*proto.JoinResponse, error) {
	game.GetRoom().Join(request.PlayerId, request.Session)
	res := &proto.JoinResponse{
		Code:     proto.Code_CODE_SUCCESS,
		Chips:    game.GetRoom().Chips,
		BetZones: []int32{0, 1},
	}
	return res, nil
}

func (g GameService) Bet(_ context.Context, request *proto.BetRequest) (*proto.BetResponse, error) {
	result := game.GetRoom().Bet(request.Session, request.BetZone, request.GetBetChip())
	res := &proto.BetResponse{}
	if result {
		res.Code = proto.Code_CODE_SUCCESS
		return res, nil
	}
	res.Code = proto.Code_CODE_FAILURE
	return res, nil
}

func (g GameService) Leave(_ context.Context, request *proto.LeaveRequest) (*proto.LeaveResponse, error) {
	game.GetRoom().Leave(request.Session)
	return &proto.LeaveResponse{
		Code: proto.Code_CODE_SUCCESS,
	}, nil
}
