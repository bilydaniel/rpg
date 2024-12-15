package world

import "bilydaniel/rpg/assets"

type World struct {
	CurrentLevel int
	Levels       []*Level
}

func InitWorld(assets assets.Assets) (*World, error) {
	return nil, nil
}
