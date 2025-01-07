package entities

import "bilydaniel/rpg/utils"

type Level interface {
	OccupiedTile(node *utils.Node) bool
	WalkableTile(node *utils.Node) bool
	SetTileOccupied(sprite Sprite, x, y int)
}
