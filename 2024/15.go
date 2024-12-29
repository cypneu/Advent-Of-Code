package main

import (
	"bufio"
	"bytes"
	"fmt"
	"slices"
	"strings"
)

func Solution15A(input string) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	var grid [][]byte
	var moves []byte
	readMap := true
	robotPos := [2]int{-1, -1}
	i := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			readMap = false
			continue
		}

		if readMap {
			idx := bytes.Index(line, []byte("@"))
			if idx >= 0 {
				robotPos[0], robotPos[1] = i, idx
			}
			grid = append(grid, append([]byte{}, line...))
		} else {
			moves = append(moves, line...)
		}
		i++
	}

	dirsByMove := map[byte][2]int{
		'<': {0, -1},
		'>': {0, 1},
		'^': {-1, 0},
		'v': {1, 0},
	}
	for _, move := range moves {
		moveRobot(grid, &robotPos, dirsByMove[move])
	}

	gpsTotal := 0
	for i, row := range grid {
		for j, ch := range row {
			if ch == 'O' {
				gpsTotal += 100*i + j
			}
		}
	}
	fmt.Println(gpsTotal)
}

func moveRobot(grid [][]byte, robotPos *[2]int, dir [2]int) {
	currPos := [2]int{robotPos[0], robotPos[1]}
	for grid[currPos[0]+dir[0]][currPos[1]+dir[1]] == 'O' {
		currPos[0], currPos[1] = currPos[0]+dir[0], currPos[1]+dir[1]
	}

	if grid[currPos[0]+dir[0]][currPos[1]+dir[1]] == '.' {
		for currPos[0] != robotPos[0] || currPos[1] != robotPos[1] {
			grid[currPos[0]+dir[0]][currPos[1]+dir[1]] = grid[currPos[0]][currPos[1]]
			currPos[0], currPos[1] = currPos[0]-dir[0], currPos[1]-dir[1]
		}
		grid[currPos[0]+dir[0]][currPos[1]+dir[1]] = grid[currPos[0]][currPos[1]]
		grid[currPos[0]][currPos[1]] = '.'
		robotPos[0], robotPos[1] = currPos[0]+dir[0], currPos[1]+dir[1]
	}
}

func Solution15B(input string) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	var grid [][]byte
	var moves []byte
	readMap := true
	var robotPos [2]int
	i := 0
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			readMap = false
			continue
		}

		if readMap {
			extendedLine := make([]byte, 2*len(line))
			for j, ch := range line {
				if ch == 'O' {
					extendedLine[2*j], extendedLine[2*j+1] = '[', ']'
				} else if ch == '@' {
					robotPos = [2]int{i, 2 * j}
					extendedLine[2*j], extendedLine[2*j+1] = '@', '.'
				} else {
					extendedLine[2*j], extendedLine[2*j+1] = line[j], line[j]
				}
			}
			grid = append(grid, extendedLine)
		} else {
			moves = append(moves, line...)
		}
		i++
	}

	dirsByMove := map[byte][2]int{
		'<': {0, -1},
		'>': {0, 1},
		'^': {-1, 0},
		'v': {1, 0},
	}
	for _, move := range moves {
		dir := dirsByMove[move]
		if dir[1] != 0 {
			moveRobotHorizontally(grid, &robotPos, dir)
		} else {
			moveRobotVertically(grid, &robotPos, dir)
		}
	}

	gpsTotal := 0
	for i, row := range grid {
		for j, ch := range row {
			if ch == '[' {
				gpsTotal += 100*i + j
			}
		}
	}
	fmt.Println(gpsTotal)
}

func moveRobotHorizontally(grid [][]byte, robotPos *[2]int, dir [2]int) {
	currPos := *robotPos
	for slices.Contains([]byte{'[', ']'}, grid[currPos[0]+dir[0]][currPos[1]+dir[1]]) {
		currPos[0], currPos[1] = currPos[0]+dir[0], currPos[1]+dir[1]
	}

	if grid[currPos[0]+dir[0]][currPos[1]+dir[1]] == '.' {
		for currPos != *robotPos {
			grid[currPos[0]+dir[0]][currPos[1]+dir[1]] = grid[currPos[0]][currPos[1]]
			currPos[0], currPos[1] = currPos[0]-dir[0], currPos[1]-dir[1]
		}
		grid[currPos[0]+dir[0]][currPos[1]+dir[1]] = grid[currPos[0]][currPos[1]]
		grid[currPos[0]][currPos[1]] = '.'
		robotPos[0], robotPos[1] = currPos[0]+dir[0], currPos[1]+dir[1]
	}
}

func moveRobotVertically(grid [][]byte, robotPos *[2]int, dir [2]int) {
	posesToMove := [][2]int{*robotPos}
	stack := [][2]int{*robotPos}
	visited := map[[2]int]struct{}{*robotPos: {}}

	for len(stack) > 0 {
		currPos := stack[0]
		stack = stack[1:]
		posesToMove = append(posesToMove, currPos)

		for _, next := range getNeighbors(currPos, dir, grid) {
			if grid[next[0]][next[1]] == '#' {
				return
			}

			if _, ok := visited[next]; !ok {
				visited[next] = struct{}{}
				stack = append(stack, next)
			}
		}
	}

	for i := len(posesToMove) - 1; i >= 0; i-- {
		currPos := posesToMove[i]
		grid[currPos[0]+dir[0]][currPos[1]+dir[1]] = grid[currPos[0]][currPos[1]]
		grid[currPos[0]][currPos[1]] = '.'
	}
	robotPos[0], robotPos[1] = robotPos[0]+dir[0], robotPos[1]+dir[1]
}

func getNeighbors(pos [2]int, dir [2]int, grid [][]byte) [][2]int {
	neighbors := [][2]int{}
	newPos := [2]int{pos[0] + dir[0], pos[1] + dir[1]}

	switch grid[newPos[0]][newPos[1]] {
	case '[':
		neighbors = append(neighbors, newPos, [2]int{newPos[0], newPos[1] + 1})
	case ']':
		neighbors = append(neighbors, newPos, [2]int{newPos[0], newPos[1] - 1})
	case '#':
		neighbors = append(neighbors, newPos)
	}

	return neighbors
}
