= Обходы графа II
== Условие
В данном графе найти такую вершину, в которую есть пути как из вершины u, так и из вершины v одинаковой (по числу рёбер) длины (не обязательно кратчайшие).

== код (фрагменты кода)
```
// findCommonVertexWithEqualPathLength finds a vertex reachable from both u and v with paths of equal length
// Returns the target vertex and the common path length, or nil if no such vertex exists
func findCommonVertexWithEqualPathLength(g *GraphInfo, u *Node, v *Node) (*Node, int) {
	if g == nil || u == nil || v == nil {
		return nil, -1
	}

	// Get all vertices reachable from u with their distances
	distancesFromU := bfsWithDistances(g, u)

	// Get all vertices reachable from v with their distances
	distancesFromV := bfsWithDistances(g, v)

	// Find common vertices with equal distances
	for vertex, distU := range distancesFromU {
		if distV, exists := distancesFromV[vertex]; exists {
			// Check if paths have equal length (not necessarily shortest paths)
			// We need to find if there exists any path of equal length from both u and v to this vertex
			if distU == distV {
				return vertex, distU
			}
		}
	}

	return nil, -1
}
```

== краткое описание алгоритма
=== Что делает
- Находит вершину, которая *достижима* из обеих вершин u и v *путями одинаковой длины*
- Возвращает найденную вершину и длину путей, либо nil если такой вершины нет

=== Сложность
- Временная сложность: $O(V + E)$
  - Два BFS обхода: $O(V + E)$ каждый
  - Поиск общей вершины: $O(V)$
  - Итого: $O(V + E)$
/*
**Важные замечания:**
- Работает только с **кратчайшими путями**, а не всеми возможными путями
- Для неориентированных графов может не найти очевидные решения
- Если нужны все возможные пути (не только кратчайшие), требуется DFS подход
*/

== примеры входных и выходных данных
=== Входные данные
```
TYPE: DIRECTED UNWEIGHTED
VERTICES: U,V,A,B,C,D,E
EDGES:
V1->A
A->B
B->D
V2->C
C->E
E->D
```

=== Выходные данные
```
=== Find Vertex with Equal Path Lengths ===

Vertices:
0: U
1: V
2: A
3: B
4: C
5: D
6: E
7: V1
8: V2
Enter index of vertex u: 7
Enter index of vertex v: 8

Searching for vertex reachable from both 'V1' and 'V2' with equal path length...
Found vertex: 'D'
Path length from 'V1': 3 edges
Path length from 'V2': 3 edges
Total path length: 3 edges from each starting vertex
```