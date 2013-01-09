package main

const (
	GRASS = '#'
	EMPTY = '.'
)

const (
	N = 1
	E = 2
	S = 4
	W = 8
)

var (
	// Array for iterating over simple directions
	DIRS = [4]int{N, E, S, W}
	// Mirrors simple directions
	MIR = [16]int{N: S, E: W, S: N, W: E}
	// Filter for N|E, N|W, S|E and S|W
	DIAG = [16]bool{N | E: true, N | W: true, S | E: true, S | W: true}
)

type Paper struct {
	Width  int
	Height int

	Vctr [16]int
	Crnr [16]int

	Table []int32
	Con   []int

	source []bool
	end    []int
	canSE  []bool
	canSW  []bool

	next []int
}

func NewPaper(width, height int, table []rune) *Paper {
	paper := new(Paper)

	// Pad the given table with #, to make boundery checks easier
	w, h := width+2, height+2
	paper.Width, paper.Height = w, h
	paper.Table = make([]rune, 0, w*h)
	for i := 0; i < w; i++ {
		paper.Table = append(paper.Table, GRASS)
	}
	for y := 1; y < h-1; y++ {
		paper.Table = append(paper.Table, GRASS)
		for x := 1; x < w-1; x++ {
			paper.Table = append(paper.Table, table[(y-1)*(w-2)+(x-1)])
		}
		paper.Table = append(paper.Table, GRASS)
	}
	for i := 0; i < w; i++ {
		paper.Table = append(paper.Table, GRASS)
	}

	paper.initTables()

	return paper
}

func Solve(paper *Paper) bool {
	return chooseConnection(paper, paper.Crnr[N|W])
}

var Calls = 0

func chooseConnection(paper *Paper, pos int) bool {
	Calls++

	// Final
	if pos == 0 {
		return paper.validate()
	}

	w := paper.Width
	if paper.source[pos] {
		switch paper.Con[pos] {
		// If the source is not yet connection
		case 0:
			// We can't connect E if we have a NE corner
			if paper.Con[pos-w+1] != S|W {
				if tryConnection(paper, pos, E) {
					return true
				}
			}
			// South connections can create a forced SE position
			if checkImplicitSE(paper, pos) {
				if tryConnection(paper, pos, S) {
					return true
				}
			}
		// If the source is already connected
		case N, W:
			return chooseConnection(paper, paper.next[pos])
		}
	} else {
		switch paper.Con[pos] {
		// SE
		case 0:
			// Should we check for implied N|W?
			if paper.canSE[pos] {
				return tryConnection(paper, pos, E|S)
			}
		// SW or WE
		case W:
			// Check there is a free line down to the source we are turning around
			if paper.canSW[pos] && checkSWLane(paper, pos) && checkImplicitSE(paper, pos) {
				if tryConnection(paper, pos, S) {
					return true
				}
			}
			// Ensure we don't block of any diagonals (NE and NW don't seem very important)
			if paper.Con[pos-w+1] != S|W && paper.Con[pos-w-1] != S|E {
				return tryConnection(paper, pos, E)
			}
		// NW
		case N | W:
			// Check if the 'by others implied' turn is actually allowed
			// We don't need to check the source connection here like in N|E
			if paper.Con[pos-w-1] == (N|W) || paper.source[pos-w-1] {
				return chooseConnection(paper, paper.next[pos])
			}
		// NE or NS
		case N:
			// Check that we are either extending a corner or starting at a non-occupied source
			if paper.Con[pos-w+1] == N|E || paper.source[pos-w+1] && paper.Con[pos-w+1]&(N|E) != 0 {
				if tryConnection(paper, pos, E) {
					return true
				}
			}
			// Ensure we don't block of any diagonals
			if paper.Con[pos-w+1] != S|W && paper.Con[pos-w-1] != S|E && checkImplicitSE(paper, pos) {
				return tryConnection(paper, pos, S)
			}
		}
	}
	return false
}

// Check that a SW line of corners, starting at pos, will not intersect a SE or NW line
func checkSWLane(paper *Paper, pos int) bool {
	for ; !paper.source[pos]; pos += paper.Width - 1 {
		// Con = 0 means we are crossing a SE line, N|W means a NW
		if paper.Con[pos] != W {
			return false
		}
	}
	return true
}

// Check that a south connection at pos won't create a forced, illegal SE corner at pos+1
// Somethine like: │└
//                 │   <-- Forced SE corner
func checkImplicitSE(paper *Paper, pos int) bool {
	return !(paper.Con[pos+1] == 0) || paper.canSE[pos+1] || paper.Table[pos+1] != EMPTY
}

