package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CLI struct {
	graphs           []*GraphInfo
	activeGraphIndex int
}

func NewCLI() *CLI {
	return &CLI{
		graphs:           make([]*GraphInfo, 0),
		activeGraphIndex: -1,
	}
}

func (c *CLI) printMainMenu() {
	fmt.Println("\n=== Main Menu ===")
	fmt.Println("1. Select graph to work with")
	fmt.Println("2. Add a graph")
	fmt.Println("3. Exit")
	fmt.Print("Choose an option: ")
}

func (c *CLI) printGraphMenu() {
	fmt.Println("\n=== Graph Operations ===")
	fmt.Println("1. Add vertex")
	fmt.Println("2. Add edge")
	fmt.Println("3. Remove vertex")
	fmt.Println("4. Remove edge")
	fmt.Println("5. List vertices")
	fmt.Println("6. List edges")
	fmt.Println("7. Change graph type")
	fmt.Println("8. Print graph info")
	fmt.Println("9. Load from file")
	fmt.Println("10. Save to file")
	fmt.Println("11. TASK 2 Ia: List Knots")
	fmt.Println("12. TASK 3 Ia")
	fmt.Println("13. TASK 4 Ib: Remove Isolated Vertices")
	fmt.Println("14. TASK 5 II: Calculate Cyclomatic Number")
	fmt.Println("15. TASK 6 II: Find Vertex with Equal Path Lengths")
	fmt.Println("16. Adjacency List")
	fmt.Println("17. Back to main menu")
	fmt.Print("Choose an option: ")
}

func (c *CLI) addGraphMenu() {
	fmt.Println("\n=== Add Graph ===")
	fmt.Println("1. Create graph manually")
	fmt.Println("2. Load graph from file")
	fmt.Println("3. Back to main menu")
	fmt.Print("Choose an option: ")
}

func (c *CLI) selectGraph() {
	if len(c.graphs) == 0 {
		fmt.Println("No graphs available. Please add a graph first.")
		return
	}

	fmt.Println("\n=== Select Graph ===")
	for i, graph := range c.graphs {
		graphType := "Undirected"
		if graph.isOriented {
			graphType = "Directed"
		}
		weightType := "Unweighted"
		if graph.isWeighted {
			weightType = "Weighted"
		}
		fmt.Printf("%d. %s %s Graph (%d vertices, %d edges)\n",
			i+1, graphType, weightType, len(graph.nodes), c.countEdges(graph))
	}
	fmt.Printf("%d. Back to main menu\n", len(c.graphs)+1)
	fmt.Print("Choose a graph: ")

	var input string
	fmt.Scanln(&input)
	choice, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		fmt.Println("Invalid input.")
		return
	}

	if choice == len(c.graphs)+1 {
		return
	}

	if choice < 1 || choice > len(c.graphs) {
		fmt.Println("Invalid graph selection.")
		return
	}

	c.activeGraphIndex = choice - 1
	fmt.Printf("Selected graph %d\n", choice)
	c.graphOperationsMenu()
}

func (c *CLI) countEdges(graph *GraphInfo) int {
	edgeCount := 0
	if graph.connectionsList != nil {
		for _, edges := range graph.connectionsList {
			edgeCount += len(edges)
		}
	}
	return edgeCount
}

func (c *CLI) addGraph() {
	for {
		c.addGraphMenu()

		var input string
		fmt.Scanln(&input)
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			fmt.Println("Invalid input.")
			continue
		}

		switch choice {
		case 1:
			c.createGraphManually()
		case 2:
			c.loadGraphFromFile()
		case 3:
			return
		default:
			fmt.Println("Invalid option.")
		}
	}
}

func (c *CLI) createGraphManually() {
	var input string

	fmt.Print("Is the graph oriented? (y/n): ")
	fmt.Scanln(&input)
	oriented := strings.ToLower(strings.TrimSpace(input)) == "y"

	fmt.Print("Is the graph weighted? (y/n): ")
	fmt.Scanln(&input)
	weighted := strings.ToLower(strings.TrimSpace(input)) == "y"

	newGraph := GraphConstructor(oriented, weighted)
	c.graphs = append(c.graphs, newGraph)
	c.activeGraphIndex = len(c.graphs) - 1

	fmt.Printf("Created new %s %s graph\n",
		map[bool]string{true: "directed", false: "undirected"}[oriented],
		map[bool]string{true: "weighted", false: "unweighted"}[weighted])

	c.graphOperationsMenu()
}

