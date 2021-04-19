package main

import (
	"flag"
	"fmt"
	"os"
)

func handle(err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		flag.Usage()
		os.Exit(1)
	}
}

var input = flag.String("input", "", "Input path for video file.")
var output = flag.String("output", "", "Output path for video file.")
var precision = flag.Int("precision", 3, "A number from 1-4 representing image precision (lower faster but worse quality).")

func init() {
	flag.StringVar(input, "i", "", "Input path for video file.")
	flag.StringVar(output, "o", "", "Output path for video file.")
	flag.IntVar(precision, "p", 3, "A number from 1-4 representing image precision (lower faster but worse quality).")
}
func main() {
	flag.Parse()

	execute(*input, *output, *precision)
}
