package main

import (
	"bufio"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func Solution1A(input string) {
	var left []int
	var right []int

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)
		left_num, _ := strconv.Atoi(words[0])
		right_num, _ := strconv.Atoi(words[1])
		left = append(left, left_num)
		right = append(right, right_num)
	}

	sort.Sort(sort.IntSlice(left))
	sort.Sort(sort.IntSlice(right))

	var total int
	for i := 0; i < len(left); i++ {
		if left[i] > right[i] {
			total += left[i] - right[i]
		} else {
			total += right[i] - left[i]
		}
	}

	fmt.Println(total)
}

func Solution1B(input string) {
	counterLeft := make(map[string]int)
	counterRight := make(map[string]int)

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)
		counterLeft[words[0]]++
		counterRight[words[1]]++
	}

	var similarityScore int
	for numStr, count := range counterLeft {
		num, _ := strconv.Atoi(numStr)
		similarityScore += num * count * counterRight[numStr]
	}
	fmt.Println(similarityScore)
}
