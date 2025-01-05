package entities

import (
	"bilydaniel/rpg/config"

	"github.com/hajimehoshi/ebiten/v2"
)

type Npc struct {
	Sprite
	Character
}

func (npc *Npc) Update() {
	x := npc.GetX()
	if x < 0 {
		npc.Movement = npc.Speed
	} else if x > 100 {
		npc.Movement = -npc.Speed
	}
	npc.SetX(x + npc.Movement)
}

func (npc *Npc) Draw(screen *ebiten.Image, camera config.Camera) {
	opts := ebiten.DrawImageOptions{}

	camera.WorldToScreenGeom(&opts, int(npc.GetX()*config.TileSize), int(npc.GetY()*config.TileSize))
	screen.DrawImage(npc.Image(), &opts)
}
