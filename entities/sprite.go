package entities

import "github.com/hajimehoshi/ebiten/v2"

const (
	Square int = iota
	Circle
)

type Sprite struct {
	X            float64
	Y            float64
	Img          *ebiten.Image
	W            *float64
	H            *float64
	R            *float64
	R_2          *float64
	ColliderType int
}

func (s *Sprite) GetCollider() {

}

func (s *Sprite) RectCollision() {

}
