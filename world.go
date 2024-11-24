package main

type World struct {
	CurrentTilemap *Tilemap
	Tilemaps       []*Tilemap
}

func InitWorld() World {
	tilemap := InitTilemap()
	tilemap.LoadTestMap()

	return World{
		CurrentTilemap: &tilemap,
	}

}
