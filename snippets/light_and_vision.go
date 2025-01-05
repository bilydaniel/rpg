package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// LightLevel represents the amount of light at a tile (0.0 to 1.0)
type LightLevel float64

// Tile extends your existing tile structure with lighting properties
type Tile struct {
	X, Y        int
	Walkable    bool
	BlocksLight bool // Whether this tile blocks light (walls, etc.)
	LightLevel  LightLevel
	Visible     bool // Whether the tile is currently visible to the player
	Explored    bool // Whether the tile has ever been seen by the player
}

// LightSource represents a source of light in the game
type LightSource struct {
	X, Y      int
	Intensity float64 // Base brightness (1.0 = full brightness)
	Radius    float64 // How far the light spreads
}

// World stores the game state
type World struct {
	Tiles        [][]*Tile
	Width        int
	Height       int
	LightSources []LightSource
}

// CalculateVisibility determines which tiles are visible from a point
func (w *World) CalculateVisibility(viewerX, viewerY int, viewDistance float64) {
	// Reset visibility (but keep explored status)
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			w.Tiles[y][x].Visible = false
		}
	}

	// Cast rays in a 360-degree arc
	for angle := 0.0; angle < 360.0; angle += 0.5 {
		w.castRay(viewerX, viewerY, angle, viewDistance)
	}
}

// castRay performs raycasting for visibility calculation
func (w *World) castRay(startX, startY int, angle float64, maxDist float64) {
	// Convert angle to radians
	radians := angle * math.Pi / 180.0

	// Calculate direction vector
	dx := math.Cos(radians)
	dy := math.Sin(radians)

	// Start at viewer position
	x := float64(startX)
	y := float64(startY)

	// Step along the ray
	for dist := 0.0; dist < maxDist; dist += 0.5 {
		// Get current tile coordinates
		tileX := int(math.Floor(x))
		tileY := int(math.Floor(y))

		// Check bounds
		if tileX < 0 || tileX >= w.Width || tileY < 0 || tileY >= w.Height {
			break
		}

		tile := w.Tiles[tileY][tileX]
		tile.Visible = true
		tile.Explored = true

		// Stop if we hit a wall
		if tile.BlocksLight {
			break
		}

		// Move along ray
		x += dx * 0.5
		y += dy * 0.5
	}
}

// UpdateLighting calculates lighting for the entire map
func (w *World) UpdateLighting() {
	// Reset light levels
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			w.Tiles[y][x].LightLevel = 0.0
		}
	}

	// Calculate contribution from each light source
	for _, light := range w.LightSources {
		w.calculateLightSource(light)
	}
}

// calculateLightSource computes lighting for a single light source
func (w *World) calculateLightSource(light LightSource) {
	// Calculate the square of the radius for performance
	radiusSquared := light.Radius * light.Radius

	// Calculate the affected area
	minX := int(math.Max(0, float64(light.X)-light.Radius))
	maxX := int(math.Min(float64(w.Width-1), float64(light.X)+light.Radius))
	minY := int(math.Max(0, float64(light.Y)-light.Radius))
	maxY := int(math.Min(float64(w.Height-1), float64(light.Y)+light.Radius))

	// Update light levels for affected tiles
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			// Calculate distance to light source
			dx := float64(x) - float64(light.X)
			dy := float64(y) - float64(light.Y)
			distSquared := dx*dx + dy*dy

			if distSquared <= radiusSquared {
				// Check if light is blocked
				if !w.isLightBlocked(light.X, light.Y, x, y) {
					// Calculate light intensity using inverse square law
					intensity := light.Intensity * (1.0 - math.Sqrt(distSquared)/light.Radius)
					w.Tiles[y][x].LightLevel += LightLevel(intensity)

					// Clamp light level to [0, 1]
					if w.Tiles[y][x].LightLevel > 1.0 {
						w.Tiles[y][x].LightLevel = 1.0
					}
				}
			}
		}
	}
}

// isLightBlocked checks if there are any blocking tiles between two points
func (w *World) isLightBlocked(x1, y1, x2, y2 int) bool {
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance == 0 {
		return false
	}

	// Normalize direction
	dx /= distance
	dy /= distance

	// Step along the line
	x := float64(x1)
	y := float64(y1)

	for i := 0.0; i < distance; i += 0.5 {
		tileX := int(math.Floor(x))
		tileY := int(math.Floor(y))

		if w.Tiles[tileY][tileX].BlocksLight {
			return true
		}

		x += dx * 0.5
		y += dy * 0.5
	}

	return false
}

// Example usage with Ebiten
type Game struct {
	world  *World
	player struct {
		X, Y           int
		VisionDistance float64
	}
}

func (g *Game) Update() error {
	// Update player position based on input...

	// Update visibility from player's position
	g.world.CalculateVisibility(g.player.X, g.player.Y, g.player.VisionDistance)

	// Update lighting (this could be done less frequently if performance is a concern)
	g.world.UpdateLighting()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw tiles based on visibility and light level
	for y := 0; y < g.world.Height; y++ {
		for x := 0; x < g.world.Width; x++ {
			tile := g.world.Tiles[y][x]

			if !tile.Explored {
				// Don't draw unexplored tiles
				continue
			}

			if tile.Visible {
				// Tile is currently visible - draw with lighting
				lightLevel := float64(tile.LightLevel)
				// Modify your tile drawing based on lightLevel (0.0 to 1.0)
				// Example: adjust color/alpha based on light level
				op := &ebiten.DrawImageOptions{}
				op.ColorM.Scale(lightLevel, lightLevel, lightLevel, 1)
				// Draw tile...
			} else if tile.Explored {
				// Tile was previously seen but not currently visible
				// Draw with darker/grayed out appearance
				op := &ebiten.DrawImageOptions{}
				op.ColorM.Scale(0.3, 0.3, 0.3, 1)
				// Draw tile...
			}
		}
	}
}
