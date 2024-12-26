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
	Width      int            //Number of tiles
	Height     int            //Number of tiles
	Sources    map[string]int //source => firstgid
	SourceData map[string]*assets.TilesetData
	Obstacles  map[string][]assets.Object
}

func (l *Level) Draw(screen *ebiten.Image, cam *config.Camera, assets assets.Assets) {
	opts := ebiten.DrawImageOptions{}
	for y := 0; y < len(l.Grid); y++ {
		for x := 0; x < len(l.Grid[y]); x++ {
			tileid := l.Grid[y][x]
			//TODO REMOVE HARDCODE
			image := assets.GetTileImage(l.SourceData["floors"].Name, l.SourceData["floors"].Image, l.SourceData["floors"].Columns, tileid.ID)
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
	x := int(point.X / config.TileSize)
	y := int(point.Y / config.TileSize)

	y = int(math.Max(float64(y), 0))
	y = int(math.Min(float64(y), float64(l.Height)))

	x = int(math.Max(float64(x), 0))
	x = int(math.Min(float64(x), float64(l.Width)))

	tile := l.Grid[y][x]
	return &tile.Node
}

func (l *Level) LoadLevel() error {
	if l.Sources == nil {
		l.Sources = map[string]int{}
	}
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
				}
			}
		}
	}
	return nil
}

type Tile struct {
	ID int
	utils.Node
	G, H, F  float64
	Parent   *Tile
	Walkable bool //TODO change to something more complex, gonna need to check for building, enemies, etc.
}

// TODO put somewhere else
type PathFinder struct {
	CollisionShapes []utils.CollisionShape
}

func (pf *PathFinder) Distance(start utils.Node, end utils.Node) float64 {
	dx := start.X - end.X
	dy := start.Y - end.Y

	return math.Hypot(float64(dx), float64(dy))
}

func (pf *PathFinder) ReconstructPath(node *Tile) []utils.Node {
	currentNode := node
	path := []utils.Node{}

	for currentNode != nil {
		path = append(path, currentNode.Node)
		currentNode = currentNode.Parent
	}
	return path
}
func (level *Level) GetNeighbors(node utils.Node) []*Tile {
	neighbors := []*Tile{}

	offsets := [][]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for _, offset := range offsets {
		offsetx := node.X + offset[0]
		offsety := node.Y + offset[1]

		if offsetx >= 0 && offsety >= 0 && offsetx < level.Width && offsety < level.Height {
			neighbors = append(neighbors, level.Grid[offsety][offsetx])
		}
	}
	return neighbors
}

func (pf *PathFinder) AlfaStar(level Level, start utils.Node, end utils.Node) []utils.Node {
	startNode := level.Grid[start.Y][start.X]
	endNode := level.Grid[end.Y][end.X]

	fmt.Printf("START: %+v\n", startNode)
	fmt.Printf("END: %+v\n", endNode)

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

		neighbors := level.GetNeighbors(current.Node)
		for _, neighbor := range neighbors {
			//TODO probably gonna need something more complex than walkable??
			if closedSet[neighbor] || !neighbor.Walkable {
				continue
			}

			tentativeG := current.G + pf.Distance(current.Node, neighbor.Node)
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
				neighbor.H = pf.Distance(neighbor.Node, endNode.Node)
				neighbor.F = neighbor.G + neighbor.H

				if !sliceContains {
					openSet = append(openSet, neighbor)
				}
			}
		}
	}
	return nil
}
