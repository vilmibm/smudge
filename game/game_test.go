package game

import (
	"reflect"
	"testing"
)

func Test_NewRay(t *testing.T) {
	tests := []struct {
		name     string
		args     []Point
		expected []Point
	}{
		{
			name:     "0 ray",
			args:     []Point{{0, 0}, {0, 0}},
			expected: []Point{},
		},
		{
			name:     "1 right",
			args:     []Point{{0, 0}, {1, 0}},
			expected: []Point{{0, 0}, {1, 0}},
		},
		{
			name:     "1 left",
			args:     []Point{{1, 0}, {0, 0}},
			expected: []Point{{1, 0}, {0, 0}},
		},
		{
			name:     "1 up",
			args:     []Point{{0, 0}, {0, 1}},
			expected: []Point{{0, 0}, {0, 1}},
		},
		{
			name:     "1 down",
			args:     []Point{{0, 1}, {0, 0}},
			expected: []Point{{0, 1}, {0, 0}},
		},
		{
			name:     "1 nw",
			args:     []Point{{1, 1}, {0, 0}},
			expected: []Point{{1, 1}, {0, 0}},
		},
		{
			name:     "1 ne",
			args:     []Point{{1, 1}, {2, 0}},
			expected: []Point{{1, 1}, {2, 0}},
		},
		{
			name:     "1 se",
			args:     []Point{{1, 1}, {2, 2}},
			expected: []Point{{1, 1}, {2, 2}},
		},
		{
			name:     "1 sw",
			args:     []Point{{1, 1}, {0, 2}},
			expected: []Point{{1, 1}, {0, 2}},
		},
		{
			name: "n",
			args: []Point{{0, 10}, {0, 0}},
			expected: []Point{
				{0, 10},
				{0, 9},
				{0, 8},
				{0, 7},
				{0, 6},
				{0, 5},
				{0, 4},
				{0, 3},
				{0, 2},
				{0, 1},
				{0, 0}},
		},
		{
			name: "s",
			args: []Point{{0, 0}, {0, 10}},
			expected: []Point{
				{0, 0},
				{0, 1},
				{0, 2},
				{0, 3},
				{0, 4},
				{0, 5},
				{0, 6},
				{0, 7},
				{0, 8},
				{0, 9},
				{0, 10},
			},
		},
		{
			name: "w",
			args: []Point{{0, 0}, {10, 0}},
			expected: []Point{
				{0, 0},
				{1, 0},
				{2, 0},
				{3, 0},
				{4, 0},
				{5, 0},
				{6, 0},
				{7, 0},
				{8, 0},
				{9, 0},
				{10, 0},
			},
		},
		{
			name: "e",
			args: []Point{{10, 0}, {0, 0}},
			expected: []Point{
				{10, 0},
				{9, 0},
				{8, 0},
				{7, 0},
				{6, 0},
				{5, 0},
				{4, 0},
				{3, 0},
				{2, 0},
				{1, 0},
				{0, 0},
			},
		},
		{
			name: "ne",
			args: []Point{{10, 10}, {20, 0}},
			expected: []Point{
				{10, 10},
				{11, 9},
				{12, 8},
				{13, 7},
				{14, 6},
				{15, 5},
				{16, 4},
				{17, 3},
				{18, 2},
				{19, 1},
				{20, 0},
			},
		},
		{
			name: "se",
			args: []Point{{0, 0}, {10, 10}},
			expected: []Point{
				{0, 0},
				{1, 1},
				{2, 2},
				{3, 3},
				{4, 4},
				{5, 5},
				{6, 6},
				{7, 7},
				{8, 8},
				{9, 9},
				{10, 10},
			},
		},
		{
			name: "sw",
			args: []Point{{10, 0}, {0, 10}},
			expected: []Point{
				{10, 0},
				{9, 1},
				{8, 2},
				{7, 3},
				{6, 4},
				{5, 5},
				{4, 6},
				{3, 7},
				{2, 8},
				{1, 9},
				{0, 10},
			},
		},
		{
			name: "nw",
			args: []Point{{10, 10}, {0, 0}},
			expected: []Point{
				{10, 10},
				{9, 9},
				{8, 8},
				{7, 7},
				{6, 6},
				{5, 5},
				{4, 4},
				{3, 3},
				{2, 2},
				{1, 1},
				{0, 0},
			},
		},
		{
			name: "flat line",
			args: []Point{{0, 0}, {5, 1}},
			expected: []Point{
				{0, 0},
				{1, 1},
				{2, 1},
				{3, 1},
				{4, 1},
				{5, 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRay(tt.args[0], tt.args[1])
			if !reflect.DeepEqual(r.Points, tt.expected) {
				t.Errorf("%v != %v", tt.expected, r.Points)
			}
		})
	}
}

// TODO other tests
