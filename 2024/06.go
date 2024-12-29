package main

import (
	"bufio"
	"fmt"
	"strings"
)

func Solution6A(input string) {
	var area [][]rune

	scanner := bufio.NewScanner(strings.NewReader(input))
	i, j := 0, 0
	var dir string

	for scanner.Scan() {
		line := scanner.Text()
		runeLine := []rune(line)

		if dir == "" {
			for col, c := range line {
				if dir == "" {
					j = col
				}
				if c == '^' {
					dir = "up"
				} else if c == 'v' {
					dir = "down"
				} else if c == '<' {
					dir = "left"
				} else if c == '>' {
					dir = "right"
				}
			}
			if dir == "" {
				i++
			}
		}

		area = append(area, runeLine)
	}

	visited := make(map[int]struct{})
	var total int
	for i >= 0 && i < len(area) && j >= 0 && j < len(area[i]) {
		if _, exists := visited[i*len(area[i])+j]; !exists {
			visited[i*len(area[i])+j] = struct{}{}
			total++
		}
		i, j, dir = nextPosition(i, j, dir, area)
	}

	fmt.Println(total)
}

func Solution6B(input string) {
	var area [][]rune

	i, j := 0, 0
	var dir string
	var total int

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		runeLine := []rune(line)

		if dir == "" {
			for col, c := range runeLine {
				if dir == "" {
					j = col
				}
				if c == '^' {
					dir = "up"
				} else if c == 'v' {
					dir = "down"
				} else if c == '<' {
					dir = "left"
				} else if c == '>' {
					dir = "right"
				}
			}
			if dir == "" {
				i++
			}
		}

		area = append(area, runeLine)
	}

	new_i, new_j, new_dir := i, j, dir
	obstacles := make(map[struct{ i, j int }]struct{})
	for i >= 0 && i < len(area) && j >= 0 && j < len(area[i]) {
		obstI, obstJ := placeObstacle(i, j, dir)
		if obstI >= 0 && obstI < len(area) && obstJ >= 0 && obstJ < len(area[0]) && area[obstI][obstJ] != '#' && area[obstI][obstJ] != '^' && area[obstI][obstJ] != 'v' && area[obstI][obstJ] != '<' && area[obstI][obstJ] != '>' {
			obstacles[struct{ i, j int }{obstI, obstJ}] = struct{}{}
		}

		i, j, dir = nextPosition(i, j, dir, area)
	}

	i, j, dir = new_i, new_j, new_dir
	for obst := range obstacles {
		area[obst.i][obst.j] = '#'

		if loopExists(i, j, dir, obst, area) {
			total++
		}

		area[obst.i][obst.j] = '.'
	}

	fmt.Println(total)
}

func loopExists(i, j int, dir string, obst struct{ i, j int }, area [][]rune) bool {
	visited := make(map[struct {
		i, j int
		dir  string
	}]struct{})
	for i >= 0 && i < len(area) && j >= 0 && j < len(area[i]) {
		key := struct {
			i, j int
			dir  string
		}{i, j, dir}
		visited[key] = struct{}{}

		i, j, dir = nextPosition(i, j, dir, area)
		key = struct {
			i, j int
			dir  string
		}{i, j, dir}
		if _, exists := visited[key]; exists {
			return true
		}
	}
	return false
}

func placeObstacle(i, j int, dir string) (int, int) {
	if dir == "up" {
		return i - 1, j
	} else if dir == "down" {
		return i + 1, j
	} else if dir == "left" {
		return i, j - 1
	} else {
		return i, j + 1
	}
}

func nextPosition(i, j int, dir string, area [][]rune) (int, int, string) {
	if dir == "up" {
		if i-1 >= 0 && area[i-1][j] == '#' {
			dir = "right"
		} else {
			i--
		}
	} else if dir == "down" {
		if i+1 < len(area) && area[i+1][j] == '#' {
			dir = "left"
		} else {
			i++
		}
	} else if dir == "left" {
		if j-1 >= 0 && area[i][j-1] == '#' {
			dir = "up"
		} else {
			j--
		}
	} else if dir == "right" {
		if j+1 < len(area[i]) && area[i][j+1] == '#' {
			dir = "down"
		} else {
			j++
		}
	}
	return i, j, dir
}
