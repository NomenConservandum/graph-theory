package main

import (
	"fmt"
	"math"
)

// floydWarshallSimple реализует упрощённую версию алгоритма Флойда-Уоршелла
// Возвращает матрицу кратчайших расстояний и флаг наличия отрицательных циклов
func floydWarshallSimple(g *GraphInfo) (map[*Node]map[*Node]float64, bool) {
	// Инициализация матрицы расстояний
	dist := make(map[*Node]map[*Node]float64)

	// Инициализация для всех вершин
	for _, u := range g.nodes {
		dist[u] = make(map[*Node]float64)
		for _, v := range g.nodes {
			if u == v {
				dist[u][v] = 0
			} else {
				dist[u][v] = math.Inf(1) // Бесконечность
			}
		}
	}

	// Заполнение начальных расстояний из рёбер графа
	for u, edges := range g.connectionsList {
		for _, edge := range edges {
			v := edge.List[1]
			weight := edge.Weight

			// Если граф невзвешенный, используем вес 1
			if !g.isWeighted {
				weight = 1
			}

			if weight < dist[u][v] {
				dist[u][v] = weight
			}
		}
	}

	// Основной алгоритм Флойда-Уоршелла
	for _, k := range g.nodes {
		for _, i := range g.nodes {
			for _, j := range g.nodes {
				if dist[i][k] < math.Inf(1) && dist[k][j] < math.Inf(1) {
					if dist[i][j] > dist[i][k]+dist[k][j] {
						dist[i][j] = dist[i][k] + dist[k][j]
					}
				}
			}
		}
	}

	// Проверка на отрицательные циклы
	hasNegativeCycle := false
	for _, i := range g.nodes {
		if dist[i][i] < 0 {
			hasNegativeCycle = true
			break
		}
	}

	return dist, hasNegativeCycle
}

// printDistanceMatrix выводит матрицу расстояний в простом формате
func printDistanceMatrix(nodes []*Node, dist map[*Node]map[*Node]float64) {
	fmt.Println("Shortest Paths Matrix:")
	fmt.Print("     ")
	for _, node := range nodes {
		fmt.Printf("%5v", fmt.Sprintf("%v", node.Value))
	}
	fmt.Println()

	for _, u := range nodes {
		fmt.Printf("%5v", fmt.Sprintf("%v", u.Value))
		for _, v := range nodes {
			d := dist[u][v]
			if math.IsInf(d, 1) {
				fmt.Printf("%5s", "INF")
			} else {
				fmt.Printf("%5.1f", d)
			}
		}
		fmt.Println()
	}
}
