package handler

import (
	"encoding/json"
	tools "github.com/duel80003/my-tools"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/google/uuid"
	"main-service/repositories"
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
		defer conn.Close()
		defer tools.Logger.Debugf("connectoin close")
		defer handler.deleteConn(conn)
		session := handler.addConn(conn)
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

func (handler *WSHandler) Broadcast(topic, gameId string, data interface{}) {
	resp := Response{
		Topic: topic,
		Data:  data,
	}
	handler.sidMapConn.Range(func(key, value any) bool {
		conn := value.(net.Conn)
		sid := key.(string)
		gid := repositories.GetPlayerGameId(sid)
		tools.Logger.Infof("sid; %s, gid: %s", sid, gid)

		if gid != gameId {
			return true
		}
		w := wsutil.NewWriter(conn, ws.StateServerSide, ws.OpText)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(&resp); err != nil {
			tools.Logger.Errorf("encode error: %s", err)
			return true
		}
		if err := w.Flush(); err != nil {
			tools.Logger.Errorf("flush error: %s", err)
			return true
		}
		return true
	})
}
