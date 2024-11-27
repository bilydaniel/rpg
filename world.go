package main

import "bilydaniel/rpg/assets"

type World struct {
	CurrentTilemap *Tilemap
	Tilemaps       []*Tilemap
}

func InitWorld(assets assets.Assets) (*World, error) {
	tilemap := InitTilemap()
	err := tilemap.LoadTestMap("floor", assets)
	if err != nil {
		return nil, err
	}

	return &World{
		CurrentTilemap: &tilemap,
	}, nil
}
