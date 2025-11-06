package main

// task5Func calculates the cyclomatic number (cycle rank) of a graph
// Formula: mu = e - v + p
// where: e = number of edges, v = number of vertices, p = number of connected components
func task5Func(g *GraphInfo) int {
	if len(g.nodes) == 0 {
		return 0
	}

	// Count edges
	e := countEdges(g)

	// Count vertices
	v := len(g.nodes)

	// Count connected components
	p := countConnectedComponents(g)

	// Calculate cyclomatic number
	cyclomaticNumber := e - v + p

	// Ensure non-negative result
	if cyclomaticNumber < 0 {
		return 0
	}

	return cyclomaticNumber
}

// countEdges returns the total number of edges in the graph
func countEdges(g *GraphInfo) int {
	edgeCount := 0
	if g.connectionsList != nil {
		for _, edges := range g.connectionsList {
			edgeCount += len(edges)
		}
	}

	// For undirected graphs, each edge is stored twice, so divide by 2
	if !g.isOriented {
		edgeCount = edgeCount / 2
	}

	return edgeCount
}

// countConnectedComponents returns the number of connected components in the graph
func countConnectedComponents(g *GraphInfo) int {
	if len(g.nodes) == 0 {
		return 0
	}

	visited := make(map[*Node]bool)
	componentCount := 0

	for _, node := range g.nodes {
		if !visited[node] {
			componentCount++
			if g.isOriented {
				// For directed graphs, use BFS/DFS that follows edges in both directions
				// to find weakly connected components
				bfsWeaklyConnected(g, node, visited)
			} else {
				// For undirected graphs, regular BFS/DFS works
				bfs(g, node, visited)
			}
		}
	}

	return componentCount
}

// bfs performs BFS for undirected graphs
func bfs(g *GraphInfo, start *Node, visited map[*Node]bool) {
	queue := []*Node{start}
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Visit all neighbors
		for _, edge := range g.connectionsList[current] {
			neighbor := edge.List[1]
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
}

// bfsWeaklyConnected performs BFS for directed graphs considering weakly connected components
// (ignoring edge directions)
func bfsWeaklyConnected(g *GraphInfo, start *Node, visited map[*Node]bool) {
	queue := []*Node{start}
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Build adjacency list ignoring directions
		neighbors := getNeighborsIgnoringDirection(g, current)

		for _, neighbor := range neighbors {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
}

// getNeighborsIgnoringDirection returns all neighbors regardless of edge direction
func getNeighborsIgnoringDirection(g *GraphInfo, node *Node) []*Node {
	neighbors := make(map[*Node]bool)

	// Outgoing edges
	for _, edge := range g.connectionsList[node] {
		neighbors[edge.List[1]] = true
	}

	// Incoming edges (need to search all edges)
	for fromNode, edges := range g.connectionsList {
		for _, edge := range edges {
			if edge.List[1] == node {
				neighbors[fromNode] = true
			}
		}
	}

	// Convert map to slice
	result := make([]*Node, 0, len(neighbors))
	for neighbor := range neighbors {
		result = append(result, neighbor)
	}

	return result
}
