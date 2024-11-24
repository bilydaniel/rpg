package assets

import "github.com/hajimehoshi/ebiten"

type Assets struct {
	Tileset map[string]*ebiten.Image
}

func (a *Assets) LoadTileSet(name string) {

}
