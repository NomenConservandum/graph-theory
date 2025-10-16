package main

func isSlipKnot(e *Edge) bool {
	return e.List[0] == e.List[1]
}

func knots(g *GraphInfo) []*Edge {
	var knots []*Edge
	for _, lst := range g.connectionsList {
		for _, tmp := range lst {
			if isSlipKnot(tmp) {
				knots = append(knots, tmp)
			}
		}
	}
	return knots
}
