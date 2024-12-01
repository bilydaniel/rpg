package config

type Camera struct {
	X, Y  float64
	Scale float64
}

func (c *Camera) ToWorld(x, y float64) (worldx float64, worldy float64) {
	worldx = (x / c.Scale) + c.X
	worldy = (y / c.Scale) + c.Y

	return
}
