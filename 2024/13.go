package main

import (
	"bufio"
	"fmt"
	"strings"
)

func Solution13A(input string) {
	totalCost := calculateTotalCost(input, 0, 0)
	fmt.Println(totalCost)
}

func Solution13B(input string) {
	const prizeOffset = 10000000000000
	totalCost := calculateTotalCost(input, prizeOffset, prizeOffset)
	fmt.Println(totalCost)
}

func calculateTotalCost(input string, xOffset, yOffset int) int {
	totalCost := 0
	costA, costB := 3, 1
	scanner := bufio.NewScanner(strings.NewReader(input))
	A, B, prize := [2]int{}, [2]int{}, [2]int{}

	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "Button A: X+%d, Y+%d", &A[0], &A[1])
		scanner.Scan()
		fmt.Sscanf(scanner.Text(), "Button B: X+%d, Y+%d", &B[0], &B[1])
		scanner.Scan()
		fmt.Sscanf(scanner.Text(), "Prize: X=%d, Y=%d", &prize[0], &prize[1])
		scanner.Scan()

		prize[0] += xOffset
		prize[1] += yOffset

		totalCost += solve13(A, B, prize, costA, costB)
	}
	return totalCost
}

func solve13(A, B, prize [2]int, costA, costB int) int {
	numeratorB := A[0]*prize[1] - A[1]*prize[0]
	denominatorB := B[1]*A[0] - A[1]*B[0]

	if numeratorB%denominatorB == 0 {
		b := numeratorB / denominatorB
		numeratorA := prize[0] - B[0]*b

		if numeratorA%A[0] == 0 {
			a := numeratorA / A[0]
			return costA*a + costB*b
		}
	}
	return 0
}
