# Water Container System

## The Problem

We have multiple water containers. We can connect any two containers together. When containers connect, the water redistributes so that all connected containers have the same water level.

We can also:
- Add water to any container (the water spreads to all connected containers)
- Disconnect containers (they stop sharing water)
- Connect containers that are already part of larger groups (all groups merge and equalize)

## How I Solved It

I built a graph-based system where each container is a node. Connections between containers are edges in an undirected graph.

- **Adding a container**: Just add a new node to the graph
- **Connecting containers**: Add an edge between two nodes, then use BFS to find all connected containers and redistribute water equally
- **Adding water**: Find the connected group using BFS, divide the water among all members
- **Disconnecting**: Remove the edge between two nodes

The adjacency list (`adj`) tracks connections. BFS (`getComponent`) finds all containers in a group.

## Running the Program

```bash
go run main.go
```

## Writing New Test Cases

In `main()`, you can use these methods:

```go
ws := NewWaterSystem()

// Add containers
ws.AddContainer(1)
ws.AddContainer(2)

// Add water (container id, amount in liters)
ws.AddWater(1, 10.0)

// Connect two containers
ws.Connect(1, 2)

// Disconnect two containers
ws.Disconnect(1, 2)

// Print current status
ws.PrintStatus()
```

All methods return errors if something goes wrong (container doesn't exist, already connected, etc.).
