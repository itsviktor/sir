package utils

import (
	"fmt"
	"log"
)

var isDebug = false

func Debug(msg string) {
	if isDebug {
		log.Print(msg)
	}
}

func Debugf(format string, v ...any) {
	Debug(fmt.Sprintf(format, v...))
}

func SetDebug(f bool) {
	isDebug = f
}
