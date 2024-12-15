package pathfinding

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	gridSize     = 40
	cellSize     = 20
)

type Game struct {
	pathFinder   *PathFinder
	movement     *Movement
	start, goal  Point
	showGrid     bool
	pathFound    bool
	smoothedPath []Point
}

func NewGame() *Game {
	pf := NewPathFinder(gridSize, gridSize)

	// Add some random obstacles
	for i := 0; i < 100; i++ {
		x := int(math.Rand(float64(gridSize)))
		y := int(math.Rand(float64(gridSize)))
		pf.Grid[y][x].Walkable = false
	}

	return &Game{
		pathFinder: pf,
		showGrid:   true,
	}
}

func (g *Game) Update() error {
	// Handle mouse input for start and goal
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		gridX, gridY := x/cellSize, y/cellSize

		if gridX < gridSize && gridY < gridSize {
			if g.start == (Point{}) {
				g.start = Point{X: float64(gridX), Y: float64(gridY)}
			} else if g.goal == (Point{}) {
				g.goal = Point{X: float64(gridX), Y: float64(gridY)}

				// Check if start and goal are valid
				if g.pathFinder.Grid[int(g.goal.Y)][int(g.goal.X)].Walkable &&
					g.pathFinder.Grid[int(g.start.Y)][int(g.start.X)].Walkable {
					// Find path
					path := g.pathFinder.ThetaStar(g.start, g.goal)

					if path != nil {
						g.smoothedPath = g.pathFinder.SmoothPath(path)
						g.movement = &Movement{
							CurrentPos:   g.start,
							TargetPos:    g.goal,
							Path:         g.smoothedPath,
							Speed:        0.2,
							PathProgress: 0,
						}
						g.pathFound = true
					}
				}
			}
		}
	}

	// Reset if right mouse button pressed
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.start = Point{}
		g.goal = Point{}
		g.pathFound = false
		g.movement = nil
	}

	// Toggle grid display
	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		g.showGrid = !g.showGrid
	}

	// Update movement if path exists
	if g.movement != nil {
		g.movement.UpdateMovement()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw grid and obstacles
	if g.showGrid {
		for y := 0; y < gridSize; y++ {
			for x := 0; x < gridSize; x++ {
				node := g.pathFinder.Grid[y][x]

				// Draw cell background
				ebitenutil.DrawRect(screen,
					float64(x*cellSize), float64(y*cellSize),
					cellSize-1, cellSize-1,
					color.RGBA{200, 200, 200, 255},
				)

				// Color non-walkable cells
				if !node.Walkable {
					ebitenutil.DrawRect(screen,
						float64(x*cellSize), float64(y*cellSize),
						cellSize-1, cellSize-1,
						color.RGBA{100, 100, 100, 255},
					)
				}
			}
		}
	}

	// Draw start point
	if g.start != (Point{}) {
		ebitenutil.DrawRect(screen,
			g.start.X*cellSize, g.start.Y*cellSize,
			cellSize-1, cellSize-1,
			color.RGBA{0, 255, 0, 255},
		)
	}

	// Draw goal point
	if g.goal != (Point{}) {
		ebitenutil.DrawRect(screen,
			g.goal.X*cellSize, g.goal.Y*cellSize,
			cellSize-1, cellSize-1,
			color.RGBA{255, 0, 0, 255},
		)
	}

	// Draw path
	if g.smoothedPath != nil {
		// Draw original path
		for i := 0; i < len(g.smoothedPath)-1; i++ {
			ebitenutil.DrawLine(screen,
				g.smoothedPath[i].X*cellSize+cellSize/2,
				g.smoothedPath[i].Y*cellSize+cellSize/2,
				g.smoothedPath[i+1].X*cellSize+cellSize/2,
				g.smoothedPath[i+1].Y*cellSize+cellSize/2,
				color.RGBA{0, 0, 255, 100},
			)
		}
	}

	// Draw current position
	if g.movement != nil {
		ebitenutil.DrawRect(screen,
			g.movement.CurrentPos.X*cellSize,
			g.movement.CurrentPos.Y*cellSize,
			cellSize-1, cellSize-1,
			color.RGBA{0, 255, 255, 255},
		)
	}

	// Draw instructions
	ebitenutil.DebugPrint(screen,
		"Left Click: Set Start/Goal\n"+
			"Right Click: Reset\n"+
			"G: Toggle Grid\n",
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Theta* Pathfinding Visualization")

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// Include the entire PathFinder implementation from the previous artifact here
// (Copy the entire PathFinder struct and methods from the previous code)
