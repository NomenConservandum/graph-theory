package main

func task4Func(graph *GraphInfo) []*Node {
	var isIsolatedList = make(map[*Node]bool, len(graph.connectionsList))
	for _, node := range graph.nodes {
		isIsolatedList[node] = false
	}

	for _, edges := range graph.connectionsList {
		for _, edge := range edges {
			isIsolatedList[edge.List[1]] = true
		}
	}

	var nodes []*Node

	for node, degree := range isIsolatedList {
		if !degree {
			nodes = append(nodes, node)
		}
	}
	return nodes
}
