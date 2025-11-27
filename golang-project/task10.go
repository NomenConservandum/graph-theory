package main

import (
	"math"
)

// BellmanFordResult представляет результат алгоритма Беллмана-Форда
type BellmanFordResult struct {
	Distances          map[*Node]float64 // Кратчайшие расстояния от стартовой вершины
	Predecessors       map[*Node]*Node   // Предшественники для восстановления путей
	HasNegativeCycle   bool              // Флаг наличия достижимого цикла отрицательного веса
	NegativeCycleNodes []*Node           // Вершины, достижимые из цикла отрицательного веса
}

// bellmanFord реализует алгоритм Беллмана-Форда для нахождения кратчайших путей из одной вершины
func bellmanFord(g *GraphInfo, start *Node) *BellmanFordResult {
	result := &BellmanFordResult{
		Distances:          make(map[*Node]float64),
		Predecessors:       make(map[*Node]*Node),
		HasNegativeCycle:   false,
		NegativeCycleNodes: make([]*Node, 0),
	}

	if len(g.nodes) == 0 {
		return result
	}

	// Инициализация расстояний
	for _, node := range g.nodes {
		if node == start {
			result.Distances[node] = 0
		} else {
			result.Distances[node] = math.Inf(1)
		}
		result.Predecessors[node] = nil
	}

	// Собираем все рёбра графа
	edges := getAllEdges(g)

	// Фаза релаксации: |V| - 1 итераций
	for i := 0; i < len(g.nodes)-1; i++ {
		changed := false
		for _, edge := range edges {
			u := edge.List[0]
			v := edge.List[1]
			weight := edge.Weight

			if !g.isWeighted {
				weight = 1
			}

			if result.Distances[u] < math.Inf(1) {
				newDist := result.Distances[u] + weight
				if newDist < result.Distances[v] {
					result.Distances[v] = newDist
					result.Predecessors[v] = u
					changed = true
				}
			}
		}
		// Если на текущей итерации не было изменений, можно завершить раньше
		if !changed {
			break
		}
	}

	// Проверка на циклы отрицательного веса
	result.HasNegativeCycle = false
	for _, edge := range edges {
		u := edge.List[0]
		v := edge.List[1]
		weight := edge.Weight

		if !g.isWeighted {
			weight = 1
		}

		if result.Distances[u] < math.Inf(1) {
			if result.Distances[u]+weight < result.Distances[v] {
				result.HasNegativeCycle = true
				// Помечаем вершины, достижимые из цикла отрицательного веса
				result.markNodesReachableFromNegativeCycle(g, v)
				break
			}
		}
	}

	return result
}

// getAllEdges возвращает все рёбра графа
func getAllEdges(g *GraphInfo) []*Edge {
	edges := make([]*Edge, 0)
	for _, edgeList := range g.connectionsList {
		edges = append(edges, edgeList...)
	}
	return edges
}

// markNodesReachableFromNegativeCycle помечает все вершины, достижимые из цикла отрицательного веса
func (r *BellmanFordResult) markNodesReachableFromNegativeCycle(g *GraphInfo, start *Node) {
	visited := make(map[*Node]bool)
	queue := []*Node{start}
	visited[start] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		r.NegativeCycleNodes = append(r.NegativeCycleNodes, current)
		r.Distances[current] = math.Inf(-1) // Помечаем как -∞

		// Добавляем всех соседей
		for _, edge := range g.connectionsList[current] {
			neighbor := edge.List[1]
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
}

// reconstructPath восстанавливает путь от стартовой вершины до целевой
func (r *BellmanFordResult) reconstructPath(target *Node) []*Node {
	if math.IsInf(r.Distances[target], -1) {
		return nil // Путь через цикл отрицательного веса
	}

	if r.Predecessors[target] == nil && r.Distances[target] > 0 {
		return nil // Путь не существует
	}

	path := make([]*Node, 0)
	current := target

	// Восстанавливаем путь в обратном порядке
	for current != nil {
		path = append([]*Node{current}, path...)
		current = r.Predecessors[current]
	}

	return path
}
