package main

// Веса IV а

// findVerticesWithinDistance находит все вершины ориентированного графа,
// расстояние от которых до заданной вершины не более N
func findVerticesWithinDistance(g *GraphInfo, start *Node, maxDistance int) []*Node {
	if !g.isOriented {
		return nil // Только для ориентированных графов
	}

	if start == nil || maxDistance < 0 {
		return nil
	}

	// Используем BFS для нахождения расстояний
	distances := make(map[*Node]int)
	visited := make(map[*Node]bool)
	queue := []*Node{start}

	distances[start] = 0
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		currentDistance := distances[current]

		// Если достигли максимального расстояния, не идём дальше
		if currentDistance >= maxDistance {
			continue
		}

		// Обрабатываем всех соседей (исходящие рёбра)
		for _, edge := range g.connectionsList[current] {
			neighbor := edge.List[1]
			if !visited[neighbor] {
				visited[neighbor] = true
				distances[neighbor] = currentDistance + 1
				queue = append(queue, neighbor)
			}
		}
	}

	// Собираем все вершины с расстоянием <= maxDistance (кроме стартовой)
	result := make([]*Node, 0)
	for vertex, distance := range distances {
		if distance <= maxDistance && vertex != start {
			result = append(result, vertex)
		}
	}

	return result
}
