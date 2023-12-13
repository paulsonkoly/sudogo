package cell

import "math/bits"

type ValT uint8  // value of a cell, 0 empty, 1-9 otherwise
type canT uint16 // bitmap of what cell can be 0-8 bits used to indicate a cell can take ix+1 as value

// everything is possible
const everything = canT(0x1ff)

// nothing is possible
const none = canT(0)

// empty cell
const empty = ValT(0)

// a pair of values, holding a digit 1-9 or 0 indicating unsolved cell
// and a bitmask that is set '1' for each possible digit for the cell
type Cell struct {
	Value ValT // value of the cell
	can   canT // possibilities for the cell
}

type possibilityIterator struct {
	can  canT
	init bool
}

// a cell with Value v and Possibilites 0
func New(v ValT) Cell { return Cell{Value: v} }

// is the cell empty? (Val: 0)
func (c Cell) IsEmpty() bool { return c.Value == empty }

// calls f with all the possibilities of the cell
func (c Cell) Possibilities() possibilityIterator {
	return possibilityIterator{can: c.can}
}

// is there a next possibility
func (p *possibilityIterator) Next() bool {
	if p.init {
		p.can &= p.can - 1
		return p.can != 0
	} else {
		p.init = true
		return true
	}
}

// value yielded by the iterator
func (p possibilityIterator) Value() ValT {
	return ValT(bits.TrailingZeros16(uint16(p.can)) + 1)
}

// set all digits possible in the cell
func (c *Cell) SetAll() { c.can = everything }

// drops v as a possibility
func (c *Cell) Drop(v ValT) { c.can &= (^(1 << (v - 1))) }

// does the cell hold a single possibility?
func (c Cell) IsSingle() bool {
	return c.can != none && c.can&(c.can-1) == none
}

// The first possible value for the cell
func (c Cell) FirstPossibility() ValT { return ValT(bits.TrailingZeros16(uint16(c.can)) + 1) }

// Is v possible in the cell c
func (c Cell) IsPossible(v ValT) bool { return c.can&(1<<(v-1)) != none }

// count the possible digits for the cell
func (c Cell) PossibilityCount() int { return bits.OnesCount16(uint16(c.can)) }
