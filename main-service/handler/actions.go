package handler

import (
	"context"
	tools "github.com/duel80003/my-tools"
	"main-service/drivers"
	proto "main-service/proto/gen/v1"
)

func sendToGs(request *Request, session string) (resp Response) {
	switch request.Topic {
	case Join:
		resp = join(request, session)
	case Bet:
		resp = bet(request, session)
	case Leave:
		resp = leave(request, session)
	default:
		resp.Topic = request.Topic
		resp.Data = &ErrorRes{
			Code: InvalidRequest,
			Msg:  "invalid request",
		}
	}
	return
}

func join(request *Request, session string) (resp Response) {
	tools.Logger.Infof("join req: %+v", request)
	resp.Topic = request.Topic
	pid, ok := request.Data["playerId"].(string)
	if !ok {
		resp.Data = invalidRes
		return
	}
	client := proto.NewGameProcessServiceClient(drivers.GetGsGRpcConn())
	response, err := client.Join(context.TODO(), &proto.JoinRequest{
		PlayerId: pid,
		Session:  session,
	})
	if err != nil {
		resp.Data = internalErrRes
		return
	}
	resp.Data = response
	return
}

func bet(request *Request, session string) (resp Response) {
	tools.Logger.Debugf("bet req: %+v", request)
	resp.Topic = request.Topic
	betZone, ok := request.Data["betZone"].(float64)
	if !ok {
		resp.Data = invalidRes
		return
	}

	chip, ok := request.Data["chip"].(float64)
	if !ok {
		resp.Data = invalidRes
		return
	}
	client := proto.NewGameProcessServiceClient(drivers.GetGsGRpcConn())
	response, err := client.Bet(context.TODO(), &proto.BetRequest{
		Session: session,
		BetZone: int32(betZone),
		BetChip: int32(chip),
	})
	if err != nil {
		tools.Logger.Errorf("bet request error: %s", err)
		resp.Data = internalErrRes
		return
	}
	resp.Data = response
	return
}

func leave(request *Request, session string) (resp Response) {
	tools.Logger.Infof("leave req: %+v", request)
	resp.Topic = request.Topic
	pid, ok := request.Data["playerId"].(string)
	if !ok {
		resp.Data = &ErrorRes{
			Code: InvalidRequest,
			Msg:  "invalid request",
		}
		return
	}
	client := proto.NewGameProcessServiceClient(drivers.GetGsGRpcConn())
	response, err := client.Leave(context.TODO(), &proto.LeaveRequest{
		PlayerId: pid,
		Session:  session,
	})
	if err != nil {

		return
	}
	if err != nil {
		resp.Data = internalErrRes
		return
	}
	resp.Data = response
	return
}
