package main

import (
	"container/heap"
	"math"
)

// Веса IV а

// DistanceItem представляет элемент для приоритетной очереди в алгоритме Дейкстры
type DistanceItem struct {
	Node     *Node
	Distance float64
	Index    int // Индекс в куче
}

// DistancePriorityQueue реализует heap.Interface для DistanceItem
type DistancePriorityQueue []*DistanceItem

func (pq DistancePriorityQueue) Len() int { return len(pq) }

func (pq DistancePriorityQueue) Less(i, j int) bool {
	return pq[i].Distance < pq[j].Distance
}

func (pq DistancePriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *DistancePriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*DistanceItem)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *DistancePriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

// afindVerticesWithinDistance находит все вершины ориентированного графа,
// расстояние от которых до заданной вершины не более N (с учётом весов рёбер)
func findVerticesWithinDistance(g *GraphInfo, start *Node, maxDistance float64) []*Node {
	if !g.isOriented {
		return nil // Только для ориентированных графов
	}

	if start == nil || maxDistance < 0 {
		return nil
	}

	// Инициализация расстояний
	distances := make(map[*Node]float64)
	for _, node := range g.nodes {
		distances[node] = math.Inf(1)
	}
	distances[start] = 0

	// Инициализация приоритетной очереди
	pq := make(DistancePriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &DistanceItem{
		Node:     start,
		Distance: 0,
	})

	for pq.Len() > 0 {
		// Извлекаем вершину с минимальным расстоянием
		currentItem := heap.Pop(&pq).(*DistanceItem)
		current := currentItem.Node
		currentDist := currentItem.Distance

		// Если текущее расстояние больше сохранённого, пропускаем
		if currentDist > distances[current] {
			continue
		}

		// Если достигли максимального расстояния, не обрабатываем соседей
		if currentDist > maxDistance {
			continue
		}

		// Обрабатываем всех соседей
		for _, edge := range g.connectionsList[current] {
			neighbor := edge.List[1]
			weight := edge.Weight

			// Если граф невзвешенный, используем вес 1
			if !g.isWeighted {
				weight = 1
			}

			newDist := currentDist + weight

			// Если нашли более короткий путь и он в пределах maxDistance
			if newDist < distances[neighbor] && newDist <= maxDistance {
				distances[neighbor] = newDist
				heap.Push(&pq, &DistanceItem{
					Node:     neighbor,
					Distance: newDist,
				})
			}
		}
	}

	// Собираем только список вершин (без расстояний)
	result := make([]*Node, 0)
	for node, dist := range distances {
		if node != start && dist <= maxDistance && !math.IsInf(dist, 1) {
			result = append(result, node)
		}
	}

	return result
}
