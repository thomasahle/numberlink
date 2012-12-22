package main

import "fmt"

type ParseError struct {
	Line    int
	Problem string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("ParseError: '%s' at line %d", e.Problem, e.Line)
}

func Parse(width int, height int, lines []string) (*Paper, *ParseError) {
	if width*height == 0 {
		return nil, &ParseError{0, "width and height cannot be 0"}
	}
	if width != len(lines[0]) || height != len(lines) {
		return nil, &ParseError{1, "width and height must match puzzle size"}
	}

	table := make([]int32, 0, width*height)
	for _, line := range lines {
		for _, c := range line {
			table = append(table, c)
		}
	}

	return NewPaper(width, height, table), nil
}
