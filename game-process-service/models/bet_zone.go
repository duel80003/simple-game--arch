package models

type BetZones struct {
	BetZoneHigher int32 `redis:"bet_zone_higher"`
	BetZoneLower  int32 `redis:"bet_zone_lower"`
}
