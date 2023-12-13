// priority queue for cells
package cqueue

import "github.com/phaul/sudoku/coord"

// priority queue for coordinates based on the amount of candidates
type PrioCoord struct {
	Count int
	Coord coord.Coord
}

type Queue []PrioCoord

func New() Queue { return make(Queue, 0, 16) }

func (q Queue) Len() int           { return len(q) }
func (q Queue) Less(i, j int) bool { return q[i].Count < q[j].Count }
func (q Queue) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func (q *Queue) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*q = append(*q, x.(PrioCoord))
}

func (q *Queue) Pop() any {
	old := *q
	n := len(old)
	x := old[n-1]
	*q = old[0 : n-1]
	return x
}
