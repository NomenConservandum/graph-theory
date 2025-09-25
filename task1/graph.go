package main

// граф: ориентированный / неориентированный, взвешанный / невзвешанный
// методы внутри и снаружи графа
// методы чтения из файла и в файл
// Хранить рёбра как отдельный класс в списке. Так будет легко удалять рёбра с обеих сторон (много ко многому?)
// Ссылка на общий список! Наверное.
// Можно использовать set.

// Метод: хранить список рёбер (вес (int32?) + направление (тернарное значение)) + список вершин. У каждого ребра список ссылок на вершины.

type Node[T any] struct {
	// Id    uint32
	Value interface{}
	// Connection *edges
}

func NodeEmptyConstructor[T any]() *Node[T] {
	return &Node[T]{}
}

func NodeConstructor[T any](value T) *Node[T] {
	return &Node[T]{Value: value}
}

type Edge[T any] struct {
	List   [2]*Node[T] // first - from, second - to
	Weight float64
}

func EdgeConstructor[T any](from *Node[T], to *Node[T], weight float64) *Edge[T] {
	var E = Edge[T]{Weight: weight}
	E.List[0] = from
	E.List[1] = to
	return &E
}

type GraphInfo[T any] struct {
	nodes           []Node[T] // here are the nodes with their id being number in the array
	connectionsList map[Node[T]][]Edge[T]
	isOriented      bool
	isWeighted      bool
}

// returns an empty graph
func GraphEmptyConstructor[T any]() *GraphInfo[T] {
	return &GraphInfo[T]{}
}

// returns the graph derived from a file
func GraphFromFileConstructor[T any](path string) *GraphInfo[T] {
	return &GraphInfo[T]{} // temp
}

func GraphConstructor[T any](isOriented bool, isWeighted bool) *GraphInfo[T] {
	var G = GraphInfo[T]{}
	G.isOriented = isOriented
	G.isWeighted = isWeighted
	return &G
}

// METHODS
