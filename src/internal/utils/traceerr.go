package utils

import (
	"fmt"
	"os"
	"strings"
)

type Position struct {
	Line   int
	Column int
}

// TraceErr prints formatted error message and points it's origin in the file, then calls os.Exit(1).
func TraceErr(filepath string, pos Position, format string, v ...any) {
	msg := fmt.Sprintf(format, v...)

	data, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("%s:%d:%d: %s\n", filepath, pos.Line, pos.Column+1, msg)
		return
	}

	lines := strings.Split(string(data), "\n")

	const contextLines = 2

	start := max(0, pos.Line-1-contextLines)
	end := min(len(lines), pos.Line+contextLines)

	fmt.Printf("%s:%d:%d: %s\n\n",
		filepath,
		pos.Line,
		pos.Column+1,
		msg,
	)

	for i := start; i < end; i++ {
		lineNumber := i + 1

		fmt.Printf("%4d | %s\n", lineNumber, lines[i])

		if lineNumber == pos.Line {
			fmt.Printf("     | %s^\n", strings.Repeat(" ", pos.Column))
		}
	}

	os.Exit(1)
}
