package main

import (
	"io"
	"os"
	"sort"
	"strings"
	"testing"
	"testing/quick"

	"github.com/google/go-cmp/cmp"
)

func TestDistances(t *testing.T) {
	testCases := []struct {
		name  string
		input []coord
		want  []int
	}{
		{
			"3x3 diag",
			[]coord{{0, 0}, {1, 1}, {2, 2}},
			[]int{2, 2, 8},
		},
		{
			"[x  ][ xx][   ]",
			[]coord{{0, 0}, {1, 1}, {2, 1}},
			[]int{1, 2, 5},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := sqDistances(tc.input)
			sort.Ints(got)
			if d := cmp.Diff(got, tc.want); d != "" {
				t.Errorf("Diff -got +want:\n%s", d)
			}
		})
	}
}

func TestAllUnique(t *testing.T) {
	testCases := []struct {
		input []int
		want  bool
	}{
		{[]int{4, 3, 2, 1}, true},
		{[]int{1, 3, 2, 1}, false},
		{[]int{1, 3, 3, 1}, false},
		{[]int{1, 3, 3, 3}, false},
		{[]int{1, 2, 3, 4}, true},
		{[]int{3}, true},
		{[]int{}, true},
	}

	for _, tc := range testCases {
		got := allUnique(tc.input)
		if got != tc.want {
			t.Errorf("allUnique(%v) = %t; want %t", tc.input, got, tc.want)
		}
	}
}

func TestBoardN(t *testing.T) {
	testCases := []struct {
		size  int
		input int64
		want  []coord
	}{
		{
			3,
			0,
			[]coord{{0, 0}, {0, 0}, {0, 0}},
		},
		{
			3,
			250,
			[]coord{{1, 2}, {0, 0}, {0, 1}},
		},
	}

	for _, tc := range testCases {
		got := boardN(tc.size, tc.input)
		if d := cmp.Diff(got, tc.want); d != "" {
			t.Errorf("Diff -got +want:\n%s", d)
		}
	}
}

func TestNumBoards(t *testing.T) {
	testCases := []struct {
		input int
		want  int64
	}{
		{1, 1},
		{2, 16},
		{3, 729},
		{4, 65536},
	}

	for _, tc := range testCases {
		got := numBoards(tc.input)
		if got != tc.want {
			t.Errorf("numBoards(%d) = %d; want %d", tc.input, got, tc.want)
		}
	}
}

func TestFindUnique(t *testing.T) {
	skipSlowTests := true

	testCases := []struct {
		name   string
		config findUniqueConfig
		want   int64
	}{
		{"3x3", findUniqueConfig{boardSize: 3}, 240},
		{"4x4", findUniqueConfig{boardSize: 4}, 4416},
		{"SLOW 5x5", findUniqueConfig{boardSize: 5}, 33600},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if skipSlowTests && strings.HasPrefix(tc.name, "SLOW") {
				t.Skip("Skipping slow test")
			}
			got := findUnique(nil, tc.config)
			if got != tc.want {
				t.Errorf("got: %d; want %d", got, tc.want)
			}
		})
	}
}

func TestUnusualDistance(t *testing.T) {
	getBoard := func(nn uint64) []coord {
		n := int64(nn % 9000)
		return boardN(4, n)
	}
	f := func(nn uint64) bool {
		board := getBoard(nn)
		ds := sqDistances(board)
		for _, x := range ds {
			if x > 18 {
				return false
			}
		}
		return true
	}
	if err := quick.Check(f, nil); err != nil {
		eiei := err.(*quick.CheckError)
		n := eiei.In[0].(uint64)
		board := getBoard(n)
		ds := sqDistances(board)
		printBoard(os.Stdout, board, ds)
		t.Error(err)
	}
}

func BenchmarkFindUnique(b *testing.B) {
	benches := []struct {
		f    func(io.Writer, findUniqueConfig) int64
		name string
	}{
		{findUnique, "findUnique"},
		{findUniqueParallel, "findUniqueParallel"},
	}

	for _, bm := range benches {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bm.f(nil, findUniqueConfig{
					boardSize: 4,
					printAll:  false,
					quitAfter: 0,
				})
			}
		})
	}
}
