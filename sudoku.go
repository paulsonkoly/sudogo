package main

import (
	"container/heap"
	"fmt"
	"math/bits"

	"github.com/phaul/sudoku/coord"
)

type cellVal uint8  // value of a cell, 0 empty, 1-9 otherwise
type cellCan uint16 // bitmap of what cell can be 0-8 bits used to indicate a cell can take ix+1 as value

type cell struct {
	val cellVal // value of the cell
	can cellCan // possibilities for the cell
}

type board [9 * 9]cell // a sudoku board

// address a board with x, y 0-8 coordinates. 0, 0 is the top left corner and 8, 0 is the top right
// errors if coordinates are out of bounds
func (b *board) at(c coord.Coord) *cell {
	return &b[coord.Ctoi(c)]
}

// sets all cells to all possible
func (b *board) all_possible() {
	i := coord.All()

	for i.Next() {
		b.at(i.Value().(coord.Coord)).can = 0x1ff
	}
}

// fill a cell in the board at c with v
func (b *board) fill(c coord.Coord, v cellVal) {
	*b.at(c) = cell{val: v, can: 0}

	i := coord.Composed(coord.Composed(coord.Row(c), coord.Column(c)), coord.Box(c))

	for i.Next() {
		c = i.Value().(coord.Coord)
		*b.at(c) = cell{b.at(c).val, b.at(c).can & (^(1 << (v - 1)))}
	}
}

// look for a cell that is single possible and fill
// return true if any were found or false otherwise
func (b *board) single_possible() bool {
	r := false
	i := coord.All()

	for i.Next() {
		co := i.Value().(coord.Coord)
		c := b.at(co)

		if c.can != 0 && c.can&(c.can-1) == 0 {
			b.fill(co, cellVal(bits.TrailingZeros16(uint16(c.can))+1))
			r = true
		}
	}
	return r
}

// finds a digit that can only go in one place, and fills it in
// returns true if one found
func (b *board) only_place() bool {
	i := coord.Composed(coord.Composed(coord.AllRows(), coord.AllColumns()), coord.AllBoxes())

	for i.Next() {
		r := i.Value().(coord.Iterator)
		counts := [9]int{}

		for r.Next() {
			can := b.at(r.Value().(coord.Coord)).can
			for j := 0; j < 9; j++ {
				if can&(1<<j) != 0 {
					counts[j] += 1
				}
			}
		}
		r.Reset()
		for r.Next() {
			co := r.Value().(coord.Coord)
			for j := 0; j < 9; j++ {
				if b.at(co).can&(1<<j) != 0 && counts[j] == 1 {
					b.fill(co, cellVal(j+1))
					return true
				}
			}
		}
	}

	// for i := 0; i < 9; i++ {
	// 	for _, c := range (row(coord{0, dim(i)})) {
	// 	}
	// 	for ix, cnt := range counts {
	// 		if cnt == 1 {
	// 			for _, c := range (row(coord{0, dim(i)})) {
	// 					b.fill(c, cellVal(ix+1))
	// 					return true
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	//
	// for i := 0; i < 9; i++ {
	// 	counts := [9]int{}
	// 	for _, c := range (column(coord{dim(i), 0})) {
	// 		can := b.at(c).can
	// 		for j := 0; j < 9; j++ {
	// 			if can&(1<<j) != 0 {
	// 				counts[j] += 1
	// 			}
	// 		}
	// 	}
	// 	for ix, cnt := range counts {
	// 		if cnt == 1 {
	// 			for _, c := range (column(coord{dim(i), 0})) {
	// 				if b.at(c).can&(1<<ix) != 0 {
	// 					b.fill(c, cellVal(ix+1))
	// 					return true
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	//
	// boxes := [...]coord{
	// 	{0, 0}, {3, 0}, {6, 0},
	// 	{0, 3}, {3, 3}, {6, 3},
	// 	{0, 6}, {3, 6}, {6, 6},
	// }
	//
	// for _, tl := range boxes {
	// 	counts := [9]int{}
	// 	for _, c := range box(tl) {
	// 		can := b.at(c).can
	// 		for j := 0; j < 9; j++ {
	// 			if can&(1<<j) != 0 {
	// 				counts[j] += 1
	// 			}
	// 		}
	// 	}
	// 	for ix, cnt := range counts {
	// 		if cnt == 1 {
	// 			for _, c := range box(tl) {
	// 				if b.at(c).can&(1<<ix) != 0 {
	// 					b.fill(c, cellVal(ix+1))
	// 					return true
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	return false
}

func (b *board) iterate() {
	for maxDepth := 3; true; maxDepth++ {
		if b.solve(0, maxDepth, max(maxDepth/3, 2)) {
			return
		}
	}
}

func (b *board) solve(depth, maxDepth, maxWidth int) bool {
	if depth >= maxDepth {
		return false
	}
	for b.single_possible() || b.only_place() {
	}
	if b.solved() {
		return true
	}
	if b.contradicts() {
		return false
	}
	return b.try(depth, maxDepth, maxWidth)
}

func (b *board) solved() bool {
	i := coord.All()

	for i.Next() {
		if b.at(i.Value().(coord.Coord)).val == 0 {
			return false
		}
	}
	return true
}

// priority queue for coordinates based on the amount of candidates
type prioCoord struct {
	count int
	coord coord.Coord
}

type queue []prioCoord

func (q queue) Len() int           { return len(q) }
func (q queue) Less(i, j int) bool { return q[i].count < q[j].count }
func (q queue) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func (q *queue) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*q = append(*q, x.(prioCoord))
}

func (q *queue) Pop() any {
	old := *q
	n := len(old)
	x := old[n-1]
	*q = old[0 : n-1]
	return x
}

