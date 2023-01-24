package handler

import (
	"context"
	"encoding/json"
	tools "github.com/duel80003/my-tools"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
	"main-service/drivers"
	proto "main-service/proto/gen/v1"
	"net"
	"sync"
)

var (
	handler     *WSHandler
	handlerOnce sync.Once
)

type WSHandler struct {
	sidMapConn sync.Map
	connMapSid sync.Map
}

func GetWsHandler() *WSHandler {
	handlerOnce.Do(func() {
		handler = new(WSHandler)
	})
	return handler
}

func (handler *WSHandler) addConn(conn net.Conn) string {
	sid := uuid.NewString()
	handler.sidMapConn.Store(sid, conn)
	handler.connMapSid.Store(conn, sid)
	return sid
}

func (handler *WSHandler) deleteConn(conn net.Conn) {
	sidI, ok := handler.connMapSid.LoadAndDelete(conn)
	if !ok {
		return
	}
	sid, ok := sidI.(string)
	if !ok {
		return
	}
	handler.sidMapConn.Delete(sid)
}

func (handler *WSHandler) Run(conn net.Conn) {
	go func() {
		session := handler.addConn(conn)
		defer conn.Close()
		defer tools.Logger.Debugf("connectoin close")
		defer handler.deleteConn(conn)
		r := wsutil.NewReader(conn, ws.StateServerSide)
		w := wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText)
		decoder := json.NewDecoder(r)
		encoder := json.NewEncoder(w)
		for {
			hdr, err := r.NextFrame()
			if err != nil {
				tools.Logger.Errorf("reat next message error: %s", err)
				return
			}
			if hdr.OpCode == ws.OpClose {
				return
			}
			var req Request
			if err := decoder.Decode(&req); err != nil {
				tools.Logger.Errorf("decode error: %s", err)
				return
			}
			resp := sendToGs(&req, session)
			if err := encoder.Encode(&resp); err != nil {
				tools.Logger.Errorf("encode error: %s", err)
				return
			}
			if err = w.Flush(); err != nil {
				tools.Logger.Errorf("flush error: %s", err)
				return
			}
		}
	}()
}

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
	tools.Logger.Infof("bet req: %+v", request)
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
