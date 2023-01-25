package models

import proto "game-process-service/proto/gen/v1"

type Event struct {
	Exchange string
	Router   string
	Data     *EventData
}

type EventData struct {
	Session  string `json:"session"`
	PlayerID string `json:"playerId"`
	Data     interface{}
}

type BetZoneInfos struct {
	*BetZones
	TMinus int32 `json:"tMinus"`
}

type StateInfo struct {
	State proto.State `json:"state"`
}
