= Обходы графа II
== Условие
Найти цикломатическое число графа — минимальное число рёбер, которые надо удалить, чтобы граф стал ациклическим.

== код (фрагменты кода)
```
// task5Func calculates the cyclomatic number (cycle rank) of a graph
// Formula: mu = e - v + p
// where: e = number of edges, v = number of vertices, p = number of connected components
func task5Func(g *GraphInfo) int {
	if len(g.nodes) == 0 {
		return 0
	}

	// Count edges
	e := countEdges(g)

	// Count vertices
	v := len(g.nodes)

	// Count connected components
	p := countConnectedComponents(g)

	// Calculate cyclomatic number
	cyclomaticNumber := e - v + p

	// Ensure non-negative result
	if cyclomaticNumber < 0 {
		return 0
	}

	return cyclomaticNumber
}

// countConnectedComponents returns the number of connected components in the graph
func countConnectedComponents(g *GraphInfo) int {
	if len(g.nodes) == 0 {
		return 0
	}

	visited := make(map[*Node]bool)
	componentCount := 0

	for _, node := range g.nodes {
		if !visited[node] {
			componentCount++
			if g.isOriented {
				// For directed graphs, use BFS/DFS that follows edges in both directions
				// to find weakly connected components
				bfsWeaklyConnected(g, node, visited)
			} else {
				// For undirected graphs, regular BFS/DFS works
				bfs(g, node, visited)
			}
		}
	}

	return componentCount
}
```

== краткое описание алгоритма
=== Что делает
- Вычисляет *цикломатическое число* графа - минимальное количество рёбер, которые нужно удалить, чтобы граф стал ациклическим
- Использует формулу: $μ = e - v + p$, где
  - e - количество рёбер
  - v - количество вершин  
  - p - количество компонент связности

=== Сложность
- Временная сложность: $O(V + E)$
  - countEdges: $O(E)$
  - countConnectedComponents: $O(V + E)$ (при использовании BFS/DFS)

=== Интерпретация результатов
- $μ = 0$: Граф ациклический (лес)
- $μ = 1$: Один простой цикл
- $μ > 1$: Несколько независимых циклов

== примеры входных и выходных данных
=== Входные данные
```
TYPE: DIRECTED UNWEIGHTED
VERTICES: A,B,C,D,E
EDGES:
A->B
B->C
C->D
D->E
E->A
```

=== Выходные данные
```
=== Cyclomatic Number Calculation ===
Cyclomatic number (cycle rank): 1
This is the minimum number of edges to remove to make the graph acyclic

Calculation details:
Number of edges (e): 5
Number of vertices (v): 5
Number of connected components (p): 1
Formula: mu = e - v + p = 5 - 5 + 1 = 1

You need to remove at least 1 edge(s) to make the graph acyclic
```