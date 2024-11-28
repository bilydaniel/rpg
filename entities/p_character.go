package entities

import (
	"bilydaniel/rpg/config"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
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
	fmt.Println(p.Selected)
	pcolor := color.RGBA{}
	if p.Name == "red" {
		pcolor = color.RGBA{255, 0, 0, 125}
	}
	if p.Name == "green" {
		pcolor = color.RGBA{0, 255, 0, 125}
	}
	if p.Name == "blue" {
		pcolor = color.RGBA{0, 0, 255, 125}
	}
	if p.Name == "yellow" {
		pcolor = color.RGBA{255, 255, 0, 125}
	}

	if p.Selected {
		pcolor = color.RGBA{255, 255, 255, 125}
	}

	vector.DrawFilledCircle(screen, float32(p.X)-float32(camera.X), float32(p.Y)-float32(camera.Y), 8*float32(camera.Scale), pcolor, false)
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
