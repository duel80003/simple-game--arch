package servers

import (
	proto "game-state-service/proto/gen/v1"
	"game-state-service/state"
)

type StateService struct{}

func (s StateService) State(_ *proto.StateRequest, server proto.GameStateService_StateServer) error {
	ch := state.NewChan()
	defer close(ch.Ch)
	defer state.GetChanManager().Delete(ch)
	state.GetChanManager().Add(ch)
	for gameState := range ch.Ch {
		err := server.Send(&proto.StateResponse{State: gameState.State, TMinus: gameState.Time})
		if err != nil {
			return err
		}
	}
	return nil
}
