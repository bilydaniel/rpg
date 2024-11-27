package assets

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Assets struct {
	Tileset map[string]*ebiten.Image
}

func InitAssets() (Assets, error) {
	var err error
	assets := Assets{Tileset: map[string]*ebiten.Image{}}
	assets.Tileset["floor"], _, err = ebitenutil.NewImageFromFile("assets/images/tilesets/TilesetFloor.png")
	if err != nil {
		return assets, err
	}
	return assets, nil
}

func (a *Assets) LoadTileSet(name string) {

}
