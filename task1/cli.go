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
	"fmt"
	"os"
	"strconv"
	"strings"
)

type CLI struct {
	graph *GraphInfo
}

func NewCLI() *CLI {
	return &CLI{
		graph: GraphConstructor(false, false), // default: undirected, unweighted
	}
}

func (c *CLI) addVertex() {
	var vt string
	fmt.Print("Enter vertex value: ")
	fmt.Scanln(&vt)
	value := strings.TrimSpace(vt)

	node := NodeConstructor(value)
	addVertex(c.graph, node)
	fmt.Println("Vertex '", value, "' added successfully")
}

func (c *CLI) addEdge() {
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
		err := addEdge(c.graph, node1, node2, weight)
		if err != nil {
			println(err.Error())
			return
		}
		fmt.Print("Added oriented edge from '", node1.Value, "' to '", node2.Value, "'")
		if c.graph.isWeighted {
			fmt.Print(" with weight ", weight)
		}
		fmt.Println()
	} else {
		if c.graph.isWeighted {
			err := addNonOrientedEdge(c.graph, node1, node2, weight)
			if err != nil {
				println(err.Error())
				return
			}
		} else {
			err := addNonOrientedNonWeightedEdge(c.graph, node1, node2)
			if err != nil {
				println(err.Error())
				return
			}
		}
		fmt.Printf("Added non-oriented edge between '%v' and '%v'", node1.Value, node2.Value)
		if c.graph.isWeighted {
			fmt.Printf(" with weight %.2f", weight)
		}
		fmt.Println()
	}
}

func (c *CLI) removeVertex() {
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

func (c *CLI) removeEdge() {
	if c.graph.connectionsList == nil || len(c.graph.connectionsList) == 0 {
		fmt.Println("No edges to remove")
		return
	}

	err := c.listEdges(true)
	if err != nil {
		return
	}
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

func (c *CLI) listVertices() {
	fmt.Println("\nVertices:")
	if len(c.graph.nodes) == 0 {
		fmt.Println("No vertices")
		return
	}

	for i, node := range c.graph.nodes {
		fmt.Printf("%d: %v\n", i, node.Value)
	}
}

func (c *CLI) listEdges(mode bool) error {
	fmt.Println("\nEdges:")
	if c.graph.connectionsList == nil || len(c.graph.connectionsList) == 0 {
		fmt.Println("No edges")
		return fmt.Errorf("")
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
		return fmt.Errorf("")
	}
	return nil
}

func (c *CLI) changeGraphType() {
	var vt string
	fmt.Print("Is the graph oriented? (y/n): ")
	fmt.Scanln(&vt)
	orientedStr := strings.TrimSpace(vt)
	oriented := strings.ToLower(orientedStr) == "y"

	fmt.Print("Is the graph weighted? (y/n): ")
	fmt.Scanln(&vt)
	weightedStr := strings.TrimSpace(vt)
	weighted := strings.ToLower(weightedStr) == "y"

	// Create new graph with new type but keep existing nodes
	nodes := c.graph.nodes
	c.graph = GraphConstructor(oriented, weighted)
	c.graph.nodes = nodes

	fmt.Printf("Graph type changed: oriented=%v, weighted=%v\n", oriented, weighted)
}

func (c *CLI) printGraphInfo() {
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

func (c *CLI) loadFromFile() {
	var vt string
	fmt.Print("Enter file path: ")
	fmt.Scanln(&vt)
	path := strings.TrimSpace(vt)

	if path == "" {
		fmt.Println("No file path provided")
		return
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("File '", path, "' does not exist")
		return
	}

	newGraph := GraphFromFileConstructor(path)
	if newGraph != nil {
		c.graph = newGraph
		fmt.Println("Graph successfully loaded from", path)
	} else {
		fmt.Println("Failed to load graph from", path)
	}
}

func (c *CLI) saveToFile() {
	var vt string
	fmt.Print("Enter file path: ")
	fmt.Scanln(&vt)
	path := strings.TrimSpace(vt)
	path = strings.TrimSpace(path)

	if path == "" {
		fmt.Println("No file path provided")
		return
	}

	var saveErr = WriteToFile(c.graph, path)

	if saveErr != nil {
		fmt.Printf("Error saving file: %v\n", saveErr)
	} else {
		fmt.Println("File saved successfully!")
	}
}

func (c *CLI) listKnots() error {
	fmt.Println("\nEdges:")
	if c.graph.connectionsList == nil || len(c.graph.connectionsList) == 0 {
		fmt.Println("No edges")
		return fmt.Errorf("")
	}

	var knots = knots(c.graph)

	for i, edges := range knots {
		fmt.Println(i, ". Vertex with value ", edges.List[0], "is a knot")
	}
	return nil
}

func (c *CLI) printMainMenu() {
	fmt.Println("\n=== Welcome! ===\n1. Add graph\n2. Select graph\n3. Exit\nChoose an option: ")
}

func (c *CLI) addGraphMenu() {
	fmt.Println("\n=== Add graph ===\n1. Load from file\n2. Write manually\n3. Exit\nChoose an option: ")
}

func (c *CLI) printMenu() {
	fmt.Println("\n=== Graph CLI ===\n1. Add vertex\n2. Add edge\n3. Remove vertex\n4. Remove edge\n5. List vertices\n6. List edges\n7. Change graph type\n8. Print graph info\n9. Load from File\n10. Save to file\n11. List Knots\n12. Exit\nChoose an option: ")
}

func (c *CLI) Run() {
	fmt.Println("Welcome to Graph CLI!")

	for {
		c.printMenu()

		var vt string
		fmt.Scanln(&vt)
		choiceStr := strings.TrimSpace(vt)
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
		case 9:
			c.loadFromFile()
		case 10:
			c.saveToFile()
		case 11:
			c.listKnots()
		case 12:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid option. Please choose 1-9.")
		}
	}
}
