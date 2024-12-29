package main

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func Solution17(input string) {
	lines := strings.Split(input, "\n")
	registers := [3]int{}

	fmt.Sscanf(lines[0], "Register A: %d", &registers[0])
	fmt.Sscanf(lines[1], "Register B: %d", &registers[1])
	fmt.Sscanf(lines[2], "Register C: %d", &registers[2])

	programStr := strings.Split(strings.Split(lines[4], " ")[1], ",")
	program := make([]int, len(programStr))
	for i, str := range programStr {
		num, _ := strconv.Atoi(str)
		program[i] = num
	}

	result := eval(program, registers)
	resStr := make([]string, len(result))
	for i, num := range result {
		resStr[i] = strconv.Itoa(num) // Convert int to string
	}
	fmt.Println(strings.Join(resStr, ","))

	res := math.MaxInt64
	solve17Manual(program, len(program)-1, 0, &res)
	fmt.Println(res)

	registers[0] = 0
	res = math.MaxInt64
	solve17(program, registers, 0, &res)
	fmt.Println(res)
}

func solve17(program []int, registers [3]int, i int, sol *int) {
	res := eval(program, registers)
	if reflect.DeepEqual(res, program) {
		*sol = min(*sol, registers[0])
	} else if i == 0 || reflect.DeepEqual(res, program[len(program)-i:]) {
		for x := 0; x < 8; x++ {
			solve17(program, [3]int{8*registers[0] + x, registers[1], registers[2]}, i+1, sol)
		}
	}
}

func solve17Manual(program []int, i, A int, res *int) {
	if i == -1 {
		*res = min(*res, A/8)
		return
	}

	for x := 0; x < 8; x++ {
		newA := A + x
		B := newA % 8
		B = B ^ 1
		C := newA >> B
		B = B ^ C
		B = B ^ 4
		if B%8 == program[i] {
			solve17Manual(program, i-1, newA*8, res)
		}
	}
}

func eval(program []int, registers [3]int) []int {
	print := []int{}
	ip := 0
	var moveIp bool
	for ip < len(program) {
		opcode := program[ip]
		switch opcode {
		case 0:
			moveIp = advInstruction(program[ip+1], &registers)
		case 1:
			moveIp = bxlInstruction(program[ip+1], &registers)
		case 2:
			moveIp = bstInstruction(program[ip+1], &registers)
		case 3:
			moveIp = jnzInstruction(program[ip+1], &registers, &ip)
		case 4:
			moveIp = bxcInstruction(program[ip+1], &registers)
		case 5:
			moveIp = outInstruction(program[ip+1], &registers, &print)
		case 6:
			moveIp = bdvInstruction(program[ip+1], &registers)
		case 7:
			moveIp = cdvInstruction(program[ip+1], &registers)
		}
		if moveIp {
			ip += 2
		}
	}
	return print
}

func literalOperand(operand int) int {
	return operand
}

func comboOperand(operand int, registers *[3]int) int {
	if operand <= 3 {
		return operand
	} else {
		return registers[operand-4]
	}
}

func advInstruction(operand int, registers *[3]int) bool {
	registers[0] = registers[0] / int(math.Pow(2, float64(comboOperand(operand, registers))))
	return true
}

func bxlInstruction(operand int, registers *[3]int) bool {
	registers[1] = registers[1] ^ literalOperand(operand)
	return true
}

func bstInstruction(operand int, registers *[3]int) bool {
	registers[1] = comboOperand(operand, registers) % 8
	return true
}

func jnzInstruction(operand int, registers *[3]int, ip *int) bool {
	if registers[0] != 0 {
		*ip = literalOperand(operand)
		return false
	}
	return true
}

func bxcInstruction(operand int, registers *[3]int) bool {
	registers[1] = registers[1] ^ registers[2]
	return true
}

func outInstruction(operand int, registers *[3]int, print *[]int) bool {
	*print = append(*print, comboOperand(operand, registers)%8)
	return true
}

func bdvInstruction(operand int, registers *[3]int) bool {
	registers[1] = registers[0] / int(math.Pow(2, float64(comboOperand(operand, registers))))
	return true
}

func cdvInstruction(operand int, registers *[3]int) bool {
	registers[2] = registers[0] / int(math.Pow(2, float64(comboOperand(operand, registers))))
	return true
}
