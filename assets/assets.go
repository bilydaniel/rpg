package assets

import (
	"bilydaniel/rpg/config"
	"fmt"
	"image"
	"io/fs"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Assets struct {
	Tileset map[string]*ebiten.Image
	Audio   AudioAssets
	Video   VideoAssets
}

type AudioAssets struct {
}

type VideoAssets struct {
	Images    map[string]map[string]*ebiten.Image //tileset type(floors) => tileset name()TilesetFloor => image
	Tilecashe map[int]*ebiten.Image
}

func InitAssets() (*Assets, error) {
	assets := &Assets{
		Video: VideoAssets{
			Images:    map[string]map[string]*ebiten.Image{},
			Tilecashe: map[int]*ebiten.Image{},
		},
	}
	err := LoadAllAssets(assets)

	if err != nil {
		return nil, err
	}
	return assets, nil
}

func LoadAllAssets(assets *Assets) error {
	if assets == nil {
		return fmt.Errorf("assets is nil")
	}

	root := "assets/tilesets"
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			//skip the root
			return nil
		}
		if filepath.Ext(path) != ".png" {
			//skip non-png files
			return nil
		}

		dirname := filepath.Base(filepath.Dir(path))
		if assets.Video.Images[dirname] == nil {
			assets.Video.Images[dirname] = map[string]*ebiten.Image{}
		}

		assets.Video.Images[dirname][d.Name()], _, err = ebitenutil.NewImageFromFile(path)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func (a *Assets) GetImage(setname string, filename string) *ebiten.Image {
	return a.Video.Images[setname][filename]
}
func (a *Assets) GetTileImage(setname string, filename string, columns int, tileid int) *ebiten.Image {
	//TODO make a cashe tileid => image
	// return cashe[tileid]
	tile, ok := a.Video.Tilecashe[tileid]
	if ok {
		return tile
	}
	x0 := ((tileid - 1) % columns) * 16
	y0 := ((tileid - 1) / columns) * 16

	a.Video.Tilecashe[tileid] = a.Video.Images[setname][filename].SubImage(image.Rect(x0, y0, x0+config.TileSize, y0+config.TileSize)).(*ebiten.Image)
	return a.Video.Tilecashe[tileid]
}
