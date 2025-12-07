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
