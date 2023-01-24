package models

type Event struct {
	Exchange string
	Router   string
	Data     interface{}
}

type BetZoneInfos struct {
	*BetZones
	TMinus int32 `json:"tMinus"`
}
