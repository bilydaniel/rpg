package utils

type Point struct {
	X, Y float64
}

type Node struct {
	X, Y int
}

// TODO disconected from the rest, kinda reimplementation,
type CollisionShape interface {
	Intersects(Point) bool
	IntersectsLine(Point, Point) bool
}

type RectangleCollision struct {
	Minx, Miny, Maxx, Maxy float64
}
