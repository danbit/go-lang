package main

import (
	"fmt"

	"github.com/looplab/fsm"
)

const (
	ON_MENU State = iota
	PLAYING
	PAUSED
	GAME_OVER
	GAME_WIN
)

type State int

type GameState struct {
	To  string
	FSM *fsm.FSM
}

func (gs *GameState) enterState(e *fsm.Event) {
	fmt.Printf("The gameState to %s is %s\n", gs.To, e.Dst)

	if e.Src == ON_MENU.value() && e.Dst == PLAYING.value() {
		snake, food, score = createLevel()
	}
}

func (gs *GameState) leaveState(e *fsm.Event) {
	fmt.Printf("Living fom state %s and enter to %s\n", e.Src, e.Dst)

	if e.Src == PLAYING.value() && (e.Dst == GAME_OVER.value() || e.Dst == GAME_WIN.value()) {
		if score > highScore {
			saveHighscore(score)
			highScore = score
		}
	}
}

func (gs *GameState) changeState(state State) {
	err := gs.FSM.Event(state.value())
	if err != nil {
		fmt.Printf("Failed to change state: %s\n", err)
	}
}

func (s State) value() string {
	return [...]string{"on_menu", "playing", "paused", "game_over", "game_win"}[s]
}