// coordinates to try in the order of least amount of possible candidates to most
func (b *board) tries(maxWidth int) queue {
	q := make(queue, 0, 16)
	i := coord.All()

	for i.Next() {
		c := i.Value().(coord.Coord)
		cell := b.at(c)
		// fmt.Printf("%v %b %d - %d\n", c, cell.can, bits.OnesCount16(uint16(cell.can)), maxWidth)
		if cell.can != 0 && bits.OnesCount16(uint16(cell.can)) <= maxWidth {
			cnt := bits.OnesCount16(uint16(cell.can))
			heap.Push(&q, prioCoord{count: cnt, coord: c})
		}
	}

	return q
}

func (b *board) try(depth, maxDepth, maxWidth int) bool {
	// look for the lowest bitcount candidate
	for q := b.tries(maxWidth); q.Len() > 0; {
		c := heap.Pop(&q).(prioCoord).coord

		// for all candidates of the cell
		for can := b.at(c).can; can != 0; can &= can - 1 {
			bb := board{}
			copy(bb[:], b[:])

			v := cellVal(bits.TrailingZeros16(uint16(can&-can)) + 1)

			bb.fill(c, v)
			if bb.solve(depth+1, maxDepth, maxWidth) {
				copy(b[:], bb[:])
				return true
			}
		}
	}
	return false
}

// there is a cell that has no possible values but also not filled in
func (b *board) contradicts() bool {
	i := coord.All()

	for i.Next() {
		c := b.at(i.Value().(coord.Coord))

		if c.val == 0 && c.can == 0 {
			return true
		}
	}
	return false
}

func (b board) print() {
	i := coord.All()

	for i.Next() {
		c := i.Value().(coord.Coord)
		if c.Y%3 == 0 && c.X == 0 {
			fmt.Println("+---+---+---")
		}
		if c.X%3 == 0 {
			fmt.Print("|")
		}
		if b.at(c).val == 0 {
			fmt.Print(" ")
		} else {
			fmt.Print(b.at(c).val)
		}
		if c.X == 8 {
			fmt.Println("|")
		}
	}
}

func main() {
	b := board{}
	b.all_possible()
	b.fill(coord.Coord{X: 0, Y: 0}, 8)
	b.fill(coord.Coord{X: 2, Y: 1}, 3)
	b.fill(coord.Coord{X: 3, Y: 1}, 6)
	b.fill(coord.Coord{X: 1, Y: 2}, 7)
	b.fill(coord.Coord{X: 4, Y: 2}, 9)
	b.fill(coord.Coord{X: 6, Y: 2}, 2)
	b.fill(coord.Coord{X: 1, Y: 3}, 5)
	b.fill(coord.Coord{X: 5, Y: 3}, 7)
	b.fill(coord.Coord{X: 4, Y: 4}, 4)
	b.fill(coord.Coord{X: 5, Y: 4}, 5)
	b.fill(coord.Coord{X: 6, Y: 4}, 7)
	b.fill(coord.Coord{X: 3, Y: 5}, 1)
	b.fill(coord.Coord{X: 7, Y: 5}, 3)
	b.fill(coord.Coord{X: 2, Y: 6}, 1)
	b.fill(coord.Coord{X: 7, Y: 6}, 6)
	b.fill(coord.Coord{X: 8, Y: 6}, 8)
	b.fill(coord.Coord{X: 2, Y: 7}, 8)
	b.fill(coord.Coord{X: 3, Y: 7}, 5)
	b.fill(coord.Coord{X: 7, Y: 7}, 1)
	b.fill(coord.Coord{X: 1, Y: 8}, 9)
	b.fill(coord.Coord{X: 6, Y: 8}, 4)

	// b.fill(coord{2, 0}, 1)
	// b.fill(coord{5, 0}, 2)
	// b.fill(coord{8, 0}, 4)
	// b.fill(coord{2, 1}, 7)
	// b.fill(coord{3, 1}, 5)
	// b.fill(coord{5, 1}, 9)
	// b.fill(coord{6, 1}, 6)
	// b.fill(coord{0, 2}, 4)
	// b.fill(coord{3, 2}, 8)
	// b.fill(coord{4, 2}, 3)
	// b.fill(coord{7, 2}, 5)
	// b.fill(coord{8, 2}, 7)
	// b.fill(coord{0, 3}, 9)
	// b.fill(coord{1, 3}, 4)
	// b.fill(coord{5, 3}, 7)
	// b.fill(coord{7, 3}, 3)
	// b.fill(coord{8, 3}, 2)
	// b.fill(coord{3, 4}, 3)
	// b.fill(coord{4, 4}, 9)
	// b.fill(coord{5, 4}, 6)
	// b.fill(coord{8, 4}, 5)
	// b.fill(coord{1, 5}, 7)
	// b.fill(coord{2, 5}, 3)
	// b.fill(coord{4, 5}, 8)
	// b.fill(coord{6, 5}, 1)
	// b.fill(coord{0, 6}, 7)
	// b.fill(coord{1, 6}, 3)
	// b.fill(coord{2, 6}, 4)
	// b.fill(coord{0, 7}, 8)
	// b.fill(coord{3, 7}, 7)
	// b.fill(coord{6, 7}, 4)
	// b.fill(coord{7, 7}, 2)
	// b.fill(coord{8, 7}, 9)
	// b.fill(coord{1, 8}, 9)
	// b.fill(coord{2, 8}, 2)
	// b.fill(coord{3, 8}, 4)
	// b.fill(coord{5, 8}, 5)
	// b.fill(coord{6, 8}, 3)
	b.print()
	fmt.Println("=========================")
	b.iterate()
	b.print()
}
