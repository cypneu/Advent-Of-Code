package main

import (
	"fmt"
	"sort"
	"strings"
)

type Graph map[string]map[string]struct{}

func Solution23(input string) {
	graph := Graph{}

	connections := strings.Split(strings.TrimSpace(input), "\n")
	for _, connection := range connections {
		connectionMembers := strings.Split(connection, "-")

		if _, exists := graph[connectionMembers[0]]; !exists {
			graph[connectionMembers[0]] = map[string]struct{}{}
		}
		if _, exists := graph[connectionMembers[1]]; !exists {
			graph[connectionMembers[1]] = map[string]struct{}{}
		}

		graph[connectionMembers[0]][connectionMembers[1]] = struct{}{}
		graph[connectionMembers[1]][connectionMembers[0]] = struct{}{}
	}

	part1 := 0
	for node := range graph {
		for adjacentNode := range graph[node] {
			for anotherAdjacentNode := range graph[node] {
				if adjacentNode < anotherAdjacentNode && node < adjacentNode && (strings.HasPrefix(node, "t") || strings.HasPrefix(adjacentNode, "t") || strings.HasPrefix(anotherAdjacentNode, "t")) {
					if _, seen := graph[adjacentNode][anotherAdjacentNode]; seen {
						part1++
					}
				}
			}
		}
	}
	fmt.Println(part1)
	fmt.Println(strings.Join(largestClique(graph), ","))
}

func largestClique(g Graph) []string {
	nodes := make([]string, 0, len(g))
	for node := range g {
		nodes = append(nodes, node)
	}

	best := []string{}
	backtrack(g, nodes, []string{}, &best)

	sort.Strings(best)
	return best
}

func backtrack(g Graph, candidates, clique []string, best *[]string) {
	for i, node := range candidates {
		if isConnectedToAll(g, node, clique) {
			newClique := append(clique, node)
			backtrack(g, candidates[i+1:], newClique, best)
			if len(newClique) > len(*best) {
				*best = append([]string{}, newClique...)
			}
		}
	}
}

func isConnectedToAll(g Graph, node string, clique []string) bool {
	for _, c := range clique {
		if _, ok := g[node][c]; !ok {
			return false
		}
	}
	return true
}
