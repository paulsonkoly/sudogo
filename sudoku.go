package main

import (
	"container/heap"
	"fmt"

	"github.com/phaul/sudoku/cell"
	"github.com/phaul/sudoku/coord"
	"github.com/phaul/sudoku/cqueue"
)

type board [9 * 9]cell.Cell // a sudoku board

// address a board with x, y 0-8 coordinates. 0, 0 is the top left corner and 8, 0 is the top right
func (b *board) at(c coord.Coord) *cell.Cell {
	return &b[coord.Ctoi(c)]
}

// sets all cells to all 9 digits are possible
func (b *board) allPossible() {
	i := coord.All()

	for i.Next() {
		b.at(i.Value().(coord.Coord)).SetAll()
	}
}

// fill a cell in the board at c with v
func (b *board) fill(c coord.Coord, v cell.ValT) {
	*b.at(c) = cell.New(v)

	i := coord.Composed(coord.Composed(coord.Row(c), coord.Column(c)), coord.Box(c))

	for i.Next() {
		c = i.Value().(coord.Coord)
		b.at(c).Drop(v)
	}
}

// look for a cell that has a single possibility and fill
//
// return true if any were found or false otherwise
func (b *board) singlePossible() bool {
	r := false
	i := coord.All()

	for i.Next() {
		co := i.Value().(coord.Coord)
		c := b.at(co)

		if c.IsSingle() {
			b.fill(co, c.FirstPossibility())
			r = true
		}
	}
	return r
}

// find a digit that can only go in one place, and fill it in
//
// returns true if one found
func (b *board) onlyPlace() bool {
	i := coord.Composed(coord.Composed(coord.AllRows(), coord.AllColumns()), coord.AllBoxes())

	for i.Next() {
		r := i.Value().(coord.Iterator)
		counts := [9]int{}

		for r.Next() {
			c := b.at(r.Value().(coord.Coord))
			for j := 1; j <= 9; j++ {
				if c.IsPossible(cell.ValT(j)) {
					counts[j-1] += 1
				}
			}
		}
		r.Reset()
		for r.Next() {
			co := r.Value().(coord.Coord)
			for j := 1; j <= 9; j++ {
				if b.at(co).IsPossible(cell.ValT(j)) && counts[j-1] == 1 {
					b.fill(co, cell.ValT(j))
					return true
				}
			}
		}
	}

	return false
}

// wrapper for solving with iterative deepening
// tune constants here for performance
// maxDepth limits the number of guesses allowed before solve returns with false
// maxWidth limits where guesses can happen, don't guess a cell if it has more possiblities than maxWidth
func (b *board) iterate() {
	for maxDepth := 3; true; maxDepth++ {
		if b.solve(0, maxDepth, max(maxDepth/3, 2)) {
			return
		}
	}
}

// tries to do a solve
// first it fills in what we know for sure
// then checks if solved or has a contradiction due to incorrect guess
// then tries the easiest guess
func (b *board) solve(depth, maxDepth, maxWidth int) bool {
	// fmt.Printf("%d / %d\n", depth, maxDepth)
	if depth >= maxDepth {
		return false
	}
	for b.singlePossible() || b.onlyPlace() {
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
		if b.at(i.Value().(coord.Coord)).IsEmpty() {
			return false
		}
	}
	return true
}


// coordinates to try in the order of least amount of possible candidates to most
func (b *board) tries(maxWidth int) cqueue.Queue {
  q := cqueue.New()
	i := coord.All()

	for i.Next() {
		c := i.Value().(coord.Coord)
		cell := b.at(c)
		p := cell.PossibilityCount()
		if 0 < p && p <= maxWidth {
			heap.Push(&q, cqueue.PrioCoord{Count: p, Coord: c})
		}
	}

	return q
}

func (b *board) try(depth, maxDepth, maxWidth int) bool {
	// look for the lowest bitcount candidate
	for q := b.tries(maxWidth); q.Len() > 0; {
		c := heap.Pop(&q).(cqueue.PrioCoord).Coord
		i := b.at(c).Possibilities()

		// for all candidates for the cell
		for i.Next() {
			v := i.Value()
			bb := board{}
			copy(bb[:], b[:])

			bb.fill(c, v)
			if bb.solve(depth+1, maxDepth, maxWidth) {
				copy(b[:], bb[:])
				return true
			}
		}
	}
	return false
}

// there is a cell that has no possible value left but also not filled in
func (b *board) contradicts() bool {
	i := coord.All()

	for i.Next() {
		c := b.at(i.Value().(coord.Coord))

		if c.Value == 0 && c.PossibilityCount() == 0 {
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
		if b.at(c).Value == 0 {
			fmt.Print(" ")
		} else {
			fmt.Print(b.at(c).Value)
		}
		if c.X == 8 {
			fmt.Println("|")
		}
	}
}

func main() {
	b := board{}
	b.allPossible()
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

	b.print()
	fmt.Println("=========================")
	b.iterate()
	b.print()
}
