package main

import (
	"fmt"
	"regexp"
	"strconv"
)

func Solution3A(input string) {
	var total int

	re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	matches := re.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		num1, _ := strconv.Atoi(match[1])
		num2, _ := strconv.Atoi(match[2])
		total += num1 * num2
	}
	fmt.Println(total)
}

func Solution3B(input string) {
	var total int

	re := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)|do\(\)|don't\(\)`)
	matches := re.FindAllStringSubmatch(input, -1)

	do := true
	for _, match := range matches {
		if match[0] == "do()" {
			do = true
		} else if match[0] == "don't()" {
			do = false
		} else if do {
			num1, _ := strconv.Atoi(match[1])
			num2, _ := strconv.Atoi(match[2])
			total += num1 * num2
		}
	}
	fmt.Println(total)
}
