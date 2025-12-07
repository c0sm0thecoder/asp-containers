package main

import (
	"fmt"
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
	ws.adj[a][b] = true

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
