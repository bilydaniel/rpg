package config

const (
	ScreenW = 640 //TODO figure out better resolution
	ScreenH = 360

	WindowW = ScreenW * 2
	WindowH = ScreenH * 2

	GameName = "RPG"

	TileSize  = 16
	Tolerance = 8
)

var PlayableCharacters map[int]string

func init() {
	PlayableCharacters = map[int]string{
		0: "red",
		1: "green",
		2: "blue",
	}
}
