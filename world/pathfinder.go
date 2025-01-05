package world

import (
	"bilydaniel/rpg/utils"
	"math"
	"sort"
)

type Tile struct {
	ID int
	utils.Node
	G, H, F  float64
	Parent   *Tile
	Walkable bool //TODO change to something more complex, gonna need to check for building, enemies, etc.
	Occupied bool
}

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
