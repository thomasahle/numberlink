package main

import "testing"

var papertests = []struct {
	width int
	height int
	n int
	out int
}{
	{1, 4, 2, 2},
	{1, 5, 2, 4},
	{2, 1, 2, 0},
	{2, 2, 2, 4},
	{2, 3, 2, 14},
	{2, 4, 2, 18},
	{2, 5, 2, 18},
	{3, 3, 2, 24},
	{3, 4, 2, 32},
	{3, 5, 2, 36},
	{4, 4, 2, 24},
	{4, 5, 2, 44},
	{4, 6, 2, 44},
	{5, 5, 2, 48},
/*	{6, 6, 2, 72},
	{7, 7, 2, 96},
	{8, 8, 2, 96},
	{9, 9, 2, 144},
	{10, 10, 2, 240},*/
}

func TestSolve(t *testing.T) {
	for _, tt := range papertests {
		if tt.n != 2 {
			t.Errorf("Currently only n=2 is supported")
			continue
		}

		al := choose2(tt.width*tt.height)
		bl := choose2(tt.width*tt.height-2)
		count := 0
		as := []int{0,1}
		for i := 0; i < al; i++ {
			bs := []int{0,1}
			for j := 0; j < bl; j++ {
				a1, a2 := as[0], as[1]
				b1, b2 := bs[0], bs[1]
				if a1 <= b1 {
					b1++
				}
				if a2 <= b1 {
					b1++
				}
				if a1 <= b2 {
					b2++
				}
				if a2 <= b2 {
					b2++
				}
				p := create([][2]int{[2]int{a1,a2}, [2]int{b1,b2}}, tt.width, tt.height)
				res := Solve(p)
				if res {
					count++
				}

				if j+1 < bl {
					nextCombination(bs, tt.width*tt.height-2)
				}
			}
			if i+1 < al {
				nextCombination(as, tt.width*tt.height)
			}
		}

		if count != tt.out {
			t.Errorf("Expected %d, got %d for %x", tt.out, count, tt)
		}
	}
}
