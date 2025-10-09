/*
package main

import "fmt"

	func mode() int {
		var md = 0
		fmt.Println("Select the mode:\n1. Read from file\n2. Enter by hand")
		fmt.Scanln(md)
		return md
	}

	func choice() int {
		var ch = 0
		fmt.Println("1. Enter a Vertex\n2. Enter an Edge\n3. Remove a Vertex\n4. Remove an Edge\n5. Print out the Graph\n6. Export the Graph")
		fmt.Scanln(ch)
		return ch
	}

	func vertexEnter() int {
		var ch = 0
		fmt.Println("1. Enter Vertex's value\n2. Enter an Edge\n3. Remove a Vertex\n4. Remove an Edge\n5. Print out the Graph\n6. Export the Graph")
		fmt.Scanln(ch)
		return ch
	}
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ClientSideGraph struct {
	graph *GraphInfo
}

func CLI() *ClientSideGraph {
	return &ClientSideGraph{
		graph: GraphConstructor(false, false), // default: undirected, unweighted
	}
}

func (c *ClientSideGraph) printMenu() {
	fmt.Println("\n=== Graph ClientSideGraph ===\n1. Add vertex\n2. Add edge\n3. Remove vertex\n4. Remove edge\n5. List vertices\n6. List edges\n7. Change graph type\n8. Print graph info\n9. Load from File\n10.Exit\nChoose an option: ")
}

func (c *ClientSideGraph) addVertex() {
	var vt string
	fmt.Print("Enter vertex value: ")
	fmt.Scanln(&vt)
	value := strings.TrimSpace(vt)

	node := NodeConstructor(value)
	addVertex(c.graph, node)
	fmt.Println("Vertex '", value, "' added successfully")
}

func (c *ClientSideGraph) addEdge() {
	if len(c.graph.nodes) < 2 {
		fmt.Println("Need at least 2 vertices to add an edge")
		return
	}

	c.listVertices()

	var vt string
	fmt.Print("Enter first vertex index: ")
	fmt.Scanln(&vt)

	idx1Str := strings.TrimSpace(vt)
	idx1, err := strconv.Atoi(idx1Str)
	if err != nil || idx1 < 0 || idx1 >= len(c.graph.nodes) {
		fmt.Println("Invalid vertex index")
		return
	}

	fmt.Print("Enter second vertex index: ")
	fmt.Scanln(&vt)
	idx2Str := strings.TrimSpace(vt)
	idx2, err := strconv.Atoi(idx2Str)
	if err != nil || idx2 < 0 || idx2 >= len(c.graph.nodes) {
		fmt.Println("Invalid vertex index")
		return
	}

	var weight float64 = 0
	if c.graph.isWeighted {
		fmt.Print("Enter edge weight: ")
		fmt.Scanln(&vt)
		weightStr := strings.TrimSpace(vt)
		weight, err = strconv.ParseFloat(weightStr, 64)
		if err != nil {
			fmt.Println("Invalid weight, using 0")
			weight = 0
		}
	}

	node1 := c.graph.nodes[idx1]
	node2 := c.graph.nodes[idx2]

	if c.graph.isOriented {
		addEdge(c.graph, node1, node2, weight)
		fmt.Printf("Added oriented edge from '%v' to '%v'", node1.Value, node2.Value)
		if c.graph.isWeighted {
			fmt.Printf(" with weight %.2f", weight)
		}
		fmt.Println()
	} else {
		if c.graph.isWeighted {
			addNonOrientedEdge(c.graph, node1, node2, weight)
		} else {
			addNonOrientedNonWeightedEdge(c.graph, node1, node2)
		}
		fmt.Printf("Added non-oriented edge between '%v' and '%v'", node1.Value, node2.Value)
		if c.graph.isWeighted {
			fmt.Printf(" with weight %.2f", weight)
		}
		fmt.Println()
	}
}

func (c *ClientSideGraph) removeVertex() {
	if len(c.graph.nodes) == 0 {
		fmt.Println("No vertices to remove")
		return
	}

	c.listVertices()

	var vt string
	fmt.Print("Enter vertex index to remove: ")
	fmt.Scanln(&vt)
	idxStr := strings.TrimSpace(vt)
	idx, err := strconv.Atoi(idxStr)
	if err != nil || idx < 0 || idx >= len(c.graph.nodes) {
		fmt.Println("Invalid vertex index")
		return
	}

	node := c.graph.nodes[idx]
	removeVertex(c.graph, node)
	fmt.Printf("Vertex '%v' removed successfully\n", node.Value)
}

var edgeLst []*Edge

func (c *ClientSideGraph) removeEdge() {
	if c.graph.connectionsList == nil || len(c.graph.connectionsList) == 0 {
		fmt.Println("No edges to remove")
		return
	}

	c.listEdges(true)

	var vt string
	fmt.Print("Enter edge index to remove: ")
	fmt.Scanln(&vt)
	idxStr := strings.TrimSpace(vt)
	idx, err := strconv.Atoi(idxStr)
	if err != nil || idx < 0 || idx >= len(c.graph.connectionsList) { // change a bit
		fmt.Println("Invalid vertex index")
		return
	}

	edge := edgeLst[idx]
	removeEdge(c.graph, edge)
	fmt.Println("Edge from '", edge.List[0].Value, "' to '", edge.List[1].Value, "' has been removed successfully")
}

func (c *ClientSideGraph) listVertices() {
	fmt.Println("\nVertices:")
	if len(c.graph.nodes) == 0 {
		fmt.Println("No vertices")
		return
	}

	for i, node := range c.graph.nodes {
		fmt.Printf("%d: %v\n", i, node.Value)
	}
}

func (c *ClientSideGraph) listEdges(mode bool) {
	fmt.Println("\nEdges:")
	if c.graph.connectionsList == nil || len(c.graph.connectionsList) == 0 {
		fmt.Println("No edges")
		return
	}

	edgeCount := 0
	for node, edges := range c.graph.connectionsList {
		for _, edge := range edges {
			if mode {
				fmt.Print(edgeCount, ". ")
			}
			fmt.Print("From '", node.Value, "' to '", edge.List[1].Value, "'")
			if c.graph.isWeighted {
				fmt.Printf(" (weight: %.2f)", edge.Weight)
			}
			fmt.Println()
			edgeCount++
			edgeLst = append(edgeLst, edge)
		}
	}

	if edgeCount == 0 {
		fmt.Println("No edges")
	}
}

func (c *ClientSideGraph) changeGraphType() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Is the graph oriented? (y/n): ")
	orientedStr, _ := reader.ReadString('\n')
	oriented := strings.ToLower(strings.TrimSpace(orientedStr)) == "y"

	fmt.Print("Is the graph weighted? (y/n): ")
	weightedStr, _ := reader.ReadString('\n')
	weighted := strings.ToLower(strings.TrimSpace(weightedStr)) == "y"

	// Create new graph with new type but keep existing nodes
	nodes := c.graph.nodes
	c.graph = GraphConstructor(oriented, weighted)
	c.graph.nodes = nodes

	fmt.Printf("Graph type changed: oriented=%v, weighted=%v\n", oriented, weighted)
}

func (c *ClientSideGraph) printGraphInfo() {
	fmt.Println("\nGraph Information:")
	fmt.Printf("Type: %s, %s\n",
		map[bool]string{true: "Oriented", false: "Non-oriented"}[c.graph.isOriented],
		map[bool]string{true: "Weighted", false: "Non-weighted"}[c.graph.isWeighted])
	fmt.Printf("Number of vertices: %d\n", len(c.graph.nodes))

	edgeCount := 0
	if c.graph.connectionsList != nil {
		for _, edges := range c.graph.connectionsList {
			edgeCount += len(edges)
		}
	}
	fmt.Printf("Number of edges: %d\n", edgeCount)
}

func (c *ClientSideGraph) loadFromFile() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter file path: ")
	path, _ := reader.ReadString('\n')
	path = strings.TrimSpace(path)

	if path == "" {
		fmt.Println("No file path provided")
		return
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("File '%s' does not exist\n", path)
		return
	}

	newGraph := GraphFromFileConstructor(path)
	if newGraph != nil {
		c.graph = newGraph
		fmt.Printf("Graph successfully loaded from %s\n", path)
	} else {
		fmt.Printf("Failed to load graph from %s\n", path)
	}
}

func (c *ClientSideGraph) Run() {
	fmt.Println("Welcome to Graph ClientSideGraph!")

	for {
		c.printMenu()

		reader := bufio.NewReader(os.Stdin)
		choiceStr, _ := reader.ReadString('\n')
		choice, err := strconv.Atoi(strings.TrimSpace(choiceStr))

		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		switch choice {
		case 1:
			c.addVertex()
		case 2:
			c.addEdge()
		case 3:
			c.removeVertex()
		case 4:
			c.removeEdge()
		case 5:
			c.listVertices()
		case 6:
			c.listEdges(false)
		case 7:
			c.changeGraphType()
		case 8:
			c.printGraphInfo()
		case 9: // New case
			c.loadFromFile()
		case 10: // Updated exit
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid option. Please choose 1-9.")
		}
	}
}
