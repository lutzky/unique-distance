package board

import (
	"bytes"
	"fmt"
	"io"
)

func main() {
	fmt.Println("vim-go")
}

// Board represents a board with markers
type Board struct {
	Markers []Coord
	Size    int
	ID      int64
}

// Coord represents a possible coordinate of a marker on a board
type Coord struct {
	X, Y int
}

// MaxDistance is the maximal possible distance between two markers on b
func (b *Board) MaxDistance() int {
	return 2 * (b.Size - 1) * (b.Size - 1)
}

func (b *Board) String() string {
	var buf bytes.Buffer
	b.Print(&buf)
	return buf.String()
}

// Print formats b (with marker distances) into w
func (b *Board) Print(w io.Writer) {
	if w == nil {
		return
	}
	rows := make([][]int, b.Size)

	for i := 0; i < b.Size; i++ {
		rows[i] = make([]int, b.Size)
		for _, c := range b.Markers {
			if c.Y == i {
				rows[i][c.X] = 1
			}
		}
	}

	for i, row := range rows {
		fmt.Fprintf(w, "[")
		for _, col := range row {
			if col == 0 {
				fmt.Fprintf(w, ".")
			} else {
				fmt.Fprintf(w, "o")
			}
		}
		fmt.Fprintf(w, "]")
		if i == 0 {
			fmt.Fprintf(w, " %v", b.SquareDistances())
		}
		fmt.Fprintln(w)
	}
}

// Generate generates a board of the given size. A given id will always
// return the same board.
func Generate(size int, id int64) Board {
	result := make([]Coord, size)
	for i := 0; i < size; i++ {
		var c Coord
		c.X = int(id % int64(size))
		id /= int64(size)
		c.Y = int(id % int64(size))
		id /= int64(size)
		result[i] = c
	}

	return Board{
		Markers: result,
		Size:    size,
		ID:      id,
	}
}

// SquareDistances returns the squares of all the pairwise distances between markers on b
func (b *Board) SquareDistances() []int {
	result := make([]int, 0, b.Size*(1+b.Size)/2)
	for i := 0; i < b.Size-1; i++ {
		for j := i + 1; j < b.Size; j++ {
			result = append(result, b.Markers[i].SquareDistance(b.Markers[j]))
		}
	}
	return result
}

// SquareDistance is the square of the distance between c and o
func (c Coord) SquareDistance(o Coord) int {
	dx := o.X - c.X
	dy := o.Y - c.Y
	return dx*dx + dy*dy
}

// Amount is the amount of different boards of size n
func Amount(n int) int64 {
	/* n ^ 2n */
	result := int64(1)

	for i := 0; i < 2*n; i++ {
		result *= int64(n)
	}

	return result
}
