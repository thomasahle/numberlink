package main

import (
	"reflect"
	"testing"
)

var flattentests = []struct {
	in  [][]int
	out [][]int
}{
	{[][]int{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}}, [][]int{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}}},
	{[][]int{{0, 1, 2}, {0, 1, 2}, {0, 1, 2}}, [][]int{{0, 1, 2}, {0, 1, 2}, {0, 1, 2}}},
	{[][]int{{3, 4, 5}, {3, 4, 5}, {3, 4, 5}}, [][]int{{0, 1, 2}, {0, 1, 2}, {0, 1, 2}}},
	{[][]int{{3, 3}, {3, 3}}, [][]int{{0, 0}, {0, 0}}},
	{[][]int{{3, 4, 5}, {3, 5, 5}}, [][]int{{0, 1, 2}, {0, 2, 2}}},
}

func TestFlatten(t *testing.T) {
	for _, tt := range flattentests {
		flatten(tt.in)
		if !reflect.DeepEqual(tt.in, tt.out) {
			t.Errorf("Expected %x, got %x", tt.out, tt.in)
		}
	}
}
