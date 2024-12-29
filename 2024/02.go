package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func Solution2A(input string) {
	var safe int

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)

		levels := make([]int, len(words))

		for i, numStr := range words {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				fmt.Println("Error converting string to int:", err)
				return
			}
			levels[i] = num
		}

		isSafe1 := true
		isSafe2 := true
		for i := 1; i < len(levels); i++ {
			if isSafe1 && !(levels[i-1]-levels[i] > 0 && levels[i-1]-levels[i] <= 3) {
				isSafe1 = false
			}

			if isSafe2 && !(levels[i]-levels[i-1] > 0 && levels[i]-levels[i-1] <= 3) {
				isSafe2 = false
			}
		}

		if isSafe1 || isSafe2 {
			safe++
		}
	}

	fmt.Println(safe)
}

func Solution2B(input string) {
	var safe int

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)

		levels := make([]int, len(words))

		for i, numStr := range words {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				fmt.Println("Error converting string to int:", err)
				return
			}
			levels[i] = num
		}

		checkedInc, checkedDec := false, false
		safeInc, safeDec := true, true
		for i := 1; i < len(levels); i++ {
			diffInc := levels[i] - levels[i-1]
			if !checkedInc && !(diffInc > 0 && diffInc <= 3) {
				safeInc = checkSafety(levels, i-1, isIncreasing) || checkSafety(levels, i, isIncreasing)
				checkedInc = true
			}

			diffDesc := levels[i-1] - levels[i]
			if !checkedDec && !(diffDesc > 0 && diffDesc <= 3) {
				safeDec = checkSafety(levels, i-1, isDecreasing) || checkSafety(levels, i, isDecreasing)
				checkedDec = true
			}

			if checkedInc && checkedDec {
				break
			}
		}

		if safeInc || safeDec {
			safe++
		}
	}

	fmt.Println(safe)
}

func checkSafety(levels []int, i int, check func(int, int) bool) bool {
	prev := -1
	for j := 0; j < len(levels); j++ {
		if j != i {
			if prev >= 0 && !check(levels[prev], levels[j]) {
				return false
			}
			prev = j
		}
	}
	return true
}

func isIncreasing(prev, curr int) bool {
	diff := curr - prev
	return diff > 0 && diff <= 3
}

func isDecreasing(prev, curr int) bool {
	diff := prev - curr
	return diff > 0 && diff <= 3
}
