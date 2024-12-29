package main

import (
	"bufio"
	"fmt"
	"strings"
)

func Solution18(input string) {
	spaceSize := 71

	corruptedMemory := [][2]int{}
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		var pos [2]int
		fmt.Sscanf(scanner.Text(), "%d,%d", &pos[0], &pos[1])
		corruptedMemory = append(corruptedMemory, pos)
	}

	fmt.Println(findExit(memoryMap(corruptedMemory, 1024), spaceSize))

	// Part 2
	// Can be solved even better by using union-find going from reverse
	// until startPos and endPos are in the same set
	left, right := 1025, len(corruptedMemory)-1
	for left < right {
		mid := left + (right-left)/2
		if findExit(memoryMap(corruptedMemory, mid), spaceSize) == -1 {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	fmt.Println(corruptedMemory[left])
}

func memoryMap(corruptedMemory [][2]int, kilobytes int) map[[2]int]struct{} {
	memoryMap := map[[2]int]struct{}{}
	for _, pos := range corruptedMemory[:kilobytes] {
		memoryMap[pos] = struct{}{}
	}
	return memoryMap
}

func findExit(corruptedMemory map[[2]int]struct{}, spaceSize int) int {
	dist := 0
	startPos, endPos := [2]int{0, 0}, [2]int{spaceSize - 1, spaceSize - 1}
	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	visited := map[[2]int]struct{}{startPos: {}}

	queue := [][2]int{startPos}
	for len(queue) > 0 {
		for levelSize := len(queue); levelSize > 0; levelSize-- {
			pos := queue[0]
			queue = queue[1:]

			if pos == endPos {
				return dist
			}

			for _, dir := range dirs {
				newPos := [2]int{pos[0] + dir[0], pos[1] + dir[1]}
				if inRange(newPos[0], newPos[1], spaceSize, spaceSize) {
					if _, ok := corruptedMemory[newPos]; !ok {
						if _, ok := visited[newPos]; !ok {
							visited[newPos] = struct{}{}
							queue = append(queue, newPos)
						}
					}
				}
			}
		}
		dist += 1
	}
	return -1
}
