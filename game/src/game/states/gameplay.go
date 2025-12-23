package states

import (
	"github.com/adm87/flinch/data"
	"github.com/adm87/flinch/engine/encoding"
	"github.com/adm87/flinch/engine/flinch"
	"github.com/adm87/flinch/storage/images"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	TargetWidth  = 918
	TargetHeight = 515
)

var GameplayID uint64 = encoding.HashType[Gameplay]()

type Gameplay struct {
}

func NewGameplay() *Gameplay {
	return &Gameplay{}
}

func (g *Gameplay) Enter(ctx *flinch.Context) error {
	testImageLoader := data.Assets.CreateBatch(images.NewLoader(
		data.SampleA,
	))
	if err := testImageLoader.Execute(ctx); err != nil {
		ctx.Logger().Error("failed to load test images", "error", err.Error())
	}
	return nil
}

func (g *Gameplay) Exit(ctx *flinch.Context) error {
	return nil
}

func (g *Gameplay) Update(ctx *flinch.Context) (uint64, error) {
	return 0, nil
}

func (g *Gameplay) Draw(ctx *flinch.Context, screen *ebiten.Image) {
	if img, exists := images.Get(data.SampleA); exists {
		screen.DrawImage(img, nil)
	}
}

func (g *Gameplay) Layout(outsideWidth, outsideHeight int) (int, int) {
	return TargetWidth, TargetHeight
}
