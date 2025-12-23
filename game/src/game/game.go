package game

import (
	"github.com/adm87/flinch/engine/flinch"
	"github.com/adm87/flinch/game/src/game/states"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	ctx    *flinch.Context
	states *GameStateMachine
}

func NewGame(ctx *flinch.Context) *Game {
	ebiten.SetWindowTitle("Flinch")
	ebiten.SetWindowSize(1280, 720)
	return &Game{
		ctx: ctx,
		states: NewGameStateMachine(map[uint64]GameState{
			states.SplashScreenID: states.NewSplashScreen(),
			states.GameplayID:     states.NewGameplay(),
		}, states.SplashScreenID),
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	width, height := g.states.Layout(g.ctx, outsideWidth, outsideHeight)

	// Ensure the context's screen buffer is the correct size.
	g.ctx.Screen().SetSize(width, height)

	return width, height
}

func (g *Game) Update() error {
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

	// Update the game state.
	return g.states.Update(g.ctx)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ctx.Screen().Buffer().Clear()

	// Draw the current game state.
	g.states.Draw(g.ctx)

	// Draw the game context's screen buffer to the application's screen.
	screen.DrawImage(g.ctx.Screen().Buffer(), nil)
}
