package models

import proto "main-service/proto/gen/v1"

type Event struct {
	Exchange string
	Router   string
	Data     interface{}
}

type BetZoneInfos struct {
	*BetZones
	TMinus int32 `json:"tMinus"`
}

type StateInfo struct {
	State proto.State `json:"state"`
}
