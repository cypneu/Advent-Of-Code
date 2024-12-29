package main

import (
	"bufio"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func Solution24(input string) {
	values := map[string]bool{}

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}

		wireValues := strings.Split(line, ": ")
		wireValue, _ := strconv.Atoi(wireValues[1])
		values[wireValues[0]] = wireValue == 1
	}
	highestZ := fmt.Sprintf("z%d", (len(values) / 2))

	graph := map[string]map[string][2]string{}
	operations := [][4]string{}
	for scanner.Scan() {
		gate := strings.Split(scanner.Text(), " ")
		if _, ok := graph[gate[0]]; !ok {
			graph[gate[0]] = map[string][2]string{}
		}
		if _, ok := graph[gate[2]]; !ok {
			graph[gate[2]] = map[string][2]string{}
		}

		operations = append(operations, [4]string{gate[0], gate[1], gate[2], gate[4]})
		graph[gate[0]][gate[4]] = [2]string{gate[1], gate[2]}
		graph[gate[2]][gate[4]] = [2]string{gate[1], gate[0]}
	}

	part1 := solvePart1(graph, values)
	part2 := solvePart2(operations, highestZ)
	fmt.Println(part1)
	fmt.Println(part2)
}

func solvePart1(graph map[string]map[string][2]string, values map[string]bool) int64 {
	visited := map[string]struct{}{}
	evalInOrder := [][4]string{}
	for node := range graph {
		if _, seen := visited[node]; !seen {
			doSolvePart1(values, graph, visited, node, &evalInOrder)
		}
	}

	for i := len(evalInOrder) - 1; i >= 0; i-- {
		gateInfo := evalInOrder[i]
		if gateInfo[1] == "XOR" {
			values[gateInfo[3]] = values[gateInfo[0]] != values[gateInfo[2]]
		} else if gateInfo[1] == "AND" {
			values[gateInfo[3]] = values[gateInfo[0]] && values[gateInfo[2]]
		} else {
			values[gateInfo[3]] = values[gateInfo[0]] || values[gateInfo[2]]
		}

	}

	zValues := []string{}
	for node := range values {
		if strings.HasPrefix(node, "z") {
			zValues = append(zValues, node)
		}
	}

	sort.Strings(zValues)

	var resBuilder strings.Builder
	for i := len(zValues) - 1; i >= 0; i-- {
		if values[zValues[i]] {
			resBuilder.WriteString("1")
		} else {
			resBuilder.WriteString("0")
		}
	}
	result := resBuilder.String()
	res, _ := strconv.ParseInt(result, 2, 0)
	return res
}

func solvePart2(operations [][4]string, highestZ string) string {
	incorrect := map[string]struct{}{}
	for _, gate := range operations {
		// All gates that output z must be XOR
		if gate[1] != "XOR" && gate[3][0] == 'z' && gate[3] != highestZ {
			incorrect[gate[3]] = struct{}{}
		}

		// All XOR gates must output z expect input x,y gates
		if gate[1] == "XOR" && gate[3][0] != 'z' && gate[0][0] != 'x' && gate[0][0] != 'y' {
			incorrect[gate[3]] = struct{}{}
		}

		// Output from AND gates must be supplied to OR gates
		if gate[1] == "AND" && gate[0] != "x00" && gate[2] != "x00" {
			for _, gateAlt := range operations {
				if (gate[3] == gateAlt[0] || gate[3] == gateAlt[2]) && gateAlt[1] != "OR" {
					incorrect[gate[3]] = struct{}{}
				}
			}
		}

		// Output from XOR gates mustn't be supplied to OR gates
		if gate[1] == "XOR" {
			for _, gateAlt := range operations {
				if (gate[3] == gateAlt[0] || gate[3] == gateAlt[2]) && gateAlt[1] == "OR" {
					incorrect[gate[3]] = struct{}{}
				}
			}
		}
	}

	res := []string{}
	for operator := range incorrect {
		res = append(res, operator)
	}
	sort.Strings(res)
	return strings.Join(res, ",")
}

func doSolvePart1(values map[string]bool, graph map[string]map[string][2]string, visited map[string]struct{}, node string, evalInOrder *[][4]string) {
	visited[node] = struct{}{}

	for nextNode := range graph[node] {
		if _, seen := visited[nextNode]; !seen {
			doSolvePart1(values, graph, visited, nextNode, evalInOrder)
		}
	}

	for nextNode, gateInfo := range graph[node] {
		*evalInOrder = append(*evalInOrder, [4]string{node, gateInfo[0], gateInfo[1], nextNode})
	}
}
