package main

import (
	"bufio"
	"fmt"
	"strings"
)

func Solution12A(input string) {
	totalPrice := 0
	var garden []string

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		garden = append(garden, scanner.Text())
	}

	dirs := []Position{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	visited := make([]bool, len(garden)*len(garden[0]))
	for i := 0; i < len(garden); i++ {
		for j := 0; j < len(garden[i]); j++ {
			if !visited[i*len(garden[i])+j] {
				visited[i*len(garden[0])+j] = true
				nodes, plots := bfs(garden, i, j, visited, dirs)
				totalPrice += nodes * plots
			}
		}
	}
	fmt.Println(totalPrice)
}

func bfs(garden []string, i, j int, visited []bool, dirs []Position) (int, int) {
	stack := []Position{{i, j}}
	plots, nodes := 0, 0
	for len(stack) > 0 {
		n := len(stack)
		currPos := stack[n-1]
		stack = stack[:n-1]
		nodes++

		for _, dir := range dirs {
			newPos := Position{currPos.row + dir.row, currPos.col + dir.col}
			if inBounds(newPos, len(garden), len(garden[i])) && garden[newPos.row][newPos.col] == garden[currPos.row][currPos.col] {
				if !visited[newPos.row*len(garden[0])+newPos.col] {
					visited[newPos.row*len(garden[0])+newPos.col] = true
					stack = append(stack, newPos)
				}
			} else {
				plots++
			}
		}
	}
	return nodes, plots
}

func Solution12B(input string) {
	totalPrice := 0
	var garden []string
	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		garden = append(garden, scanner.Text())
	}

	dirs := []Position{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	visited := make([]bool, len(garden)*len(garden[0]))
	for i := 0; i < len(garden); i++ {
		for j := 0; j < len(garden[i]); j++ {
			if !visited[i*len(garden[i])+j] {
				visited[i*len(garden[0])+j] = true
				nodes, plots := bfs2(garden, i, j, visited, dirs)
				totalPrice += nodes * plots
			}
		}
	}
	fmt.Println(totalPrice)
}

func bfs2(garden []string, i, j int, visited []bool, dirs []Position) (int, int) {
	stack := []Position{{i, j}}
	n, m := len(garden), len(garden[0])
	nodes, plots := 0, 0
	for len(stack) > 0 {
		currPos := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		nodes++

		for k, dir := range dirs {
			if isInnerCorner(garden, currPos, dir, dirs[(k+1)%4], n, m) {
				plots++
			}

			newPos := Position{currPos.row + dir.row, currPos.col + dir.col}
			if inBounds(newPos, len(garden), len(garden[i])) && garden[newPos.row][newPos.col] == garden[currPos.row][currPos.col] {
				if !visited[newPos.row*len(garden[0])+newPos.col] {
					visited[newPos.row*len(garden[0])+newPos.col] = true
					stack = append(stack, newPos)
				}
			} else {
				first := Position{currPos.row + dir.row, currPos.col + dir.col}
				second := Position{currPos.row + dirs[(k+1)%4].row, currPos.col + dirs[(k+1)%4].col}
				if isOuterCorner(garden, garden[currPos.row][currPos.col], first, second, n, m) {
					plots++
				}
			}
		}
	}

	return nodes, plots
}

func isInnerCorner(garden []string, pos, dir Position, nextDir Position, n, m int) bool {
	newPos1 := Position{pos.row + dir.row, pos.col + dir.col}
	newPos2 := Position{pos.row + nextDir.row, pos.col + nextDir.col}
	newPos3 := Position{pos.row + dir.row + nextDir.row, pos.col + dir.col + nextDir.col}
	if !inBounds(newPos1, n, m) || !inBounds(newPos2, n, m) || !inBounds(newPos3, n, m) {
		return false
	}

	return garden[newPos1.row][newPos1.col] == garden[pos.row][pos.col] && garden[newPos2.row][newPos2.col] == garden[pos.row][pos.col] && garden[newPos3.row][newPos3.col] != garden[pos.row][pos.col]
}

func isOuterCorner(garden []string, color byte, first, second Position, n, m int) bool {
	firstDifferent := first.row < 0 || first.row >= n || first.col < 0 || first.col >= m || garden[first.row][first.col] != color
	secondDifferent := second.row < 0 || second.row >= n || second.col < 0 || second.col >= m || garden[second.row][second.col] != color
	if firstDifferent && secondDifferent {
		return true
	}
	return false
}
