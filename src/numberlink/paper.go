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
	// Mirrors directions so W becomes E and SE becomes NW
	MIR = [16]int{0, S, W, S | W, N, N | S, N | W, N | S | W, E, E | S, E | W, E | S | W, N | E, N | E | S, N | E | W, N | E | S | W}
	// Filter for N|E, N|W, S|E and S|W
	DIAG = [16]bool{false, false, false, true, false, false, true, false, false, true, false, false, true, false, false, false}
)

type Paper struct {
	Width  int
	Height int
	Table  []int32
	Con    []int

	source []bool
	end    []int
	canSE  []bool
	canSW  []bool
	next   []int
	hist   []entry

	Vctr [16]int
	Crnr [16]int
}

type entry struct {
	ptr []int
	key int
	val int
}

func NewPaper(width, height int, table []rune) *Paper {
	paper := new(Paper)

	// Pad the given table with #, to make boundery checks easier
	paper.Width, paper.Height = width+2, height+2
	paper.Table = make([]rune, 0, paper.Width*paper.Height)
	for i := 0; i < paper.Width; i++ {
		paper.Table = append(paper.Table, GRASS)
	}
	for y := 1; y < paper.Height-1; y++ {
		paper.Table = append(paper.Table, GRASS)
		for x := 1; x < paper.Width-1; x++ {
			paper.Table = append(paper.Table, table[(y-1)*(paper.Width-2)+(x-1)])
		}
		paper.Table = append(paper.Table, GRASS)
	}
	for i := 0; i < paper.Width; i++ {
		paper.Table = append(paper.Table, GRASS)
	}

	initTables(paper)

	return paper
}

func Solve(paper *Paper) bool {
	return solve(paper, paper.Width+1)
}

var Calls = 0

func solve(paper *Paper, pos int) bool {
	Calls++

	// Final
	if pos == paper.Crnr[S|E] {
		return paper.validate()
	}

	// Connect
	histSize := len(paper.hist)
	res := chooseConnections(paper, pos)
	if !res {
		paper.goBack(histSize)
	}
	return res
}

func chooseConnections(paper *Paper, pos int) bool {
	w := paper.Width
	// Instead of long if-else chains: http://tour.golang.org/#43
	histSize := len(paper.hist)
	if paper.source[pos] {
		switch paper.Con[pos] {
		// If the source is not yet connection
		case 0:
			// It will be a long time before we come back to this row, so do an extra check
			if paper.Con[pos-w+1] != S|W {
				if paper.connect(pos, E) && solve(paper, paper.next[pos]) {
					return true
				}
				paper.goBack(histSize)
			}
			// South connections can create a forced SE position
			if checkImplicitSE(paper, pos) {
				if paper.connect(pos, S) && solve(paper, paper.next[pos]) {
					return true
				}
			}
		// If the source is already connected
		case N, W:
			return solve(paper, paper.next[pos])
		}
	} else {
		switch paper.Con[pos] {
		// SE
		case 0:
			// Should we check for implied N|W?
			if paper.canSE[pos] {
				return paper.connect(pos, E) && paper.connect(pos, S) && solve(paper, paper.next[pos])
			}
		// SW or WE
		case W:
			// Check there is a free line down to the source we are turning around
			if paper.canSW[pos] && checkSWLane(paper, pos) && checkImplicitSE(paper, pos) {
				if paper.connect(pos, S) && solve(paper, paper.next[pos]) {
					return true
				}
				paper.goBack(histSize)
			}
   			// Ensure we don't block of any diagonals (NE and NW don't seem very important)
			if paper.Con[pos-w+1] != S|W && paper.Con[pos-w-1] != S|E {
				return paper.connect(pos, E) && solve(paper, paper.next[pos])
			}
		// NW
		case N | W:
			// Check if the 'by others implied' turn is actually allowed
			// We don't need to check the source connection here like in N|E
			if paper.Con[pos-w-1] == (N|W) || paper.source[pos-w-1] {
				return solve(paper, paper.next[pos])
			}
		// NE or NS
		case N:
			// Check that we are either extending a corner or starting at a non-occupied source
			if paper.Con[pos-w+1] == N|E || paper.source[pos-w+1] && paper.Con[pos-w+1]&(N|E) != 0 {
				if paper.connect(pos, E) && solve(paper, paper.next[pos]) {
					return true
				}
				paper.goBack(histSize)
			}
			// Ensure we don't block of any diagonals
			if paper.Con[pos-w+1] != S|W && paper.Con[pos-w-1] != S|E && checkImplicitSE(paper, pos) {
				return paper.connect(pos, S) && solve(paper, paper.next[pos])
			}
		}
	}
	return false
}

