package utils

import (
	"bilydaniel/rpg/config"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Drag struct {
	Startx   int
	Starty   int
	Endx     int
	Endy     int
	Dragging bool
}

func (drag *Drag) Draw(screen *ebiten.Image, camera *config.Camera) {
	if drag.Dragging {
		vector.StrokeRect(screen, float32(drag.Startx), float32(drag.Starty), float32(drag.Endx-drag.Startx), float32(drag.Endy-drag.Starty), 0.5, color.RGBA{0, 255, 0, 125}, true)

	}
}
