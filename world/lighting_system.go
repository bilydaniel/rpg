package world

import (
	"bilydaniel/rpg/config"

	"github.com/hajimehoshi/ebiten"
)

const shaderConst = `

	`

type LightSource struct {
	X, Y, R   float32
	Intensity float32
}
type LightingSystem struct {
	Shader       *ebiten.Shader
	Lights       []LightSource
	ShaderLights []float32
	ScreenW      float32
	ScreenH      float32
}

func NewLightingSystem(W, H int) (*LightingSystem, error) {
	shader, err := ebiten.NewShader([]byte(shaderConst))
	if err != nil {
		return nil, err
	}

	return &LightingSystem{
		Shader:       shader,
		Lights:       []LightSource{},
		ShaderLights: []float32{},
		ScreenW:      float32(W),
		ScreenH:      float32(H),
	}, nil
}

func (l *LightingSystem) AddLight(x, y int) {
	l.Lights = append(l.Lights, LightSource{
		X:         float32(x * config.TileSize),
		Y:         float32(y * config.TileSize),
		R:         64.0,
		Intensity: 0.75,
	})

}
