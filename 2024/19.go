package main

import (
	"bufio"
	"fmt"
	"strings"
)

func Solution19(input string) {
	scanner := bufio.NewScanner(strings.NewReader(input))

	var patterns, designs []string
	for scanner.Scan() {
		if len(patterns) == 0 {
			patterns = strings.Split(scanner.Text(), ", ")
			scanner.Scan()
		} else {
			designs = append(designs, scanner.Text())
		}
	}

	var ways func(string) int
	cache := map[string]int{}

	ways = func(design string) (n int) {
		if n, ok := cache[design]; ok {
			return n
		}

		if design == "" {
			return 1
		}

		for _, pattern := range patterns {
			if strings.HasPrefix(design, pattern) {
				n += ways(design[len(pattern):])
			}
		}

		cache[design] = n
		return n
	}

	part1, part2 := 0, 0
	for _, design := range designs {
		if res := ways(design); res > 0 {
			part1++
			part2 += res
		}
	}

	fmt.Println(part1)
	fmt.Println(part2)
}
