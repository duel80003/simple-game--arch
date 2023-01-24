package models

type BetZones struct {
	BetZoneHigher int32 `redis:"bet_zone_higher" json:"betZoneHigher"`
	BetZoneLower  int32 `redis:"bet_zone_lower" json:"betZoneLower"`
}
