package main

import (
	"bufio"
	"fmt"
	"strings"
)

type Position struct {
	row, col int
}

func Solution8A(input string) {
	total := solve8(input, strategyA)
	fmt.Println(total)
}

func Solution8B(input string) {
	total := solve8(input, strategyB)
	fmt.Println(total)
}

func solve8(input string, strategy func(Position, Position, *[]bool, int, int) int) int {
	total := 0
	antennas, n, m := parseAntennas(input)
	taken := make([]bool, n*m)
	for _, a := range antennas {
		for i := 0; i < len(a); i++ {
			for j := 0; j < len(a); j++ {
				if i == j {
					continue
				}
				total += strategy(a[i], a[j], &taken, n, m)
			}
		}
	}
	return total
}

func strategyA(a1, a2 Position, taken *[]bool, n, m int) int {
	new_j := a1.col + (a1.col - a2.col)
	new_i := a1.row + (a1.row - a2.row)
	if inRange(new_i, new_j, n, m) && markPosition(taken, Position{new_i, new_j}, n, m) {
		return 1
	}
	return 0
}

func parseAntennas(input string) (map[rune][]Position, int, int) {
	antennas := make(map[rune][]Position)
	row, cols := 0, 0

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		cols = len(line)
		for col, char := range line {
			if char != '.' {
				antennas[char] = append(antennas[char], Position{row, col})
			}
		}
		row++
	}

	return antennas, row, cols
}

func strategyB(a1, a2 Position, taken *[]bool, n, m int) int {
	total := 0
	new_i, new_j := a1.row, a1.col
	for k := 0; inRange(new_i, new_j, n, m); k++ {
		if markPosition(taken, Position{new_i, new_j}, n, m) {
			total++
		}

		new_i = a1.row + (k+1)*(a1.row-a2.row)
		new_j = a1.col + (k+1)*(a1.col-a2.col)
	}
	return total
}

func markPosition(taken *[]bool, pos Position, n, m int) bool {
	index := pos.row*m + pos.col
	if !((*taken)[index]) {
		(*taken)[index] = true
		return true
	}
	return false
}

func inRange(i, j, n, m int) bool {
	return i >= 0 && i < n && j >= 0 && j < m
}
