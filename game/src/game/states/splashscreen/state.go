package splashscreen

import (
	"errors"

	"github.com/adm87/flinch/data"
	"github.com/adm87/flinch/engine/flinch"
	"github.com/adm87/flinch/game/src/state"
	"github.com/adm87/flinch/storage/images"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	SplashscreenComplete state.StateExitCondition = iota + 1
)

var ()

type State struct {
	img     *ebiten.Image
	op      *ebiten.DrawImageOptions
	opacity float64
}

func New() state.State[flinch.Context] {
	return &State{
		op: &ebiten.DrawImageOptions{
			Filter: ebiten.FilterLinear,
		},
		opacity: 0.0,
	}
}

func (s *State) Enter(ctx *flinch.Context) error {
	loadingOp := data.Static.CreateBatch(
		images.NewLoader(data.Splash1920x1080Black),
	)
	if err := loadingOp.Execute(ctx); err != nil {
		return err
	}

	img, ok := images.Get(data.Splash1920x1080Black)
	if !ok {
		return errors.New("failed to load splashscreen")
	}
	s.img = img

	return nil
}

func (s *State) Exit(ctx *flinch.Context) error {
	images.Delete(data.Splash1920x1080Black)
	s.img = nil

	return nil
}

func (s *State) Process(ctx *flinch.Context) (state.StateExitCondition, error) {
	return state.NilExitCondition, nil
}

func (s *State) Draw(ctx *flinch.Context) {
	w, h := ctx.Screen().Size()
	sx, sy, sw, sh := transformScreen(
		s.img.Bounds().Dx(),
		s.img.Bounds().Dy(),
		w,
		h,
	)

	s.op.GeoM.Reset()
	s.op.GeoM.Scale(sw, sh)
	s.op.GeoM.Translate(sx, sy)

	ctx.Screen().Buffer().DrawImage(s.img, s.op)
}

func transformScreen(sw, sh, tw, th int) (float64, float64, float64, float64) {
	scale := float64(th) / float64(sh)
	x := (float64(tw) - float64(sw)*scale) / 2
	return x, 0.0, scale, scale
}
