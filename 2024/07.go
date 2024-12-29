package main

import (
	"bufio"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func Solution7A(input string) {
	fmt.Println(solve(input, false))
}

func Solution7B(input string) {
	fmt.Println(solve(input, true))
}

func solve(input string, checkConcat bool) int {
	var total int

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		res, nums := getResultAndNumbers(line)

		if matchTestValue(len(nums)-1, res, nums, checkConcat) {
			total += res
		}
	}
	return total
}

func matchTestValue(i, target int, nums []int, checkConcat bool) bool {
	if i == 0 && target == nums[i] {
		return true
	}

	if target < 0 || i == 0 {
		return false
	}

	if target%nums[i] == 0 && matchTestValue(i-1, target/nums[i], nums, checkConcat) {
		return true
	}

	if checkConcat {
		x := int(math.Pow(10, float64(digits(nums[i]))))
		if (target-nums[i])%x == 0 && matchTestValue(i-1, target/x, nums, checkConcat) {
			return true
		}
	}

	return matchTestValue(i-1, target-nums[i], nums, checkConcat)
}

func digits(a int) int {
	return int(math.Log10(float64(a))) + 1
}

func getResultAndNumbers(line string) (int, []int) {
	split := strings.Split(line, ": ")
	res, _ := strconv.Atoi(split[0])

	numsStr := strings.Split(split[1], " ")
	nums := make([]int, len(numsStr))
	for i, numStr := range numsStr {
		nums[i], _ = strconv.Atoi(numStr)
	}
	return res, nums
}
