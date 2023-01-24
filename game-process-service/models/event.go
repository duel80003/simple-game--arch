package models

import proto "game-process-service/proto/gen/v1"

type Event struct {
	Exchange string
	Router   string
	Data     interface{}
}

type BetZoneInfos struct {
	*BetZones
	GameID string `json:"gameId"`
	TMinus int32  `json:"tMinus"`
}

type StateInfo struct {
	GameID string      `json:"gameId"`
	State  proto.State `json:"state"`
}
