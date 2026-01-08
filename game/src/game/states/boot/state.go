package boot

import (
	"github.com/adm87/flinch/engine/flinch"
	"github.com/adm87/flinch/game/src/state"
)

const (
	BootSuccess state.StateExitCondition = iota + 1
)

type State struct {
}

func New() state.State[flinch.Context] {
	return &State{}
}

func (s *State) Enter(ctx *flinch.Context) error {
	return nil
}

func (s *State) Exit(ctx *flinch.Context) error {
	return nil
}

func (s *State) Process(ctx *flinch.Context) (state.StateExitCondition, error) {
	return BootSuccess, nil
}
