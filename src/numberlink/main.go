package main

import "fmt"
import "flag"
import "os"
import "runtime/pprof"
import "strings"
import "strconv"
import "bufio"

var (
	colorsFlag   = flag.Bool("colors", false, "Make the output more readable with colors")
	tubesFlag    = flag.Bool("tubes", false, "Draw lines between sources")
	callsFlag    = flag.Bool("calls", false, "Count number of recursive calls")
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

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			break
		}
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}
		parts := strings.Split(line, " ")
		bad := len(parts) != 2
		var w, h int
		if !bad {
			var err1, err2 error
			w, err1 = strconv.Atoi(parts[0])
			h, err2 = strconv.Atoi(parts[1])
			bad = bad || err1 != nil || err2 != nil
		}
		if bad {
			fmt.Fprintf(os.Stderr, "Error: Expected 'width height' got '%s'", line)
			os.Exit(1)
		}

		// We use 0 0 as a visible eof mark
		if w == 0 && h == 0 {
			break
		}
		lines := make([]string, 0, w*h)
		for i := 0; i < h; i++ {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
				os.Exit(1)
			}
			line = strings.TrimSpace(line)
			lines = append(lines, line)
		}

		// Done parsing stuff, time for the fun part
		for i := 0; i < 1; i++ {
			p, err := Parse(w, h, lines)
			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
				os.Exit(1)
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

			if *callsFlag {
				fmt.Printf("Called %d times\n", Calls)
				Calls = 0
			}
			fmt.Println()
		}
	}
}
