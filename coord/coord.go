// A 9x9 sudoku coordinate system with various iterations and iterator compositions
//
// Example:
//
// iterating a row that contains a certain coordinate:
//
// l := coord.Coord{X: 5, Y: 3}
// c := coord.Row(l)
//
//	for c.Next() {
//	  fmt.Print(c.Value())
//	}
//
// Example 2:
//
// iterating all cells by rows
//
// c := coord.AllRows()
//
//	for c.Next() {
//	  r := c.Value() // row iterator
//	  for r.Next() {
//	    fmt.Print(r.Value())
//	  }
//	}
//
// Example 3:
//
// # Composing 2 iterators into an iterator that calls the constituents in order
//
// l := coord.Coord{X: 5, Y: 3}
// c := coord.Composed(coord.Row(l), coord.Column(l))
//
//	for c.Next() {
//	  fmt.Print(c.Value())
//	}
package coord

type dim int8
type Coord struct {
	X, Y dim // X,Y coordinates on a sudoku board
}

// coordinate to integer
func Ctoi(c Coord) int {
	return int(c.Y*9 + c.X)
}

// composed iterator iterating first a then b
func Composed(a, b Iterator) Iterator { return &composed{a: a, b: b} }

// iterates all coordinates row by row
func All() *allIterator { return &allIterator{i: -1} }

// iterating same row as c
func Row(c Coord) *rowIterator { return &rowIterator{base: c, i: -1} }

// iterating same column as c
func Column(c Coord) *columnIterator { return &columnIterator{base: c, i: -1} }

// coordinates for the cells in the same 3x3 box
func Box(c Coord) *boxIterator {
	i := boxIterator{base: c, i: -1}

	sx := i.base.X - i.base.X%3
	sy := i.base.Y - i.base.Y%3
	n := 0
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			i.coords[n] = Coord{dim(x) + sx, dim(y) + sy}
			n++
		}
	}
	return &i
}

// iterator that yields row iterators, one for each column
func AllRows() *allRowsIterator { return &allRowsIterator{i: -1} }

// iterator that yields column iterators, one for each row
func AllColumns() *allColumnsIterator { return &allColumnsIterator{i: -1} }

// iterator that yields box iterators, one for each 3x3 box of sudoku
func AllBoxes() *allBoxesIterator { return &allBoxesIterator{i: -1} }

type any interface{}

// iterator
type Iterator interface {
	Next() bool // iterator Next
	Value() any // iterator Value
	Reset()     // reset iterator
}

type composed struct {
	a, b Iterator
	bRun bool
}

func (i *composed) Next() bool {
	if i.bRun {
		return i.b.Next()
	} else {
		i.bRun = !i.a.Next()
		if i.bRun {
			return i.b.Next()
		}
		return true
	}
}

func (i composed) Value() any {
	if i.bRun {
		return i.b.Value()
	} else {
		return i.a.Value()
	}
}

func (i *composed) Reset() {
	i.a.Reset()
	i.b.Reset()
	i.bRun = false
}

type allIterator struct {
	i dim
}

func (i *allIterator) Next() bool {
	i.i++
	return i.i < 81
}

func (i allIterator) Value() any {
	return Coord{i.i % 9, i.i / 9}
}

func (i *allIterator) Reset() {
	i.i = -1
}

type rowIterator struct {
	base Coord
	i    dim
}

func (i *rowIterator) Next() bool {
	i.i++
	return i.i < 9
}

func (i rowIterator) Value() any {
	return Coord{i.i, i.base.Y}
}

func (i *rowIterator) Reset() {
	i.i = -1
}

type columnIterator struct {
	base Coord
	i    dim
}

func (i *columnIterator) Next() bool {
	i.i++
	return i.i < 9
}

func (i columnIterator) Value() any {
	return Coord{i.base.X, i.i}
}

func (i *columnIterator) Reset() {
	i.i = -1
}

type boxIterator struct {
	base   Coord
	i      dim
	coords [9]Coord
}

func (i *boxIterator) Next() bool {
	i.i++
	return i.i < 9
}

func (i boxIterator) Value() any {
	return i.coords[i.i]
}

func (i *boxIterator) Reset() {
	i.i = -1
}

type allRowsIterator struct {
	i dim
}

func (i *allRowsIterator) Next() bool {
	i.i++
	return i.i < 9
}

func (i allRowsIterator) Value() any {
	return Row(Coord{0, i.i})
}

func (i *allRowsIterator) Reset() {
	i.i = -1
}

type allColumnsIterator struct {
	i dim
}

func (i *allColumnsIterator) Next() bool {
	i.i++
	return i.i < 9
}

func (i allColumnsIterator) Value() any {
	return Column(Coord{i.i, 0})
}

func (i *allColumnsIterator) Reset() {
	i.i = -1
}

type allBoxesIterator struct{ i dim }

func (i *allBoxesIterator) Next() bool {
	i.i++
	return i.i < 9
}

func (i allBoxesIterator) Value() any {
	bx, by := i.i%3, i.i/3
	return Box(Coord{bx * 3, by * 3})
}

func (i *allBoxesIterator) Reset() {
	i.i = -1
}
