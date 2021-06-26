package main

import (
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

func newBar(max int64) *progressbar.ProgressBar {
	return progressbar.NewOptions64(
		max,
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetItsString("f"),
	)
}
