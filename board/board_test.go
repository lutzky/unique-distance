package board

import (
	"fmt"
	"sort"
	"testing"
	"testing/quick"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSquareDistances(t *testing.T) {
	testCases := []struct {
		name  string
		input Board
		want  []int
	}{
		{
			"3x3 diag",
			Board{Size: 3, Markers: []Coord{{0, 0}, {1, 1}, {2, 2}}},
			[]int{2, 2, 8},
		},
		{
			"[x  ][ xx][   ]",
			Board{Size: 3, Markers: []Coord{{0, 0}, {1, 1}, {2, 1}}},
			[]int{1, 2, 5},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.input.SquareDistances()
			sort.Ints(got)
			if d := cmp.Diff(got, tc.want, boardCmpOpt); d != "" {
				t.Errorf("Board:\n%s\nDiff -got +want:\n%s", tc.input, d)
			}
		})
	}
}

var boardCmpOpt = cmpopts.SortSlices(func(a, b Coord) bool {
	if a.X != b.X {
		return a.X < b.X
	}
	return a.Y < b.Y
})

func TestGenerate(t *testing.T) {
	testCases := []struct {
		size  int
		input int64
		want  Board
	}{
		{
			3,
			0,
			Board{Size: 3, Markers: []Coord{{0, 0}, {0, 0}, {0, 0}}},
		},
		{
			3,
			250,
			Board{Size: 3, Markers: []Coord{{1, 2}, {0, 0}, {0, 1}}},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Gen%d", tc.input), func(t *testing.T) {
			got := Generate(tc.size, tc.input)
			if d := cmp.Diff(tc.want, got, boardCmpOpt); d != "" {
				t.Errorf("Want:\n%v\nGot:\n%v", tc.want, got)
				t.Errorf("Diff -want +got:\n%s", d)
			}
		})
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
		got := Amount(tc.input)
		if got != tc.want {
			t.Errorf("numBoards(%d) = %d; want %d", tc.input, got, tc.want)
		}
	}
}

func TestUnusualDistance(t *testing.T) {
	getBoard := func(nn uint64) Board {
		n := int64(nn % 9000)
		return Generate(4, n)
	}
	f := func(nn uint64) bool {
		board := getBoard(nn)
		ds := board.SquareDistances()
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
		b := getBoard(n)
		t.Errorf("%v:\n%s", err, b)
	}
}
