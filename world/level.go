package world

// !!!!!
// !!!!!
// !!!!!
//TODO MADE BY AI, GO THROUGH EVERYTHING AND FIND ALL THAT IS WRONG
// !!!!!
// !!!!!
// !!!!!

import (
	"bilydaniel/rpg/utils"
	"math"
	"sort"
)

type Level struct {
	Grid   [][]*Node
	Width  int //Number of tiles
	Height int //Number of tiles

}

type Node struct {
	utils.Point
	G, H, F  float64
	Parent   *Node
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

func (pf *PathFinder) ReconstructPath(node *Node) []utils.Point {
	currentNode := node
	path := []utils.Point{}

	for currentNode != nil {
		path = append(path, currentNode.Point)
		currentNode = currentNode.Parent
	}
	return path
}
func (level *Level) GetNeighbors(point utils.Point) []*Node {
	neighbors := []*Node{}

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

	openSet := []*Node{}
	closedSet := map[*Node]bool{}

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
