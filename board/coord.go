package board

import "fmt"

// Coord represents a possible coordinate of a marker on a board
type Coord struct {
	X, Y int
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

// SquareDistance is the square of the distance between c and o
func (c Coord) SquareDistance(o Coord) int {
	dx := o.X - c.X
	dy := o.Y - c.Y
	return dx*dx + dy*dy
}

func (c Coord) sortCompare(o Coord) bool {
	if c.X != o.X {
		return c.X < o.X
	}
	return c.Y < o.Y
}
