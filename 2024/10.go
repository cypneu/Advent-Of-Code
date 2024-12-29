package main

import (
	"bufio"
	"fmt"
	"strings"
)

type ScoringStrategy func(topoMap []string, startPos Position, dirs []Position, n, m int) int

func Solution10A(input string) {
	solve10(input, trailheadScore())
}

func Solution10B(input string) {
	solve10(input, trailheadRating())
}

func solve10(input string, scoringStrategy ScoringStrategy) {
	var startPositions []Position
	var topoMap []string

	scanner := bufio.NewScanner(strings.NewReader(input))
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		for j, c := range line {
			if c == '0' {
				startPositions = append(startPositions, Position{i, j})
			}
		}
		topoMap = append(topoMap, line)
		i++
	}

	n, m := len(topoMap), len(topoMap[0])
	dirs := []Position{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	totalScore := 0
	for _, startPos := range startPositions {
		totalScore += scoringStrategy(topoMap, startPos, dirs, n, m)
	}

	fmt.Println(totalScore)
}

func trailheadScore() ScoringStrategy {
	dp := make(map[Position]map[Position]struct{})

	var countTrailheadScore func([]string, Position, []Position, int, int) map[Position]struct{}
	countTrailheadScore = func(topoMap []string, startPos Position, dirs []Position, n, m int) map[Position]struct{} {
		if _, ok := dp[startPos]; ok {
			return dp[startPos]
		}

		if topoMap[startPos.row][startPos.col] == '9' {
			dp[startPos] = map[Position]struct{}{startPos: {}}
			return dp[startPos]
		}

		uniquePositions := make(map[Position]struct{})
		for _, dir := range dirs {
			newPos := Position{startPos.row + dir.row, startPos.col + dir.col}
			if inBounds(newPos, n, m) && topoMap[startPos.row][startPos.col]+1 == topoMap[newPos.row][newPos.col] {
				res := countTrailheadScore(topoMap, newPos, dirs, n, m)
				for key := range res {
					uniquePositions[key] = struct{}{}
				}
			}
		}

		dp[startPos] = uniquePositions
		return uniquePositions
	}

	return func(topoMap []string, startPos Position, dirs []Position, n, m int) int {
		return len(countTrailheadScore(topoMap, startPos, dirs, n, m))
	}
}

func trailheadRating() ScoringStrategy {
	dp := make(map[Position]int)

	var countTrailheadRating func([]string, Position, []Position, int, int) int
	countTrailheadRating = func(topoMap []string, startPos Position, dirs []Position, n, m int) int {
		if score, ok := dp[startPos]; ok {
			return score
		}

		if topoMap[startPos.row][startPos.col] == '9' {
			dp[startPos] = 1
			return 1
		}

		score := 0
		for _, dir := range dirs {
			newPos := Position{startPos.row + dir.row, startPos.col + dir.col}
			if inBounds(newPos, n, m) && topoMap[startPos.row][startPos.col]+1 == topoMap[newPos.row][newPos.col] {
				score += countTrailheadRating(topoMap, newPos, dirs, n, m)
			}
		}
		dp[startPos] = score
		return score
	}

	return func(topoMap []string, startPos Position, dirs []Position, n, m int) int {
		return countTrailheadRating(topoMap, startPos, dirs, n, m)
	}
}

func inBounds(newPos Position, n, m int) bool {
	return newPos.row >= 0 && newPos.row < n && newPos.col >= 0 && newPos.col < m
}
