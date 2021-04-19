package main

import (
	_ "embed"
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
	"gocv.io/x/gocv"
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

func execute(file string, out string, precision int) {
	var kind gocv.MatType
	switch precision {
	case 0:
		kind = gocv.MatTypeCV32FC3
	case 1:
		kind = gocv.MatTypeCV8UC3
	case 2:
		kind = gocv.MatTypeCV16UC3
	case 3:
		kind = gocv.MatTypeCV32FC3
	case 4:
		kind = gocv.MatTypeCV64FC4
	default:
		fmt.Println("Precision must be from 1-4.")
		return
	}

	video, err := gocv.VideoCaptureFile(file)
	handle(err)
	defer video.Close()

	img := gocv.NewMat()
	video.Read(&img)
	img.ConvertTo(&img, kind)

	writer, err := gocv.VideoWriterFile(out, video.CodecString(), video.Get(gocv.VideoCaptureFPS), img.Cols(), img.Rows(), true)
	handle(err)
	defer writer.Close()

	output := img.Clone()
	writable := gocv.NewMat()

	var time float64 = 1
	bar := newBar(int64(video.Get(gocv.VideoCaptureFrameCount)))
	for video.Read(&img) {
		if img.Empty() {
			continue
		}
		img.ConvertTo(&img, kind)

		gocv.AddWeighted(output, (time-1)/time, img, 1/time, 0, &output)

		output.ConvertTo(&writable, gocv.MatTypeCV8UC3)
		err = writer.Write(writable)
		handle(err)

		time++
		err = bar.Add(1)
		handle(err)
	}
	bar.Finish()
}
