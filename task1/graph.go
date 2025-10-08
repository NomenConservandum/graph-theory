package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseStringToType(valueStr string) (interface{}, error) {

	bl, er1 := strconv.ParseBool(valueStr)
	fl, er2 := strconv.ParseFloat(valueStr, 64)
	in, er3 := strconv.ParseInt(valueStr, 10, 32)
	if er1 == nil {
		return bl, nil
	} else if er2 == nil {
		return fl, nil
	} else if er3 == nil {
		return in, nil
	} else if er1 != nil && er2 != nil && er3 != nil {
		return valueStr, nil
	}
	return nil, fmt.Errorf("sth went wrong")
}

// граф: ориентированный / неориентированный, взвешанный / невзвешанный
// методы внутри и снаружи графа
// методы чтения из файла и в файл
// Хранить рёбра как отдельный класс в списке. Так будет легко удалять рёбра с обеих сторон (много ко многому?)
// Ссылка на общий список! Наверное.
// Можно использовать set.

// Метод: хранить список рёбер (вес (int32?) + направление (тернарное значение)) + список вершин. У каждого ребра список ссылок на вершины.

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

// METHODS

func addVertex(g *GraphInfo, n *Node) {
	g.nodes = append(g.nodes, n)
}

func addEdge(g *GraphInfo, n1 *Node, n2 *Node, weight float64) {
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

// Removes a_k element by criteria. Returns [a_1, a_2, ..., a_k-1, a_k+1, ..., a_n]
func findElementArrayByFunc[T any](el *T, l []*T, f func(*T, *T) bool) int {
	for i, tmp := range l {
		if f(tmp, el) {
			return i
		}
	}
	return -1
}

// a little workaround. The second element should have node as 'from' (first) element of the list
func edgeHasNode(el1 *Edge, el2 *Edge) bool {
	return eqByAdress(el1.List[0], el2.List[0]) || eqByAdress(el1.List[1], el2.List[0])
}

func removeEdge(g *GraphInfo, e *Edge) {
	g.connectionsList[e.List[0]] = removeElementArrayByFunc(e, g.connectionsList[e.List[0]], eqByAdress)
}

// Removes vertex both from the nodes list and the map: the key and all of the appearances of the vertex in values.
func removeVertex(g *GraphInfo, n *Node) {
	g.nodes = removeElementArrayByFunc(n, g.nodes, eqByAdress)

	var conLstTmp map[*Node][]*Edge = make(map[*Node][]*Edge)

	var tmpEdge = EdgeConstructor(n, nil, 0)
	for tmp, lst := range g.connectionsList {
		if eqByAdress(tmp, n) {
			continue
		}
		conLstTmp[tmp] = removeElementArrayByFunc[Edge](tmpEdge, lst, edgeHasNode)
	}

	g.connectionsList = conLstTmp
}

func addNonWeightedEdge(g *GraphInfo, n1 *Node, n2 *Node) {
	addEdge(g, n1, n2, 0)
}

func addNonOrientedEdge(g *GraphInfo, n1 *Node, n2 *Node, weight float64) {
	addEdge(g, n1, n2, weight)
	addEdge(g, n2, n1, weight)
}

func addNonOrientedNonWeightedEdge(g *GraphInfo, n1 *Node, n2 *Node) {
	addEdge(g, n1, n2, 0)
	addEdge(g, n2, n1, 0)
}

// TODO: change it

func GraphFromFileConstructor(path string) *GraphInfo {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	defer file.Close()

	graph := GraphEmptyConstructor()
	scanner := bufio.NewScanner(file)
	var lineNumber int

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		lineNumber++

		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		// Parse graph type
		if strings.HasPrefix(line, "TYPE:") {
			parseGraphType(graph, line)
			continue
		}

		// Parse vertices
		if strings.HasPrefix(line, "VERTICES:") {
			parseVertices(graph, line)
			continue
		}

		// Parse edges section start
		if line == "EDGES:" {
			continue // Next lines will be edges
		}

		// Parse edges in various formats
		if strings.Contains(line, "-") || strings.Contains(line, "->") {
			parseEdge(graph, line)
			continue
		}

		// Simple edge list format (fallback)4
		fields := strings.Fields(line)
		//val1, er1 := parseStringToType()
		if len(fields) == 2 {
			parseSimpleEdge(graph, fields[0], fields[1], 0)
		} else if len(fields) == 3 {
			if weight, err := strconv.ParseFloat(fields[2], 64); err == nil {
				parseSimpleEdge(graph, fields[0], fields[1], weight)
			} else {
				parseSimpleEdge(graph, fields[0], fields[1], 0)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return nil
	}

	fmt.Printf("Graph loaded from %s: %d vertices, ", path, len(graph.nodes))

	// Count edges
	edgeCount := 0
	if graph.connectionsList != nil {
		for _, edges := range graph.connectionsList {
			edgeCount += len(edges)
		}
	}
	fmt.Printf("%d edges\n", edgeCount)

	return graph
}

func parseGraphType(graph *GraphInfo, line string) {
	parts := strings.Split(line, ":")
	if len(parts) < 2 {
		return
	}

	typeStr := strings.ToUpper(strings.TrimSpace(parts[1]))

	if strings.Contains(typeStr, "DIRECTED") {
		graph.isOriented = true
	} else if strings.Contains(typeStr, "UNDIRECTED") {
		graph.isOriented = false
	}

	if strings.Contains(typeStr, "WEIGHTED") {
		graph.isWeighted = true
	} else if strings.Contains(typeStr, "UNWEIGHTED") {
		graph.isWeighted = false
	}
}

func parseVertices(graph *GraphInfo, line string) {
	parts := strings.Split(line, ":")
	if len(parts) < 2 {
		return
	}

	verticesStr := strings.TrimSpace(parts[1])
	vertices := strings.Split(verticesStr, ",")

	for _, v := range vertices {
		v = strings.TrimSpace(v)
		if v != "" {
			addVertex(graph, NodeConstructor(v))
		}
	}
}

func parseEdge(graph *GraphInfo, line string) {
	line = strings.TrimSpace(line)

	// Handle weighted directed: "1->2: 5.0"
	if strings.Contains(line, "->") && strings.Contains(line, ":") {
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			edgePart := strings.TrimSpace(parts[0])
			weightStr := strings.TrimSpace(parts[1])

			nodes := strings.Split(edgePart, "->")
			if len(nodes) == 2 {
				from := strings.TrimSpace(nodes[0])
				to := strings.TrimSpace(nodes[1])
				weight, _ := strconv.ParseFloat(weightStr, 64)

				addEdgeBetweenNodes(graph, from, to, weight)
				return
			}
		}
	}

	// Handle unweighted directed: "1->2"
	if strings.Contains(line, "->") {
		nodes := strings.Split(line, "->")
		if len(nodes) == 2 {
			from := strings.TrimSpace(nodes[0])
			to := strings.TrimSpace(nodes[1])
			addEdgeBetweenNodes(graph, from, to, 0)
			return
		}
	}

	// Handle weighted undirected: "A-B: 3.0"
	if strings.Contains(line, "-") && strings.Contains(line, ":") {
		parts := strings.Split(line, ":")
		if len(parts) == 2 {
			edgePart := strings.TrimSpace(parts[0])
			weightStr := strings.TrimSpace(parts[1])

			nodes := strings.Split(edgePart, "-")
			if len(nodes) == 2 {
				from := strings.TrimSpace(nodes[0])
				to := strings.TrimSpace(nodes[1])
				weight, _ := strconv.ParseFloat(weightStr, 64)

				addNonOrientedEdgeBetweenNodes(graph, from, to, weight)
				return
			}
		}
	}

	// Handle simple undirected: "A-B"
	if strings.Contains(line, "-") {
		nodes := strings.Split(line, "-")
		if len(nodes) == 2 {
			from := strings.TrimSpace(nodes[0])
			to := strings.TrimSpace(nodes[1])
			addNonOrientedEdgeBetweenNodes(graph, from, to, 0)
			return
		}
	}
}

func parseSimpleEdge(graph *GraphInfo, fromStr, toStr interface{}, weight float64) {
	if graph.isOriented {
		addEdgeBetweenNodes(graph, fromStr, toStr, weight)
	} else {
		addNonOrientedEdgeBetweenNodes(graph, fromStr, toStr, weight)
	}
}

func addEdgeBetweenNodes(graph *GraphInfo, fromStr, toStr interface{}, weight float64) {
	fromNode := findNodeByValue(graph, fromStr)
	toNode := findNodeByValue(graph, toStr)

	if fromNode == nil {
		fromNode = NodeConstructor(fromStr)
		addVertex(graph, fromNode)
	}
	if toNode == nil {
		toNode = NodeConstructor(toStr)
		addVertex(graph, toNode)
	}

	addEdge(graph, fromNode, toNode, weight)
}

func addNonOrientedEdgeBetweenNodes(graph *GraphInfo, node1Str, node2Str interface{}, weight float64) {
	node1 := findNodeByValue(graph, node1Str)
	node2 := findNodeByValue(graph, node2Str)

	if node1 == nil {
		node1 = NodeConstructor(node1Str)
		addVertex(graph, node1)
	}
	if node2 == nil {
		node2 = NodeConstructor(node2Str)
		addVertex(graph, node2)
	}

	if graph.isWeighted {
		addNonOrientedEdge(graph, node1, node2, weight)
	} else {
		addNonOrientedNonWeightedEdge(graph, node1, node2)
	}
}

func findNodeByValue(graph *GraphInfo, value interface{}) *Node {
	for _, node := range graph.nodes {
		if fmt.Sprintf("%v", node.Value) == fmt.Sprintf("%v", value) {
			return node
		}
	}
	return nil
}
