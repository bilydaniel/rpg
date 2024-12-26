package config

const (
	ScreenW = 400 //TODO figure out better resolution
	ScreenH = 240

	WindowW = 1280
	WindowH = 960

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
		3: "yellow",
	}
}
