package entities

import (
	"bilydaniel/rpg/config"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	Square int = iota
	Circle
)

type Sprite interface {
	Position() (float64, float64)
	Center() (float64, float64)
	GetX() float64
	GetY() float64
	Size() float64
	Image() *ebiten.Image
	SetPosition(x float64, y float64)
	SetX(x float64)
	SetY(y float64)
}

type SquareSprite struct {
	X   float64
	Y   float64
	W   float64
	H   float64
	Img *ebiten.Image
}

func (ss *SquareSprite) Top() float64 {
	return ss.Y
}
func (ss *SquareSprite) Bottom() float64 {
	return ss.Y + ss.H
}
func (ss *SquareSprite) Left() float64 {
	return ss.X
}
func (ss *SquareSprite) Right() float64 {
	return ss.X + ss.W
}
func (ss *SquareSprite) Center() (float64, float64) {
	return ss.X + config.TileSize/2, ss.Y + config.TileSize/2
}
func (ss *SquareSprite) Centerx() float64 {
	return ss.X + config.TileSize/2
}
func (ss *SquareSprite) Centery() float64 {
	return ss.Y + config.TileSize/2
}

func (ss *SquareSprite) Position() (float64, float64) {
	return ss.X, ss.Y
}

func (ss *SquareSprite) GetX() float64 {
	return ss.X
}

func (ss *SquareSprite) GetY() float64 {
	return ss.Y
}

func (ss *SquareSprite) Size() float64 {
	return ss.W
}

func (ss *SquareSprite) Image() *ebiten.Image {
	return ss.Img
}

func (ss *SquareSprite) SetPosition(x float64, y float64) {
	ss.X = x
	ss.Y = y
}
func (ss *SquareSprite) SetX(x float64) {
	ss.X = x
}
func (ss *SquareSprite) SetY(y float64) {
	ss.Y = y
}

type CircleSprite struct {
	X   float64
	Y   float64
	R   float64
	R_2 float64
	Img *ebiten.Image
}

func (cs *CircleSprite) Position() (float64, float64) {
	return cs.X, cs.Y
}

func (cs *CircleSprite) GetX() float64 {
	return cs.X
}

func (cs *CircleSprite) GetY() float64 {
	return cs.Y
}

func (cs *CircleSprite) Size() float64 {
	return cs.R_2
}

func (cs *CircleSprite) Image() *ebiten.Image {
	return cs.Img
}

func (cs *CircleSprite) SetPosition(x float64, y float64) {
	cs.X = x
	cs.Y = y
}
func (cs *CircleSprite) SetX(x float64) {
	cs.X = x
}
func (cs *CircleSprite) SetY(y float64) {
	cs.Y = y
}
func (cs *CircleSprite) Center() (x float64, y float64) {
	return cs.X, cs.Y
}
