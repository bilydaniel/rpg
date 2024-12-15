package utils

type Point struct {
	X, Y float64 // tiles
}

// TODO disconected from the rest, kinda reimplementation,
// put it together with sprite
// need to do this part alone first
type CollisionShape interface {
	Intersects(Point) bool
	IntersectsLine(Point, Point) bool
}

type RectangleCollision struct {
	Minx, Miny, Maxx, Maxy float64
}
