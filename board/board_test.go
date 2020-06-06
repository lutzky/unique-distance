package board

import (
	"fmt"
	"sort"
	"strconv"
	"testing"
	"testing/quick"
	"unicode"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func fromString(size int, s string) Board {
	result := Board{Size: size}
	var i int64
	for _, r := range s {
		if unicode.IsSpace(r) {
			continue
		}
		n := int64(0)
		if r == '.' {
			r = '0'
		}
		if r == 'o' {
			r = '1'
		}
		n, err := strconv.ParseInt(string(r), 16, 16)
		if err != nil {
			panic(fmt.Sprintf("Board has invalid marker '%c':\n%s", r, s))
		}
		x := i % int64(size)
		y := i / int64(size)
		for j := int64(0); j < n; j++ {
			result.Markers = append(result.Markers, Coord{int(x), int(y)})
		}
		i++
	}

	return result
}

func TestSquareDistances(t *testing.T) {
	testCases := []struct {
		name  string
		input Board
		want  []int
	}{
		{
			"3x3 diag",
			fromString(3, `
			o..
			.o.
			..o
			`),
			[]int{2, 2, 8},
		},
		{
			"[x  ][ xx][   ]",
			fromString(3, `
			o..
			.oo
			...
			`),
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
			fromString(3, `
			3..
			...
			...
			`),
		},
		{
			3,
			250,
			fromString(3, `
			o..
			o..
			.o.
			`),
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
