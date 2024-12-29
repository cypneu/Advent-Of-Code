package main

import (
	"bufio"
	"fmt"
	"strings"
)

type Vector struct {
	X int
	Y int
}

func Solution4A(input string) {
	pattern := "XMAS"
	var counter int

	grid := inputIntoGrid(input)
	directions := []Vector{{1, 0}, {0, 1}, {-1, 0}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	isInRange := func(i, j int) bool { return i >= 0 && j >= 0 && i < len(grid) && j < len(grid[i]) }

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			for _, dir := range directions {
				if checkPattern(grid, pattern, i, j, dir, isInRange) {
					counter++
				}
			}
		}
	}
	fmt.Println(counter)
}

func Solution4B(input string) {
	var counter int
	grid := inputIntoGrid(input)

	for i := 0; i < len(grid)-2; i++ {
		for j := 0; j < len(grid[i])-2; j++ {
			areDiagCornersOk1 := (grid[i][j] == 'M' && grid[i+2][j+2] == 'S') || (grid[i][j] == 'S' && grid[i+2][j+2] == 'M')
			areDiagCornersOk2 := (grid[i][j+2] == 'M' && grid[i+2][j] == 'S') || (grid[i][j+2] == 'S' && grid[i+2][j] == 'M')

			if grid[i+1][j+1] == 'A' && areDiagCornersOk1 && areDiagCornersOk2 {
				counter++
			}
		}
	}
	fmt.Println(counter)
}

func inputIntoGrid(input string) []string {
	var grid []string
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		grid = append(grid, line)
	}
	return grid
}

func checkPattern(grid []string, pattern string, i, j int, dir Vector, isInRange func(int, int) bool) bool {
	for k := 0; k < len(pattern); k++ {
		new_i, new_j := i+k*dir.Y, j+k*dir.X
		if !isInRange(new_i, new_j) || grid[new_i][new_j] != pattern[k] {
			return false
		}
	}
	return true
}
