package utils

type Item[T any] struct {
	Value    T
	Priority []float64
	index    int
}

type PriorityQueue[T any] []*Item[T]

func (pq PriorityQueue[T]) Len() int { return len(pq) }

func (pq PriorityQueue[T]) Less(i, j int) bool {
	p1 := pq[i].Priority
	p2 := pq[j].Priority

	minLen := len(p1)
	if len(p2) < minLen {
		minLen = len(p2)
	}

	for k := 0; k < minLen; k++ {
		if p1[k] > p2[k] {
			return true
		} else if p1[k] < p2[k] {
			return false
		}
	}

	return len(p1) > len(p2)
}

func (pq PriorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue[T]) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item[T])
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[T]) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func NewPriorityQueue[T any]() PriorityQueue[T] {
	return make(PriorityQueue[T], 0)
}
