package game

import (
	"sync"
)

// Node represents a tile in the world grid
type Node struct {
	X, Y     int
	Walkable bool
	// Base properties that don't change
}

// PathNode extends Node with A* specific data
type PathNode struct {
	*Node
	G, H, F float64
	Parent  *PathNode
	mu      sync.Mutex // Mutex for concurrent access
}

// World represents the game world grid
type World struct {
	Nodes [][]*Node
	// Other world properties
}

// PathfindingSystem manages concurrent pathfinding requests
type PathfindingSystem struct {
	world     *World
	nodePool  sync.Pool        // Pool of PathNodes to reuse
	workQueue chan PathRequest // Channel for queuing pathfinding requests
}

// PathRequest represents a single pathfinding request
type PathRequest struct {
	Start, End *Node
	ResultChan chan []Node
}

// NewPathfindingSystem creates a new concurrent pathfinding system
func NewPathfindingSystem(world *World, workers int) *PathfindingSystem {
	ps := &PathfindingSystem{
		world: world,
		nodePool: sync.Pool{
			New: func() interface{} {
				return &PathNode{}
			},
		},
		workQueue: make(chan PathRequest, 100),
	}

	// Start worker goroutines
	for i := 0; i < workers; i++ {
		go ps.worker()
	}

	return ps
}

// RequestPath initiates a pathfinding request
func (ps *PathfindingSystem) RequestPath(start, end *Node) chan []Node {
	resultChan := make(chan []Node, 1)
	ps.workQueue <- PathRequest{
		Start:      start,
		End:        end,
		ResultChan: resultChan,
	}
	return resultChan
}

// worker processes pathfinding requests
func (ps *PathfindingSystem) worker() {
	for req := range ps.workQueue {
		path := ps.findPath(req.Start, req.End)
		req.ResultChan <- path
	}
}

// findPath performs A* pathfinding with isolated PathNodes
func (ps *PathfindingSystem) findPath(start, end *Node) []Node {
	// Create a temporary map of PathNodes for this calculation
	nodeMap := make(map[*Node]*PathNode)

	// Get or create PathNode wrapper for the start node
	startNode := ps.getPathNode(start)
	startNode.G = 0
	startNode.H = heuristic(start, end)
	startNode.F = startNode.G + startNode.H

	// Standard A* implementation using the temporary nodes...
	// ... (A* implementation details)

	// Clean up and return to pool
	for _, pn := range nodeMap {
		ps.nodePool.Put(pn)
	}

	return path
}

// getPathNode creates or gets a PathNode from the pool
func (ps *PathfindingSystem) getPathNode(n *Node) *PathNode {
	pn := ps.nodePool.Get().(*PathNode)
	pn.Node = n
	pn.G = 0
	pn.H = 0
	pn.F = 0
	pn.Parent = nil
	return pn
}

// Entity represents a game entity that needs pathfinding
type Entity struct {
	pathfinder *PathfindingSystem
	// Other entity properties
}

// RequestMove initiates an asynchronous movement request
func (e *Entity) RequestMove(targetX, targetY int) {
	start := e.getCurrentNode()
	end := e.pathfinder.world.Nodes[targetY][targetX]

	// Request path asynchronously
	pathChan := e.pathfinder.RequestPath(start, end)

	// Handle the result in a goroutine
	go func() {
		path := <-pathChan
		e.followPath(path)
	}()
}
