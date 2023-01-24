package game

import (
	"context"
	"game-process-service/drivers"
	proto "game-process-service/proto/gen/v1"
	tools "github.com/duel80003/my-tools"
	"io"
)

func StartGameFlow() {
	client := proto.NewGameStateServiceClient(drivers.GetStateGRpcConn())
	stream, err := client.State(context.TODO(), &proto.StateRequest{})
	if err != nil {
		tools.Logger.Fatalf("client.State failed: %v", err)
	}
	for {
		state, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			tools.Logger.Errorf("state machine failed: %v", err)
			return
		}
		//tools.Logger.Infof("state info: %+v", state)
		switch state.State {
		case proto.State_STATE_GAME_START:
			GetRoom().State = proto.State_STATE_GAME_START
		case proto.State_STATE_START_BET:
			GetRoom().State = proto.State_STATE_START_BET
		case proto.State_STATE_STOP_BET:
			GetRoom().State = proto.State_STATE_START_BET
		case proto.State_STATE_AWARD:
			GetRoom().State = proto.State_STATE_AWARD
		case proto.State_STATE_END:
			GetRoom().State = proto.State_STATE_END
			GetRoom().Reset()
		}
	}
}
