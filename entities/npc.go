package entities

import (
	"bilydaniel/rpg/config"

	"github.com/hajimehoshi/ebiten/v2"
)

type Npc struct {
	Sprite
	Character
	LevelName string
}

func (npc *Npc) Update(level Level) {
	x := npc.GetX()

	if x < 0 {
		npc.Movement = npc.Speed
	} else if x > 100 {
		npc.Movement = -npc.Speed
	}
	if int(x+npc.Movement) > int(x) {
		level.SetTileOccupied(npc, int(x+npc.Movement), int(npc.GetY()))
		level.SetTileOccupied(nil, int(x), int(npc.GetY()))
	}
	//npc.SetX(x + npc.Movement)
}

func (npc *Npc) Draw(screen *ebiten.Image, camera config.Camera) {
	opts := ebiten.DrawImageOptions{}

	camera.WorldToScreenGeom(&opts, int(npc.GetX()*config.TileSize), int(npc.GetY()*config.TileSize))
	screen.DrawImage(npc.Image(), &opts)
}
