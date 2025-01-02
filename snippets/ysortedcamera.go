package game

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

// Renderable represents any object that can be drawn
type Renderable interface {
	GetY() float64
	Draw(screen *ebiten.Image, camera *Camera)
	IsStatic() bool
}

// Camera handles view transformation and sorting
type Camera struct {
	X, Y          float64
	Width, Height int
	Scale         float64

	// Separate static and dynamic objects
	staticObjects  []Renderable
	dynamicObjects []Renderable

	// Cache for sorted static objects
	sortedStaticCache []Renderable

	// Viewport bounds for culling
	bounds struct {
		minX, maxX float64
		minY, maxY float64
	}
}

// NewCamera creates a new camera instance
func NewCamera(width, height int) *Camera {
	return &Camera{
		Width:          width,
		Height:         height,
		Scale:          1.0,
		staticObjects:  make([]Renderable, 0, 1000),
		dynamicObjects: make([]Renderable, 0, 100),
	}
}

// AddObject adds an object to the appropriate list
func (c *Camera) AddObject(obj Renderable) {
	if obj.IsStatic() {
		c.staticObjects = append(c.staticObjects, obj)
		// Invalidate cache when adding new static objects
		c.sortedStaticCache = nil
	} else {
		c.dynamicObjects = append(c.dynamicObjects, obj)
	}
}

// RemoveObject removes an object from the appropriate list
func (c *Camera) RemoveObject(obj Renderable) {
	if obj.IsStatic() {
		for i, o := range c.staticObjects {
			if o == obj {
				c.staticObjects = append(c.staticObjects[:i], c.staticObjects[i+1:]...)
				c.sortedStaticCache = nil
				break
			}
		}
	} else {
		for i, o := range c.dynamicObjects {
			if o == obj {
				c.dynamicObjects = append(c.dynamicObjects[:i], c.dynamicObjects[i+1:]...)
				break
			}
		}
	}
}

// UpdateBounds updates the camera's visible area for culling
func (c *Camera) UpdateBounds() {
	c.bounds.minX = c.X - float64(c.Width)/2/c.Scale
	c.bounds.maxX = c.X + float64(c.Width)/2/c.Scale
	c.bounds.minY = c.Y - float64(c.Height)/2/c.Scale
	c.bounds.maxY = c.Y + float64(c.Height)/2/c.Scale
}

// IsVisible checks if an object is within the camera's view
func (c *Camera) IsVisible(x, y, width, height float64) bool {
	return x+width >= c.bounds.minX &&
		x <= c.bounds.maxX &&
		y+height >= c.bounds.minY &&
		y <= c.bounds.maxY
}

// GetTransform returns the camera transform for the current frame
func (c *Camera) GetTransform() ebiten.GeoM {
	geom := ebiten.GeoM{}
	geom.Scale(c.Scale, c.Scale)
	geom.Translate(-c.X*c.Scale+float64(c.Width)/2, -c.Y*c.Scale+float64(c.Height)/2)
	return geom
}

// Draw handles rendering all objects with proper Y-sorting
func (c *Camera) Draw(screen *ebiten.Image) {
	c.UpdateBounds()

	// Initialize or update static cache if needed
	if c.sortedStaticCache == nil {
		c.sortedStaticCache = make([]Renderable, len(c.staticObjects))
		copy(c.sortedStaticCache, c.staticObjects)
		sort.Slice(c.sortedStaticCache, func(i, j int) bool {
			return c.sortedStaticCache[i].GetY() < c.sortedStaticCache[j].GetY()
		})
	}

	// Sort dynamic objects
	sort.Slice(c.dynamicObjects, func(i, j int) bool {
		return c.dynamicObjects[i].GetY() < c.dynamicObjects[j].GetY()
	})

	// Merge sorted static and dynamic objects while drawing
	staticIdx := 0
	dynamicIdx := 0

	for staticIdx < len(c.sortedStaticCache) || dynamicIdx < len(c.dynamicObjects) {
		var nextObject Renderable

		// Determine which object to draw next
		if staticIdx >= len(c.sortedStaticCache) {
			nextObject = c.dynamicObjects[dynamicIdx]
			dynamicIdx++
		} else if dynamicIdx >= len(c.dynamicObjects) {
			nextObject = c.sortedStaticCache[staticIdx]
			staticIdx++
		} else if c.sortedStaticCache[staticIdx].GetY() <= c.dynamicObjects[dynamicIdx].GetY() {
			nextObject = c.sortedStaticCache[staticIdx]
			staticIdx++
		} else {
			nextObject = c.dynamicObjects[dynamicIdx]
			dynamicIdx++
		}

		// Draw the object
		nextObject.Draw(screen, c)
	}
}

// Example implementation of a renderable object
type GameObject struct {
	X, Y   float64
	Image  *ebiten.Image
	Static bool
}

func (g *GameObject) GetY() float64 {
	return g.Y
}

func (g *GameObject) IsStatic() bool {
	return g.Static
}

func (g *GameObject) Draw(screen *ebiten.Image, camera *Camera) {
	// Skip if not visible
	if !camera.IsVisible(g.X, g.Y, float64(g.Image.Bounds().Dx()), float64(g.Image.Bounds().Dy())) {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM = camera.GetTransform()
	op.GeoM.Translate(g.X, g.Y)
	screen.DrawImage(g.Image, op)
}
