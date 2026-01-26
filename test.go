package main

import "fmt"

func choose2(n int) int {
	return n*(n-1)/2
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func nextCombination(as []int, n int) {
	r := len(as)
	i := r - 1
	for as[i] == n - r + i {
		i--
	}
	as[i]++
	for j := i + 1; j < r; j++ {
		as[j] = as[j-1] + 1
	}
}

func create(sources [][2]int, w, h int) *Paper {
	table := make([]rune, w*h)
	for i := 0; i < w*h; i++ {
		table[i] = '.'
	}
	for i, pair := range sources {
		table[pair[0]] = SIGMA[i]
		table[pair[1]] = SIGMA[i]
	}
	return NewPaper(w, h, table)
}

func printThemAll(w, h int) {
	al := choose2(w*h)
	bl := choose2(w*h-2)
	fmt.Println(al, bl)
	count := 0
	max := 0
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
			p := create([][2]int{[2]int{a1,a2}, [2]int{b1,b2}}, w, h)
			res := Solve(p)
			if res {
				PrintTubes(p, true)
				fmt.Println(Calls)
				max = imax(max, Calls)
				count++
			}
			Calls = 0

			if j+1 < bl {
				nextCombination(bs, w*h-2)
			}
		}
		if i+1 < al {
			nextCombination(as, w*h)
		}
	}
	fmt.Println("Count:", count)
	fmt.Println("Max:", max)
}
