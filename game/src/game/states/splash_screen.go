package states

import (
	"github.com/adm87/flinch/data"
	"github.com/adm87/flinch/engine/encoding"
	"github.com/adm87/flinch/engine/flinch"
	"github.com/adm87/flinch/storage/images"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

var SplashScreenID uint64 = encoding.HashType[SplashScreen]()

type SplashScreen struct {
	sequence *flinch.Sequence
	opacity  float32
	width    int
	height   int
}

func NewSplashScreen() *SplashScreen {
	return &SplashScreen{}
}

func (ss *SplashScreen) Enter(ctx *flinch.Context) error {
	ctx.Logger().Info("Entering SplashScreen state")

	splashScreenLoader := data.Static.CreateBatch(images.NewLoader(
		data.Splash1920x1080Black,
	))
	if err := splashScreenLoader.Execute(ctx); err != nil {
		ctx.Logger().Error("failed to load splash screen", "error", err.Error())
	}

	if img, exists := images.Get(data.Splash1920x1080Black); exists {
		ss.width = img.Bounds().Dx()
		ss.height = img.Bounds().Dy()
	}

	ss.sequence = flinch.NewSequence(
		ss.fadeSplashScreenIn(ctx, 0.5),
		ss.holdSplashScreen(ctx, 1.0),
		ss.fadeSplashScreenOut(ctx, 0.5),
	)
	ctx.Scripts().AddSequence(ss.sequence)

	return nil
}

func (ss *SplashScreen) Exit(ctx *flinch.Context) error {
	ctx.Scripts().RemoveSequence(ss.sequence)
	ss.sequence = nil

	return nil
}

func (ss *SplashScreen) Update(ctx *flinch.Context) (uint64, error) {
	if !ss.sequence.Started() {
		ss.sequence.Start()
	}
	if ss.sequence.Complete() {
		return GameplayID, nil
	}
	return 0, nil
}

func (ss *SplashScreen) Draw(ctx *flinch.Context, screen *ebiten.Image) {
	if img, exists := images.Get(data.Splash1920x1080Black); exists {
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(ss.opacity)
		screen.DrawImage(img, op)
	}
}

func (ss *SplashScreen) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ss.width, ss.height
}

func (ss *SplashScreen) fadeSplashScreenIn(ctx *flinch.Context, duration float32) flinch.Action {
	tween := gween.New(0, 1, duration, ease.Linear)
	return func(dt float64) bool {
		curr, done := tween.Update(float32(dt))
		ss.opacity = curr
		return done
	}
}

func (ss *SplashScreen) fadeSplashScreenOut(ctx *flinch.Context, duration float32) flinch.Action {
	tween := gween.New(1, 0, duration, ease.Linear)
	return func(dt float64) bool {
		curr, done := tween.Update(float32(dt))
		ss.opacity = curr
		return done
	}
}

func (ss *SplashScreen) holdSplashScreen(ctx *flinch.Context, duration float32) flinch.Action {
	timer := flinch.NewTimer(duration, false)
	return func(dt float64) bool {
		timer.Update(float32(dt))
		return timer.Completed()
	}
}
