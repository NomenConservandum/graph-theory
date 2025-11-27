= Каркас III
== Условие
Дан взвешенный неориентированный граф из N вершин и M ребер. Требуется найти в нем каркас минимального веса.

Алгоритм, который необходимо реализовать для решения задачи (Прима или Краскала), выдает преподаватель.

== код (фрагменты кода)
```
// prim implements Prim's algorithm for Minimum Spanning Tree using your graph structure
func prim(g *GraphInfo, start *Node) *PrimResult {
	result := &PrimResult{
		MSTEdges:    make([]*Edge, 0),
		TotalWeight: 0.0,
	}

	// Prim's algorithm only works for undirected graphs
	if g.isOriented {
		result.IsConnected = false
		return result
	}

	if len(g.nodes) == 0 {
		result.IsConnected = true
		return result
	}

	// If no start node provided, use first node
	if start == nil {
		start = g.nodes[0]
	}

	// Track visited nodes
	visited := make(map[*Node]bool)

	// Initialize priority queue
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// Start with the initial node
	visited[start] = true

	// Add all edges from start node to the priority queue
	for _, edge := range g.connectionsList[start] {
		heap.Push(&pq, &EdgeItem{
			From:   edge.List[0],
			To:     edge.List[1],
			Weight: edge.Weight,
		})
	}

	// Continue until we have v-1 edges or priority queue is empty
	targetEdges := len(g.nodes) - 1

	for pq.Len() > 0 && len(result.MSTEdges) < targetEdges {
		// Get the minimum weight edge
		minEdgeItem := heap.Pop(&pq).(*EdgeItem)

		// If the destination node is already visited, skip
		if visited[minEdgeItem.To] {
			continue
		}

		// Create the actual edge to add to MST
		mstEdge := EdgeConstructor(minEdgeItem.From, minEdgeItem.To, minEdgeItem.Weight)

		// Add this edge to MST
		result.MSTEdges = append(result.MSTEdges, mstEdge)
		result.TotalWeight += minEdgeItem.Weight

		// Mark the node as visited
		visited[minEdgeItem.To] = true

		// Add all edges from the new node to unvisited nodes
		for _, edge := range g.connectionsList[minEdgeItem.To] {
			neighbor := edge.List[1]
			if !visited[neighbor] {
				heap.Push(&pq, &EdgeItem{
					From:   edge.List[0],
					To:     edge.List[1],
					Weight: edge.Weight,
				})
			}
		}
	}

	// Check if the graph is connected (we should have visited all nodes)
	result.IsConnected = len(visited) == len(g.nodes)

	return result
}
```

== краткое описание алгоритма
=== Что делает
- Строит *минимальное остовное дерево* (MST) для неориентированного взвешенного графа
- Использует *жадную стратегию*: на каждом шаге добавляет минимальное по весу ребро, соединяющее дерево с новой вершиной

=== Сложность
- Временная сложность: $O(E log V)$
  - $V$ - количество вершин
  - $E$ - количество рёбер
  - Каждое ребро обрабатывается один раз: $O(E)$
  - Операции с приоритетной очередью: $O(log V)$ на операцию

/*
**Ключевые шаги:**
1. **Инициализация** с начальной вершиной
2. **Добавление всех рёбер** из начальной вершины в приоритетную очередь
3. **Пока не построено MST** (V-1 рёбер):
   - Извлечь **минимальное ребро** из очереди
   - Если вершина уже в MST - пропустить
   - Добавить ребро в MST
   - Добавить все рёбра из новой вершины в очередь

**Особенности:**
- Работает **только с неориентированными графами**
- **Жадный алгоритм** - всегда выбирает локально оптимальное решение
- Строит MST **постепенно**, начиная с произвольной вершины
- **Проверяет связность** графа в конце
*/
== примеры входных и выходных данных
=== Входные данные
```
TYPE: UNDIRECTED WEIGHTED
VERTICES: A,B,C,D,E
EDGES:
A-B: 4
A-C: 2
B-C: 1
B-D: 5
C-D: 8
C-E: 10
D-E: 2
```

=== Выходные данные
```
=== Prim's Algorithm - Minimum Spanning Tree ===

Vertices:
0: A
1: B
2: C
3: D
4: E
Enter starting vertex index (or press Enter for automatic): 0
Starting from vertex: A
Minimum Spanning Tree found!
Total weight: 10.00
Number of edges in MST: 4

Edges in Minimum Spanning Tree:
1. A -- C (weight: 2.00)
2. C -- B (weight: 1.00)
3. B -- D (weight: 5.00)
4. D -- E (weight: 2.00)

Original graph: 7 edges
MST reduction: 3 edges removed
```