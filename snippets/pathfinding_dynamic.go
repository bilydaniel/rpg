package game

import (
	"math"
)

// Position represents a 2D coordinate
type Position struct {
	X, Y int
}

// Character represents a moving entity in the game
type Character struct {
	Pos    Position
	Path   []Position
	Moving bool
	Speed  float64 // tiles per update
	Size   int     // collision size in tiles
}

// World handles game state and collision detection
type World struct {
	StaticObstacles map[Position]bool
	DynamicEntities []*Character
	Width, Height   int
}

// IsPositionOccupied checks if a position is blocked by static or dynamic obstacles
func (w *World) IsPositionOccupied(pos Position) bool {
	// Check static obstacles
	if w.StaticObstacles[pos] {
		return true
	}

	// Check dynamic entities
	for _, entity := range w.DynamicEntities {
		if entity.Pos.X == pos.X && entity.Pos.Y == pos.Y {
			return true
		}
	}

	return false
}

// UpdateCharacterMovement handles character movement with collision detection
func (w *World) UpdateCharacterMovement(char *Character) {
	if !char.Moving || len(char.Path) == 0 {
		return
	}

	nextPos := char.Path[0]

	// Check if next position is occupied by dynamic obstacle
	if w.IsPositionOccupied(nextPos) {
		// Recalculate path avoiding current dynamic obstacles
		newPath := w.RecalculatePath(char.Pos, char.Path[len(char.Path)-1])
		if len(newPath) > 0 {
			char.Path = newPath
		} else {
			// No valid path found, wait
			return
		}
	}

	// Move towards next position
	dx := float64(nextPos.X - char.Pos.X)
	dy := float64(nextPos.Y - char.Pos.Y)
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance <= char.Speed {
		// Reached next position
		char.Pos = nextPos
		char.Path = char.Path[1:]
		if len(char.Path) == 0 {
			char.Moving = false
		}
	} else {
		// Move towards next position
		moveX := dx / distance * char.Speed
		moveY := dy / distance * char.Speed
		char.Pos.X += int(moveX)
		char.Pos.Y += int(moveY)
	}
}

// RecalculatePath uses A* to find a new path avoiding current obstacles
func (w *World) RecalculatePath(start, end Position) []Position {
	// Implement A* pathfinding here, considering both static and dynamic obstacles
	// This should use your existing A* implementation but with IsPositionOccupied()
	// for checking valid moves
	return nil // Placeholder
}

// Update is called every game tick
func (w *World) Update() error {
	// Update all characters
	for _, char := range w.DynamicEntities {
		w.UpdateCharacterMovement(char)
	}
	return nil
}
