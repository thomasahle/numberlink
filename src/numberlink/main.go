package main

import "fmt"
import "flag"
import "os"
import "runtime/pprof"
import "strings"
import "strconv"

var (
	colorsFlag   = flag.Bool("colors", false, "Make the output more readable with colors")
	tubesFlag    = flag.Bool("tubes", false, "Draw lines between sources")
	countFlag    = flag.Bool("count", false, "Count number of nodes visited")
	profileFlag  = flag.String("profile", "", "Write profiling data to file")
	generateFlag = flag.String("generate", "", "Generate a puzzle of a certain size. Usage: --generate=5x5")
)

func main() {
	flag.Parse()

	if *profileFlag != "" {
		f, err := os.Create(*profileFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *generateFlag != "" {
		size := strings.Split(*generateFlag, "x")
		if len(size) != 2 {
			fmt.Fprintf(os.Stderr, "Error: Must have exactly two arguments to --generate")
			os.Exit(1)
		}
		width, err1 := strconv.Atoi(size[0])
		height, err2 := strconv.Atoi(size[1])
		if err1 != nil || err2 != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to parse arguments to --generate")
			os.Exit(1)
		}
		err := Generate(width, height)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		return
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
					PrintTubes(p, *colorsFlag)
				default:
					PrintSimple(p, *colorsFlag)
				}
			} else {
				fmt.Println("No solutions")
			}
			if *countFlag {
				fmt.Printf("Called %d times\n", Calls)
				Calls = 0
			}
			fmt.Println()
		}
	}
}