func (c *CLI) loadGraphFromFile() {
	var input string
	fmt.Print("Enter file path: ")
	fmt.Scanln(&input)
	path := strings.TrimSpace(input)

	if path == "" {
		fmt.Println("No file path provided")
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("File '%s' does not exist\n", path)
		return
	}

	newGraph := GraphFromFileConstructor(path)
	if newGraph != nil {
		c.graphs = append(c.graphs, newGraph)
		c.activeGraphIndex = len(c.graphs) - 1
		fmt.Printf("Graph successfully loaded from %s\n", path)
		c.graphOperationsMenu()
	} else {
		fmt.Printf("Failed to load graph from %s\n", path)
	}
}

func (c *CLI) graphOperationsMenu() {
	if c.activeGraphIndex == -1 {
		fmt.Println("No active graph selected.")
		return
	}

	currentGraph := c.graphs[c.activeGraphIndex]

	for {
		c.printGraphMenu()

		var input string
		fmt.Scanln(&input)
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		switch choice {
		case 1:
			c.addVertex(currentGraph)
		case 2:
			c.addEdge(currentGraph)
		case 3:
			c.removeVertex(currentGraph)
		case 4:
			c.removeEdge(currentGraph)
		case 5:
			c.listVertices(currentGraph)
		case 6:
			c.listEdges(currentGraph, false)
		case 7:
			c.changeGraphType(currentGraph)
		case 8:
			c.printGraphInfo(currentGraph)
		case 9:
			c.loadFromFile(currentGraph)
		case 10:
			c.saveToFile(currentGraph)
		case 11:
			c.listKnots(currentGraph)
		case 12:
			c.task3(currentGraph)
		case 13:
			c.task4(currentGraph)
		case 14:
			c.task5(currentGraph)
		case 15:
			c.findCommonVertexWithEqualPaths(currentGraph)
		case 16:
			c.adjacencyList(currentGraph)
		case 17:
			c.activeGraphIndex = -1
			return
		default:
			fmt.Println("Invalid option. Please choose 1-17.")
		}
	}
}

// Updated methods to accept graph as parameter
func (c *CLI) addVertex(graph *GraphInfo) {
	var input string
	fmt.Print("Enter vertex value: ")
	fmt.Scanln(&input)
	value := strings.TrimSpace(input)

	node := NodeConstructor(value)
	addVertex(graph, node)
	fmt.Printf("Vertex '%s' added successfully\n", value)
}

