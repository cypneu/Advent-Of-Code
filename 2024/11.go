package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func Solution11A(input string) {
	stonesStr := strings.Fields(input)
	stones := make([]int, len(stonesStr))
	for i, stone := range stonesStr {
		stone, _ := strconv.Atoi(stone)
		stones[i] = stone
	}

	for i := 0; i < 25; i++ {
		var newStones []int
		for j := 0; j < len(stones); j++ {
			if stones[j] == 0 {
				newStones = append(newStones, 1)
			} else if digits(stones[j])%2 == 0 {
				div := int(math.Pow(10, float64(digits(stones[j])/2)))
				leftPart := stones[j] / div
				newStones = append(newStones, leftPart)
				newStones = append(newStones, stones[j]%(leftPart*div))
			} else {
				newStones = append(newStones, stones[j]*2024)
			}
		}
		stones = newStones
	}

	fmt.Println(len(stones))
}

func Solution11B(input string) {
	stonesStr := strings.Fields(input)
	stones := make([]int, len(stonesStr))
	for i, stone := range stonesStr {
		stone, _ := strconv.Atoi(stone)
		stones[i] = stone
	}

	dp := make(map[struct{ num, iterLeft int }]int)
	total := 0
	for _, stone := range stones {
		total += solve11(stone, dp, 75)
	}
	fmt.Println(total)
}

func solve11(num int, dp map[struct{ num, iterLeft int }]int, iterLeft int) int {
	if res, ok := dp[struct{ num, iterLeft int }{num, iterLeft}]; ok {
		return res
	}

	if iterLeft == 0 {
		return 1
	}

	var res int
	if num == 0 {
		res = solve11(1, dp, iterLeft-1)
	} else {
		digitsNum := digits(num)
		if digitsNum%2 == 0 {
			div := int(math.Pow(10, float64(digitsNum/2)))
			leftPart := num / div
			res = solve11(leftPart, dp, iterLeft-1) + solve11(num%(leftPart*div), dp, iterLeft-1)
		} else {
			res = solve11(num*2024, dp, iterLeft-1)
		}
	}

	dp[struct{ num, iterLeft int }{num, iterLeft}] = res
	return res
}
