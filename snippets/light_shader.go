// lighting.go
package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	lightingShader = `
package main

// Uniforms for light properties
var Lights[16]vec4     // x, y, radius, intensity for each light
var NumLights int      // Number of active lights
var PlayerPos vec2     // Player position for FOV calculation
var ViewDistance float // Maximum view distance
var ScreenSize vec2    // Screen dimensions

func Fragment(position vec4, texCoord vec2, color vec4) vec4 {
    // Convert texture coordinates to world position
    worldPos := texCoord * ScreenSize

    // Calculate distance to player for FOV
    distToPlayer := distance(worldPos, PlayerPos)
    if distToPlayer > ViewDistance {
        // Outside view distance - render as dark/unexplored
        return vec4(0.0, 0.0, 0.0, 1.0)
    }

    // Start with ambient light
    totalLight := 0.2

    // Accumulate light from all sources
    for i := 0; i < NumLights; i++ {
        light := Lights[i]
        lightPos := light.xy
        radius := light.z
        intensity := light.w

        dist := distance(worldPos, lightPos)
        if dist < radius {
            // Inverse square law falloff
            contribution := intensity * (1.0 - (dist * dist) / (radius * radius))
            totalLight += contribution
        }
    }

    // Clamp final light value
    totalLight = clamp(totalLight, 0.0, 1.0)

    // Sample the original texture
    original := imageSrc0At(texCoord)

    // Apply lighting to the original color
    return vec4(original.rgb * totalLight, original.a)
}
`
)

// LightSource represents a light in the game world
type LightSource struct {
	X, Y      float32
	Radius    float32
	Intensity float32
}

// LightingSystem manages the shader-based lighting
type LightingSystem struct {
	shader    *ebiten.Shader
	lights    []LightSource
	offscreen *ebiten.Image // Render target for the game world
	lightmap  []float32     // Buffer for light data
	screenW   float32
	screenH   float32
}

func NewLightingSystem(screenWidth, screenHeight int) (*LightingSystem, error) {
	shader, err := ebiten.NewShader([]byte(lightingShader))
	if err != nil {
		return nil, err
	}

	return &LightingSystem{
		shader:    shader,
		lights:    make([]LightSource, 0, 16),
		offscreen: ebiten.NewImage(screenWidth, screenHeight),
		lightmap:  make([]float32, 16*4), // Space for 16 lights x 4 components
		screenW:   float32(screenWidth),
		screenH:   float32(screenHeight),
	}, nil
}

func (ls *LightingSystem) AddLight(x, y, radius, intensity float32) {
	if len(ls.lights) < 16 {
		ls.lights = append(ls.lights, LightSource{x, y, radius, intensity})
	}
}

func (ls *LightingSystem) RemoveLight(index int) {
	if index < len(ls.lights) {
		ls.lights = append(ls.lights[:index], ls.lights[index+1:]...)
	}
}

func (ls *LightingSystem) updateLightBuffer() {
	// Update the light data buffer
	for i, light := range ls.lights {
		baseIndex := i * 4
		ls.lightmap[baseIndex] = light.X
		ls.lightmap[baseIndex+1] = light.Y
		ls.lightmap[baseIndex+2] = light.Radius
		ls.lightmap[baseIndex+3] = light.Intensity
	}
}

func (ls *LightingSystem) Draw(screen *ebiten.Image, worldImage *ebiten.Image, playerX, playerY, viewDistance float32) {
	// Update light positions and properties
	ls.updateLightBuffer()

	// Set up shader uniforms
	op := &ebiten.DrawRectShaderOptions{}
	op.Uniforms = map[string]interface{}{
		"Lights":       ls.lightmap,
		"NumLights":    len(ls.lights),
		"PlayerPos":    []float32{playerX, playerY},
		"ViewDistance": viewDistance,
		"ScreenSize":   []float32{ls.screenW, ls.screenH},
	}

	// The world image is our texture input
	op.Images[0] = worldImage

	// Draw the world with lighting applied
	screen.DrawRectShader(
		screen.Bounds().Dx(),
		screen.Bounds().Dy(),
		ls.shader,
		op,
	)
}
