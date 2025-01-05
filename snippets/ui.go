package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

// Widget represents a UI element with position and size
type Widget interface {
	Update() error
	Draw(screen *ebiten.Image)
	SetPosition(x, y float64)
	GetBounds() (x, y, width, height float64)
	HandleInput(x, y float64, pressed bool) bool
}

// BaseWidget provides common functionality for all widgets
type BaseWidget struct {
	X, Y          float64
	Width, Height float64
	Visible       bool
	Parent        Widget
	Children      []Widget
}

// Panel is a container for other widgets
type Panel struct {
	BaseWidget
	BackgroundColor color.Color
	BorderColor     color.Color
	BorderWidth     float64
}

func NewPanel(x, y, width, height float64) *Panel {
	return &Panel{
		BaseWidget: BaseWidget{
			X:        x,
			Y:        y,
			Width:    width,
			Height:   height,
			Visible:  true,
			Children: make([]Widget, 0),
		},
		BackgroundColor: color.RGBA{40, 40, 40, 200},
		BorderColor:     color.RGBA{60, 60, 60, 255},
		BorderWidth:     2,
	}
}

func (p *Panel) Draw(screen *ebiten.Image) {
	if !p.Visible {
		return
	}

	// Draw background
	vector.DrawFilledRect(screen, float32(p.X), float32(p.Y), float32(p.Width), float32(p.Height), p.BackgroundColor, true)

	// Draw border
	if p.BorderWidth > 0 {
		vector.StrokeRect(screen, float32(p.X), float32(p.Y), float32(p.Width), float32(p.Height), float32(p.BorderWidth), p.BorderColor, true)
	}

	// Draw children
	for _, child := range p.Children {
		child.Draw(screen)
	}
}

// ProgressBar represents a bar that can show health, mana, etc.
type ProgressBar struct {
	BaseWidget
	Value        float64
	MaxValue     float64
	FillColor    color.Color
	BackColor    color.Color
	BorderColor  color.Color
	BorderWidth  float64
	ShowText     bool
	ValueChanged func(float64)
}

func NewProgressBar(x, y, width, height float64) *ProgressBar {
	return &ProgressBar{
		BaseWidget: BaseWidget{
			X:       x,
			Y:       y,
			Width:   width,
			Height:  height,
			Visible: true,
		},
		Value:       100,
		MaxValue:    100,
		FillColor:   color.RGBA{0, 255, 0, 255},
		BackColor:   color.RGBA{60, 60, 60, 255},
		BorderColor: color.RGBA{200, 200, 200, 255},
		BorderWidth: 1,
		ShowText:    true,
	}
}

func (b *ProgressBar) Draw(screen *ebiten.Image) {
	if !b.Visible {
		return
	}

	// Draw background
	vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height), b.BackColor, true)

	// Draw fill
	fillWidth := (b.Value / b.MaxValue) * b.Width
	vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(fillWidth), float32(b.Height), b.FillColor, true)

	// Draw border
	if b.BorderWidth > 0 {
		vector.StrokeRect(screen, float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height), float32(b.BorderWidth), b.BorderColor, true)
	}
}

// Button represents a clickable button
type Button struct {
	BaseWidget
	Text            string
	TextColor       color.Color
	BackgroundColor color.Color
	HoverColor      color.Color
	PressedColor    color.Color
	BorderColor     color.Color
	BorderWidth     float64
	Font            font.Face
	OnClick         func()
	isHovered       bool
	isPressed       bool
}

func NewButton(x, y, width, height float64, text string, font font.Face) *Button {
	return &Button{
		BaseWidget: BaseWidget{
			X:       x,
			Y:       y,
			Width:   width,
			Height:  height,
			Visible: true,
		},
		Text:            text,
		TextColor:       color.White,
		BackgroundColor: color.RGBA{60, 60, 60, 255},
		HoverColor:      color.RGBA{80, 80, 80, 255},
		PressedColor:    color.RGBA{40, 40, 40, 255},
		BorderColor:     color.RGBA{200, 200, 200, 255},
		BorderWidth:     1,
		Font:            font,
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	if !b.Visible {
		return
	}

	// Choose color based on state
	bgColor := b.BackgroundColor
	if b.isPressed {
		bgColor = b.PressedColor
	} else if b.isHovered {
		bgColor = b.HoverColor
	}

	// Draw background
	vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height), bgColor, true)

	// Draw border
	if b.BorderWidth > 0 {
		vector.StrokeRect(screen, float32(b.X), float32(b.Y), float32(b.Width), float32(b.Height), float32(b.BorderWidth), b.BorderColor, true)
	}

	// Draw text centered
	bounds := text.BoundString(b.Font, b.Text)
	textX := b.X + (b.Width-float64(bounds.Dx()))/2
	textY := b.Y + (b.Height+float64(bounds.Dy()))/2
	text.Draw(screen, b.Text, b.Font, int(textX), int(textY), b.TextColor)
}

// HUD manages all UI elements
type HUD struct {
	Widgets []Widget
}

func NewHUD() *HUD {
	return &HUD{
		Widgets: make([]Widget, 0),
	}
}

func (h *HUD) Update() error {
	for _, widget := range h.Widgets {
		if err := widget.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (h *HUD) Draw(screen *ebiten.Image) {
	for _, widget := range h.Widgets {
		widget.Draw(screen)
	}
}

// Example usage in your game
type Game struct {
	hud        *HUD
	healthBar  *ProgressBar
	manaBar    *ProgressBar
	inventory  *Panel
	playerHP   float64
	playerMana float64
}

func NewGame() *Game {
	g := &Game{
		hud:        NewHUD(),
		playerHP:   100,
		playerMana: 100,
	}

	// Create health bar
	g.healthBar = NewProgressBar(10, 10, 200, 20)
	g.healthBar.FillColor = color.RGBA{255, 0, 0, 255}
	g.hud.Widgets = append(g.hud.Widgets, g.healthBar)

	// Create mana bar
	g.manaBar = NewProgressBar(10, 40, 200, 20)
	g.manaBar.FillColor = color.RGBA{0, 0, 255, 255}
	g.hud.Widgets = append(g.hud.Widgets, g.manaBar)

	// Create inventory panel
	g.inventory = NewPanel(10, 70, 300, 200)
	g.hud.Widgets = append(g.hud.Widgets, g.inventory)

	return g
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw game world first
	// ...

	// Draw HUD on top
	g.hud.Draw(screen)
}

func (g *Game) Update() error {
	// Update game logic
	// ...

	// Update HUD elements
	g.healthBar.Value = g.playerHP
	g.manaBar.Value = g.playerMana

	return g.hud.Update()
}
