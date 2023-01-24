package handler

const (
	Join  = "join"
	Bet   = "bet"
	Leave = "leave"
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
