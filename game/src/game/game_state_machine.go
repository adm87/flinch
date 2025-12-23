package game

import (
	"github.com/adm87/flinch/engine/flinch"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameState interface {
	Enter(ctx *flinch.Context) error
	Exit(ctx *flinch.Context) error
	Update(ctx *flinch.Context) (uint64, error)
	Draw(ctx *flinch.Context, screen *ebiten.Image)
	Layout(outsideWidth, outsideHeight int) (int, int)
}

type GameStateMachine struct {
	transitioning bool

	current      uint64
	next         uint64
	screenBuffer *ebiten.Image

	states map[uint64]GameState
}

func NewGameStateMachine(states map[uint64]GameState, initial uint64) *GameStateMachine {
	return &GameStateMachine{
		states: states,
		next:   initial,
	}
}

func (gsm *GameStateMachine) Update(ctx *flinch.Context) error {
	if gsm.next != 0 {
		current := gsm.states[gsm.current]
		if current != nil {
			if err := current.Exit(ctx); err != nil {
				return err
			}
		}

		gsm.current = gsm.next
		gsm.next = 0
		gsm.transitioning = true

		next := gsm.states[gsm.current]
		if next != nil {
			if err := next.Enter(ctx); err != nil {
				return err
			}
		}
	}

	current := gsm.states[gsm.current]
	if current != nil && !gsm.transitioning {
		nextStateID, err := current.Update(ctx)
		if err != nil {
			return err
		}
		if nextStateID != 0 {
			gsm.next = nextStateID
		}
	}
	return nil
}

func (gsm *GameStateMachine) Draw(ctx *flinch.Context) {
	current := gsm.states[gsm.current]
	if current != nil && !gsm.transitioning {
		current.Draw(ctx, ctx.Screen().Buffer())
	}
}

func (gsm *GameStateMachine) Layout(ctx *flinch.Context, outsideWidth, outsideHeight int) (int, int) {
	if current := gsm.states[gsm.current]; current != nil {
		gsm.transitioning = false
		return current.Layout(outsideWidth, outsideHeight)
	}
	return outsideWidth, outsideHeight
}
