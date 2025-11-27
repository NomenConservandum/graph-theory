= Максимальный поток V
== Условие
Решить задачу на нахождение максимального потока любым алгоритмом. Подготовить примеры, демонстрирующие работу алгоритма в разных случаях.

== код (фрагменты кода)
```
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
```

== краткое описание алгоритма
=== Что делает
- Находит *максимальный поток* из истока (source) в сток (sink) в потоковой сети
- Определяет *минимальный разрез* (min-cut) сети
- Использует *кратчайшие увеличивающие пути* по числу рёбер

=== Алгоритм

1. Поиск увеличивающего пути (BFS)
  - Ищем путь от истока к стоку в *остаточной сети*
  - Остаточная пропускная способность = capacity - flow

2. Увеличение потока
  - Находим *минимальную остаточную способность* на пути
  - Увеличиваем поток вдоль всего пути на эту величину

3. Минимальный разрез
  - Выполняем BFS от истока в остаточной сети
  - Разрез = рёбра из достижимых в недостижимые вершины

=== Сложность
- Временная сложность: $O(V × E^2)$
  - $V$ - количество вершин
  - $E$ - количество рёбер
  - Количество итераций: $O(V × E)$
  - BFS за итерацию: $O(E)$
  - Итого: $O(V × E × E) = O(V × E^2)$
/*
- **Пространственная сложность:** O(V + E)
  - Хранение потока: O(E)
  - BFS структуры: O(V)
  - Минимальный разрез: O(E)

**Ключевые особенности:**
1. **Специальный случай Форда-Фалкерсона** с BFS вместо DFS
2. **Гарантированная полиномиальная сложность**
3. **Находит кратчайший по рёбрам увеличивающий путь** на каждой итерации
4. **Определяет минимальный разрез** через достижимые вершины
*/

/*
**Преимущества:**
- Проще в реализации, чем Диниц
- Гарантированная сходимость
- Находит минимальный разрез

**Недостатки:**
- Медленнее, чем Диниц на практике
- Много итераций для некоторых графов

**Теорема:**
- **Максимальный поток = пропускная способность минимального разреза**
- Алгоритм гарантированно находит оптимальное решение

**Применение:**
- Транспортные сети
- Сетевые потоки
- Задачи назначения ресурсов
*/
== примеры входных и выходных данных
=== Входные данные
```
TYPE: DIRECTED WEIGHTED
VERTICES: s,a,b,c,t
EDGES:
s->a: 10
s->b: 5
a->b: 15
a->c: 10
b->c: 10
b->t: 10
c->t: 10
```

=== Выходные данные
```
=== Edmond-Karp algorithm - Max Flow ===

Vertices:
0: s
1: a
2: b
3: c
4: t
Enter source vertex index: 0
Enter sink vertex index: 1

=== Max Flow Search Results ===
Source: 's'
Sink: 'a'
Max Flow Value: 10.00

Edge flow distribution:
  From s to a: 10.00 / 10.00

Minimal cut of an edge:
  from s to a (capacity: 10.00)
Cut's total capacity: 10.00
Probe: Flow value is equal to minimal cut's capacity (10.00 = 10.00)
```