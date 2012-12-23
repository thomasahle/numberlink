package main

import "container/list"
import "fmt"

const (
	RESET = "\x1b[0m"
	BOLD  = "\x1b[1m"

	BLACK   = "\x1b[30m"
	RED     = "\x1b[31m"
	GREEN   = "\x1b[32m"
	YELLOW  = "\x1b[33m"
	BLUE    = "\x1b[34m"
	MAGENTA = "\x1b[35m"
	CYAN    = "\x1b[36m"
	WHITE   = "\x1b[37m"
)

var (
	TUBE = [16]rune{' ', '╵', '╶', '└', '╷', '│', '┌', '├', '╴', '┘', '─', '┴', '┐', '┤', '┬', '┼'}
)

func PrintSimple(paper *Paper, color bool) {
	colors := makeColorTable(paper, !color)
	table := fillTable(paper)
	fmt.Println(paper.Width-2, paper.Height-2)
	for y := 1; y < paper.Height-1; y++ {
		for x := 1; x < paper.Width-1; x++ {
			pos := y*paper.Width + x
			if col := colors[pos]; col == "" {
				fmt.Printf("%c", table[pos])
			} else {
				fmt.Printf("%s%c%s", col, table[pos], RESET)
			}
		}
		fmt.Println()
	}
}

func PrintTubes(paper *Paper, color bool) {
	colors := makeColorTable(paper, !color)
	for y := 1; y < paper.Height-1; y++ {
		for x := 1; x < paper.Width-1; x++ {
			pos := y*paper.Width + x
			val := paper.Table[pos]
			var c rune
			if val == EMPTY {
				c = TUBE[paper.Con[pos]]
			} else {
				c = val
			}
			if col := colors[pos]; col == "" {
				fmt.Printf("%c", c)
			} else {
				fmt.Printf("%s%c%s", col, c, RESET)
			}
		}
		fmt.Println()
	}
}

// Assigns a terminal color code to every position on the paper
func makeColorTable(paper *Paper, empty bool) []string {
	color := make([]string, paper.Width*paper.Height)
	if !empty {
		table := fillTable(paper)

		next := 0
		available := []string{RED, GREEN, YELLOW, BLUE, MAGENTA, CYAN, BLACK, WHITE}
		for _, c := range available {
			available = append(available, BOLD+c)
		}
		var mapping = make(map[rune]string)

		for y := 1; y < paper.Height-1; y++ {
			for x := 1; x < paper.Width-1; x++ {
				c := table[y*paper.Width+x]
				if _, found := mapping[c]; !found {
					if len(available) >= 1 {
						mapping[c] = available[next]
						next = (next + 1) % len(available)
					} else {
						mapping[c] = BLACK
					}
				}
				color[y*paper.Width+x] = mapping[c]
			}
		}
	}
	return color
}

// Does a bfs search on every source, filling out its connected nodes
func fillTable(paper *Paper) []rune {
	w, h := paper.Width, paper.Height
	table := make([]rune, w*h)
	copy(table, paper.Table)
	for pos := 0; pos < w*h; pos++ {
		if paper.source[pos] {
			queue := list.New()
			queue.PushBack(pos)
			for queue.Len() != 0 {
				pos := queue.Remove(queue.Front()).(int)
				paint := table[pos]
				for _, dir := range DIRS {
					next := pos + paper.Vctr[dir]
					if paper.Con[pos]&dir != 0 && table[next] == EMPTY {
						table[next] = paint
						queue.PushBack(next)
					}
				}
			}
		}
	}
	return table
}
