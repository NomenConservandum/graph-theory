package main

// knots returns all vertices that have self-loops (edges from a vertex to itself)
func knots(graph *GraphInfo) []*Node {
	var knots []*Node

	// Only makes sense for directed graphs
	if !graph.isOriented {
		return knots // Return empty for undirected graphs
	}

	// Check each vertex for self-loops
	for node, edges := range graph.connectionsList {
		for _, edge := range edges {
			// Check if this edge is a loop (from node to itself)
			if edge.List[0] == edge.List[1] {
				knots = append(knots, node)
				break // Found at least one loop for this vertex, move to next
			}
		}
	}

	return knots
}
