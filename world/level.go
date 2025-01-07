package world

// !!!!!
// !!!!!
// !!!!!
//TODO MADE BY AI, GO THROUGH EVERYTHING AND FIND ALL THAT IS WRONG
// !!!!!
// !!!!!
// !!!!!

import (
	"bilydaniel/rpg/assets"
	"bilydaniel/rpg/config"
	"bilydaniel/rpg/entities"
	"bilydaniel/rpg/utils"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

type Level struct {
	Name       string
	Grid       [][]*Tile
	Occupancy  [][]entities.Sprite //TODO probably something else than sprite
	Width      int                 //Number of tiles
	Height     int                 //Number of tiles
	Sources    map[string]int      //source => firstgid
	SourceData map[string]*assets.TilesetData
	Obstacles  map[string][]assets.Object
}

func InitLevel() Level {
	l := Level{}

	if l.Sources == nil {
		l.Sources = map[string]int{}
	}
	if l.SourceData == nil {
		l.SourceData = map[string]*assets.TilesetData{}
	}
	if l.Obstacles == nil {
		l.Obstacles = map[string][]assets.Object{}
	}

	return l
}

func (l *Level) Draw(screen *ebiten.Image, cam *config.Camera, assets assets.Assets) {
	opts := ebiten.DrawImageOptions{}
	for y := 0; y < len(l.Grid); y++ {
		for x := 0; x < len(l.Grid[y]); x++ {
			tile := l.Grid[y][x]
			//TODO REMOVE HARDCODE
			image := assets.GetTileImage(l.SourceData["floors"].Name, l.SourceData["floors"].Image, l.SourceData["floors"].Columns, tile.ID)
			if image != nil {
				opts.GeoM.Reset()
				cam.WorldToScreenGeom(&opts, x*config.TileSize, y*config.TileSize)
				screen.DrawImage(image, &opts)
			}
		}
	}

	for _, v := range l.Obstacles["buildings"] {
		gid := v.GID
		firstgid := l.Sources["buildings.tsj"] //TODO REMOVE HARDCODE
		id := gid - firstgid

		resultimage := ""

		//TODO change from O(n) to O(1)
		for _, data := range l.SourceData["buildings"].Tiles {
			if data.ID == id {
				resultimage = data.Image
			}
		}

		image := assets.GetImage(l.SourceData["buildings"].Name, resultimage)
		if image != nil {
			opts.GeoM.Reset()
			cam.WorldToScreenGeom(&opts, v.X, v.Y)
			screen.DrawImage(image, &opts)
		}
	}

}

func (l *Level) NodeFromPoint(point utils.Point) *utils.Node {
	//TODO breaks when I zoom out really far, out of bounds, probably just put a max to the zoomout
	x := int(point.X / config.TileSize)
	y := int(point.Y / config.TileSize)

	y = int(math.Max(float64(y), 0))
	y = int(math.Min(float64(y), float64(l.Height)))

	x = int(math.Max(float64(x), 0))
	x = int(math.Min(float64(x), float64(l.Width)))

	tile := l.Grid[y][x]
	return &tile.Node
}

func (l *Level) ResetValues() {
	for i := 0; i < l.Height; i++ {
		for j := 0; j < l.Width; j++ {
			l.Grid[i][j].F = 0
			l.Grid[i][j].G = 0
			l.Grid[i][j].H = 0
			l.Grid[i][j].Parent = nil
		}
	}
}

func (l *Level) LoadLevel(name string) error {
	l.Name = name

	tilemap, err := assets.LoadTilemap(l.Name)
	if err != nil {
		return err
	}

	for _, source := range tilemap.Tilesets {
		_, ok := l.Sources[source.Source]
		if !ok {
			l.Sources[source.Source] = source.Firstgid
		}
		l.Sources[source.Source] = source.Firstgid
		path := "assets/maps/" + source.Source
		sourceFile, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		sourceData := assets.TilesetData{}
		err = json.Unmarshal(sourceFile, &sourceData)
		if err != nil {
			return err
		}

		if sourceData.Image != "" {
			//Source is tileset
			sourceData.TilesImage = true
			sourceData.Image = filepath.Base(sourceData.Image)
		} else if len(sourceData.Tiles) != 0 {
			//Source is objects
			sourceData.TilesImage = false
			for k, v := range sourceData.Tiles {
				sourceData.Tiles[k].Image = filepath.Base(v.Image)
			}
		}

		l.SourceData[sourceData.Name] = &sourceData
	}
	l.Height = tilemap.Height
	l.Width = tilemap.Width

	l.Occupancy = make([][]entities.Sprite, l.Height)
	for i := 0; i < l.Height; i++ {
		l.Occupancy[i] = make([]entities.Sprite, l.Width)
	}

	l.Grid = make([][]*Tile, l.Height)
	for _, layer := range tilemap.Layers {
		if layer.Type == "tilelayer" {
			if layer.Name == "tiles" {
				if len(layer.Data) < l.Width*l.Height {
					return fmt.Errorf("Tile layer has not enough data")
				}
				//TODO TEST WITH DIFFERENT WIDTH AND HEIGHT, both 100 now
				for i := 0; i < l.Height; i++ {
					//Y
					l.Grid[i] = make([]*Tile, l.Width)
					for j := 0; j < l.Width; j++ {
						//X
						l.Grid[i][j] = &Tile{ID: layer.Data[(i*l.Width)+j], Node: utils.Node{X: j, Y: i}, Walkable: true}
					}
				}
			}
		}

		if layer.Type == "objectgroup" {
			if layer.Name == "buildings" {
				for _, v := range layer.Objects {
					l.Obstacles[layer.Name] = append(l.Obstacles[layer.Name], v)
					//TODO solve rotation
					xgrid := v.X / 16
					ygrid := v.Y / 16

					heightgrid := v.Height / 16
					widthgrid := v.Width / 16
					l.Grid[ygrid][xgrid].Walkable = false
					for i := 0; i < heightgrid; i++ {
						l.Grid[ygrid+i][xgrid].Walkable = false
						l.Grid[ygrid+i][xgrid+widthgrid-1].Walkable = false
					}
					for j := 0; j < widthgrid; j++ {
						l.Grid[ygrid][xgrid+j].Walkable = false
						l.Grid[ygrid+heightgrid-1][xgrid+j].Walkable = false
					}
				}
			}
		}
	}
	return nil
}

func (level *Level) WalkableTile(node *utils.Node) bool {
	return level.Grid[node.Y][node.X].Walkable
}

func (level *Level) OccupiedTile(node *utils.Node) bool {
	return level.Occupancy[node.Y][node.X] != nil
}

func (level *Level) SetTileOccupied(sprite entities.Sprite, x, y int) {
	if x < 0 || y < 0 {
		return
	}
	if x >= level.Width || y >= level.Height {
		return
	}
	level.Occupancy[y][x] = sprite
}
