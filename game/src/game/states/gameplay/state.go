package gameplay

import (
	"github.com/adm87/flinch/engine/flinch"
	"github.com/adm87/flinch/game/src/state"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 1280 * 0.25
	screenHeight = 720 * 0.25
)

type State struct {
	worldBuffer *ebiten.Image
}

func New() state.State[flinch.Context] {
	return &State{
		worldBuffer: ebiten.NewImage(screenWidth, screenHeight),
	}
}

func (s *State) Enter(ctx *flinch.Context) error {
	return nil
}

func (s *State) Exit(ctx *flinch.Context) error {
	s.worldBuffer.Deallocate()
	return nil
}

func (s *State) Process(ctx *flinch.Context) (state.StateExitCondition, error) {
	return state.NilExitCondition, nil
}

func (s *State) Draw(ctx *flinch.Context) {
	ebitenutil.DebugPrint(s.worldBuffer, "Gameplay...")
}
