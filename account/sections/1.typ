= Минимальные требования для класса Граф
== Условие
Для решения всех задач курса необходимо создать класс (или иерархию классов - на усмотрение разработчика), содержащий:

1. Структуру для хранения списка смежности графа (не работать с графом через матрицы смежности, если в некоторых алгоритмах удобнее использовать список ребер - реализовать метод, создающий список рёбер на основе списка смежности);

2. Конструкторы (не менее 3-х):
  - конструктор по умолчанию, создающий пустой граф
  - конструктор, заполняющий данные графа из файла
  - конструктор-копию (аккуратно, не все сразу делают именно копию)
  - специфические конструкторы для удобства тестирования

3. Методы:

  - добавляющие вершину,
  - добавляющие ребро (дугу),
  - удаляющие вершину,
  - удаляющие ребро (дугу),
  - выводящие список смежности в файл (в том числе в пригодном для чтения конструктором формате).
  - Не выполняйте некорректные операции, сообщайте об ошибках.

4. Должны поддерживаться как ориентированные, так и неориентированные графы. Заранее предусмотрите возможность добавления меток и\или весов для дуг. Поддержка мультиграфа не требуется.

5. Добавьте минималистичный консольный интерфейс пользователя (не смешивая его с реализацией!), позволяющий добавлять и удалять вершины и рёбра (дуги) и просматривать текущий список смежности графа.

6. Сгенерируйте не менее 4 входных файлов с разными типами графов (балансируйте на комбинации ориентированность-взвешенность) для тестирования класса в этом и последующих заданиях. Графы должны содержать не менее 7-10 вершин, в том числе петли и изолированные вершины.



Замечание:

В зависимости от выбранного способа хранения графа могут появиться дополнительные трудности при удалении-добавлении, например, необходимость переименования вершин, если граф хранится списком $($например, vector C++, List C#$)$. Этого можно избежать, если хранить в списке пару (имя вершины, список смежных вершин), или хранить в другой структуре (например, Dictionary C#$,$ map в С++, при этом список смежности вершины может также храниться в виде словаря с ключами - смежными вершинами и значениями - весами соответствующих ребер). Идеально, если в качестве вершины реализуется обобщенный тип (generic), но достаточно использовать строковый тип или свой класс.

== код (фрагменты кода)

```go
type Node struct {
	// Id    uint32
	Value interface{}
	// Connection *edges
}

func NodeEmptyConstructor() *Node {
	return &Node{}
}

func NodeConstructor(value interface{}) *Node {
	return &Node{Value: value}
}

type Edge struct {
	List [2]*Node // first - from, second - to
	// Why do we store 'from?' It is easy to remove edge this way:
	// you don't have to walk through the whole map to find the edge
	Weight float64
}

func EdgeConstructor(from *Node, to *Node, weight float64) *Edge {
	var E = Edge{Weight: weight}
	E.List[0] = from
	E.List[1] = to
	return &E
}

type GraphInfo struct {
	nodes           []*Node // here are the nodes with their id being number in the array
	connectionsList map[*Node][]*Edge
	isOriented      bool
	isWeighted      bool
}

// returns an empty graph
func GraphEmptyConstructor() *GraphInfo {
	return &GraphInfo{
		connectionsList: make(map[*Node][]*Edge),
		nodes:           make([]*Node, 0),
	}
}

func GraphConstructor(isOriented bool, isWeighted bool) *GraphInfo {
	var G = GraphInfo{
		isOriented:      isOriented,
		isWeighted:      isWeighted,
		connectionsList: make(map[*Node][]*Edge), // Initialize here
		nodes:           make([]*Node, 0),
	}
	return &G
}
```

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

=== Удаление вершины
```
Edges:
From '1' to '1'
From '1' to '2'
From '2' to '3'
From '3' to '4'
From '4' to '5'
From '5' to '5'

...

Vertices:
0: 1
1: 2
2: 3
3: 4
4: 5
Enter vertex index to remove: 0
Vertex '1' removed successfully

...

Edges:
From '2' to '3'
From '3' to '4'
From '4' to '5'
From '5' to '5'
```

=== Добавление ребра
```
Edges:
From '2' to '3'
From '3' to '4'
From '4' to '5'
From '5' to '5'

Vertices:
0: 2
1: 3
2: 4
3: 5
Enter first vertex index: 0
Enter second vertex index: 3
Added oriented edge from '2' to '5'

...

Edges:
From '2' to '3'
From '2' to '5'
From '3' to '4'
From '4' to '5'
From '5' to '5'
```

=== Пример ввода некорректного индекса

```
Edges:
0. From '2' to '3'
1. From '2' to '5'
2. From '3' to '4'
3. From '4' to '5'
4. From '5' to '5'
Enter edge index to remove: 5
Invalid edge index
```