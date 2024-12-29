package main

import (
	"fmt"
	"strconv"
	"strings"
)

func Solution22(input string) {
	secretNumbers := strings.Split(input, "\n")

	part1, part2 := 0, 0
	counter := map[[4]int]int{}
	for _, num := range secretNumbers {
		secret, _ := strconv.Atoi(num)

		seen := map[[4]int]struct{}{}
		last4, delta := []int{}, 0
		for x := 0; x < 2000; x++ {
			secret, delta = nextSecretValue(secret)
			last4 = append(last4, delta)

			if len(last4) >= 4 {
				four := [4]int(last4[len(last4)-4:])
				if _, ok := seen[four]; !ok {
					counter[four] += (secret % 10)
					seen[four] = struct{}{}
				}
			}
		}
		part1 += secret
	}

	for _, v := range counter {
		part2 = max(part2, v)
	}

	fmt.Println(part1)
	fmt.Println(part2)
}

func nextSecretValue(secret int) (int, int) {
	prev := secret % 10
	secret = prune(mix(64*secret, secret))
	secret = prune(mix(secret/32, secret))
	secret = prune(mix(2048*secret, secret))
	return secret, (secret%10 - prev)
}

func mix(a, b int) int {
	return a ^ b
}

func prune(a int) int {
	return a % 16777216
}
