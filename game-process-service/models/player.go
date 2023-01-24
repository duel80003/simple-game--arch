package models

type Player struct {
	ID string `redis:"player_id"`
	BetZones
}

func NewPlayer(id string) *Player {
	return &Player{
		ID: id,
	}
}
