package main

import (
	"bilydaniel/rpg/config"
	"bilydaniel/rpg/entities"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	PCharacters []*entities.PCharacter
	Camera      config.Camera
	World       World
	Assets      Assets
}

func initGame() Game {
	return Game{
		PCharacters: entities.InitPCharacters(),
		World:       InitWorld(),
	}
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.World.CurrentTilemap.Draw(screen, g.Camera)

	for _, character := range g.PCharacters {
		if character != nil {
			character.Draw(screen, g.Camera)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.ScreenW, config.ScreenH
}

func main() {
	ebiten.SetWindowSize(config.WindowW, config.WindowH)
	ebiten.SetWindowTitle(config.GameName)
	game := initGame()
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
