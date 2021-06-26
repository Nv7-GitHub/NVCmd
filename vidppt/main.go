package main

import (
	"embed"
	"flag"
	"fmt"
	"image"
	"io"
	"os"

	"github.com/Nv7-Github/pptx"
	"github.com/nfnt/resize"
	"gocv.io/x/gocv"
)

var outFilename = *flag.String("output", "out.pptx", "the output PowerPoint file")
var in = *flag.String("input", "", "the input video file")

const (
	width  = 960
	height = 720
)

//go:embed blank.pptx
var blankFS embed.FS

func handle(err error) {
	if err != nil {
		flag.Usage()
		fmt.Println("error: ", err.Error())
		os.Exit(2)
	}
}

func main() {
	// Open video
	vid, err := gocv.VideoCaptureFile(in)
	handle(err)

	// Create output file with blank PPTX
	blank, err := blankFS.Open("blank.pptx")
	handle(err)
	defer blank.Close()

	outFile, err := os.Create(outFilename)
	handle(err)
	defer outFile.Close()

	io.Copy(outFile, blank)

	// Open as PPTX
	out, err := pptx.Open(outFilename)
	handle(err)

	// Write to file
	img := gocv.NewMat()
	var im image.Image
	var s pptx.Slide

	// Loop
	bar := newBar(int64(vid.Get(gocv.VideoCaptureFrameCount)))
	for vid.Read(&img) {
		// Decode image
		if img.Empty() {
			continue
		}
		im, err = img.ToImage()
		handle(err)

		// Resize image
		im = resize.Resize(width, height, im, resize.Bicubic)

		// Save image to PPT
		s = pptx.Slide{
			Images: []pptx.Image{
				{
					X: 0,
					Y: 0,
					Image: pptx.GoImage{
						Image: im,
					},
				},
			},
		}
		err = out.Add(s)
		handle(err)

		err = bar.Add(1)
		handle(err)
	}

	err = out.Close()
	handle(err)
	err = bar.Finish()
	handle(err)
}
