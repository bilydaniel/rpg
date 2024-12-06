package main

import (
	"bilydaniel/rpg/assets"
	"bilydaniel/rpg/config"
	"bilydaniel/rpg/entities"
	"bilydaniel/rpg/utils"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	PCharacters []*entities.PCharacter
	Camera      config.Camera
	World       *World
	Assets      assets.Assets
	Drag        utils.Drag
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
		Drag:        utils.Drag{},
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

	//TODO gonna need to change clicking, think it through
	//Probably should split it up and not generalize
	// make a menu first so i have an idea about the rest clicking stuff???
	// probably gonna need some soft of ID system

	// SELECT
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		for _, pchar := range g.PCharacters {
			pchar.Selected = false
		}
		//TODO how to handle other clickable stuff?
		for _, pchar := range g.PCharacters {
			mx, my := ebiten.CursorPosition()
			if pchar.ClickCollision(mx, my, g.Camera) {
				pchar.OnClick()
				break
			}
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.Drag.Dragging = true
		g.Drag.Startx, g.Drag.Starty = ebiten.CursorPosition()
	}

	if g.Drag.Dragging {
		g.Drag.Endx, g.Drag.Endy = ebiten.CursorPosition()
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && g.Drag.Dragging {
		g.Drag.Dragging = false
		//TODO SELECT ALL THE CHARACTERS IN SQUARE
		for _, pchar := range g.PCharacters {
			if pchar.RectCollision(g.Drag.Startx, g.Drag.Starty, g.Drag.Endx, g.Drag.Endy, g.Camera) {
				pchar.Selected = true
			}
		}

	}
	// MOVEMENT
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		for _, pchar := range g.PCharacters {
			mx, my := ebiten.CursorPosition()
			pchar.SetDestination(mx, my, g.Camera)
		}
	}

	for _, pchar := range g.PCharacters {
		pchar.Update()
	}

	return nil
}

// TODO CHECK ALL THE NILLS

func (g *Game) Draw(screen *ebiten.Image) {
	//TODO load all the assets only once
	if g.World != nil && g.World.CurrentTilemap != nil {
		g.World.CurrentTilemap.Draw(screen, g.Camera, g.Assets)
	}

	for _, character := range g.PCharacters {
		if character != nil {
			character.Draw(screen, g.Camera)
		}
	}

	g.Drag.Draw(screen, &g.Camera)
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
