package main

import (
	"bufio"
	"fmt"
	"strings"
)

func Solution25(input string) {
	locks, keys := [][5]int{}, [][5]int{}
	schema := [7]string{}
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		for i := 0; i < 7; i++ {
			schema[i] = scanner.Text()
			scanner.Scan()
		}

		if isLock(schema) {
			locks = append(locks, countPinHeights(schema))
		} else {
			keys = append(keys, countPinHeights(schema))
		}
	}

	part1 := 0
	for _, lock := range locks {
		for _, key := range keys {
			if doesKeyFitLock(lock, key) {
				part1++
			}
		}
	}
	fmt.Println(part1)
}

func doesKeyFitLock(lock, key [5]int) bool {
	for i := 0; i < len(lock); i++ {
		if lock[i]+key[i] > 5 {
			return false
		}
	}
	return true
}

func countPinHeights(schema [7]string) [5]int {
	pinHeights := [5]int{}
	for j := 0; j < 5; j++ {
		for i := 1; i < 6; i++ {
			if schema[i][j] == '#' {
				pinHeights[j]++
			}
		}
	}
	return pinHeights
}

func isLock(schema [7]string) bool {
	return schema[0][0] == '#'
}
