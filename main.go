package main

import (
	"bilydaniel/rpg/assets"
	"bilydaniel/rpg/config"
	"bilydaniel/rpg/entities"
	"bilydaniel/rpg/utils"
	"bilydaniel/rpg/world"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	PCharacters []*entities.PCharacter
	Camera      *config.Camera
	World       *world.World
	Drag        *utils.Drag
	Assets      *assets.Assets
}

func initGame() (*Game, error) {
	world, err := world.InitWorld()
	if err != nil {
		return nil, err
	}

	assets, err := assets.InitAssets()
	if err != nil {
		return nil, err
	}
	return &Game{
		PCharacters: entities.InitPCharacters(),
		World:       world,
		Camera:      &config.Camera{X: 0, Y: 0, Scale: 1.0}, //TODO make an init function
		Drag:        &utils.Drag{},
		Assets:      assets,
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
			if pchar.ClickCollision(mx, my, *g.Camera) {
				pchar.OnClick()
				break
			}
		}
	}

	// DRAGING
	// TODO combine with selecting
	if inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 5 && !g.Drag.Dragging {
		g.Drag.Dragging = true
		g.Drag.Startx, g.Drag.Starty = ebiten.CursorPosition()
	}

	if g.Drag.Dragging {
		g.Drag.Endx, g.Drag.Endy = ebiten.CursorPosition()
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && g.Drag.Dragging {
		g.Drag.Dragging = false
		for _, pchar := range g.PCharacters {
			if pchar.RectCollision(g.Drag.Startx, g.Drag.Starty, g.Drag.Endx, g.Drag.Endy, *g.Camera) {
				pchar.Selected = true
			}
		}

	}
	// MOVEMENT
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		for _, pchar := range g.PCharacters {
			mx, my := ebiten.CursorPosition()
			pchar.SetDestination(mx, my, *g.Camera)
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
	debug := "tps"
	if debug == "tps" {
		ebitenutil.DebugPrint(screen, strconv.Itoa(int(ebiten.ActualTPS())))
	}
	if debug == "fps" {
		ebitenutil.DebugPrint(screen, strconv.Itoa(int(ebiten.ActualFPS())))
	}
	if g.World != nil && g.World.CurrentLevel != nil {
		g.World.CurrentLevel.Draw(screen, g.Camera, *g.Assets)
	}

	/*
		for _, character := range g.PCharacters {
			if character != nil {
				character.Draw(screen, g.Camera)
			}
		}

		g.Drag.Draw(screen, &g.Camera)
	*/
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
