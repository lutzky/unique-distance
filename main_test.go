package main

import (
	"io"
	"strings"
	"testing"
)

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
		got := allUnique(tc.input, make([]bool, 17))
		if got != tc.want {
			t.Errorf("allUnique(%v) = %t; want %t", tc.input, got, tc.want)
		}
	}
}

func TestFindUnique(t *testing.T) {
	testCases := []struct {
		name   string
		config findUniqueConfig
		want   int64
	}{
		{"3x3", findUniqueConfig{boardSize: 3}, 5},
		{"4x4", findUniqueConfig{boardSize: 4}, 23},
		{"SLOW 5x5", findUniqueConfig{boardSize: 5}, 35},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if testing.Short() && strings.HasPrefix(tc.name, "SLOW") {
				t.Skip("Skipping slow test")
			}
			got := findUnique(nil, tc.config)
			if got != tc.want {
				t.Errorf("got: %d; want %d", got, tc.want)
			}
		})
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
