= Список смежности Ia
== Условие
Вывести те вершины орграфа, в которых есть петли.

== код (фрагменты кода)
```
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
```

== краткое описание алгоритма
Алгоритм поиска вершин с петлями (self-loops)

=== Что делает
- Находит все вершины ориентированного графа, которые имеют петли (рёбра из вершины в саму себя)
- Возвращает список таких вершин

== Сложность
- Временная сложность: $O(V + E)$
  - $V$ - количество вершин
  - $E$ - количество рёбер
  - Каждое ребро проверяется ровно один раз

== примеры входных и выходных данных

=== Входные данные
```
TYPE: DIRECTED UNWEIGHTED
VERTICES: 1,2,3,4,5
EDGES:
1->1
1->2
2->3
3->4
4->5
5->5
```

=== Выходные данные
```
Vertices with loops (knots):
1. Vertex '5' has self-loop(s)
2. Vertex '1' has self-loop(s)
```
