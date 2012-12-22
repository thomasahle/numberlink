package main

import "flag"

var (
	widthFlag  = flag.Int("width", 4, "Width of the generated paper")
	heightFlag = flag.Int("height", 4, "Height of the generated paper")
)
/*

import "strings"
import "math/rand"

func canGo(paper *ppr.Paper, alpha rune, from int, dir int) bool {
	next := from + paper.vctr[dir]
	if paper.table[next] != EMPTY {
		return false
	}
	// Check no self touch (from won't be set to alpha yet)
	neighbours := 0
	for _, dir := range [4]int{N,E,S,W} {
		if paper.table[next + paper.vctr[dir]] == alpha {
			neighbours += 1
		}
	}
	if neighbours != 0 {
		return false
	}
	// Check we don't leave isolated points around from
	for _, dir := range [4]int{N,E,S,W} {
		isopos := from + paper.vctr[dir]
		if paper.table[isopos] == EMPTY && isopos != next {
			anyneighs := false
			for _, dir := range[4]int{N,E,S,W} {
				neigh := isopos + paper.vctr[dir]
				if paper.table[neigh] == EMPTY && neigh != from {
					anyneighs = true
				}
			}
			if !anyneighs {
				return false
			}
		}
	}
	// Check we don't leave isolated points around next
	// Its okay for a point to be isolated if it doesn't
	// have any neighbours in our color
	for _, dir := range [4]int{N,E,S,W} {
		isopos := next + paper.vctr[dir]
		if paper.table[isopos] == EMPTY && isopos != from {
			anyneighs := false
			alpneighs := false
			for _, dir := range[4]int{N,E,S,W} {
				neigh := isopos + paper.vctr[dir]
				if paper.table[neigh] == EMPTY && neigh != next {
					anyneighs = true
				}
				if paper.table[neigh] == alpha {
					alpneighs = true
				}
			}
			if alpneighs && !anyneighs {
				return false
			}
		}
	}
	// These 

	// nej
	//aaa bb
	//a  fn
	//aaaxcc

	// ja
	//aaa bb
	//a  n
	//aaafcc

	// nej
	//xx
	// x
	//nf

	// Fuck
	// .
	// ..
	// .

	return true
}

func generate(w, h int) {
	table := make([]rune, w*h)
	for i := 0; i < w*h; i++ {
		table[i] = '.'
	}
	paper := ppr.NewPaper(w, h, table)
	w, h = paper.width, paper.height

	alpha := 0
	empties := (w-2)*(h-2)
	for empties != 0 {
		for pos := 0; pos < w*h; pos++ {
			if paper.table[pos] == EMPTY {
				//maxLen = rand.Intn(empties)
				p := pos
				changed := true
				for changed {
					changed = false
					for _, i := range rand.Perm(4) {
						if dir := []int{N,E,S,W}[i]; canGo(paper, SIGMA[alpha], p, dir) {
							paper.table[p] = SIGMA[alpha]
							p = p + paper.vctr[dir]
							changed = true
							break
						}
					}
				}
				if p != pos {
					paper.table[p] = SIGMA[alpha]
					alpha += 1
				}

				for y := 1; y < paper.height-1; y++ {
					for x := 1; x < paper.width-1; x++ {
						fmt.Printf("%c", paper.table[y*paper.width+x])
					}
					fmt.Println()
				}
				fmt.Println()
			}
		}
	}
}
func main() {
	flag.Parse()
	fmt.Println(*widthFlag, *heightFlag)
}

*/
