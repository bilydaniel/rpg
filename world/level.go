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
	"bilydaniel/rpg/utils"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type Level struct {
	Name       string
	Grid       [][]*Tile
	Width      int //Number of tiles
	Height     int //Number of tiles
	SourceData map[string]*assets.TilesetData
	Obstacles  map[string][]assets.Object
}

func (l *Level) Draw(screen *ebiten.Image, cam *config.Camera, assets assets.Assets) {
	for y := 0; y < len(l.Grid); y++ {
		for x := 0; x < len(l.Grid[y]); x++ {
			tileid := l.Grid[y][x]
			//TODO REMOVE HARDCODE
			image := assets.GetTileImage(l.SourceData["floors"].Name, l.SourceData["floors"].Image, l.SourceData["floors"].Columns, tileid.ID)
			if image != nil {
				opts := ebiten.DrawImageOptions{}
				opts.GeoM.Translate(float64(x*config.TileSize), float64(y*config.TileSize))
				opts.GeoM.Translate(-cam.X*cam.Speed, -cam.Y*cam.Speed)
				opts.GeoM.Scale(cam.Scale, cam.Scale)
				screen.DrawImage(image, &opts)
			}
		}
	}

	for _, v := range l.Obstacles["buildings"] {
		image := assets.GetImage(l.SourceData["buildings"].Name, l.SourceData["buildings"].Image, l.SourceData["floors"].Columns, v.ID)
		if image != nil {
			opts := ebiten.DrawImageOptions{}
			/*
				opts.GeoM.Translate(float64(x*config.TileSize), float64(y*config.TileSize))
				opts.GeoM.Translate(-cam.X*cam.Speed, -cam.Y*cam.Speed)
				opts.GeoM.Scale(cam.Scale, cam.Scale)
			*/
			screen.DrawImage(image, &opts)
		}
	}
}

func (l *Level) LoadLevel() error {
	if l.SourceData == nil {
		l.SourceData = map[string]*assets.TilesetData{}
	}
	if l.Obstacles == nil {
		l.Obstacles = map[string][]assets.Object{}
	}
	tilemap, err := assets.LoadTilemap(l.Name)
	if err != nil {
		return err
	}

	//TODO transofrm the tilemap data into my own
	//TODO add things to transform when needed
	//TODO minimal for now
	for _, source := range tilemap.Tilesets {
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
		/*
			for _, y := range l.SourceData {
				//fmt.Printf("%+v\n", y)
			}
		*/
	}

	l.Height = tilemap.Height
	l.Width = tilemap.Width
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
						l.Grid[i][j] = &Tile{ID: layer.Data[(i*l.Width)+j], Point: utils.Point{X: float64(j), Y: float64(i)}}
					}
				}
			}
		}

		if layer.Type == "objectgroup" {
			if layer.Name == "buildings" {
				for _, v := range layer.Objects {
					fmt.Printf("builds: %+v\n", l.SourceData["buildings"].Tiles)
					//TODO FINISH MAPPING source ID TO NAME
					l.Obstacles[layer.Name] = append(l.Obstacles[layer.Name], v)
				}

			}
		}
	}
	return nil
}

type Tile struct {
	ID int
	utils.Point
	G, H, F  float64
	Parent   *Tile
	Walkable bool //TODO change to something more complex, gonna need to check for building, enemies, etc.
}

// TODO put somewhere else
type PathFinder struct {
	CollisionShapes []utils.CollisionShape
}

func (pf *PathFinder) Distance(start utils.Point, end utils.Point) float64 {
	dx := start.X - end.X
	dy := start.Y - end.Y

	return math.Hypot(dx, dy)
}

func (pf *PathFinder) ReconstructPath(node *Tile) []utils.Point {
	currentNode := node
	path := []utils.Point{}

	for currentNode != nil {
		path = append(path, currentNode.Point)
		currentNode = currentNode.Parent
	}
	return path
}
func (level *Level) GetNeighbors(point utils.Point) []*Tile {
	neighbors := []*Tile{}

	offsets := [][]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for _, offset := range offsets {
		offsetx := int(point.X) + offset[0]
		offsety := int(point.Y) + offset[1]

		if offsetx >= 0 && offsety >= 0 && offsetx < level.Width && offsety < level.Height {
			neighbors = append(neighbors, level.Grid[offsety][offsetx])
		}
	}
	return neighbors
}

func (pf *PathFinder) AlfaStar(level Level, start utils.Point, end utils.Point) []utils.Point {
	startNode := level.Grid[int(start.Y)][int(start.X)]
	endNode := level.Grid[int(end.Y)][int(end.X)]

	openSet := []*Tile{}
	closedSet := map[*Tile]bool{}

	startNode.G = 0
	startNode.H = pf.Distance(start, end)
	startNode.F = startNode.G + startNode.H
	startNode.Parent = nil

	openSet = append(openSet, startNode)
	for len(openSet) > 0 {
		sort.Slice(openSet, func(i, j int) bool {
			return openSet[i].F < openSet[j].F
		})

		current := openSet[0]
		if current == endNode {
			return pf.ReconstructPath(current)
		}

		openSet = openSet[1:]
		closedSet[current] = true

		neighbors := level.GetNeighbors(current.Point)
		for _, neighbor := range neighbors {
			//TODO probably gonna need something more complex than walkable??
			if closedSet[neighbor] || !neighbor.Walkable {
				continue
			}

			tentativeG := current.G + pf.Distance(current.Point, neighbor.Point)
			// POSSIBLE UPGRADE FROM A* TO THETA*, DOESENT SEEM NEEDED
			/*
				if current.Parent != nil && pf.LineOfSight(current.Parent.Pos, neighbor.Pos) {
					newG := current.Parent.G + pf.Distance(current.Parent.Pos, neighbor.Pos)
					if newG < neighbor.G {
						neighbor.Parent = current.Parent // Rewire the parent
						neighbor.G = newG
						neighbor.F = neighbor.G + neighbor.H
					}
				}
			*/

			sliceContains := utils.SliceContains(openSet, neighbor)
			if !sliceContains || tentativeG < neighbor.G {
				neighbor.Parent = current
				neighbor.G = tentativeG
				neighbor.H = pf.Distance(neighbor.Point, endNode.Point)
				neighbor.F = neighbor.G + neighbor.H

				if !sliceContains {
					openSet = append(openSet, neighbor)
				}
			}
		}
	}
	return nil
}
