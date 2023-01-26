package models

const (
	NotifyState       EventType = "notify_state"
	NotifyBetZoneInfo EventType = "notify_bet_zone_info"
	PlayerJoin        EventType = "player_join"
	PlayerLeave       EventType = "player_leave"
)

type EventType string

type NotificationEvent struct {
	Type   EventType
	SID    string
	PID    string
	TMinus int32
	Event  *Event
}
