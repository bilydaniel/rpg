package config

import "github.com/hajimehoshi/ebiten/v2"

type Camera struct {
	X, Y  float64
	Scale float64
	Speed float64
}

//TODO implement Y sorted
// https://claude.ai/chat/4f32ea21-3789-45e8-900c-45b31c46dd77

func (c *Camera) ScreenToWorld(x, y float64) (worldx float64, worldy float64) {
	worldx = (x / c.Scale) + c.X*c.Speed
	worldy = (y / c.Scale) + c.Y*c.Speed

	return
}

func (c *Camera) WorldToScreen(x, y float64) (worldx float64, worldy float64) {
	worldx = (x / c.Scale) + c.X*c.Speed
	worldy = (y / c.Scale) + c.Y*c.Speed

	return
}

func (c *Camera) WorldToScreenGeom(opts *ebiten.DrawImageOptions, x int, y int) {
	if opts != nil {
		opts.GeoM.Translate(float64(x), float64(y))
		opts.GeoM.Translate(-c.X*c.Speed, -c.Y*c.Speed)
		opts.GeoM.Scale(c.Scale, c.Scale)
	}
}
