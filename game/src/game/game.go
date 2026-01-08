package game

import (
	"image/color"

	"github.com/adm87/flinch/engine/flinch"
	"github.com/adm87/flinch/game/src/game/states/boot"
	"github.com/adm87/flinch/game/src/game/states/gameplay"
	"github.com/adm87/flinch/game/src/game/states/splashscreen"
	"github.com/adm87/flinch/game/src/state"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	TargetWidth  = 1280
	TargetHeight = 720
)

var (
	ClearColor = color.RGBA{100, 149, 237, 255}
)

// Game States
var (
	fsm = state.NewFSM[flinch.Context]()

	bootStateID         = state.Register[boot.State](fsm, boot.New)
	gameplayID          = state.Register[gameplay.State](fsm, gameplay.New)
	splashscreenStateID = state.Register[splashscreen.State](fsm, splashscreen.New)

	transitions = state.Transitions{
		bootStateID: {
			boot.BootSuccess: splashscreenStateID,
		},
		splashscreenStateID: {
			splashscreen.SplashscreenComplete: gameplayID,
		},
	}
)

type Drawable interface {
	Draw(ctx *flinch.Context)
}

type ggame struct {
	ctx *flinch.Context
	op  *ebiten.DrawImageOptions
}

func Run(ctx *flinch.Context) error {
	ebiten.SetWindowSize(TargetWidth, TargetHeight)
	ebiten.SetWindowTitle("Flinch")

	ctx.Screen().SetSize(TargetWidth, TargetHeight)

	fsm.SetNext(bootStateID)
	fsm.SetTransitions(transitions)

	return ebiten.RunGame(&ggame{
		ctx: ctx,
		op: &ebiten.DrawImageOptions{
			Filter: ebiten.FilterLinear,
		},
	})
}

func (g *ggame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return TargetWidth, TargetHeight
}

func (g *ggame) Update() error {
	// Debug: Exit the game when the Escape key is pressed.
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// Debug: Toggle fullscreen mode when F11 is pressed.
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	// Update the game context.
	g.ctx.Update()

	// Process the FSM.
	return fsm.Process(g.ctx)
}

func (g *ggame) Draw(screen *ebiten.Image) {
	buffer := g.ctx.Screen().Buffer()

	if !fsm.IsTransitioning() {
		if drawable, ok := fsm.State().(Drawable); ok {
			buffer.Fill(ClearColor)
			drawable.Draw(g.ctx)
		}
	}

	screen.DrawImage(buffer, g.op)
}
