package config

import "sync"

var (
	Clickers   []*Clicker
	clickersMu sync.RWMutex
)

type Clicker interface {
	OnClick()
	ClickCollision(x int, y int, camera Camera) bool
}

func AddClicker(clicker Clicker) {
	clickersMu.Lock()
	Clickers = append(Clickers, &clicker)
	clickersMu.Unlock()
}
