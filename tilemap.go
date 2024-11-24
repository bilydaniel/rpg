package main

import (
	"bilydaniel/rpg/config"
	"encoding/json"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tilemap struct {
	Layers []TilemapLayer `json:"layers"`
}

func InitTilemap() Tilemap {
	return Tilemap{}
}

type TilemapLayer struct {
	Data   []int `json:"data"`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

func (t *Tilemap) LoadTestMap() error {
	jsonmap, err := os.ReadFile("assets/maps/test_map.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonmap, t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Tilemap) Draw(screen *ebiten.Image, camera config.Camera) {
	for _, layer := range t.Layers {
		for idx, id := range layer.Data {
			//TODO what can I cashe??
			x := idx % layer.Width
			y := idx / layer.Width

			x *= config.TileSize
			y *= config.TileSize

			tile := t.GetTile()

		}
	}
}
