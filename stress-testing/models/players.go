package models

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Player struct {
	PlayerID string
	Conn     *websocket.Conn
	BetMux   sync.Mutex
}
