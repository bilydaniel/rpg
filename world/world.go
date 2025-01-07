package world

import (
	"bilydaniel/rpg/entities"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type World struct {
	//TODO level switching for player and npc
	CurrentLevel *Level
	Levels       map[string]*Level
	Npcs         map[string]*entities.Npc
}

func InitWorld() (*World, error) {
	currentLevel := InitLevel()
	world := World{
		CurrentLevel: &currentLevel,
		Levels:       map[string]*Level{},
	}

	err := currentLevel.LoadLevel("level_1")
	if err != nil {
		return nil, err
	}

	world.Levels[world.CurrentLevel.Name] = &currentLevel
	npcs := map[string]*entities.Npc{}
	//TODO put into assets
	image, _, err := ebitenutil.NewImageFromFile("assets/images/greenchar.png")
	if err != nil {
		return nil, err
	}

	for i := 0; i < 100; i++ {
		id := strconv.Itoa(i)
		npc := &entities.Npc{
			Sprite: &entities.CircleSprite{
				X:   float64(i),
				Y:   float64(i),
				Img: image,
			},
			Character: entities.Character{
				Speed:    1.0 / 30,
				Movement: 1.0 / 30,
			},
		}
		world.CurrentLevel.Occupancy[i][i] = *npc

		npcs[id] = npc

	}
	world.Npcs = npcs

	return &world, nil
}
