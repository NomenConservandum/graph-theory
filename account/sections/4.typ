= Список смежности Iб: несколько графов
== Условие
Построить орграф, полученный из исходного удалением изолированных вершин.

== код (фрагменты кода)
```
func task4Func(graph *GraphInfo) []*Node {
	var isIsolatedList = make(map[*Node]bool, len(graph.connectionsList))
	for _, node := range graph.nodes {
		isIsolatedList[node] = false
	}

	for _, edges := range graph.connectionsList {
		for _, edge := range edges {
			isIsolatedList[edge.List[1]] = true
			isIsolatedList[edge.List[0]] = true
		}
	}

	var nodes []*Node

	for node, degree := range isIsolatedList {
		if !degree {
			nodes = append(nodes, node)
		}
	}
	return nodes
}
```

== краткое описание алгоритма
=== Что делает
- Находит все *изолированные вершины* в графе (вершины без инцидентных рёбер)
- Возвращает список таких вершин

=== Сложность
- Временная сложность: $O(V + E)$
  - $V$ - количество вершин
  - $E$ - количество рёбер
  - Инициализация: $O(V)$
  - Обход рёбер: $O(E)$
  - Фильтрация: $O(V)$

== примеры входных и выходных данных
=== Входные данные
```
TYPE: UNDIRECTED UNWEIGHTED
VERTICES: A,B,C,D,E
EDGES:
B-C
C-D
D-E
```

=== Выходные данные
```
Removing vertex 'A'
```