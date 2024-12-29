package main

import (
	"bufio"
	"fmt"
	"math"
	"strings"
)

func Solution14A(input string) {
	n, m := 103, 101
	safetyFactor := 1

	q := [4]int{}
	scanner := bufio.NewScanner(strings.NewReader(input))
	var px, py, vx, vy int

	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "p=%d,%d v=%d,%d", &px, &py, &vx, &vy)

		px = (px + 100*vx + 100*m) % m
		py = (py + 100*vy + 100*n) % n

		if px != m/2 && py != n/2 {
			q[py/((n+1)/2)*2+px/((m+1)/2)]++
		}
	}
	for i := 0; i < 4; i++ {
		safetyFactor *= q[i]
	}
	fmt.Println(safetyFactor)
}

func Solution14B(input string) {
	n, m := 103, 101

	scanner := bufio.NewScanner(strings.NewReader(input))
	var positions, initPos, velocities [][2]int

	var px, py, vx, vy int
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "p=%d,%d v=%d,%d", &px, &py, &vx, &vy)
		positions = append(positions, [2]int{px, py})
		initPos = append(initPos, [2]int{px, py})
		velocities = append(velocities, [2]int{vx, vy})
	}

	minVar, minVarIdx := math.MaxInt64, -1
	for i := 0; i < n*m; i++ {
		center := [2]int{0, 0}
		for j := 0; j < len(positions); j++ {
			positions[j][0] = (positions[j][0] + velocities[j][0] + m) % m
			positions[j][1] = (positions[j][1] + velocities[j][1] + n) % n
			center[0] += positions[j][0]
			center[1] += positions[j][1]
		}
		center[0] /= len(positions)
		center[1] /= len(positions)

		variance := 0
		for _, pos := range positions {
			variance += (pos[0]-center[0])*(pos[0]-center[0]) + (pos[1]-center[1])*(pos[1]-center[1])
		}

		if variance < minVar {
			minVar = variance
			minVarIdx = i + 1
		}
	}
	fmt.Println(minVarIdx)
	printChristmasTree(initPos, velocities, minVarIdx, n, m)
}

func printChristmasTree(initPos, velocities [][2]int, minVarIdx, n, m int) {
	posMap := make(map[[2]int]struct{})
	for i := 0; i < len(initPos); i++ {
		initPos[i][0] = (initPos[i][0] + velocities[i][0]*minVarIdx + minVarIdx*m) % m
		initPos[i][1] = (initPos[i][1] + velocities[i][1]*minVarIdx + minVarIdx*n) % n
		posMap[initPos[i]] = struct{}{}
	}

	for y := 0; y < n; y++ {
		for x := 0; x < m; x++ {
			if _, ok := posMap[[2]int{x, y}]; ok {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
