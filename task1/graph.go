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
	List [2]*Node[T] // first - from, second - to
	// Why do we store 'from?' It is easy to remove edge this way:
	// you don't have to walk through the whole map to find the edge
	Weight float64
}

func EdgeConstructor[T any](from *Node[T], to *Node[T], weight float64) *Edge[T] {
	var E = Edge[T]{Weight: weight}
	E.List[0] = from
	E.List[1] = to
	return &E
}

type GraphInfo[T any] struct {
	nodes           []*Node[T] // here are the nodes with their id being number in the array
	connectionsList map[*Node[T]][]*Edge[T]
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

func addVertex[T any](g *GraphInfo[T], n *Node[T]) {
	g.nodes = append(g.nodes, n)
}

func addEdge[T any](g *GraphInfo[T], n1 *Node[T], n2 *Node[T], weight float64) {
	g.connectionsList[n1] = append(g.connectionsList[n1], EdgeConstructor(n1, n2, weight))
}

func eqByAdress[T any](el1 *T, el2 *T) bool {
	return el1 == el2
}

// Removes a_k element by criteria. Returns [a_1, a_2, ..., a_k-1, a_k+1, ..., a_n]
func removeElementArrayByFunc[T any](el *T, l []*T, f func(*T, *T) bool) []*T {
	var lst []*T
	for _, tmp := range l {
		if f(tmp, el) {
			continue
		}
		lst = append(lst, el)
	}
	return lst
}

// a little workaround. The second element should have node as 'from' (first) element of the list
func edgeHasNode[T any](el1 *Edge[T], el2 *Edge[T]) bool {
	return eqByAdress(el1.List[0], el2.List[0]) || eqByAdress(el1.List[1], el2.List[0])
}

func removeEdge[T any](g *GraphInfo[T], e *Edge[T]) {
	g.connectionsList[e.List[0]] = removeElementArrayByFunc(e, g.connectionsList[e.List[0]], eqByAdress)
}

// Removes vertex both from the nodes list and the map: the key and all of the appearances of the vertex in values.
func removeVertex[T any](g *GraphInfo[T], n *Node[T]) {
	g.nodes = removeElementArrayByFunc(n, g.nodes, eqByAdress)

	var conLstTmp map[*Node[T]][]*Edge[T]

	var tmpEdge = EdgeConstructor(n, nil, 0)
	for tmp, lst := range g.connectionsList {
		if eqByAdress(tmp, n) {
			continue
		}
		conLstTmp[tmp] = removeElementArrayByFunc(tmpEdge, lst, edgeHasNode)
	}

	g.connectionsList = conLstTmp
}
