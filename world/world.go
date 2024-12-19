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
	level := Level{}

	err := level.LoadLevel(world.CurrentLevel.Name)
	if err != nil {
		return nil, err
	}
	world.Levels[world.CurrentLevel.Name] = &level
	return &world, nil
}
