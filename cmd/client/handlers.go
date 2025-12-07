package main

import (
	"fmt"

	"github.com/jcourtney5/peril/internal/gamelogic"
	"github.com/jcourtney5/peril/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
	}
}
