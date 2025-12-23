package images

import (
	"bytes"
	"sync"

	"github.com/adm87/flinch/engine/flinch"
	"github.com/adm87/flinch/engine/resources"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var (
	cache = make(map[resources.Asset]*ebiten.Image)
	mu    = sync.RWMutex{}
)

func Get(asset resources.Asset) (*ebiten.Image, bool) {
	mu.RLock()
	defer mu.RUnlock()

	img, exists := cache[asset]
	return img, exists
}

func Set(asset resources.Asset, img *ebiten.Image) {
	mu.Lock()
	defer mu.Unlock()

	cache[asset] = img
}

// NewLoader creates a new LoadingTask that loads the specified assets into the image cache.
func NewLoader(assets ...resources.Asset) resources.LoadingTask {
	return func(ctx *flinch.Context, rs *resources.ResourceSystem, batchID uint64) error {
		for _, asset := range assets {
			if err := loadImage(rs, asset, batchID); err != nil {
				return err
			}
		}
		return nil
	}
}

// loadImage is a helper to maintain concurrent image loading safety.
func loadImage(rs *resources.ResourceSystem, asset resources.Asset, batchID uint64) error {
	lock := rs.LockAsset(batchID, asset)
	defer lock.Release()

	data, err := rs.ReadBytes(asset)
	if err != nil {
		return err
	}

	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(data))
	if err != nil {
		return err
	}

	Set(asset, img)

	return nil
}
