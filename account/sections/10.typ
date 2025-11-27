= Веса IV с
== Условие
Вывести кратчайшие пути из вершины u во все остальные вершины.

== код (фрагменты кода)
```
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
```

== краткое описание алгоритма
=== Что делает
- Находит *кратчайшие пути из одной стартовой вершины* во все остальные
- *Обнаруживает циклы отрицательного веса*, достижимые из стартовой вершины
- Возвращает расстояния, предшественников и информацию об отрицательных циклах

=== Сложность
- Временная сложность: $O(V × E)$
  - $V$ - количество вершин
  - $E$ - количество рёбер
  - Фаза релаксации: $O(V × E)$
  - Проверка циклов: $O(E)$
  
/*
- **Пространственная сложность:** O(V)
  - Хранение расстояний и предшественников: O(V)
  - Список рёбер: O(E)

**Ключевые особенности:**
1. **Работает с отрицательными весами** рёбер
2. **Обнаруживает отрицательные циклы**, достижимые из стартовой вершины
3. **Оптимизация**: ранний выход при отсутствии изменений
4. **Восстанавливает пути** через предшественников

**Основные фазы:**

**1. Фаза релаксации (V-1 итераций):**
```
for i = 1 to |V| - 1:
    for each edge (u, v) with weight w:
        if dist[u] + w < dist[v]:
            dist[v] = dist[u] + w
            pred[v] = u
```

**2. Проверка отрицательных циклов:**
```
for each edge (u, v) with weight w:
    if dist[u] + w < dist[v]:
        // Обнаружен отрицательный цикл
```

**Преимущества:**
- Обрабатывает отрицательные веса
- Обнаруживает отрицательные циклы
- Проще в реализации, чем Дейкстра

**Недостатки:**
- Медленнее Дейкстры для графов без отрицательных весов
- Не работает при недостижимых отрицательных циклах

**Оптимизации:**
- **Ранний выход**: если на итерации не было изменений
- **Очередь**: можно использовать очередь для более эффективной релаксации

**Применение:**
- Графы с отрицательными весами
- Обнаружение арбитражных возможностей
- Анализ финансовых сетей
- Когда нужно обнаружить отрицательные циклы
*/
== примеры входных и выходных данных
=== Входные данные
```
TYPE: DIRECTED WEIGHTED
VERTICES: 1,2,3,4
EDGES:
1->2: 5
1->3: 3
2->4: 2
3->2: 1
2->3: 0
3->4: 6
```

=== Выходные данные
```
=== Bellman-Ford algorithm - Shortest Paths From Vertex U ===

Vertices:
0: 1
1: 2
2: 3
3: 4
Enter starting vertex index: 0

Results for starting vertex '1':

Shortest distances to the vertex:
  From '2': 4.00 | PATH: '1' -> '3' -> '2'
  From '3': 3.00 | PATH: '1' -> '3'
  From '4': 6.00 | PATH: '1' -> '3' -> '2' -> '4'
```