package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/lutzky/unique-distance/board"
)

var (
	workers     = flag.Int64("workers", 4, "Number of workers for parallel version")
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

func findUnique(w io.Writer, config findUniqueConfig) int64 {
	var found int64
	for i := int64(0); i < board.Amount(config.boardSize); i++ {
		b := board.Generate(config.boardSize, i)
		ds := b.SquareDistances()
		if allUnique(ds, b.MaxDistance()) {
			if config.printAll {
				b.Print(w)
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

func findUniqueParallel(w io.Writer, config findUniqueConfig) int64 {
	var found int64
	boardsPerWorker := board.Amount(config.boardSize) / *workers
	ch := make(chan int64)
	for i := int64(0); i < *workers; i++ {
		go func(i int64) {
			for q := int64(0); q < boardsPerWorker; q++ {
				b := board.Generate(config.boardSize, boardsPerWorker*i+q)
				ds := b.SquareDistances()
				if allUnique(ds, b.MaxDistance()) {
					ch <- b.ID
				} else {
					ch <- -1
				}
			}
		}(i)
	}

	for i := int64(0); i < board.Amount(config.boardSize); i++ {
		bID := <-ch
		if bID != -1 {
			if config.printAll {
				b := board.Generate(config.boardSize, bID)
				b.Print(w)
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
