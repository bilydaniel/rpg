package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// TileElevation represents a tile's height information
type TileElevation struct {
	BaseHeight   int     // Base elevation level (0 = ground)
	IsRamp       bool    // Whether this tile is a ramp/stairs
	RampDir      int     // Direction of elevation change (0=N, 1=E, 2=S, 3=W)
	ShadowLength float64 // Length of shadow to cast
}

// ElevationMap manages the height data for the game world
type ElevationMap struct {
	Width, Height int
	Elevations    [][]*TileElevation
	ShadowColor   Color
}

// NewElevationMap creates a new elevation map
func NewElevationMap(width, height int) *ElevationMap {
	em := &ElevationMap{
		Width:       width,
		Height:      height,
		ShadowColor: Color{R: 0, G: 0, B: 0, A: 80}, // Semi-transparent black
	}

	em.Elevations = make([][]*TileElevation, height)
	for y := range em.Elevations {
		em.Elevations[y] = make([]*TileElevation, width)
		for x := range em.Elevations[y] {
			em.Elevations[y][x] = &TileElevation{
				BaseHeight: 0,
				IsRamp:     false,
			}
		}
	}

	return em
}

// GetElevationAt returns the elevation at given coordinates
func (em *ElevationMap) GetElevationAt(x, y int) *TileElevation {
	if x < 0 || y < 0 || x >= em.Width || y >= em.Height {
		return nil
	}
	return em.Elevations[y][x]
}

// RenderTile draws a tile with elevation effects
func (em *ElevationMap) RenderTile(screen *ebiten.Image, tileImage *ebiten.Image, x, y int, tileSize int) {
	elevation := em.GetElevationAt(x, y)
	if elevation == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}

	// Base position
	baseX := float64(x * tileSize)
	baseY := float64(y * tileSize)

	// Apply elevation offset (moves tiles up for higher elevations)
	elevationOffset := float64(elevation.BaseHeight * 4) // 4 pixels per level
	baseY -= elevationOffset

	// Position the tile
	op.GeoM.Translate(baseX, baseY)

	// Draw shadows for elevated tiles
	if elevation.BaseHeight > 0 {
		shadowImg := ebiten.NewImage(tileSize, int(elevation.ShadowLength))
		shadowImg.Fill(em.ShadowColor)

		shadowOp := &ebiten.DrawImageOptions{}
		shadowOp.GeoM.Translate(baseX, baseY+float64(tileSize))
		screen.DrawImage(shadowImg, shadowOp)
	}

	// Draw ramp overlay if this is a ramp tile
	if elevation.IsRamp {
		rampOverlay := em.createRampOverlay(tileSize, elevation.RampDir)
		rampOp := &ebiten.DrawImageOptions{}
		rampOp.GeoM.Translate(baseX, baseY)
		screen.DrawImage(rampOverlay, rampOp)
	}

	// Draw the base tile
	screen.DrawImage(tileImage, op)
}

// createRampOverlay creates a visual overlay for ramps
func (em *ElevationMap) createRampOverlay(tileSize int, direction int) *ebiten.Image {
	overlay := ebiten.NewImage(tileSize, tileSize)

	// Create a gradient effect based on direction
	for y := 0; y < tileSize; y++ {
		for x := 0; x < tileSize; x++ {
			var alpha uint8
			switch direction {
			case 0: // North
				alpha = uint8(float64(y) / float64(tileSize) * 128)
			case 1: // East
				alpha = uint8(float64(x) / float64(tileSize) * 128)
			case 2: // South
				alpha = uint8((1 - float64(y)/float64(tileSize)) * 128)
			case 3: // West
				alpha = uint8((1 - float64(x)/float64(tileSize)) * 128)
			}
			overlay.Set(x, y, Color{R: 255, G: 255, B: 255, A: alpha})
		}
	}

	return overlay
}

// GameObject extension for elevation support
type GameObject struct {
	X, Y    float64
	Z       int // Elevation level
	Image   *ebiten.Image
	ElevMap *ElevationMap
}

// Draw draws the game object with elevation considerations
func (g *GameObject) Draw(screen *ebiten.Image) {
	tileX, tileY := int(g.X), int(g.Y)
	elevation := g.ElevMap.GetElevationAt(tileX, tileY)

	if elevation == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}

	// Calculate visual position including elevation
	visualX := g.X
	visualY := g.Y - float64(g.Z*4) // 4 pixels per elevation level

	// If on a ramp, interpolate height
	if elevation.IsRamp {
		progress := 0.0
		switch elevation.RampDir {
		case 0: // North
			progress = 1.0 - (g.Y - float64(tileY))
		case 1: // East
			progress = g.X - float64(tileX)
		case 2: // South
			progress = g.Y - float64(tileY)
		case 3: // West
			progress = 1.0 - (g.X - float64(tileX))
		}
		visualY -= progress * 4 // Smooth elevation change
	}

	op.GeoM.Translate(visualX, visualY)
	screen.DrawImage(g.Image, op)
}
