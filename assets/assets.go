package assets

import (
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
	Tilesets map[string]TilesetAsset
	Objects  map[string]ObjectAsset
}

type TilesetAsset struct {
	Img *ebiten.Image
}
type ObjectAsset struct {
}

func InitAssets() (Assets, error) {
	//TODO add all assets here
	var err error
	assets := Assets{Tileset: map[string]*ebiten.Image{}}
	assets.Tileset["floor"], _, err = ebitenutil.NewImageFromFile("assets/images/tilesets/TilesetFloor.png")
	if err != nil {
		return assets, err
	}

	assets.Tileset["objects"], _, err = ebitenutil.NewImageFromFile("assets/images/tilesets/house1.png")
	if err != nil {
		return assets, err
	}
	return assets, nil
}

func (a *Assets) LoadTileSet(name string) {

}
