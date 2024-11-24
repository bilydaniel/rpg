package config

const (
	ScreenW = 640
	ScreenH = 360

	WindowW = 1024
	WindowH = 720

	GameName = "RPG"

	TileSize = 16
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
