package main

import (
	"bilydaniel/rpg/assets"
	"bilydaniel/rpg/config"
	"bilydaniel/rpg/entities"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	PCharacters []*entities.PCharacter
	Camera      config.Camera
	World       *World
	Assets      assets.Assets
}

func initGame() (*Game, error) {
	assets, err := assets.InitAssets()
	if err != nil {
		return nil, err
	}

	world, err := InitWorld(assets)
	if err != nil {
		return nil, err
	}
	return &Game{
		PCharacters: entities.InitPCharacters(),
		World:       world,
		Assets:      assets,
	}, nil
}

func (g *Game) Update() error {
	return nil
}

// TODO CHECK ALL THE NILLS

func (g *Game) Draw(screen *ebiten.Image) {
	if g.World != nil && g.World.CurrentTilemap != nil {
		g.World.CurrentTilemap.Draw(screen, g.Camera, g.Assets)
	}

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
	game, err := initGame()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
