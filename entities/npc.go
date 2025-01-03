package entities

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

func (npc *Npc) Draw() {

}
