package main

import (
	"container/heap"
)

// EdgeItem represents an edge for the priority queue in Prim's algorithm
type EdgeItem struct {
	From   *Node
	To     *Node
	Weight float64
	Index  int // Index in the heap
}

// PriorityQueue implements heap.Interface for EdgeItem
type PriorityQueue []*EdgeItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Weight < pq[j].Weight
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*EdgeItem)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// PrimResult represents the result of Prim's algorithm
type PrimResult struct {
	MSTEdges    []*Edge
	TotalWeight float64
	IsConnected bool
}

// prim implements Prim's algorithm for Minimum Spanning Tree using your graph structure
func prim(g *GraphInfo, start *Node) *PrimResult {
	result := &PrimResult{
		MSTEdges:    make([]*Edge, 0),
		TotalWeight: 0.0,
	}

	// Prim's algorithm only works for undirected graphs
	if g.isOriented {
		result.IsConnected = false
		return result
	}

	if len(g.nodes) == 0 {
		result.IsConnected = true
		return result
	}

	// If no start node provided, use first node
	if start == nil {
		start = g.nodes[0]
	}

	// Track visited nodes
	visited := make(map[*Node]bool)

	// Initialize priority queue
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// Start with the initial node
	visited[start] = true

	// Add all edges from start node to the priority queue
	for _, edge := range g.connectionsList[start] {
		heap.Push(&pq, &EdgeItem{
			From:   edge.List[0],
			To:     edge.List[1],
			Weight: edge.Weight,
		})
	}

	// Continue until we have v-1 edges or priority queue is empty
	targetEdges := len(g.nodes) - 1

	for pq.Len() > 0 && len(result.MSTEdges) < targetEdges {
		// Get the minimum weight edge
		minEdgeItem := heap.Pop(&pq).(*EdgeItem)

		// If the destination node is already visited, skip
		if visited[minEdgeItem.To] {
			continue
		}

		// Create the actual edge to add to MST
		mstEdge := EdgeConstructor(minEdgeItem.From, minEdgeItem.To, minEdgeItem.Weight)

		// Add this edge to MST
		result.MSTEdges = append(result.MSTEdges, mstEdge)
		result.TotalWeight += minEdgeItem.Weight

		// Mark the node as visited
		visited[minEdgeItem.To] = true

		// Add all edges from the new node to unvisited nodes
		for _, edge := range g.connectionsList[minEdgeItem.To] {
			neighbor := edge.List[1]
			if !visited[neighbor] {
				heap.Push(&pq, &EdgeItem{
					From:   edge.List[0],
					To:     edge.List[1],
					Weight: edge.Weight,
				})
			}
		}
	}

	// Check if the graph is connected (we should have visited all nodes)
	result.IsConnected = len(visited) == len(g.nodes)

	return result
}

// primAllStarts runs Prim's algorithm from all possible start nodes and returns the best MST
func primAllStarts(g *GraphInfo) *PrimResult {
	if len(g.nodes) == 0 {
		return &PrimResult{IsConnected: true}
	}

	var bestResult *PrimResult

	for _, start := range g.nodes {
		result := prim(g, start)

		// If this is the first valid result or better than current best
		if result.IsConnected && (bestResult == nil || result.TotalWeight < bestResult.TotalWeight) {
			bestResult = result
		}
	}

	if bestResult == nil {
		// No connected MST found, return result from first node
		return prim(g, g.nodes[0])
	}

	return bestResult
}
