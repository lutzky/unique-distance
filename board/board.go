package board

import (
	"bytes"
	"fmt"
	"io"
	"sort"
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

// MaxDistance is the maximal possible distance between two markers on a board
// of the given size.
func MaxDistance(size int) int {
	return 2 * (size - 1) * (size - 1)
}

func (b Board) String() string {
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
				rows[i][c.X]++
			}
		}
	}

	for i, row := range rows {
		fmt.Fprintf(w, "[")
		for _, col := range row {
			switch {
			case col == 0:
				fmt.Fprintf(w, ".")
			case col == 1:
				fmt.Fprintf(w, "o")
			case col < 16:
				fmt.Fprintf(w, "%x", col)
			default:
				panic("Board with more than 16 markers in same spot: " + fmt.Sprint(b.Markers))
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
//
// Note that the board's ID may be changed after several operations,
// which sort the markers. However, generating a new board from the
// modified ID will result an a board with the same markers.
func Generate(size int, id int64) Board {
	result := Board{
		Markers: markersFromID(size, id),
		Size:    size,
		ID:      id,
	}

	return result
}

func markersFromID(size int, id int64) []Coord {
	markers := make([]Coord, size)
	for i := 0; i < size; i++ {
		var c Coord
		c.X = int(id % int64(size))
		id /= int64(size)
		c.Y = int(id % int64(size))
		id /= int64(size)
		markers[i] = c
	}

	return markers
}

func (b *Board) sortMarkers() {
	sort.Slice(b.Markers, func(i, j int) bool {
		return b.Markers[i].sortCompare(b.Markers[j])
	})
}

// SquareDistances returns the squares of all the pairwise distances between markers on b
func (b *Board) SquareDistances() []int {
	if len(b.Markers) == 0 {
		return nil
	}
	result := make([]int, 0, b.Size*(1+b.Size)/2)
	for i := 0; i < b.Size-1; i++ {
		for j := i + 1; j < b.Size; j++ {
			result = append(result, b.Markers[i].SquareDistance(b.Markers[j]))
		}
	}
	return result
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
