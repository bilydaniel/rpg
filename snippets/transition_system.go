package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Portal represents a connection between two levels
type Portal struct {
	Position     Vec2
	TargetLevel  string
	TargetPortal string // ID of the destination portal
	Size         Vec2   // Collision box size
}

// Level represents a game level and its contents
type Level struct {
	ID       string
	Portals  map[string]*Portal // Map of portal ID to Portal
	Entities []Entity           // All entities in this level (players, enemies, etc)
	Width    float64
	Height   float64
}

// Entity interface for anything that can move between levels
type Entity interface {
	GetPosition() Vec2
	SetPosition(Vec2)
	GetSize() Vec2
	OnLevelTransition(oldLevel, newLevel *Level)
}

// Game struct managing levels
type Game struct {
	CurrentLevel    *Level
	Levels          map[string]*Level
	Player          *Player
	TransitionTimer float64 // For transition effects
	IsTransitioning bool
}

func NewGame() *Game {
	game := &Game{
		Levels: make(map[string]*Level),
	}

	// Initialize levels
	game.initializeLevels()

	// Set starting level
	game.CurrentLevel = game.Levels["level1"]
	return game
}

func (g *Game) initializeLevels() {
	// Create level 1
	level1 := &Level{
		ID:      "level1",
		Portals: make(map[string]*Portal),
		Width:   800,
		Height:  600,
	}

	// Add portal to level 2
	level1.Portals["to_level2"] = &Portal{
		Position:     Vec2{X: 750, Y: 300},
		TargetLevel:  "level2",
		TargetPortal: "from_level1",
		Size:         Vec2{X: 32, Y: 32},
	}

	// Create level 2 similarly...
	level2 := &Level{
		ID:      "level2",
		Portals: make(map[string]*Portal),
		Width:   800,
		Height:  600,
	}

	// Add portal back to level 1
	level2.Portals["from_level1"] = &Portal{
		Position:     Vec2{X: 50, Y: 300},
		TargetLevel:  "level1",
		TargetPortal: "to_level2",
		Size:         Vec2{X: 32, Y: 32},
	}

	g.Levels["level1"] = level1
	g.Levels["level2"] = level2
}

// Check if an entity is colliding with any portal
func (g *Game) checkPortalCollisions(entity Entity) *Portal {
	entityPos := entity.GetPosition()
	entitySize := entity.GetSize()

	for _, portal := range g.CurrentLevel.Portals {
		// Simple AABB collision check
		if entityPos.X < portal.Position.X+portal.Size.X &&
			entityPos.X+entitySize.X > portal.Position.X &&
			entityPos.Y < portal.Position.Y+portal.Size.Y &&
			entityPos.Y+entitySize.Y > portal.Position.Y {
			return portal
		}
	}
	return nil
}

// Handle level transition for an entity
func (g *Game) transitionEntity(entity Entity, portal *Portal) {
	if g.IsTransitioning {
		return
	}

	targetLevel := g.Levels[portal.TargetLevel]
	if targetLevel == nil {
		return
	}

	// Find target portal position
	targetPortal := targetLevel.Portals[portal.TargetPortal]
	if targetPortal == nil {
		return
	}

	// Remove entity from current level
	g.removeEntityFromLevel(entity, g.CurrentLevel)

	// Add entity to new level
	targetLevel.Entities = append(targetLevel.Entities, entity)

	// Position entity at target portal
	entity.SetPosition(targetPortal.Position)

	// Notify entity of level change
	entity.OnLevelTransition(g.CurrentLevel, targetLevel)

	// If this is the player, change current level
	if entity == g.Player {
		g.CurrentLevel = targetLevel
		g.IsTransitioning = true
		g.TransitionTimer = 0
	}
}

func (g *Game) removeEntityFromLevel(entity Entity, level *Level) {
	for i, e := range level.Entities {
		if e == entity {
			level.Entities = append(level.Entities[:i], level.Entities[i+1:]...)
			return
		}
	}
}

func (g *Game) Update() error {
	if g.IsTransitioning {
		g.TransitionTimer += 1.0 / 60.0
		if g.TransitionTimer >= 0.5 { // Half second transition
			g.IsTransitioning = false
		}
		return nil
	}

	// Check portal collisions for player
	if portal := g.checkPortalCollisions(g.Player); portal != nil {
		g.transitionEntity(g.Player, portal)
	}

	// Check portal collisions for all entities in current level
	for _, entity := range g.CurrentLevel.Entities {
		if entity != g.Player { // Skip player as we already checked
			if portal := g.checkPortalCollisions(entity); portal != nil {
				g.transitionEntity(entity, portal)
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw current level
	// Draw all entities in current level
	// If transitioning, apply fade effect
	if g.IsTransitioning {
		// Apply transition effect (fade, slide, etc.)
	}
}