func checkSWLane(paper *Paper, pos int) bool {
	for ; !paper.source[pos]; pos += paper.Width-1 {
		// Con = 0 means we are crossing a SE line, N|W means a NW
		if paper.Con[pos] != W {
			return false
		}
	}
	return true
}

func checkImplicitSE(paper *Paper, pos int) bool {
	return !(paper.Con[pos+1] == 0) || paper.canSE[pos+1] || paper.Table[pos+1] != EMPTY
}

// As it turns out, our smart algorithm isn't 100% able to avoid self-touching flows
// Hence we need this validation to filter out the false positives
func (paper *Paper) validate() bool {
	w, h := paper.Width, paper.Height
	vtable := make([]rune, w*h)
	for pos := 0; pos < w*h; pos++ {
		if paper.source[pos] {
			// Run throw the flow
			alpha := paper.Table[pos]
			old, next := pos, pos
			for {
				// Mark our path as we go
				vtable[pos] = alpha
				for _, dir := range DIRS {
					cand := pos + paper.Vctr[dir]
					if paper.Con[pos]&dir != 0 {
						if cand != old {
							next = cand
						}
					} else if vtable[cand] == alpha && !paper.source[cand] {
						// We aren't connected, but it has our color,
						// this is exactly what we doesn't want.
						return false
					}
				}
				// We have reached the end
				if old != pos && paper.source[pos] {
					break
				}
				old, pos = pos, next
			}
		}
	}
	return true
}

func (paper *Paper) setnrem(ptr []int, key, val int) {
	if ptr[key] != val {
		paper.hist = append(paper.hist, entry{ptr, key, ptr[key]})
		ptr[key] = val
	}
}

func (paper *Paper) goBack(histSize int) {
	for i := len(paper.hist) - 1; i >= histSize; i-- {
		entry := paper.hist[i]
		entry.ptr[entry.key] = entry.val
	}
	paper.hist = paper.hist[:histSize]
}

func (paper *Paper) connect(pos1 int, dir int) bool {
	// Asserts dir in {N, S, E, W}
	pos2 := pos1 + paper.Vctr[dir]
	// Cannot connect out of the paper
	if paper.Table[pos2] == GRASS {
		return false
	}
	// No close corners (The value of this optimization is doubtable)
	if paper.Con[pos1] != 0 {
		dir2 := paper.Con[pos1+paper.Vctr[paper.Con[pos1]]]
		dir3 := paper.Con[pos1]|dir
		if DIAG[dir2] && DIAG[dir3] && dir2&dir3 != 0 {
			return false
		}
	}
	// No loops
	end1, end2 := paper.end[pos1], paper.end[pos2]
	if end1 == pos2 && end2 == pos1 {
		return false
	}
	// Check different sources aren't connected
	if paper.Table[end1] != EMPTY && paper.Table[end2] != EMPTY &&
		paper.Table[end1] != paper.Table[end2] {
		return false
	}
	// Add the connection and a backwards connection from pos2
	paper.setnrem(paper.Con, pos1, paper.Con[pos1]|dir)
	paper.setnrem(paper.Con, pos2, paper.Con[pos2]|MIR[dir])
	// Change states of ends to connect pos1 and pos2
	paper.setnrem(paper.end, end1, end2)
	paper.setnrem(paper.end, end2, end1)
	return true
}

func initTables(paper *Paper) {
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

	// Source table
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

	// Diagonal next table
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

	// History
	paper.hist = make([]entry, 0, 4*w*h)
}

func xrange(i int, j int, step int) []int {
	slice := make([]int, 0, (j-i+step-1)/step)
	for i < j {
		slice = append(slice, i)
		i += step
	}
	return slice
}
