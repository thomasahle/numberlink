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
)

type Paper struct {
	Width  int
	Height int
	Table  []int32
	Con    []int

	source []bool
	end    []int
	canNW  []bool
	canNE  []bool
	mustNE []bool
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

func Solve(paper *Paper) bool {
	return solve(paper, paper.Width+1)
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

var Calls = 0

func solve(paper *Paper, pos int) bool {
	Calls++

	// Final
	if pos == (paper.Height-1)*paper.Width-2 {
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
			// If this connection is bad, the next solve will notice immidiately
			if paper.connect(pos, S) && solve(paper, paper.next[pos]) {
				return true
			}
		// If the source is already connected
		case N, W:
			return solve(paper, paper.next[pos])
		}
	} else {
		// In some cases SW or SE is forced, but we haven't got the needed inputs
		// Rather than putting the test in every wrong input clause, we put them here
		if paper.Con[pos] != W && (pos == 2*w-2 || paper.Con[pos-w+1] == S|W) {
			return false
		}
		if paper.Con[pos] != 0 && (pos == w+1 || paper.Con[pos-w-1] == S|E) {
			return false
		}
		if paper.Con[pos] != N && paper.mustNE[pos] {
			return false
		}

		switch paper.Con[pos] {
		// SE
		case 0:
			if paper.canNW[pos] {
				return paper.connect(pos, S|E) && solve(paper, paper.next[pos])
			}
		// SW or WE
		case W:
			if paper.canNE[pos] {
				if paper.connect(pos, S) && solve(paper, paper.next[pos]) {
					return true
				}
				paper.goBack(histSize)
			}
			// In some cases SW is forced, otherwise try a simple W-E bridge
			if paper.Con[pos-w+1] != S|W && paper.Con[pos-w-1] != S|E {
				if paper.connect(pos, E) && solve(paper, paper.next[pos]) {
					return true
				}
			}
		// NW
		case N | W:
			// Check if the 'by others implied' turn is actually allowed
			taken := pos == 2*w+2 || pos >= 2*w+2 && paper.Con[pos-2*w-2] == S|E
			if (paper.source[pos-w-1] && !taken) || paper.Con[pos-w-1] == (N|W) {
				return solve(paper, paper.next[pos])
			}
		// NE or NS
		case N:
			// 'taken' tests whether the pivot is already active in the opposite direction
			taken := pos == 3*w-3 || pos >= 2*w-2 && paper.Con[pos-2*w+2] == S|W
			if (paper.source[pos-w+1] && !taken) || paper.Con[pos-w+1] == N|E {
				if paper.connect(pos, E) && solve(paper, paper.next[pos]) {
					return true
				}
				paper.goBack(histSize)
			}
			// Prevents blocking off the SW diagonal
			if !paper.mustNE[pos] && paper.Con[pos-w+1] != S|W && paper.Con[pos-w-1] != S|E {
				return paper.connect(pos, S) && solve(paper, paper.next[pos])
			}
		}
	}
	return false
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
	// If connect is called with multiple dir bits set:
	// isolate one and recurse with the rest
	if rest := dir & (dir - 1); rest != 0 {
		if !paper.connect(pos1, rest) {
			return false
		}
		dir = dir & (-dir)
	}
	pos2 := pos1 + paper.Vctr[dir]
	// Cannot connect out of the paper
	if paper.Table[pos2] == GRASS {
		return false
	}
	// Check different sources aren't connected
	end1, end2 := paper.end[pos1], paper.end[pos2]
	if paper.source[end1] && paper.source[end2] &&
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
	paper.canNW = make([]bool, w*h)
	paper.canNE = make([]bool, w*h)
	paper.mustNE = make([]bool, w*h)
	for pos := range paper.Table {
		if paper.source[pos] {
			d := paper.Vctr[N|W]
			for p := pos + d; paper.Table[p] == EMPTY; p += d {
				paper.canNW[p] = true
			}
			d = paper.Vctr[N|E]
			for p := pos + d; paper.Table[p] == EMPTY; p += d {
				paper.canNE[p] = true
			}
		}
	}
	for pos := paper.Crnr[S|W]; paper.Table[pos] == EMPTY; pos += -w + 1 {
		paper.mustNE[pos] = true
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
}

func xrange(i int, j int, step int) []int {
	slice := make([]int, 0, (j-i+step-1)/step)
	for i < j {
		slice = append(slice, i)
		i += step
	}
	return slice
}
