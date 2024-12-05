package entities

import (
	"bilydaniel/rpg/config"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type PCharacter struct {
	Name            string
	Selected        bool
	DestinationX    *float64
	DestinationY    *float64
	DestinationDist *float64
	Sprite
	Character
}

func InitPCharacter(name string) *PCharacter {
	r := 8.0
	r_2 := r * r
	pcharacter := PCharacter{
		Name:     name,
		Selected: false,
		Sprite: Sprite{
			X:            50,
			Y:            50,
			ColliderType: Circle,
			R:            &r,
			R_2:          &r_2,
		},
		Character: Character{
			Speed: 3.0,
		},
	}
	if name == "red" {
		pcharacter.X = 0
		pcharacter.Y = 0
	} else if name == "green" {
		pcharacter.Y = 100
	} else if name == "blue" {
		pcharacter.X = 150
	} else if name == "yellow" {
		pcharacter.Y = 150
	}
	config.AddClicker(&pcharacter)
	return &pcharacter
}

func InitPCharacters() []*PCharacter {
	characters := []*PCharacter{}
	for i := 0; i < 4; i++ {
		characters = append(characters, InitPCharacter(config.PlayableCharacters[i]))
	}
	return characters
}

func (p *PCharacter) Update() {
	//TODO do FLOCKING behaviour
	if p.Selected {
		if p.DestinationX != nil && p.DestinationY != nil {
			dx := *p.DestinationX - p.X
			dy := *p.DestinationY - p.Y

			dist := math.Hypot(dx, dy)
			p.DestinationDist = &dist

			if dist > config.Tolerance {
				dxnorm := dx / dist
				dynorm := dy / dist

				p.X += dxnorm * p.Speed
				p.Y += dynorm * p.Speed
			}
		}
	}
}

func (p *PCharacter) Draw(screen *ebiten.Image, camera config.Camera) {

	//pcolor := color.RGBA{0, 255, 0, 255}
	tile := &ebiten.Image{}
	circle, _, err := ebitenutil.NewImageFromFile("assets/images/circle.png")
	//TODO ERROR HANDLE
	if err != nil {
		return
	}
	if p.Name == "red" {
		//pcolor = color.RGBA{255, 0, 0, 125}
		tile, _, err = ebitenutil.NewImageFromFile("assets/images/redchar.png")
	}
	if p.Name == "green" {
		//pcolor = color.RGBA{0, 255, 0, 125}
		tile, _, err = ebitenutil.NewImageFromFile("assets/images/greenchar.png")
	}
	if p.Name == "blue" {
		//pcolor = color.RGBA{0, 0, 255, 125}
		tile, _, err = ebitenutil.NewImageFromFile("assets/images/bluechar.png")
	}
	if p.Name == "yellow" {
		//pcolor = color.RGBA{255, 255, 0, 125}
		tile, _, err = ebitenutil.NewImageFromFile("assets/images/yellowchar.png")
	}
	//TODO ERROR HANDLE
	if err != nil {
		return
	}

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(p.X), float64(p.Y))
	opts.GeoM.Translate(-camera.X, -camera.Y)
	opts.GeoM.Scale(camera.Scale, camera.Scale)
	if p.Selected {
		screen.DrawImage(circle, &opts)

	}

	screen.DrawImage(tile, &opts)
	if p.DestinationX != nil && p.DestinationY != nil {
		if p.DestinationDist != nil && *p.DestinationDist > float64(config.Tolerance) {
			destinationImage, _, err := ebitenutil.NewImageFromFile("assets/images/target.png")
			if err == nil {
				opts.GeoM.Reset()
				//TODO make a function from this
				opts.GeoM.Translate(*p.DestinationX, *p.DestinationY)
				opts.GeoM.Translate(-camera.X-config.TileSize/2, -camera.Y-config.TileSize/2)
				opts.GeoM.Scale(camera.Scale, camera.Scale)
				screen.DrawImage(destinationImage, &opts)
			}
		}

	}
	//vector.StrokeCircle(screen, float32(p.X-camera.X+8), float32(p.Y-camera.Y+8), float32(*p.R), 1.0, color.RGBA{255, 0, 0, 125}, true)

}

func (p *PCharacter) OnClick() {
	if p.Selected {
		p.Selected = false
	} else {
		p.Selected = true
	}
}

func (p *PCharacter) ClickCollision(x int, y int, camera config.Camera) bool {
	if p.ColliderType == Square {

	} else if p.ColliderType == Circle {
		worldx := (float64(x) / camera.Scale) + camera.X
		worldy := (float64(y) / camera.Scale) + camera.Y
		dx := worldx - p.X - config.TileSize/2
		dy := worldy - p.Y - config.TileSize/2

		distance := math.Pow(dx, 2) + math.Pow(dy, 2)
		if p.R_2 != nil {
			if distance <= *p.R_2 {
				return true
			}
		}
	}

	return false
}

func (p *PCharacter) RectCollision(startx int, starty int, endx int, endy int, camera config.Camera) bool {

	return false
}

func (p *PCharacter) SetDestination(x int, y int, camera config.Camera) {
	if p.Selected {
		worldx, worldy := camera.ToWorld(float64(x), float64(y))
		p.DestinationX = &worldx
		p.DestinationY = &worldy
	}
}
