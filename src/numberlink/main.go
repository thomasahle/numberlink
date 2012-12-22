package main

import "fmt"
import "flag"
import "os"
import "log"
import "runtime/pprof"

var (
	colorFlag    = flag.Bool("color", false, "Make the output more readable with colors")
	tubesFlag    = flag.Bool("tubes", false, "Draw lines between sources")
	countFlag    = flag.Bool("count", false, "Count number of nodes visited")
	profileFlag  = flag.String("profile", "", "Write profiling data to file")
)

func main() {
	flag.Parse()

	if *profileFlag != "" {
		f, err := os.Create(*profileFlag)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var w, h int
	for n, _ := fmt.Scan(&w, &h); n != 0; n, _ = fmt.Scan(&w, &h) {
		lines := make([]string, 0, w*h)
		for i := 0; i < h; i++ {
			var line string
			fmt.Scan(&line)
			lines = append(lines, line)
		}
		for i := 0; i < 1; i++ {
			p, err := Parse(w, h, lines)
			if err != nil {
				fmt.Println(err.Error())
				break
			}

			if Solve(p) {
				fmt.Println("Found a solution!")
				switch {
				case *tubesFlag:
					PrintTubes(p, *colorFlag)
				default:
					PrintSimple(p, *colorFlag)
				}
			} else {
				fmt.Println("No solutions")
			}
			if *countFlag {
				fmt.Printf("Called %d times\n", calls)
				calls = 0
			}
			fmt.Println()
		}
	}
}

