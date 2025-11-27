package main

import (
	"container/list"
	"fmt"
	"math"
)

// FlowEdge представляет ребро в потоковой сети
type FlowEdge struct {
	From     *Node
	To       *Node
	Capacity float64
}

// FlowNetwork представляет потоковую сеть
type FlowNetwork struct {
	Nodes []*Node
	Edges []*FlowEdge
	Adj   map[*Node][]*FlowEdge
}

// MaxFlowResult представляет результат поиска максимального потока
type MaxFlowResult struct {
	MaxFlowValue float64
	Flow         map[*FlowEdge]float64
	MinCut       []*FlowEdge
	Source       *Node
	Sink         *Node
}

// createFlowNetwork создаёт потоковую сеть из обычного графа
func createFlowNetwork(g *GraphInfo) *FlowNetwork {
	network := &FlowNetwork{
		Nodes: make([]*Node, len(g.nodes)),
		Edges: make([]*FlowEdge, 0),
		Adj:   make(map[*Node][]*FlowEdge),
	}

	// Копируем вершины
	copy(network.Nodes, g.nodes)

	// Создаём рёбра потоковой сети
	for fromNode, edges := range g.connectionsList {
		for _, edge := range edges {
			capacity := edge.Weight

			// Если граф невзвешенный, используем capacity = 1
			if !g.isWeighted {
				capacity = 1
			}

			// Если capacity = 0 (невзвешенный граф без указания весов), устанавливаем 1
			if capacity == 0 {
				capacity = 1
			}

			// Проверяем на отрицательную вместимость
			if capacity < 0 {
				fmt.Printf("\033[31mWarning\033[0m: edge from %v to %v has negative capacity: %.2f\n",
					fromNode.Value, edge.List[1].Value, capacity)
				continue // Пропускаем рёбра с отрицательной пропускной способностью
			}

			flowEdge := &FlowEdge{
				From:     fromNode,
				To:       edge.List[1],
				Capacity: capacity,
			}

			network.Edges = append(network.Edges, flowEdge)
			network.Adj[fromNode] = append(network.Adj[fromNode], flowEdge)

			// Для неориентированного графа добавляем обратное ребро с той же пропускной способностью
			if !g.isOriented {
				reverseEdge := &FlowEdge{
					From:     edge.List[1],
					To:       fromNode,
					Capacity: capacity,
				}
				network.Edges = append(network.Edges, reverseEdge)
				network.Adj[edge.List[1]] = append(network.Adj[edge.List[1]], reverseEdge)
			}
		}
	}

	return network
}

// edmondsKarp реализует алгоритм Эдмондса-Карпа для поиска максимального потока
func edmondsKarp(network *FlowNetwork, source *Node, sink *Node) *MaxFlowResult {
	result := &MaxFlowResult{
		MaxFlowValue: 0,
		Flow:         make(map[*FlowEdge]float64),
		MinCut:       make([]*FlowEdge, 0),
		Source:       source,
		Sink:         sink,
	}

	// Инициализация потока
	for _, edge := range network.Edges {
		result.Flow[edge] = 0
	}

	// Пока существует увеличивающий путь
	for {
		// BFS для поиска кратчайшего увеличивающего пути
		parent := make(map[*Node]*FlowEdge)
		visited := make(map[*Node]bool) // Сбрасываем visited для каждой итерации
		queue := list.New()

		queue.PushBack(source)
		visited[source] = true

		foundPath := false
		for queue.Len() > 0 {
			current := queue.Remove(queue.Front()).(*Node)

			// Если достигли стока, прерываем BFS
			if current == sink {
				foundPath = true
				break
			}

			for _, edge := range network.Adj[current] {
				residualCapacity := edge.Capacity - result.Flow[edge]
				if residualCapacity > 0 && !visited[edge.To] {
					parent[edge.To] = edge
					visited[edge.To] = true
					queue.PushBack(edge.To)
				}
			}
		}

		// Если путь до стока не найден - завершаем
		if !foundPath {
			break
		}

		// Находим минимальную остаточную пропускную способность на пути
		pathFlow := math.Inf(1)
		for node := sink; node != source; {
			edge := parent[node]
			residualCapacity := edge.Capacity - result.Flow[edge]
			if residualCapacity < pathFlow {
				pathFlow = residualCapacity
			}
			node = edge.From
		}

		// Увеличиваем поток вдоль пути
		for node := sink; node != source; {
			edge := parent[node]
			result.Flow[edge] += pathFlow
			node = edge.From
		}

		result.MaxFlowValue += pathFlow
	}

	// Для нахождения минимального разреза нужно выполнить финальный BFS
	finalVisited := make(map[*Node]bool)
	queue := list.New()
	queue.PushBack(source)
	finalVisited[source] = true

	for queue.Len() > 0 {
		current := queue.Remove(queue.Front()).(*Node)
		for _, edge := range network.Adj[current] {
			residualCapacity := edge.Capacity - result.Flow[edge]
			if residualCapacity > 0 && !finalVisited[edge.To] {
				finalVisited[edge.To] = true
				queue.PushBack(edge.To)
			}
		}
	}

	// Находим минимальный разрез
	result.findMinCut(network, finalVisited)

	return result
}

// findMinCut находит минимальный разрез в потоковой сети
func (r *MaxFlowResult) findMinCut(network *FlowNetwork, visited map[*Node]bool) {
	// Ребра из посещённых в непосещённые образуют минимальный разрез
	for _, edge := range network.Edges {
		if visited[edge.From] && !visited[edge.To] {
			r.MinCut = append(r.MinCut, edge)
		}
	}
}
