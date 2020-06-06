package board

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMirror(t *testing.T) {
	testCases := []struct {
		name        string
		columns     bool
		input, want Board
	}{
		{
			name:    "topLeftTriangle-rows",
			columns: false,
			input: fromString(3, `
			oo.
			o..
			...
			`),
			want: fromString(3, `
			...
			o..
			oo.
			`),
		},
		{
			name:    "topLeftTriangle-columns",
			columns: true,
			input: fromString(3, `
			oo.
			o..
			...
			`),
			want: fromString(3, `
			.oo
			..o
			...
			`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Mirroring board:\n%s", tc.input)
			tc.input.Mirror(tc.columns)

			if d := cmp.Diff(tc.want, tc.input, boardCmpOpt); d != "" {
				t.Errorf("tc.input.Mirror(columns: %t) returned this board:\n%sBut wanted this one:\n%s\nDiff -want +got:\n%s",
					tc.columns, tc.input, tc.want, d)
			}
		})
	}
}

func TestRotate(t *testing.T) {
	testCases := []struct {
		name  string
		input Board
		wants []Board
	}{
		{
			name: "topLeftTriangle",
			input: fromString(3, `
			23.
			4..
			...
			`),
			wants: []Board{
				fromString(3, `
				.42
				..3
				...
				`),
				fromString(3, `
				...
				..4
				.32
				`),
				fromString(3, `
				...
				3..
				24.
				`),
				fromString(3, `
				23.
				4..
				...
				`),
			},
		},
		{
			name: "justCenter",
			input: fromString(3, `
			...
			.o.
			...
			`),
			wants: []Board{
				fromString(3, `
				...
				.o.
				...
				`),
				fromString(3, `
				...
				.o.
				...
				`),
			},
		},
		{
			name: "tetris-L",
			input: fromString(4, `
			....
			.o..
			.o..
			.oo.
			`),
			wants: []Board{
				fromString(4, `
				....
				ooo.
				o...
				....
				`),
				fromString(4, `
				.oo.
				..o.
				..o.
				....
				`),
				fromString(4, `
				....
				...o
				.ooo
				....
				`),
				fromString(4, `
				....
				.o..
				.o..
				.oo.
				`),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, want := range tc.wants {
				t.Logf("Rotating board:\n%s", tc.input)
				tc.input.Rotate()

				if d := cmp.Diff(want, tc.input, boardCmpOpt); d != "" {
					t.Errorf("tc.input.Rotate() returned this board:\n%sBut wanted this one:\n%s\nDiff -want +got:\n%s",
						tc.input, want, d)
				}
			}
		})
	}
}
