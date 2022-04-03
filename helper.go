package main

import (
	"fmt"
	"strconv"
	"time"
)

// From https://golangbyexample.com/print-output-text-color-console/

const ColorReset = "\033[0m"

const ColorRed = "\033[31m"
const ColorGreen = "\033[32m"
const ColorYellow = "\033[33m"
const ColorBlue = "\033[34m"
const ColorCyan = "\033[36m"
const ColorWhite = "\033[0m"

func PrintLogInColor(color string, format string, args ...interface{}) {
	colorMu.Lock()
	fmt.Print(color)
	fmt.Printf(strconv.Itoa(int(time.Now().UnixNano()/int64(time.Millisecond))) /*.Format("15:04:05")*/ +" - "+format, args...)
	fmt.Print(ColorReset)
	colorMu.Unlock()
}
