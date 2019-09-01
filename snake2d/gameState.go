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

func NewGameState(to string) *GameState {
	gameState := &GameState{
		To: to,
	}

	gameState.FSM = fsm.NewFSM(
		ON_MENU.value(),
		fsm.Events{
			{Name: PLAYING.value(), Src: []string{ON_MENU.value()}, Dst: PLAYING.value()},
			{Name: PAUSED.value(), Src: []string{PLAYING.value()}, Dst: PAUSED.value()},
			{Name: PAUSED.value(), Src: []string{PAUSED.value()}, Dst: PLAYING.value()},
			{Name: GAME_OVER.value(), Src: []string{PLAYING.value()}, Dst: GAME_OVER.value()},
			{Name: GAME_OVER.value(), Src: []string{GAME_OVER.value()}, Dst: ON_MENU.value()},
			{Name: ON_MENU.value(), Src: []string{GAME_OVER.value()}, Dst: ON_MENU.value()},
			{Name: ON_MENU.value(), Src: []string{GAME_WIN.value()}, Dst: ON_MENU.value()},
			{Name: GAME_WIN.value(), Src: []string{PLAYING.value()}, Dst: GAME_WIN.value()},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { gameState.enterState(e) },
			"leave_state": func(e *fsm.Event) { gameState.leaveState(e) },
		},
	)

	return gameState
}

func (gs *GameState) enterState(e *fsm.Event) {
	fmt.Printf("The gameState to %s is %s\n", gs.To, e.Dst)

	if e.Src == ON_MENU.value() && e.Dst == PLAYING.value() {
		level = NewLevel()
		snake, food, score = level.snake, level.food, level.score

		scene = &Scene{
			GameObjects: []GameObject{
				&food,
				&snake,
			},
		}
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
