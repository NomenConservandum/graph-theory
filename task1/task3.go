package main

func task3Func(graph *GraphInfo, nodeGiven *Node) []*Node {
	var degreesList = make(map[*Node]int)
	// Check each vertex for self-loops
	for _, edges := range graph.connectionsList {
		for _, edge := range edges {
			degreesList[edge.List[1]]++
		}
	}

	var nodes []*Node

	var degreeGiven = degreesList[nodeGiven]

	// Check each vertex for self-loops
	for node, degree := range degreesList {
		if degree < degreeGiven {
			nodes = append(nodes, node)
		}
	}
	return nodes
}
