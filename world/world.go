package world

type World struct {
	CurrentLevel string
	Levels       map[string]*Level
}

func InitWorld() (*World, error) {
	world := World{
		CurrentLevel: "level_1",
		Levels:       map[string]*Level{},
	}
	level := Level{}

	err := level.LoadLevel(world.CurrentLevel)
	if err != nil {
		return nil, err
	}
	world.Levels[world.CurrentLevel] = &level
	return &world, nil
}
