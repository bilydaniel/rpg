package entities

import (
	"bilydaniel/rpg/config"
	"bilydaniel/rpg/utils"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type PCharacter struct {
	Name            string
	Selected        bool
	Destination     *utils.Node
	DestinationX    *float64
	DestinationY    *float64
	DestinationDist *float64
	Path            []utils.Node
	PathProgress    int
	Sprite
	Character
}

func InitPCharacter(name string) (*PCharacter, error) {
	r := 8.0
	r_2 := r * r

	image, _, err := ebitenutil.NewImageFromFile("assets/images/cavegirl.png")

	if err != nil {
		return nil, err
	}
	pcharacter := PCharacter{
		Name:     name,
		Selected: false,
		Sprite: &CircleSprite{
			X:   0,
			Y:   0,
			R:   r,
			R_2: r_2,
			Img: image,
		},
		Character: Character{
			Speed: 1 / 30.0,
		},
		Path: []utils.Node{},
	}
	if name == "red" {
		pcharacter.SetPosition(0, 0)
	} else if name == "green" {
		pcharacter.SetY(2)
	} else if name == "blue" {
		pcharacter.SetX(4)
	}
	config.AddClicker(&pcharacter)

	return &pcharacter, nil
}

func InitPCharacters() ([]*PCharacter, error) {
	characters := []*PCharacter{}
	for i := 0; i < 3; i++ {
		character, err := InitPCharacter(config.PlayableCharacters[i])
		if err != nil {
			return nil, err
		}
		characters = append(characters, character)
	}
	return characters, nil
}

func (p *PCharacter) Update(level Level) {
	if len(p.Path) > 0 {
		if p.PathProgress > len(p.Path)-1 {
			p.Path = []utils.Node{}
			p.PathProgress = 0
			return
		}

		target := p.Path[p.PathProgress]
		if level.OccupiedTile(&target) {
			p.ResetWalking()
			return
		}

		dx := float64(target.X) - p.GetX()
		dy := float64(target.Y) - p.GetY()
		dist := math.Hypot(dx, dy)

		if dist == 0 {
			return
		}

		dxnorm := dx / dist
		dynorm := dy / dist

		p.SetPosition(p.GetX()+dxnorm*p.Speed, p.GetY()+dynorm*p.Speed)

		if math.Abs(p.GetX()-float64(target.X)) <= p.Speed && math.Abs(p.GetY()-float64(target.Y)) <= p.Speed {
			p.SetX(float64(target.X))
			p.SetY(float64(target.Y))
			p.PathProgress++
		}
	}
}

func (p *PCharacter) Draw(screen *ebiten.Image, camera config.Camera) {

	//pcolor := color.RGBA{0, 255, 0, 255}

	opts := ebiten.DrawImageOptions{}

	camera.WorldToScreenGeom(&opts, int(p.GetX()*config.TileSize), int(p.GetY()*config.TileSize))
	circle, _, err := ebitenutil.NewImageFromFile("assets/images/circle.png")
	if err != nil {
		return
	}
	//TODO use shaders for this????
	if p.Selected {
		screen.DrawImage(circle, &opts)
	}
	screen.DrawImage(p.Image(), &opts)

	if p.DestinationX != nil && p.DestinationY != nil {
		if p.DestinationDist != nil && *p.DestinationDist > float64(config.Tolerance) {
			destinationImage, _, err := ebitenutil.NewImageFromFile("assets/images/target.png")
			if err == nil {
				opts.GeoM.Reset()
				camera.WorldToScreenGeom(&opts, int(*p.DestinationX)-config.TileSize/2, int(*p.DestinationY)-config.TileSize/2)
				screen.DrawImage(destinationImage, &opts)
			}
		}

	}

	//path
	if len(p.Path) != 0 {
		opt := ebiten.GeoM{}
		opt.Translate(-camera.X*camera.Speed, -camera.Y*camera.Speed)
		opt.Scale(camera.Scale, camera.Scale)
		if p.PathProgress < 1 {

			//TODO add walkable=false
			x0, y0 := p.GetX()*config.TileSize+config.TileSize/2, p.GetY()*config.TileSize+config.TileSize/2
			x1, y1 := float64(p.Path[0].X)*config.TileSize+config.TileSize/2, float64(p.Path[0].Y)*config.TileSize+config.TileSize/2

			sx, sy := opt.Apply(x0, y0)
			ex, ey := opt.Apply(x1, y1)
			vector.StrokeLine(screen, float32(sx), float32(sy), float32(ex), float32(ey), 1, color.RGBA{255, 0, 0, 255}, false)
		}

		for i, node := range p.Path {
			if i >= p.PathProgress {
				if i < len(p.Path)-1 {
					x0, y0 := float64(node.X*config.TileSize+config.TileSize/2), float64(node.Y*config.TileSize+config.TileSize/2)
					x1, y1 := float64(p.Path[i+1].X*config.TileSize+config.TileSize/2), float64(p.Path[i+1].Y*config.TileSize+config.TileSize/2)

					sx, sy := opt.Apply(x0, y0)
					ex, ey := opt.Apply(x1, y1)

					vector.StrokeLine(screen, float32(sx), float32(sy), float32(ex), float32(ey), 1, color.RGBA{255, 0, 0, 255}, false)
				}
			}
		}
	}
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

		worldx, worldy := camera.ScreenToWorld(float64(x), float64(y))
		dx := worldx - p.GetX()*config.TileSize - config.TileSize/2
		dy := worldy - p.GetY()*config.TileSize - config.TileSize/2

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
	worldx, worldy := camera.ScreenToWorld(float64(startx), float64(starty))
	startx = int(worldx)
	starty = int(worldy)

	worldx, worldy = camera.ScreenToWorld(float64(endx), float64(endy))
	endx = int(worldx)
	endy = int(worldy)

	charx := p.GetX()*config.TileSize + config.TileSize/2
	chary := p.GetY()*config.TileSize + config.TileSize/2

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

func (p *PCharacter) ResetWalking() {
	p.Path = []utils.Node{}
	p.PathProgress = 0
}
