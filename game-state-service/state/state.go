package state

import (
	"fmt"
	proto "game-state-service/proto/gen/v1"
	tools "github.com/duel80003/my-tools"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	chanManager     *ChanManager
	chanManagerOnce sync.Once
)

type ChanManager struct {
	m   map[string]*Chan
	mux sync.RWMutex
}

func (manager *ChanManager) Add(ch *Chan) {
	manager.mux.Lock()
	defer manager.mux.Unlock()
	manager.m[ch.index] = ch
}

func (manager *ChanManager) Delete(ch *Chan) {
	manager.mux.Lock()
	defer manager.mux.Unlock()
	tools.Logger.Infof("delete chan: %s", ch.index)
	delete(manager.m, ch.index)
}

func (manager *ChanManager) Notify(state proto.State, t int32) {
	manager.mux.RLock()
	defer manager.mux.RUnlock()
	tools.Logger.Infof("notify chan: %d", len(manager.m))

	for _, value := range manager.m {
		value.Ch <- &Info{State: state, Time: t}
	}
}

type Info struct {
	State proto.State
	Time  int32
}

type Chan struct {
	Ch    chan *Info
	index string
}

func NewChan() *Chan {
	return &Chan{
		Ch:    make(chan *Info, 1),
		index: fmt.Sprintf("%d", time.Now().UnixMilli()),
	}
}

func GetChanManager() *ChanManager {
	chanManagerOnce.Do(func() {
		chanManager = new(ChanManager)
		chanManager.m = make(map[string]*Chan)
	})
	return chanManager
}

func StartStateMachine() {
	stateGameStart := os.Getenv("STATE_GAME_START")
	stateStartBet := os.Getenv("STATE_START_BET")
	stateStopBet := os.Getenv("STATE_STOP_BET")
	stateAward := os.Getenv("STATE_AWARD")
	stateEnd := os.Getenv("STATE_END")
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	currentState := proto.State_STATE_INIT
	var wg sync.WaitGroup
	for {
		select {
		case <-ticker.C:
			switch currentState {
			case proto.State_STATE_INIT:
				tools.Logger.Infof("current state: %s", currentState)
				currentState = proto.State_STATE_GAME_START
			case proto.State_STATE_GAME_START:
				tools.Logger.Infof("current state: %s", currentState)
				wg.Add(1)
				GetChanManager().Notify(currentState, getTime(stateGameStart))
				time.AfterFunc(getTimeDurationByString(stateGameStart), func() {
					currentState = proto.State_STATE_START_BET
					wg.Done()
				})
				wg.Wait()
			case proto.State_STATE_START_BET:
				tools.Logger.Infof("current state: %s", currentState)
				wg.Add(1)
				GetChanManager().Notify(currentState, getTime(stateStartBet))
				time.AfterFunc(getTimeDurationByString(stateStartBet), func() {
					currentState = proto.State_STATE_STOP_BET
					wg.Done()
				})
				wg.Wait()
			case proto.State_STATE_STOP_BET:
				tools.Logger.Infof("current state: %s", currentState)
				wg.Add(1)
				GetChanManager().Notify(currentState, getTime(stateStopBet))
				time.AfterFunc(getTimeDurationByString(stateStopBet), func() {
					currentState = proto.State_STATE_AWARD
					wg.Done()
				})
				wg.Wait()
			case proto.State_STATE_AWARD:
				tools.Logger.Infof("current state: %s", currentState)
				wg.Add(1)
				GetChanManager().Notify(currentState, getTime(stateAward))
				time.AfterFunc(getTimeDurationByString(stateAward), func() {
					currentState = proto.State_STATE_END
					wg.Done()
				})
				wg.Wait()
			case proto.State_STATE_END:
				tools.Logger.Infof("current state: %s", currentState)
				wg.Add(1)
				GetChanManager().Notify(currentState, getTime(stateEnd))
				time.AfterFunc(getTimeDurationByString(stateEnd), func() {
					currentState = proto.State_STATE_INIT
					wg.Done()
				})
				wg.Wait()
			}
		}
	}
}

func getTimeDurationByString(str string) time.Duration {
	value, err := strconv.Atoi(str)
	if err != nil {
		tools.Logger.Panicf("[getTimeDurationByString] invalid time value: %s", str)
	}
	return time.Duration(value) * time.Second
}

func getTime(str string) int32 {
	value, err := strconv.Atoi(str)
	if err != nil {
		tools.Logger.Panicf("[getTime] invalid time value: %s", str)
	}
	return int32(value)
}
