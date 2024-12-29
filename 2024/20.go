package main

import (
	"bufio"
	"fmt"
	"math"
	"strings"
)

func Solution20(input string) {
	grid := map[complex128]byte{}
	var start complex128

	i := 0
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		for j, char := range scanner.Bytes() {
			pos := complex(float64(i), float64(j))
			if char != '#' {
				grid[pos] = char
			}
			if char == 'S' {
				start = complex(float64(i), float64(j))
			}
		}
		i++
	}

	positions := []complex128{start}
	dist := map[complex128]int{start: 0}
	dirs := [4]complex128{1, -1, 1i, -1i}
	i = 0
	for i < len(positions) {
		pos := positions[i]
		for _, dir := range dirs {
			newPos := pos + dir
			if _, onPath := grid[newPos]; onPath {
				if _, visited := dist[newPos]; !visited {
					dist[newPos] = dist[pos] + 1
					positions = append(positions, newPos)
				}
			}
		}
		i++
	}

	part1, part2 := 0, 0
	for pos1, dist1 := range dist {
		for pos2, dist2 := range dist {
			taxiDist := int(math.Abs(real(pos1-pos2)) + math.Abs(imag(pos1-pos2)))
			if taxiDist == 2 && dist1-dist2-taxiDist >= 100 {
				part1++
			}
			if taxiDist <= 20 && dist1-dist2-taxiDist >= 100 {
				part2++
			}
		}
	}

	fmt.Println(part1)
	fmt.Println(part2)
}
