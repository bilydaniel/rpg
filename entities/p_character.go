package entities

import (
	"bilydaniel/rpg/config"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type PCharacter struct {
	Name     string
	Selected bool
	Sprite
	Character
}

func InitPCharacter(name string) *PCharacter {
	pcharacter := PCharacter{
		Name:     name,
		Selected: false,
		Sprite: Sprite{
			X: 50,
			Y: 50,
		},
		Character: Character{},
	}
	if name == "red" {
		pcharacter.X = 20
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

func (p *PCharacter) Draw(screen *ebiten.Image, camera config.Camera) {
	pcolor := color.RGBA{0, 255, 0, 255}
	tile := &ebiten.Image{}
	var err error
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
	//TODO weird
	if err != nil {
		return
	}

	opts := ebiten.DrawImageOptions{}
	if p.Selected {
		vector.StrokeCircle(screen, float32(camera.Scale)*(float32(p.X)+8-float32(camera.X)), float32(camera.Scale)*float32(p.Y)+8-float32(camera.Y), 8*float32(camera.Scale), 1, pcolor, true)
		//opts.ColorScale.SetR(255)
	}
	opts.GeoM.Translate(float64(p.X), float64(p.Y))
	opts.GeoM.Translate(-camera.X, -camera.Y)
	opts.GeoM.Scale(camera.Scale, camera.Scale)

	screen.DrawImage(tile, &opts)
	//TODO make an image, weird
}

func (p *PCharacter) OnClick() {
	if p.Selected {
		p.Selected = false
	} else {
		p.Selected = true
	}
}

func (p *PCharacter) ClickCollision(x int, y int, camera config.Camera) bool {
	return true
}
