package main

func task4Func(graph *GraphInfo) {
	var isIsolatedList = make(map[*Node]bool, len(graph.connectionsList))
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

	for _, node := range nodes {
		removeVertex(graph, node)
	}
}
