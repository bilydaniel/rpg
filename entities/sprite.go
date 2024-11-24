package entities

import "github.com/hajimehoshi/ebiten/v2"

type Sprite struct {
	X   float64
	Y   float64
	Img *ebiten.Image
}
