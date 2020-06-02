package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	boardSize   = flag.Int("n", 3, "Board size")
	printAll    = flag.Bool("print_all", true, "Print all valid boards seen")
	quitAfter   = flag.Int64("quit_after", 0, "Quit after finding this many solutions (0 for 'all')")
	useParallel = flag.Bool("use_parallel", false, "Use parallel implementation")
)

type findUniqueConfig struct {
	boardSize int
	printAll  bool
	quitAfter int64
}

func main() {
	flag.Parse()
	config := findUniqueConfig{
		boardSize: *boardSize,
		printAll:  *printAll,
		quitAfter: *quitAfter,
	}
	var found int64
	if *useParallel {
		found = findUniqueParallel(os.Stdout, config)
	} else {
		found = findUnique(os.Stdout, config)
	}
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

func maxDistance(boardSize int) int {
	return 2 * (boardSize - 1) * (boardSize - 1)
}

func findUnique(w io.Writer, config findUniqueConfig) int64 {
	var found int64
	for i := int64(0); i < numBoards(config.boardSize); i++ {
		board := boardN(config.boardSize, i)
		ds := sqDistances(board)
		if allUnique(ds, maxDistance(config.boardSize)) {
			if config.printAll {
				printBoard(w, board, ds)
				fmt.Fprintln(w)
			}
			found++
			if config.quitAfter != 0 && found >= config.quitAfter {
				return found
			}
		}
	}
	return found
}

const Workers = 4

func findUniqueParallel(w io.Writer, config findUniqueConfig) int64 {
	var found int64
	boardsPerWorker := numBoards(config.boardSize) / Workers
	ch := make(chan []coord)
	for i := int64(0); i < Workers; i++ {
		go func(i int64) {
			for q := int64(0); q < boardsPerWorker; q++ {
				board := boardN(config.boardSize, boardsPerWorker*i+q)
				ds := sqDistances(board)
				if allUnique(ds, maxDistance(config.boardSize)) {
					ch <- board
				} else {
					ch <- nil
				}
			}
		}(i)
	}

	for i := int64(0); i < numBoards(config.boardSize); i++ {
		board := <-ch
		if board != nil {
			if config.printAll {
				ds := sqDistances(board)
				printBoard(w, board, ds)
				fmt.Fprintln(w)
			}
			found++
			if config.quitAfter != 0 && found >= config.quitAfter {
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
	result := make([]int, 0, len(board)*(1+len(board))/2)
	for i := 0; i < len(board)-1; i++ {
		for j := i + 1; j < len(board); j++ {
			result = append(result, board[i].sqDist(board[j]))
		}
	}
	return result
}

func allUnique(ns []int, max int) bool {
	if len(ns) < 2 {
		return true
	}
	found := make([]bool, max+1)
	for _, d := range ns {
		if found[d] {
			return false
		}
		found[d] = true
	}
	return true
}

func boardN(size int, input int64) []coord {
	result := make([]coord, size)
	for i := 0; i < size; i++ {
		var c coord
		c.X = int(input % int64(size))
		input /= int64(size)
		c.Y = int(input % int64(size))
		input /= int64(size)
		result[i] = c
	}
	return result
}
