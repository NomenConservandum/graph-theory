= Список смежности Ia
== Условие
Вывести те вершины, полустепень захода которых меньше, чем у заданной вершины.

== код (фрагменты кода)
```
func task3Func(graph *GraphInfo, nodeGiven *Node) []*Node {
	var degreesList = make(map[*Node]int)
	for _, edges := range graph.connectionsList {
		for _, edge := range edges {
			degreesList[edge.List[1]]++
		}
	}

	var nodes []*Node

	var degreeGiven = degreesList[nodeGiven]

	for node, degree := range degreesList {
		if degree < degreeGiven {
			nodes = append(nodes, node)
		}
	}
	return nodes
}
```

== краткое описание алгоритма
=== Что делает
- Находит все вершины графа, у которых *полустепень захода* (количество входящих рёбер) *меньше*, чем у заданной вершины
- Возвращает список таких вершин

=== Сложность
- Временная сложность: $O(V + E)$
  - $V$ - количество вершин
  - $E$ - количество рёбер
  - Первый цикл: $O(E)$ - подсчёт входящих рёбер
  - Второй цикл: $O(V)$ - сравнение степеней

== примеры входных и выходных данных
=== Входные данные
```
TYPE: DIRECTED UNWEIGHTED
VERTICES: 1,2,3,4,5,6
EDGES:
1->2
1->6
2->5
2->3
3->4
4->5
6->5
```

=== Выходные данные
```
Vertices:
0: 1
1: 2
2: 3
3: 4
4: 5
5: 6
Enter main vertex index: 4

Vertices which half-degree of entrance is lesser than that of given Vertex:
1. Vertex '2' has less half-degree of entrance than that of '5'
2. Vertex '6' has less half-degree of entrance than that of '5'
3. Vertex '3' has less half-degree of entrance than that of '5'
4. Vertex '4' has less half-degree of entrance than that of '5'
```
