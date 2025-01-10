package world

import (
	"bilydaniel/rpg/config"
	"bilydaniel/rpg/entities"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

const shaderConst = `
	package main

	const MAX_LIGHTS int = 100

	// Uniforms
	var NumLights int
	var Lights[100]vec4
	var PlayerPost vec2
	var ViewDistance float
	var ScreenSize vec2

func Fragment(position vec4, texCoord vec2, color vec4) vec4{
	//worldPos := texCoord * ScreenSize 
	return vec4(0.0, 0.0, 0.0, 0.0)
		
}

	`

type LightSource struct {
	X, Y, R   float32
	Intensity float32
}
type LightingSystem struct {
	Shader       *ebiten.Shader
	Lights       []LightSource
	ShaderLights []float32
	ViewDistance float32
	ScreenW      float32
	ScreenH      float32
}

func NewLightingSystem(W, H int) (*LightingSystem, error) {
	shader, err := ebiten.NewShader([]byte(shaderConst))
	if err != nil {
		fmt.Println("SHADER ERROR:")
		return nil, err
	}

	return &LightingSystem{
		Shader:       shader,
		Lights:       []LightSource{},
		ShaderLights: []float32{},
		ViewDistance: 64,
		ScreenW:      float32(W),
		ScreenH:      float32(H),
	}, nil
}

func (l *LightingSystem) Draw(screen *ebiten.Image, worldImage *ebiten.Image, pcharacters []*entities.PCharacter) {
	playerX := pcharacters[0].GetX() * config.TileSize
	playerY := pcharacters[0].GetY() * config.TileSize

	op := &ebiten.DrawRectShaderOptions{}
	op.Uniforms = map[string]interface{}{
		"Lights":       l.ShaderLights,
		"NumLights":    len(l.Lights),
		"PlayerPos":    []float32{float32(playerX), float32(playerY)},
		"ViewDistance": l.ViewDistance,
		"ScreenSize":   []float32{l.ScreenW, l.ScreenH},
	}

	op.Images[0] = worldImage

	screen.DrawRectShader(screen.Bounds().Dx(), screen.Bounds().Dy(), l.Shader, op)
}
func (l *LightingSystem) AddLight(x, y int) {
	l.Lights = append(l.Lights, LightSource{
		X:         float32(x * config.TileSize),
		Y:         float32(y * config.TileSize),
		R:         64.0,
		Intensity: 0.75,
	})

	l.ShaderLights = make([]float32, 400)
	for i, light := range l.Lights {
		baseIndex := i * 4
		l.ShaderLights[baseIndex] = light.X
		l.ShaderLights[baseIndex+1] = light.Y
		l.ShaderLights[baseIndex+2] = light.R
		l.ShaderLights[baseIndex+3] = light.Intensity
	}
}
