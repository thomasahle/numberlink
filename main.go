package main

import "fmt"
import "flag"
import "io"
import "os"
import "runtime/pprof"
import "strings"
import "strconv"
import "bufio"

var (
	colorsFlag    = flag.Bool("colors", false, "Make the output more readable with colors")
	tubesFlag     = flag.Bool("tubes", false, "Draw lines between sources")
	callsFlag     = flag.Bool("calls", false, "Count number of recursive calls")
	callsOnlyFlag = flag.Bool("calls-only", false, "Print only the culminative number of recursive calls")
	profileFlag   = flag.String("profile", "", "Write profiling data to file")
	generateFlag  = flag.String("generate", "", "Generate a puzzle of a certain size. Usage: --generate=5x5")
)

func main() {
	flag.Parse()

	// Generating
	if *generateFlag != "" {
		size := strings.Split(*generateFlag, "x")
		if len(size) != 2 {
			fmt.Fprintf(os.Stderr, "Error: Must have exactly two arguments to --generate\n")
			os.Exit(1)
		}
		width, err1 := strconv.Atoi(size[0])
		height, err2 := strconv.Atoi(size[1])
		if err1 != nil || err2 != nil {
			fmt.Fprintf(os.Stderr, "Error: Unable to parse arguments to --generate\n")
			os.Exit(1)
		}
		pzzl, _, err := Generate(width, height)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		} else {
			fmt.Println(len(pzzl[0]), len(pzzl))
			for _, line := range pzzl {
				fmt.Println(line)
			}
		}
		return
	}

	// Profiling
	if *profileFlag != "" {
		f, err := os.Create(*profileFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// Normal run
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Fprintln(os.Stderr, err.Error())
			}
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
			fmt.Fprintf(os.Stderr, "Error: Expected 'width height' got '%s'\n", line)
			os.Exit(1)
		}

		// We use 0 0 as an end of puzzles mark
		if w == 0 && h == 0 {
			break
		}
		lines := make([]string, 0, w*h)
		for i := 0; i < h; i++ {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Fprintln(os.Stderr, err.Error())
				}
				os.Exit(1)
			}
			line = strings.TrimSpace(line)
			lines = append(lines, line)
		}

		// Done parsing stuff, time for the fun part
		p, err := Parse(w, h, lines)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		res := Solve(p)
		if !*callsOnlyFlag {
			if res {
				switch {
				case *tubesFlag:
					PrintTubes(p, *colorsFlag)
				default:
					PrintSimple(p, *colorsFlag)
				}
			} else {
				fmt.Println("IMPOSSIBLE")
			}

			if *callsFlag {
				fmt.Printf("Called %d times\n", Calls)
				Calls = 0
			}
			fmt.Println()
		}
	}
	if *callsOnlyFlag {
		fmt.Printf("Called %d times\n", Calls)
	}
}
