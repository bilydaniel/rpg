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
		Camera:      config.Camera{X: 0, Y: 0, Scale: 1.0}, //TODO make an init function
	}, nil
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.Camera.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.Camera.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.Camera.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.Camera.Y += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.Camera.Scale -= 0.01
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		g.Camera.Scale += 0.01
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		//TODO how to handle other clickable stuff?
		for _, pchar := range g.PCharacters {

		}

	}

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
