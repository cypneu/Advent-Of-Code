package utils

type UnionFind[T comparable] struct {
	Parent map[T]T
	Size   map[T]int
}

func NewUnionFind[T comparable]() UnionFind[T] {
	return UnionFind[T]{
		Parent: map[T]T{},
		Size:   map[T]int{},
	}
}

func (uf *UnionFind[T]) Find(x T) T {
	if uf.Parent[x] != x {
		uf.Parent[x] = uf.Find(uf.Parent[x])
	}
	return uf.Parent[x]
}

func (uf *UnionFind[T]) Union(x, y T) {
	rootX, rootY := uf.Find(x), uf.Find(y)

	if rootX != rootY {
		if uf.Size[rootX] > uf.Size[rootY] {
			uf.Parent[rootY] = rootX
			uf.Size[rootX] += uf.Size[rootY]
		} else {
			uf.Parent[rootX] = rootY
			uf.Size[rootY] += uf.Size[rootX]
		}
	}
}

func (uf *UnionFind[T]) Add(x T) {
	if _, exists := uf.Parent[x]; !exists {
		uf.Parent[x] = x
		uf.Size[x] = 1
	}
}
