package pathfinding

// CollisionShape defines different types of collision shapes
type CollisionShape interface {
	Intersects(point Point) bool
	IntersectsLine(start, end Point) bool
}

// CircleCollision represents a circular collision area
type CircleCollision struct {
	Center Point
	Radius float64
}

func (c *CircleCollision) Intersects(point Point) bool {
	dx := point.X - c.Center.X
	dy := point.Y - c.Center.Y
	return dx*dx+dy*dy <= c.Radius*c.Radius
}

func (c *CircleCollision) IntersectsLine(start, end Point) bool {
	// Line-circle intersection algorithm
	x0, y0 := start.X, start.Y
	x1, y1 := end.X, end.Y
	cx, cy := c.Center.X, c.Center.Y
	r := c.Radius

	// Vector from line start to circle center
	dx := cx - x0
	dy := cy - y0

	// Line direction vector
	lineDx := x1 - x0
	lineDy := y1 - y0

	// Project circle center onto the line
	t := (dx*lineDx + dy*lineDy) / (lineDx*lineDx + lineDy*lineDy)

	// Closest point on the line
	closestX := x0 + t*lineDx
	closestY := y0 + t*lineDy

	// Clamp to line segment
	if t < 0 {
		closestX = x0
		closestY = y0
	} else if t > 1 {
		closestX = x1
		closestY = y1
	}

	// Distance from closest point to circle center
	distX := cx - closestX
	distY := cy - closestY

	return distX*distX+distY*distY <= r*r
}

// RectangleCollision represents a rectangular collision area
type RectangleCollision struct {
	MinX, MinY, MaxX, MaxY float64
}

func (r *RectangleCollision) Intersects(point Point) bool {
	return point.X >= r.MinX && point.X <= r.MaxX &&
		point.Y >= r.MinY && point.Y <= r.MaxY
}

func (r *RectangleCollision) IntersectsLine(start, end Point) bool {
	// Check if line segment intersects rectangle
	return lineIntersectsRectangle(start, end, *r)
}

// ComplexPathFinder extends the previous PathFinder with collision detection
type ComplexPathFinder struct {
	Grid            [][]*Node
	Width, Height   int
	CollisionShapes []CollisionShape
}

// NewComplexPathFinder creates a new PathFinder with collision support
func NewComplexPathFinder(width, height int) *ComplexPathFinder {
	grid := make([][]*Node, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]*Node, width)
		for x := 0; x < width; x++ {
			grid[y][x] = &Node{
				Pos:      Point{X: float64(x), Y: float64(y)},
				Walkable: true,
			}
		}
	}
	return &ComplexPathFinder{
		Grid:            grid,
		Width:           width,
		Height:          height,
		CollisionShapes: make([]CollisionShape, 0),
	}
}

// AddCollisionShape adds a collision shape to the pathfinder
func (pf *ComplexPathFinder) AddCollisionShape(shape CollisionShape) {
	pf.CollisionShapes = append(pf.CollisionShapes, shape)
}

// IsCollisionFree checks if a point or line is free of collisions
func (pf *ComplexPathFinder) IsCollisionFree(point Point) bool {
	// Check grid walkability
	gridX, gridY := int(point.X), int(point.Y)
	if gridX < 0 || gridX >= pf.Width || gridY < 0 || gridY >= pf.Height {
		return false
	}

	// Check grid walkability
	if !pf.Grid[gridY][gridX].Walkable {
		return false
	}

	// Check against all collision shapes
	for _, shape := range pf.CollisionShapes {
		if shape.Intersects(point) {
			return false
		}
	}

	return true
}

// IsLineCollisionFree checks if a line segment is free of collisions
func (pf *ComplexPathFinder) IsLineCollisionFree(start, end Point) bool {
	// Check grid boundaries
	if !pf.IsPointInGrid(start) || !pf.IsPointInGrid(end) {
		return false
	}

	// Line of sight with grid walkability
	x0, y0 := int(start.X), int(start.Y)
	x1, y1 := int(end.X), int(end.Y)

	dx := abs(x1 - x0)
	dy := abs(y1 - y0)

	sx := 1
	if x0 >= x1 {
		sx = -1
	}
	sy := 1
	if y0 >= y1 {
		sy = -1
	}

	err := dx - dy
	x, y := x0, y0

	// Bresenham's line algorithm with collision checks
	for {
		// Check grid walkability
		if !pf.Grid[y][x].Walkable {
			return false
		}

		// Check against all collision shapes
		currentPoint := Point{X: float64(x), Y: float64(y)}
		for _, shape := range pf.CollisionShapes {
			if shape.Intersects(currentPoint) {
				return false
			}
		}

		if x == x1 && y == y1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}

	// Check line intersections with collision shapes
	for _, shape := range pf.CollisionShapes {
		if shape.IntersectsLine(start, end) {
			return false
		}
	}

	return true
}

// IsPointInGrid checks if a point is within the grid
func (pf *ComplexPathFinder) IsPointInGrid(point Point) bool {
	return point.X >= 0 && point.X < float64(pf.Width) &&
		point.Y >= 0 && point.Y < float64(pf.Height)
}

// Modify the existing ThetaStar method to use collision checks
func (pf *ComplexPathFinder) ThetaStar(start, goal Point) []Point {
	// Previous ThetaStar implementation, but replace LineOfSight with IsLineCollisionFree
	// Replace walkability checks with IsCollisionFree
	// The core logic remains the same, just use the new collision methods

	// ... (rest of the previous ThetaStar implementation)
	// Replace pf.LineOfSight with pf.IsLineCollisionFree
	// Replace walkability checks with pf.IsCollisionFree
}

// Utility functions to help with collision detection
func lineIntersectsRectangle(start, end Point, rect RectangleCollision) bool {
	// Check if line segment intersects rectangle
	// Implement Cohen-Sutherland line clipping algorithm

	const (
		INSIDE = 0
		LEFT   = 1
		RIGHT  = 2
		BOTTOM = 4
		TOP    = 8
	)

	computeOutCode := func(point Point) int {
		code := INSIDE

		if point.X < rect.MinX {
			code |= LEFT
		} else if point.X > rect.MaxX {
			code |= RIGHT
		}

		if point.Y < rect.MinY {
			code |= BOTTOM
		} else if point.Y > rect.MaxY {
			code |= TOP
		}

		return code
	}

	// Compute outcodes
	outcode0 := computeOutCode(start)
	outcode1 := computeOutCode(end)

	// Trivial reject and accept
	if outcode0 == 0 || outcode1 == 0 {
		return true
	}

	if outcode0&outcode1 != 0 {
		return false
	}

	return true // Line potentially intersects
}
