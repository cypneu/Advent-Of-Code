package main

import (
	"AoC/2024/utils"
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"strings"
)

type Reindeer struct {
	Pos complex128
	Dir complex128
}

type PQItem = utils.Item[Reindeer]

type CompleteState struct {
	reindeer Reindeer
	score    int
}

func getOrMax(dict map[Reindeer]int, key Reindeer) int {
	if _, ok := dict[key]; !ok {
		return math.MaxInt64
	}
	return dict[key]
}

func Solution16(input string) {
	var grid []string
	var reindeer Reindeer
	var finishPos complex128
	i := 0

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		j := strings.Index(line, "S")
		if j >= 0 {
			reindeer = Reindeer{Pos: complex(float64(i), float64(j)), Dir: 1i}
		}
		k := strings.Index(line, "E")
		if k >= 0 {
			finishPos = complex(float64(i), float64(k))
		}
		grid = append(grid, line)
		i++
	}

	dist := make(map[Reindeer]int)
	dist[reindeer] = 0

	dirs := [3]complex128{1, 1i, -1i}
	cost := [3]int{1, 1001, 1001}

	completedStates := []CompleteState{}

	pq := utils.NewPriorityQueue[Reindeer]()
	heap.Init(&pq)
	heap.Push(&pq, &PQItem{Value: reindeer, Priority: []float64{0}})

	score := math.MaxInt64
	for pq.Len() > 0 {
		cur := heap.Pop(&pq).(*PQItem)
		reindeer, curDist := cur.Value, int(cur.Priority[0])

		if curDist > getOrMax(dist, reindeer) {
			continue
		}

		if reindeer.Pos == finishPos && curDist <= score {
			score = curDist
			completedStates = append(completedStates, CompleteState{reindeer, curDist})
		}

		for i, dir := range dirs {
			newDir := reindeer.Dir * dir
			nextReindeer := Reindeer{Pos: reindeer.Pos + newDir, Dir: newDir}
			if grid[int(real(nextReindeer.Pos))][int(imag(nextReindeer.Pos))] != '#' {
				if curDist+cost[i] < getOrMax(dist, nextReindeer) {
					dist[nextReindeer] = curDist + cost[i]
					heap.Push(&pq, &PQItem{Value: nextReindeer, Priority: []float64{float64(dist[nextReindeer])}})
				}
			}
		}
	}

	fmt.Println(score)
	visited := reconstructPath(completedStates, dist, dirs, cost)
	fmt.Println(len(visited))
}

func reconstructPath(queue []CompleteState, dist map[Reindeer]int, dirs [3]complex128, cost [3]int) map[complex128]struct{} {
	visited := map[complex128]struct{}{}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		visited[cur.reindeer.Pos] = struct{}{}

		for i, dir := range dirs {
			prevDir := cur.reindeer.Dir * dir
			prevReindeer := Reindeer{Pos: cur.reindeer.Pos - cur.reindeer.Dir, Dir: prevDir}
			if _, ok := visited[prevReindeer.Pos]; !ok && dist[prevReindeer] == cur.score-cost[i] {
				queue = append(queue, CompleteState{prevReindeer, cur.score - cost[i]})
			}
		}
	}
	return visited
}
