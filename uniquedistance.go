package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
)

var (
	boardSize = flag.Int("n", 3, "Board size")
	printAll  = flag.Bool("print_all", true, "Print all valid boards seen")
	quitAfter = flag.Int64("quit_after", 0, "Quit after finding this many solutions (0 for 'all')")
)

func main() {
	flag.Parse()
	found := findUnique(os.Stdout, *boardSize)
	fmt.Printf("Found %d solutions\n", found)
}

type coord struct {
	X, Y int
}

func printBoard(w io.Writer, board []coord, ds []int) {
	if w == nil {
		return
	}
	n := len(board)
	rows := make([][]int, n)

	for i := 0; i < n; i++ {
		rows[i] = make([]int, n)
		for _, c := range board {
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
			fmt.Fprintf(w, " %v", ds)
		}
		fmt.Fprintln(w)
	}
}

func numBoards(n int) int64 {
	/* n ^ 2n */
	result := int64(1)

	for i := 0; i < 2*n; i++ {
		result *= int64(n)
	}

	return result
}

func findUnique(w io.Writer, n int) int64 {
	var found int64
	for i := int64(0); i < numBoards(n); i++ {
		board := boardN(n, i)
		ds := sqDistances(board)
		if allUnique(ds) {
			if *printAll {
				printBoard(w, board, ds)
				fmt.Fprintln(w)
			}
			found++
			if *quitAfter != 0 && found >= *quitAfter {
				return found
			}
		}
	}
	return found
}

func (c coord) sqDist(o coord) int {
	dx := o.X - c.X
	dy := o.Y - c.Y
	return dx*dx + dy*dy
}

func sqDistances(board []coord) []int {
	var result []int
	for i := 0; i < len(board)-1; i++ {
		for j := i + 1; j < len(board); j++ {
			result = append(result, board[i].sqDist(board[j]))
		}
	}
	return result
}

func allUnique(ns []int) bool {
	if len(ns) < 2 {
		return true
	}
	sort.Ints(ns)
	for i := 0; i < len(ns)-1; i++ {
		if ns[i] == ns[i+1] {
			return false
		}
	}
	return true
}

func boardN(size int, input int64) []coord {
	var result []coord
	for i := 0; i < size; i++ {
		var c coord
		c.X = int(input % int64(size))
		input /= int64(size)
		c.Y = int(input % int64(size))
		input /= int64(size)
		result = append(result, c)
	}
	return result
}
