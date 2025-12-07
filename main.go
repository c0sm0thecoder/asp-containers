package main

import (
	"fmt"
	"sort"
	"sync"
)

// Struct to build the container network
type WaterSystem struct {
	mu           sync.RWMutex
	adj          map[int]map[int]bool
	levels       map[int]float64
	containerIDs map[int]bool
}

func NewWaterSystem() *WaterSystem {
	return &WaterSystem{
		adj:          make(map[int]map[int]bool),
		levels:       make(map[int]float64),
		containerIDs: make(map[int]bool),
	}
}

func (ws *WaterSystem) AddContainer(id int) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if _, exists := ws.containerIDs[id]; !exists {
		ws.containerIDs[id] = true
		ws.adj[id] = make(map[int]bool)
		ws.levels[id] = 0.0
		return nil
	}

	return fmt.Errorf("container with id %d already exists", id)
}

func (ws *WaterSystem) AddWater(id int, amount float64) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if !ws.containerIDs[id] {
		return fmt.Errorf("container with id %d does not exist", id)
	}

	group := ws.getComponent(id)
	levelIncrease := amount / float64(len(group))

	for _, memberID := range group {
		ws.levels[memberID] += levelIncrease
	}

	return nil
}

func (ws *WaterSystem) getComponent(startNode int) []int {
	// Do BFS to get all connected nodes to the startNode
	visited := make(map[int]bool)
	queue := []int{startNode}
	visited[startNode] = true
	component := []int{}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		component = append(component, curr)

		for neighbor := range ws.adj[curr] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}

	return component
}

func (ws *WaterSystem) Connect(a, b int) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if !ws.containerIDs[a] || !ws.containerIDs[b] {
		return fmt.Errorf("one or both containers do not exist")
	}

	if ws.adj[a][b] {
		return fmt.Errorf("containers are already connected")
	}

	ws.adj[a][b] = true
	ws.adj[b][a] = true

	group := ws.getComponent(a)

	totalWater := 0.0
	for _, id := range group {
		totalWater += ws.levels[id]
	}

	newLevel := totalWater / float64(len(group))
	for _, id := range group {
		ws.levels[id] = newLevel
	}

	return nil
}

func (ws *WaterSystem) Disconnect(a, b int) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if !ws.adj[a][b] {
		return fmt.Errorf("containers are not connected")
	}

	delete(ws.adj[a], b)
	delete(ws.adj[b], a)

	return nil
}

func (ws *WaterSystem) PrintStatus() {
	fmt.Println("--------------------- Container Status --------------------")

	ids := make([]int, 0, len(ws.containerIDs))
	for id := range ws.containerIDs {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	for _, id := range ids {
		group := ws.getComponent(id)
		fmt.Printf("Container %d: %.2fL (connected to: ", id, ws.levels[id])

		neighbors := make([]int, 0, len(ws.adj[id]))
		for neighbor := range ws.adj[id] {
			neighbors = append(neighbors, neighbor)
		}
		sort.Ints(neighbors)

		if len(neighbors) == 0 {
			fmt.Print("none")
		} else {
			for i, neighbor := range neighbors {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%d", neighbor)
			}
		}
		fmt.Printf(") [group size: %d]\n", len(group))
	}
	fmt.Println("-----------------------------------------------------------")
}

func main() {
	ws := NewWaterSystem()

	// Setup 4 containers
	fmt.Println("\nSetting up 4 containers")
	for i := 1; i <= 4; i++ {
		if err := ws.AddContainer(i); err != nil {
			fmt.Println("Error:", err)
		}
	}
	ws.PrintStatus()

	// Add 10L water to the container 1
	fmt.Println("\nAdding 10L water to container 1")
	if err := ws.AddWater(1, 10.0); err != nil {
		fmt.Println("Error:", err)
	}
	ws.PrintStatus()

	// Connect containers 1 and 2
	fmt.Println("\nConnecting containers 1 and 2")
	if err := ws.Connect(1, 2); err != nil {
		fmt.Println("Error:", err)
	}
	ws.PrintStatus()

	// Connect containers 3 and 4
	fmt.Println("\nConnecting containers 3 and 4")
	if err := ws.Connect(3, 4); err != nil {
		fmt.Println("Error:", err)
	}
	ws.PrintStatus()

	// Add 20L water to the container 3
	fmt.Println("\nAdding 20L water to container 3")
	if err := ws.AddWater(3, 20.0); err != nil {
		fmt.Println("Error:", err)
	}
	ws.PrintStatus()

	// Connect the two groups we have
	fmt.Println("\nConnecting the two groups (2 and 3)")
	if err := ws.Connect(2, 3); err != nil {
		fmt.Println("Error:", err)
	}
	ws.PrintStatus()

	// Disconnect the groups
	fmt.Println("\nDisconnecting the groups (2 and 3)")
	if err := ws.Disconnect(2, 3); err != nil {
		fmt.Println("Error:", err)
	}
	ws.PrintStatus()

	// Add water to the group 1-2
	fmt.Println("\nAdding 5L water to group 1-2")
	if err := ws.AddWater(1, 5.0); err != nil {
		fmt.Println("Error:", err)
	}
	ws.PrintStatus()
}
