package world

type World struct {
	CurrentLevel *Level
	Levels       map[string]*Level
}

func InitWorld() (*World, error) {
	currentLevel := Level{Name: "level_1"}
	world := World{
		CurrentLevel: &currentLevel,
		Levels:       map[string]*Level{},
	}

	err := currentLevel.LoadLevel()
	if err != nil {
		return nil, err
	}

	world.Levels[world.CurrentLevel.Name] = &currentLevel
	return &world, nil
}
