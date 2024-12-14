package entities

import (
	"bilydaniel/rpg/config"
	"fmt"
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
		Sprite: &CircleSprite{
			X:   50,
			Y:   50,
			R:   r,
			R_2: r_2,
		},
		Character: Character{
			Speed:          1.0,
			TurnSpeed:      0.1,
			AngleTolerance: 0.0,
		},
	}
	if name == "red" {
		pcharacter.SetPosition(0, 0)
	} else if name == "green" {
		pcharacter.SetY(100)
	} else if name == "blue" {
		pcharacter.SetX(100)
	} else if name == "yellow" {
		pcharacter.SetY(150)
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
	if p.DestinationX != nil && p.DestinationY != nil {
		dx := *p.DestinationX - p.GetX()
		dy := *p.DestinationY - p.GetY()

		dist := math.Hypot(dx, dy)
		p.DestinationDist = &dist

		if p.Angle < p.AngleDestination-p.AngleTolerance {
			p.Angle += p.TurnSpeed
		}

		if p.Angle > p.AngleDestination+p.AngleTolerance {
			p.Angle -= p.TurnSpeed
		}

		if dist > config.Tolerance {
			dxnorm := dx / dist
			dynorm := dy / dist

			p.SetPosition(p.GetX()+dxnorm*p.Speed, p.GetY()+dynorm*p.Speed)
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
	opts.GeoM.Translate(-8, -8)
	opts.GeoM.Rotate(p.Angle)
	opts.GeoM.Translate(8, 8)
	// TODO SWITCH ASSET IF IN A CERTAIN ANGLE, good for now

	opts.GeoM.Translate(float64(p.GetX()), float64(p.GetY()))
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

	switch value := p.Sprite.(type) {
	case *CircleSprite:
		worldx := (float64(x) / camera.Scale) + camera.X
		worldy := (float64(y) / camera.Scale) + camera.Y
		dx := worldx - p.GetX() - config.TileSize/2
		dy := worldy - p.GetY() - config.TileSize/2

		distance := math.Pow(dx, 2) + math.Pow(dy, 2)
		if distance <= value.R_2 {
			return true
		}
	case *SquareSprite:
	//TODO
	default:
		fmt.Errorf("Unknown collision type")
	}

	return false
}

func (p *PCharacter) RectCollision(startx int, starty int, endx int, endy int, camera config.Camera) bool {
	//TODO try to understand this algorithm a bit more, draw it

	circleCollision, ok := p.Sprite.(*CircleSprite)
	if !ok {
		fmt.Errorf("Unknown collision type")
		return false
	}

	startx = startx + int(camera.X)
	starty = starty + int(camera.Y)
	endx = endx + int(camera.X)
	endy = endy + int(camera.Y)

	charx := p.GetX() + config.TileSize/2
	chary := p.GetY() + config.TileSize/2

	rectLeft := math.Min(float64(startx), float64(endx))
	rectRight := math.Max(float64(startx), float64(endx))
	rectTop := math.Min(float64(starty), float64(endy))
	rectBottom := math.Max(float64(starty), float64(endy))

	closestx := math.Max(rectLeft, math.Min(charx, rectRight))
	closesty := math.Max(rectTop, math.Min(chary, rectBottom))

	distancex := closestx - charx
	distancey := closesty - chary

	distance := math.Hypot(distancex, distancey)

	return distance <= circleCollision.R
}

func (p *PCharacter) SetDestination(x int, y int, camera config.Camera) {
	if p.Selected {
		worldx, worldy := camera.ToWorld(float64(x), float64(y))

		p.DestinationX = &worldx
		p.DestinationY = &worldy

		dx := *p.DestinationX - p.GetX()
		dy := *p.DestinationY - p.GetY()
		p.AngleDestination = math.Atan2(-dy, dx)
	}
}
