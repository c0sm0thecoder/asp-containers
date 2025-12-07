package main

import "sync"

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

func (ws *WaterSystem) AddContainer(id int) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if _, exists := ws.containerIDs[id]; !exists {
		ws.containerIDs[id] = true
		ws.adj[id] = make(map[int]bool)
		ws.levels[id] = 0.0
	}
}
