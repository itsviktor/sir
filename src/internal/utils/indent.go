package utils

import (
	"fmt"
	"strings"
)

func IndentPrintf(indent int, format string, v ...any) {
	fmt.Printf("%s%s", strings.Repeat(" ", indent), fmt.Sprintf(format, v...))
}
