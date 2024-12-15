package pathfinding

import (
	"math"
	"sort"
)

// Point represents a 2D coordinate
type Point struct {
	X, Y float64
}

// Node represents a node in the pathfinding grid
type Node struct {
	Pos      Point
	G, H, F  float64
	Parent   *Node
	Walkable bool
}

// PathFinder manages pathfinding and smooth movement
type PathFinder struct {
	Grid          [][]*Node
	Width, Height int
}

// NewPathFinder creates a new PathFinder with a given grid size
func NewPathFinder(width, height int) *PathFinder {
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
	return &PathFinder{
		Grid:   grid,
		Width:  width,
		Height: height,
	}
}

// Distance calculates Euclidean distance between two points
func (pf *PathFinder) Distance(a, b Point) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// LineOfSight checks if there's a clear path between two points
func (pf *PathFinder) LineOfSight(start, end Point) bool {
	//TODO look at exlanation
	// Bresenham's line algorithm for visibility checking
	x0, y0 := int(start.X), int(start.Y)
	x1, y1 := int(end.X), int(end.Y)

	dx := abs(x1 - x0)
	dy := abs(y1 - y0)

	sx, sy := 1, 1
	if x0 >= x1 {
		sx = -1
	}
	if y0 >= y1 {
		sy = -1
	}

	err := dx - dy

	for {
		if !pf.Grid[y0][x0].Walkable {
			return false
		}

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}

	return true
}

// ThetaStar implements the Theta* pathfinding algorithm
func (pf *PathFinder) ThetaStar(start, goal Point) []Point {
	startNode := pf.Grid[int(start.Y)][int(start.X)]
	goalNode := pf.Grid[int(goal.Y)][int(goal.X)]

	openSet := make([]*Node, 0)
	closedSet := make(map[*Node]bool)

	startNode.G = 0
	startNode.H = pf.Distance(start, goal)
	startNode.F = startNode.H
	startNode.Parent = nil

	openSet = append(openSet, startNode)

	for len(openSet) > 0 {
		// Sort open set by F score
		sort.Slice(openSet, func(i, j int) bool {
			return openSet[i].F < openSet[j].F
		})

		current := openSet[0]

		// Reached goal
		if current == goalNode {
			return pf.ReconstructPath(current)
		}

		openSet = openSet[1:]
		closedSet[current] = true

		// Check neighboring nodes
		neighbors := pf.GetNeighbors(current)
		for _, neighbor := range neighbors {
			if closedSet[neighbor] || !neighbor.Walkable {
				continue
			}

			// Calculate potential new G score
			//TODO put in else so it doesent get computed twice
			tentativeG := current.G + pf.Distance(current.Pos, neighbor.Pos)

			// Rewire path if line of sight is clear
			if current.Parent != nil && pf.LineOfSight(current.Parent.Pos, neighbor.Pos) {
				tentativeG = current.Parent.G + pf.Distance(current.Parent.Pos, neighbor.Pos)
			}

			// Update if better path found
			if tentativeG < neighbor.G || !contains(openSet, neighbor) {
				neighbor.Parent = current
				neighbor.G = tentativeG
				neighbor.H = pf.Distance(neighbor.Pos, goal)
				neighbor.F = neighbor.G + neighbor.H

				if !contains(openSet, neighbor) {
					openSet = append(openSet, neighbor)
				}
			}
		}
	}

	return nil // No path found
}

// ReconstructPath rebuilds the path from goal to start
func (pf *PathFinder) ReconstructPath(node *Node) []Point {
	path := make([]Point, 0)
	current := node

	for current != nil {
		path = append([]Point{current.Pos}, path...)
		current = current.Parent
	}

	return path
}

// GetNeighbors returns walkable neighboring nodes
func (pf *PathFinder) GetNeighbors(node *Node) []*Node {
	neighbors := make([]*Node, 0, 8)
	directions := [][]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for _, dir := range directions {
		newX := int(node.Pos.X) + dir[0]
		newY := int(node.Pos.Y) + dir[1]

		if newX >= 0 && newX < pf.Width && newY >= 0 && newY < pf.Height {
			neighbors = append(neighbors, pf.Grid[newY][newX])
		}
	}

	return neighbors
}

// Helper functions
func contains(slice []*Node, item *Node) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// SmoothPath reduces path points while maintaining line of sight
func (pf *PathFinder) SmoothPath(path []Point) []Point {
	if len(path) <= 2 {
		return path
	}

	smoothedPath := []Point{path[0]}
	current := 0

	for current < len(path)-1 {
		next := current + 1

		// Look ahead to find furthest visible point
		for next+1 < len(path) && pf.LineOfSight(path[current], path[next+1]) {
			next++
		}

		smoothedPath = append(smoothedPath, path[next])
		current = next
	}

	return smoothedPath
}

// Movement represents the current movement state
type Movement struct {
	CurrentPos   Point
	TargetPos    Point
	Path         []Point
	Speed        float64
	PathProgress int
}

// UpdateMovement moves the entity along the smoothed path
func (m *Movement) UpdateMovement() bool {
	if m.PathProgress >= len(m.Path)-1 {
		return false // Path completed
	}

	// Calculate direction to next point
	target := m.Path[m.PathProgress+1]
	dx := target.X - m.CurrentPos.X
	dy := target.Y - m.CurrentPos.Y

	// Normalize direction
	length := math.Sqrt(dx*dx + dy*dy)
	if length == 0 {
		return true
	}

	// Move towards next point
	m.CurrentPos.X += (dx / length) * m.Speed
	m.CurrentPos.Y += (dy / length) * m.Speed

	// Check if reached or passed next point
	if math.Abs(m.CurrentPos.X-target.X) <= m.Speed &&
		math.Abs(m.CurrentPos.Y-target.Y) <= m.Speed {
		m.CurrentPos = target
		m.PathProgress++
	}

	return true
}
