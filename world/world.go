package world

type World struct {
	//TODO level switching for player and npc
	CurrentLevel *Level
	Levels       map[string]*Level
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

	return &world, nil
}
