package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func Solution5A(input string) {
	total := 0
	readDependencies := true
	pageOrderingRules := make(map[string]map[string]struct{})

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			readDependencies = false
		} else if readDependencies {
			orderings := strings.Split(line, "|")
			if _, exists := pageOrderingRules[orderings[0]]; !exists {
				pageOrderingRules[orderings[0]] = make(map[string]struct{})
			}
			pageOrderingRules[orderings[0]][orderings[1]] = struct{}{}
		} else {
			update := strings.Split(line, ",")
			middle, _ := strconv.Atoi(update[len(update)/2])
			if checkAllDependencies(update, pageOrderingRules) {
				total += middle
			}
		}
	}

	fmt.Println(total)
}

func Solution5B(input string) {
	total := 0
	readDependencies := true
	pageOrderingRules := make(map[string]map[string]struct{})

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			readDependencies = false
		} else if readDependencies {
			orderings := strings.Split(line, "|")
			if _, exists := pageOrderingRules[orderings[0]]; !exists {
				pageOrderingRules[orderings[0]] = make(map[string]struct{})
			}
			pageOrderingRules[orderings[0]][orderings[1]] = struct{}{}
		} else {
			update := strings.Split(line, ",")
			updateSet := sliceToSet(update)

			if !checkAllDependencies(update, pageOrderingRules) {
				res := topologicalSort(updateSet, pageOrderingRules)
				middle, _ := strconv.Atoi(res[len(res)/2])
				total += middle
			}
		}
	}

	fmt.Println(total)
}

func checkAllDependencies(update []string, pageOrderingRules map[string]map[string]struct{}) bool {
	for i := 0; i < len(update); i++ {
		if !checkDependenciesForPage(update, pageOrderingRules, i) {
			return false
		}
	}
	return true
}

func checkDependenciesForPage(update []string, pageOrderingRules map[string]map[string]struct{}, i int) bool {
	for j := i + 1; j < len(update); j++ {
		if _, exists := pageOrderingRules[update[i]][update[j]]; !exists {
			return false
		}
	}
	return true
}

func topologicalSort(update map[string]struct{}, pageOrderingRules map[string]map[string]struct{}) []string {
	res := []string{}
	visited := make(map[string]struct{})

	for node := range update {
		if _, exists := visited[node]; !exists {
			dfs(&res, visited, pageOrderingRules, node, update)
		}
	}

	return reverse(res)
}

func dfs(res *[]string, visited map[string]struct{}, pageOrderingRules map[string]map[string]struct{}, node string, update map[string]struct{}) {
	if _, exists := visited[node]; exists {
		return
	}

	for dependency := range pageOrderingRules[node] {
		if _, inUpdate := update[dependency]; inUpdate {
			if _, exists := visited[dependency]; !exists {
				dfs(res, visited, pageOrderingRules, dependency, update)
			}
		}
	}

	visited[node] = struct{}{}
	*res = append(*res, node)
}

func reverse(input []string) []string {
	for i, j := 0, len(input)-1; i < j; i, j = i+1, j-1 {
		input[i], input[j] = input[j], input[i]
	}
	return input
}

func sliceToSet(slice []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, item := range slice {
		set[item] = struct{}{}
	}
	return set
}
