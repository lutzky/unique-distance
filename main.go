package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/lutzky/unique-distance/board"
)

var (
	workers     = flag.Int64("workers", 64, "Number of workers for parallel version")
	useParallel = flag.Bool("use_parallel", false, "Use parallel implementation")
)

type findUniqueConfig struct {
	boardSize int
	printAll  bool
	quitAfter int
}

func registerFlags(fs *flag.FlagSet, config *findUniqueConfig) {
	fs.IntVar(&config.boardSize, "n", 3, "Board size")
	fs.BoolVar(&config.printAll, "print_all", true, "Print all valid boards seen")
	fs.IntVar(&config.quitAfter, "quit_after", 0, "Quit after finding this many solutions (0 for 'all')")
}

func main() {
	config := findUniqueConfig{}
	registerFlags(flag.CommandLine, &config)
	flag.Parse()

	var found int64
	if *useParallel {
		found = findUniqueParallel(os.Stdout, config)
	} else {
		found = findUnique(os.Stdout, config)
	}
	fmt.Printf("Found %d solutions\n", found)
}

func findUnique(w io.Writer, config findUniqueConfig) int64 {
	found := map[int64]bool{}
	workspace := make([]bool, board.MaxDistance(config.boardSize)+1)
	for i := int64(0); i < board.Amount(config.boardSize); i++ {
		b := board.Generate(config.boardSize, i)
		ds := b.SquareDistances()
		if allUnique(ds, workspace) {
			b.Normalize()
			if !found[b.ID] {
				if config.printAll {
					b.Print(w)
					fmt.Fprintln(w)
				}
				found[b.ID] = true
			}
			if config.quitAfter != 0 && len(found) >= config.quitAfter {
				return int64(len(found))
			}
		}
	}
	return int64(len(found))
}

func findUniqueParallel(w io.Writer, config findUniqueConfig) int64 {
	found := map[int64]bool{}
	var wg sync.WaitGroup
	boardsPerWorker := board.Amount(config.boardSize) / *workers
	ch := make(chan int64)
	for i := int64(0); i < *workers; i++ {
		wg.Add(1)
		go func(i int64) {
			workspace := make([]bool, board.MaxDistance(config.boardSize)+1)
			for q := int64(0); q < boardsPerWorker; q++ {
				b := board.Generate(config.boardSize, boardsPerWorker*i+q)
				ds := b.SquareDistances()
				if allUnique(ds, workspace) {
					b.Normalize()
					ch <- b.ID
				}
			}
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for bID := range ch {
		if !found[bID] {
			if config.printAll {
				b := board.Generate(config.boardSize, bID)
				b.Print(w)
				fmt.Fprintln(w)
			}
			found[bID] = true

			if config.quitAfter != 0 && len(found) >= config.quitAfter {
				return int64(len(found))
			}
		}
	}
	return int64(len(found))
}

func allUnique(ns []int, workspace []bool) bool {
	if len(ns) < 2 {
		return true
	}
	for i := range workspace {
		workspace[i] = false
	}
	for _, d := range ns {
		if workspace[d] {
			return false
		}
		workspace[d] = true
	}
	return true
}
