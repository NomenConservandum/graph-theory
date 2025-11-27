= Веса IV а
== Условие
Определить множество вершин орграфа, расстояние от которых до заданной вершины не более N.

== код (фрагменты кода)
```
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
```

== краткое описание алгоритма
=== Что делает
- Находит все вершины *ориентированного графа*, находящиеся на *расстоянии ≤ maxDistance* от стартовой вершины
- Учитывает *веса рёбер* при вычислении расстояний
- Возвращает *список вершин* (без конкретных расстояний)
=== Алгоритм
1. Инициализация расстояний ($infinity$ для всех, кроме старта $= 0$)
2. Добавление стартовой вершины в приоритетную очередь
3. Пока очередь не пуста:
   - Извлечь вершину с минимальным расстоянием
   - Пропустить, если расстояние устарело
   - Прервать, если расстояние > maxDistance
   - Обработать соседей: обновить расстояния и добавить в очередь
4. Собрать результат: вершины с расстоянием ≤ maxDistance

=== Сложность
- Временная сложность: $O((V + E) log V)$
  - $V$ - количество вершин
  - $E$ - количество рёбер
  - Каждая вершина и ребро обрабатываются: $O(V + E)$
  - Операции с приоритетной очередью: $O(log V)$

/*
**Ключевые особенности:**
1. **Только для ориентированных графов**
2. **Учитывает веса рёбер** (для невзвешенных использует вес = 1)
3. **Оптимизация**: прекращает обработку при превышении maxDistance
4. **Lazy deletion**: пропускает устаревшие элементы очереди

**Алгоритм:**
1. **Инициализация** расстояний (∞ для всех, кроме старта = 0)
2. **Добавление стартовой вершины** в приоритетную очередь
3. **Пока очередь не пуста**:
   - Извлечь вершину с минимальным расстоянием
   - **Пропустить**, если расстояние устарело
   - **Прервать**, если расстояние > maxDistance
   - **Обработать соседей**: обновить расстояния и добавить в очередь
4. **Собрать результат**: вершины с расстоянием ≤ maxDistance

**Отличие от полного Дейкстры:**
- **Ранний выход** при превышении maxDistance
- **Возвращает только вершины**, а не расстояния
- **Не строит полные кратчайшие пути** до всех вершин

**Применение:**
- Поиск "близких" вершин в социальных графах
- Анализ зоны влияния в транспортных сетях
- Поиск соседей в взвешенных графах
*/
== примеры входных и выходных данных
=== Входные данные
```
# TASK 3 TESTING FILE
TYPE: DIRECTED WEIGHTED
VERTICES: 1,2,3,4
EDGES:
1->2: 5.0
2->3: 3.2
3->4: 7.1
```

=== Выходные данные
```
=== Search Within Distance N ===

Vertices:
0: 1
1: 2
2: 3
3: 4
Enter the start vertex index: 0
Enter N: 9

Search Results:
Starting vertex: '1'
N: 9.000000
Number of vertices found: 2
Vertices within distance N:
1. '3'
2. '2'
```