func (c *CLI) addEdge(graph *GraphInfo) {
	if len(graph.nodes) < 2 {
		fmt.Println("Need at least 2 vertices to add an edge")
		return
	}

	c.listVertices(graph)

	var input string
	fmt.Print("Enter first vertex index: ")
	fmt.Scanln(&input)
	idx1, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || idx1 < 0 || idx1 >= len(graph.nodes) {
		fmt.Println("Invalid vertex index")
		return
	}

	fmt.Print("Enter second vertex index: ")
	fmt.Scanln(&input)
	idx2, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || idx2 < 0 || idx2 >= len(graph.nodes) {
		fmt.Println("Invalid vertex index")
		return
	}

	var weight float64 = 0
	if graph.isWeighted {
		fmt.Print("Enter edge weight: ")
		fmt.Scanln(&input)
		weight, err = strconv.ParseFloat(strings.TrimSpace(input), 64)
		if err != nil {
			fmt.Println("Invalid weight, using 0")
			weight = 0
		}
	}

	node1 := graph.nodes[idx1]
	node2 := graph.nodes[idx2]

	if graph.isOriented {
		err := addEdge(graph, node1, node2, weight)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("Added oriented edge from '%v' to '%v'", node1.Value, node2.Value)
		if graph.isWeighted {
			fmt.Printf(" with weight %.2f", weight)
		}
		fmt.Println()
	} else {
		if graph.isWeighted {
			err := addNonOrientedEdge(graph, node1, node2, weight)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		} else {
			err := addNonOrientedNonWeightedEdge(graph, node1, node2)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
		fmt.Printf("Added non-oriented edge between '%v' and '%v'", node1.Value, node2.Value)
		if graph.isWeighted {
			fmt.Printf(" with weight %.2f", weight)
		}
		fmt.Println()
	}
}

func (c *CLI) removeVertex(graph *GraphInfo) {
	if len(graph.nodes) == 0 {
		fmt.Println("No vertices to remove")
		return
	}

	c.listVertices(graph)

	var input string
	fmt.Print("Enter vertex index to remove: ")
	fmt.Scanln(&input)
	idx, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || idx < 0 || idx >= len(graph.nodes) {
		fmt.Println("Invalid vertex index")
		return
	}

	node := graph.nodes[idx]
	removeVertex(graph, node)
	fmt.Printf("Vertex '%v' removed successfully\n", node.Value)
}

var edgeLst []*Edge

func (c *CLI) removeEdge(graph *GraphInfo) {
	if graph.connectionsList == nil || len(graph.connectionsList) == 0 {
		fmt.Println("No edges to remove")
		return
	}

	err := c.listEdges(graph, true)
	if err != nil {
		return
	}

	var input string
	fmt.Print("Enter edge index to remove: ")
	fmt.Scanln(&input)
	idx, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || idx < 0 || idx >= len(edgeLst) {
		fmt.Println("Invalid edge index")
		return
	}

	edge := edgeLst[idx]
	removeEdge(graph, edge)
	fmt.Printf("Edge from '%v' to '%v' has been removed successfully\n", edge.List[0].Value, edge.List[1].Value)
}

func (c *CLI) listVertices(graph *GraphInfo) {
	fmt.Println("\nVertices:")
	if len(graph.nodes) == 0 {
		fmt.Println("No vertices")
		return
	}

	for i, node := range graph.nodes {
		fmt.Printf("%d: %v\n", i, node.Value)
	}
}

func (c *CLI) listEdges(graph *GraphInfo, mode bool) error {
	fmt.Println("\nEdges:")
	if graph.connectionsList == nil || len(graph.connectionsList) == 0 {
		fmt.Println("No edges")
		return fmt.Errorf("no edges")
	}

	edgeCount := 0
	edgeLst = make([]*Edge, 0) // Reset edge list

	for node, edges := range graph.connectionsList {
		for _, edge := range edges {
			if mode {
				fmt.Printf("%d. ", edgeCount)
			}
			fmt.Printf("From '%v' to '%v'", node.Value, edge.List[1].Value)
			if graph.isWeighted {
				fmt.Printf(" (weight: %.2f)", edge.Weight)
			}
			fmt.Println()
			edgeCount++
			edgeLst = append(edgeLst, edge)
		}
	}

	if edgeCount == 0 {
		fmt.Println("No edges")
		return fmt.Errorf("no edges")
	}
	return nil
}

func (c *CLI) changeGraphType(graph *GraphInfo) {
	var input string
	fmt.Print("Is the graph oriented? (y/n): ")
	fmt.Scanln(&input)
	oriented := strings.ToLower(strings.TrimSpace(input)) == "y"

	fmt.Print("Is the graph weighted? (y/n): ")
	fmt.Scanln(&input)
	weighted := strings.ToLower(strings.TrimSpace(input)) == "y"

	// Update the existing graph's properties
	graph.isOriented = oriented
	graph.isWeighted = weighted

	fmt.Printf("Graph type changed: oriented=%v, weighted=%v\n", oriented, weighted)
}

func (c *CLI) printGraphInfo(graph *GraphInfo) {
	fmt.Println("\nGraph Information:")
	fmt.Printf("Type: %s, %s\n",
		map[bool]string{true: "Oriented", false: "Non-oriented"}[graph.isOriented],
		map[bool]string{true: "Weighted", false: "Non-weighted"}[graph.isWeighted])
	fmt.Printf("Number of vertices: %d\n", len(graph.nodes))
	fmt.Printf("Number of edges: %d\n", c.countEdges(graph))
}

func (c *CLI) loadFromFile(graph *GraphInfo) {
	var input string
	fmt.Print("Enter file path: ")
	fmt.Scanln(&input)
	path := strings.TrimSpace(input)

	if path == "" {
		fmt.Println("No file path provided")
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("File '%s' does not exist\n", path)
		return
	}

	newGraph := GraphFromFileConstructor(path)
	if newGraph != nil {
		// Replace the current graph with the loaded one
		c.graphs[c.activeGraphIndex] = newGraph
		fmt.Printf("Graph successfully loaded from %s\n", path)
	} else {
		fmt.Printf("Failed to load graph from %s\n", path)
	}
}

func (c *CLI) saveToFile(graph *GraphInfo) {
	var input string
	fmt.Print("Enter file path: ")
	fmt.Scanln(&input)
	path := strings.TrimSpace(input)

	if path == "" {
		fmt.Println("No file path provided")
		return
	}

	err := WriteToFile(graph, path)
	if err != nil {
		fmt.Printf("Error saving file: %v\n", err)
	} else {
		fmt.Println("File saved successfully!")
	}
}

func (c *CLI) listKnots(graph *GraphInfo) {
	fmt.Println("\nVertices with loops (knots):")

	if !graph.isOriented {
		fmt.Println("This operation only makes sense for directed graphs")
		return
	}

	knots := knots(graph)

	if len(knots) == 0 {
		fmt.Println("No vertices with loops found")
		return
	}

	for i, node := range knots {
		fmt.Printf("%d. Vertex '%v' has self-loop(s)\n", i+1, node.Value)
	}
}

func (c *CLI) task3(graph *GraphInfo) {
	if len(graph.nodes) == 0 {
		fmt.Println("No vertices to examine")
		return
	}

	c.listVertices(graph)

	var input string
	fmt.Print("Enter main vertex index: ")
	fmt.Scanln(&input)
	idx, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || idx < 0 || idx >= len(graph.nodes) {
		fmt.Println("Invalid vertex index")
		return
	}

	nodes := task3Func(graph, graph.nodes[idx])

	if len(nodes) == 0 {
		fmt.Println("No vertices with half-degree of entrance lesser than that of given Vertex")
		return
	} else {
		fmt.Println("\nVertices which half-degree of entrance is lesser than that of given Vertex:")
	}

	for i, node := range nodes {
		fmt.Printf("%d. Vertex '%v' has less half-degree of entrance than that of '%v'\n", i+1, node.Value, graph.nodes[idx].Value)
	}
}

func (c *CLI) task4(graph *GraphInfo) {
	var nodes = task4Func(graph)

	if len(nodes) == 0 {
		fmt.Printf("No vertices to remove\n")
	}

	for _, node := range nodes {
		fmt.Printf("Removing vertex '%v'\n", node.Value)
		removeVertex(graph, node)
	}
}

func (c *CLI) adjacencyList(graph *GraphInfo) {
	fmt.Println("\nAdjacency List:")

	if len(graph.connectionsList) == 0 {
		fmt.Println("No edges")
		return
	}

	for nd, edges := range graph.connectionsList {
		fmt.Printf("%v: ", nd.Value)
		for _, v := range edges {
			fmt.Printf("%v ", v.List[1].Value)
		}
		print("\n")
	}
}

func (c *CLI) task5(graph *GraphInfo) {
	fmt.Println("\n=== Cyclomatic Number Calculation ===")

	if len(graph.nodes) == 0 {
		fmt.Println("Graph is empty")
		return
	}

	cyclomaticNumber := task5Func(graph)

	fmt.Printf("Cyclomatic number (cycle rank): %d\n", cyclomaticNumber)
	fmt.Printf("This is the minimum number of edges to remove to make the graph acyclic\n")

	// Additional information
	e := countEdges(graph)
	v := len(graph.nodes)
	p := countConnectedComponents(graph)

	fmt.Printf("\nCalculation details:\n")
	fmt.Printf("Number of edges (e): %d\n", e)
	fmt.Printf("Number of vertices (v): %d\n", v)
	fmt.Printf("Number of connected components (p): %d\n", p)
	fmt.Printf("Formula: mu = e - v + p = %d - %d + %d = %d\n", e, v, p, cyclomaticNumber)

	if cyclomaticNumber == 0 {
		fmt.Println("\nThe graph is already acyclic!")
	} else {
		fmt.Printf("\nYou need to remove at least %d edge(s) to make the graph acyclic\n", cyclomaticNumber)
	}
}

func (c *CLI) findCommonVertexWithEqualPaths(graph *GraphInfo) {
	fmt.Println("\n=== Find Vertex with Equal Path Lengths ===")

	if len(graph.nodes) < 3 {
		fmt.Println("Need at least 3 vertices for this operation")
		return
	}

	c.listVertices(graph)

	var input string

	// Get vertex u
	fmt.Print("Enter index of vertex u: ")
	fmt.Scanln(&input)
	idxU, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || idxU < 0 || idxU >= len(graph.nodes) {
		fmt.Println("Invalid vertex index")
		return
	}

	// Get vertex v
	fmt.Print("Enter index of vertex v: ")
	fmt.Scanln(&input)
	idxV, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || idxV < 0 || idxV >= len(graph.nodes) {
		fmt.Println("Invalid vertex index")
		return
	}

	if idxU == idxV {
		fmt.Println("Vertices u and v must be different")
		return
	}

	u := graph.nodes[idxU]
	v := graph.nodes[idxV]

	fmt.Printf("\nSearching for vertex reachable from both '%v' and '%v' with equal path length...\n", u.Value, v.Value)

	// Try with shortest paths first (more efficient)
	target, length := findCommonVertexWithEqualPathLength(graph, u, v)

	if target == nil {
		// If no result with shortest paths, try with all possible paths
		fmt.Println("No vertex found with shortest paths of equal length. Trying all possible paths...")
		target, length = findCommonVertexWithEqualPathLengthAllPaths(graph, u, v)
	}

	if target != nil {
		fmt.Printf("✓ Found vertex: '%v'\n", target.Value)
		fmt.Printf("Path length from '%v': %d edges\n", u.Value, length)
		fmt.Printf("Path length from '%v': %d edges\n", v.Value, length)
		fmt.Printf("Total path length: %d edges from each starting vertex\n", length)
	} else {
		fmt.Println("✗ No vertex found that is reachable from both u and v with paths of equal length")

		// Provide some diagnostic information
		distancesFromU := bfsWithDistances(graph, u)
		distancesFromV := bfsWithDistances(graph, v)

		fmt.Println("\nDiagnostic info:")
		fmt.Printf("Vertices reachable from '%v': %d\n", u.Value, len(distancesFromU))
		fmt.Printf("Vertices reachable from '%v': %d\n", v.Value, len(distancesFromV))

		// Find common reachable vertices (regardless of path length)
		commonVertices := 0
		for vertex := range distancesFromU {
			if _, exists := distancesFromV[vertex]; exists {
				commonVertices++
			}
		}
		fmt.Printf("Common reachable vertices: %d\n", commonVertices)
	}
}

func (c *CLI) exitProgram() {
	var input string
	fmt.Print("Do you want to exit? All of your data will be lost, if not saved. (y/n): ")
	fmt.Scanln(&input)

	if strings.ToLower(strings.TrimSpace(input)) == "y" {
		// Count unsaved graphs
		unsavedCount := len(c.graphs)
		if unsavedCount > 0 {
			fmt.Printf("Warning: %d graphs will be lost.\n", unsavedCount)
			fmt.Print("Are you sure? (y/n): ")
			fmt.Scanln(&input)
			if strings.ToLower(strings.TrimSpace(input)) == "y" {
				fmt.Println("Goodbye!")
				os.Exit(0)
			}
		} else {
			fmt.Println("Goodbye!")
			os.Exit(0)
		}
	}
}

func (c *CLI) Run() {
	fmt.Println("Welcome to Graph CLI!")

	for {
		c.printMainMenu()

		var input string
		fmt.Scanln(&input)
		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		switch choice {
		case 1:
			c.selectGraph()
		case 2:
			c.addGraph()
		case 3:
			c.exitProgram()
		default:
			fmt.Println("Invalid option. Please choose 1-3.")
		}
	}
}