func tryConnection(paper *Paper, pos1 int, dirs int) bool {
	// Extract the (last) bit which we will process in this call
	dir := dirs & -dirs
	pos2 := pos1 + paper.Vctr[dir]
	end1, end2 := paper.end[pos1], paper.end[pos2]

	// Cannot connect out of the paper
	if paper.Table[pos2] == GRASS {
		return false
	}
	// Check different sources aren't connected
	if paper.Table[end1] != EMPTY && paper.Table[end2] != EMPTY &&
		paper.Table[end1] != paper.Table[end2] {
		return false
	}
	// No loops
	if end1 == pos2 && end2 == pos1 {
		return false
	}
	// No tight corners (Just an optimization)
	if paper.Con[pos1] != 0 {
		dir2 := paper.Con[pos1+paper.Vctr[paper.Con[pos1]]]
		dir3 := paper.Con[pos1] | dir
		if DIAG[dir2] && DIAG[dir3] && dir2&dir3 != 0 {
			return false
		}
	}

	// Add the connection and a backwards connection from pos2
	old1, old2 := paper.Con[pos1], paper.Con[pos2]
	paper.Con[pos1] |= dir
	paper.Con[pos2] |= MIR[dir]
	// Change states of ends to connect pos1 and pos2
	old3, old4 := paper.end[end1], paper.end[end2]
	paper.end[end1] = end2
	paper.end[end2] = end1

	// Remove the done bit and recurse if nessecary
	dir2 := dirs &^ dir
	res := false
	if dir2 == 0 {
		res = chooseConnection(paper, paper.next[pos1])
	} else {
		res = tryConnection(paper, pos1, dir2)
	}

	// Recreate the state, but not if a solution was found,
	// since we'll let it bubble all the way to the caller
	if !res {
		paper.Con[pos1] = old1
		paper.Con[pos2] = old2
		paper.end[end1] = old3
		paper.end[end2] = old4
	}

	return res
}

// As it turns out, though our algorithm avoids must self-touching flows, it
// can be tricked to allow some. Hence we need this validation to filter out
// the false positives
func (paper *Paper) validate() bool {
	w, h := paper.Width, paper.Height
	vtable := make([]rune, w*h)
	for pos := 0; pos < w*h; pos++ {
		if paper.source[pos] {
			// Run throw the flow
			alpha := paper.Table[pos]
			p, old, next := pos, pos, pos
			for {
				// Mark our path as we go
				vtable[p] = alpha
				for _, dir := range DIRS {
					cand := p + paper.Vctr[dir]
					if paper.Con[p]&dir != 0 {
						if cand != old {
							next = cand
						}
					} else if vtable[cand] == alpha {
						// We aren't connected, but it has our color,
						// this is exactly what we doesn't want.
						return false
					}
				}
				// We have reached the end
				if old != p && paper.source[p] {
					break
				}
				old, p = p, next
			}
		}
	}
	return true
}

func (paper *Paper) initTables() {
	w, h := paper.Width, paper.Height

	// Direction vector table
	for dir := 0; dir < 16; dir++ {
		if dir&N != 0 {
			paper.Vctr[dir] += -w
		}
		if dir&E != 0 {
			paper.Vctr[dir] += 1
		}
		if dir&S != 0 {
			paper.Vctr[dir] += w
		}
		if dir&W != 0 {
			paper.Vctr[dir] -= 1
		}
	}

	// Positions of the four corners inside the grass
	paper.Crnr[N|W] = w + 1
	paper.Crnr[N|E] = 2*w - 2
	paper.Crnr[S|E] = h*w - w - 2
	paper.Crnr[S|W] = h*w - 2*w + 1

	// Table to easily check if a position is a source
	paper.source = make([]bool, w*h)
	for pos := 0; pos < w*h; pos++ {
		paper.source[pos] = paper.Table[pos] != EMPTY && paper.Table[pos] != GRASS
	}

	// Pivot tables
	paper.canSE = make([]bool, w*h)
	paper.canSW = make([]bool, w*h)
	for pos := range paper.Table {
		if paper.source[pos] {
			d := paper.Vctr[N|W]
			for p := pos + d; paper.Table[p] == EMPTY; p += d {
				paper.canSE[p] = true
			}
			d = paper.Vctr[N|E]
			for p := pos + d; paper.Table[p] == EMPTY; p += d {
				paper.canSW[p] = true
			}
		}
	}

	// Diagonal 'next' table
	paper.next = make([]int, w*h)
	last := 0
	for _, pos := range append(
		xrange(paper.Crnr[N|W], paper.Crnr[N|E], 1),
		xrange(paper.Crnr[N|E], paper.Crnr[S|E]+1, w)...) {
		for paper.Table[pos] != GRASS {
			paper.next[last] = pos
			last = pos
			pos = pos + w - 1
		}
	}

	// 'Where is the other end' table
	paper.end = make([]int, w*h)
	for pos := 0; pos < w*h; pos++ {
		paper.end[pos] = pos
	}

	// Connection table
	paper.Con = make([]int, w*h)
}

// Makes a slice of the interval [i, i+step, i+2step, ..., j)
func xrange(i int, j int, step int) []int {
	slice := make([]int, 0, (j-i+step-1)/step)
	for i < j {
		slice = append(slice, i)
		i += step
	}
	return slice
}
