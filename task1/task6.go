package main

// findCommonVertexWithEqualPathLength finds a vertex reachable from both u and v with paths of equal length
// Returns the target vertex and the common path length, or nil if no such vertex exists
func findCommonVertexWithEqualPathLength(g *GraphInfo, u *Node, v *Node) (*Node, int) {
	if g == nil || u == nil || v == nil {
		return nil, -1
	}

	// Get all vertices reachable from u with their distances
	distancesFromU := bfsWithDistances(g, u)

	// Get all vertices reachable from v with their distances
	distancesFromV := bfsWithDistances(g, v)

	// Find common vertices with equal distances
	for vertex, distU := range distancesFromU {
		if distV, exists := distancesFromV[vertex]; exists {
			// Check if paths have equal length (not necessarily shortest paths)
			// We need to find if there exists any path of equal length from both u and v to this vertex
			if distU == distV {
				return vertex, distU
			}
		}
	}

	return nil, -1
}

// bfsWithDistances performs BFS and returns a map of vertices to their shortest distances
func bfsWithDistances(g *GraphInfo, start *Node) map[*Node]int {
	distances := make(map[*Node]int)
	visited := make(map[*Node]bool)
	queue := []*Node{start}

	distances[start] = 0
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		currentDistance := distances[current]

		// Visit all neighbors
		for _, edge := range g.connectionsList[current] {
			neighbor := edge.List[1]
			if !visited[neighbor] {
				visited[neighbor] = true
				distances[neighbor] = currentDistance + 1
				queue = append(queue, neighbor)
			}
		}
	}

	return distances
}

// Alternative version that considers all possible path lengths (not just shortest)
func findCommonVertexWithEqualPathLengthAllPaths(g *GraphInfo, u *Node, v *Node) (*Node, int) {
	if g == nil || u == nil || v == nil {
		return nil, -1
	}

	// Get all possible path lengths from u to each vertex
	allPathsFromU := findAllPathLengths(g, u)

	// Get all possible path lengths from v to each vertex
	allPathsFromV := findAllPathLengths(g, v)

	// Find common vertices with at least one equal path length
	for vertex, uLengths := range allPathsFromU {
		if vLengths, exists := allPathsFromV[vertex]; exists {
			// Check if there's any common path length
			for lengthU := range uLengths {
				if vLengths[lengthU] {
					return vertex, lengthU
				}
			}
		}
	}

	return nil, -1
}

// findAllPathLengths uses DFS to find all possible path lengths from start to every reachable vertex
func findAllPathLengths(g *GraphInfo, start *Node) map[*Node]map[int]bool {
	result := make(map[*Node]map[int]bool)

	var dfs func(current *Node, visited map[*Node]bool, pathLength int)
	dfs = func(current *Node, visited map[*Node]bool, pathLength int) {
		// Initialize the map for this vertex if needed
		if result[current] == nil {
			result[current] = make(map[int]bool)
		}

		// Record this path length
		result[current][pathLength] = true

		// Continue to neighbors
		for _, edge := range g.connectionsList[current] {
			neighbor := edge.List[1]

			// Create a copy of visited map for this path
			newVisited := make(map[*Node]bool)
			for k, v := range visited {
				newVisited[k] = v
			}

			if !newVisited[neighbor] {
				newVisited[neighbor] = true
				dfs(neighbor, newVisited, pathLength+1)
			}
		}
	}

	initialVisited := make(map[*Node]bool)
	initialVisited[start] = true
	dfs(start, initialVisited, 0)

	return result
}
