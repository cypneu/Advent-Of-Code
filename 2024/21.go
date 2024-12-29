package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

type CacheKey struct {
	iter int
	code string
}

func Solution21(input string) {
	codes := []string{}
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		codes = append(codes, scanner.Text())
	}

	numericKeypad := [][]byte{{'7', '8', '9'}, {'4', '5', '6'}, {'1', '2', '3'}, {0, '0', 'A'}}
	directionalKeypad := [][]byte{{0, '^', 'A'}, {'<', 'v', '>'}}

	numericTransitions := map[byte]map[byte]string{}
	fillTransitions(numericTransitions, numericKeypad)
	numericTransitions['7']['A'] = ">>vvv"
	numericTransitions['4']['A'] = ">>vv"
	numericTransitions['1']['A'] = ">>v"
	numericTransitions['7']['0'] = ">vvv"
	numericTransitions['4']['0'] = ">vv"
	numericTransitions['1']['0'] = ">v"

	numericTransitions['A']['7'] = "^^^<<"
	numericTransitions['A']['4'] = "^^<<"
	numericTransitions['A']['1'] = "^<<"
	numericTransitions['0']['7'] = "^^^<"
	numericTransitions['0']['4'] = "^^<"
	numericTransitions['0']['1'] = "^<"

	directionalTransitions := map[byte]map[byte]string{}
	fillTransitions(directionalTransitions, directionalKeypad)
	directionalTransitions['<']['A'] = ">>^"
	directionalTransitions['<']['^'] = ">^"
	directionalTransitions['A']['<'] = "v<<"
	directionalTransitions['^']['<'] = "v<"

	cache := map[CacheKey]int{}
	var countMoveSequence func(iter, maxIter int, code string, numericTransitions, directionalTransitions map[byte]map[byte]string) int
	countMoveSequence = func(iter, maxIter int, code string, numericTransitions, directionalTransitions map[byte]map[byte]string) (n int) {
		key := CacheKey{iter, code}
		if val, ok := cache[key]; ok {
			return val
		}

		if iter == maxIter+1 {
			return len(code)
		}

		transitions := directionalTransitions
		if iter == 0 {
			transitions = numericTransitions
		}

		var prev byte = 'A'
		for _, char := range code {
			nextCode := transitions[prev][byte(char)] + "A"
			n += countMoveSequence(iter+1, maxIter, nextCode, numericTransitions, directionalTransitions)
			prev = byte(char)
		}

		cache[key] = n
		return n
	}

	complexities := 0
	for _, code := range codes {
		minSeq := countMoveSequence(0, 25, code, numericTransitions, directionalTransitions)
		codeNum, _ := strconv.Atoi(code[:len(code)-1])
		complexities += codeNum * minSeq
	}
	fmt.Println(complexities)
}

func fillTransitions(transitions map[byte]map[byte]string, keypad [][]byte) {
	for i := 0; i < len(keypad); i++ {
		for j := 0; j < len(keypad[i]); j++ {
			if _, ok := transitions[keypad[i][j]]; !ok {
				transitions[keypad[i][j]] = map[byte]string{}
			}
			for k := 0; k < len(keypad); k++ {
				for l := 0; l < len(keypad[k]); l++ {
					transitions[keypad[i][j]][keypad[k][l]] = generateMoves(l-j, k-i)
				}
			}
		}
	}
}

func generateMoves(diffX, diffY int) string {
	moveX, moveY := ">", "v"
	if diffX < 0 {
		diffX, moveX = -diffX, "<"
	}
	if diffY < 0 {
		diffY, moveY = -diffY, "^"
	}

	horizontalMoves := strings.Repeat(moveX, diffX)
	verticalMoves := strings.Repeat(moveY, diffY)
	if diffX != 0 && diffY != 0 && moveX == ">" {
		return verticalMoves + horizontalMoves
	}
	return horizontalMoves + verticalMoves
}
