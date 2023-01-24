package handler

const (
	Join         = "join"
	Bet          = "bet"
	Leave        = "leave"
	State        = "state"
	BetZoneInfos = "betZoneInfos"
)

var (
	invalidRes = &ErrorRes{
		Code: InvalidRequest,
		Msg:  "invalid request",
	}

	internalErrRes = &ErrorRes{
		Code: InternalErr,
	}
)
