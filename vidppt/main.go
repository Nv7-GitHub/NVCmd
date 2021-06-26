package main

import (
	"embed"
	"fmt"
	"image"
	"io"
	"os"

	"github.com/Nv7-Github/pptx"
	"github.com/Nv7-Github/vidego"
	arg "github.com/alexflint/go-arg"
	"github.com/nfnt/resize"
)

type Args struct {
	Input  string `help:"input video file" arg:"-i"`
	Output string `help:"output video file" arg:"-o"`
}

const (
	width  = 960
	height = 720
)

//go:embed blank.pptx
var blankFS embed.FS

var p *arg.Parser

func handle(err error) {
	if err != nil {
		p.Fail(err.Error())
	}
}

func main() {
	var args Args
	p = arg.MustParse(&args)

	// Check if file exists
	_, err := os.Stat(args.Input)
	if os.IsNotExist(err) {
		fmt.Println(args.Input)
		handle(fmt.Errorf("error: file does not exist"))
	}

	// Open video
	vid, err := vidego.NewDecoder(args.Input)
	handle(err)
	defer vid.Free()

	// Create output file with blank PPTX
	blank, err := blankFS.Open("blank.pptx")
	handle(err)
	defer blank.Close()

	outFile, err := os.Create(args.Output)
	handle(err)
	defer outFile.Close()

	io.Copy(outFile, blank)

	// Open as PPTX
	out, err := pptx.Open(args.Output)
	handle(err)

	// Write to file
	var s pptx.Slide

	// Loop
	bar := newBar(int64(vid.FrameCount()))
	var imgs []image.Image
	var cont bool
	for {
		// Decode frames
		cont, imgs, err = vid.GetNextFrame()
		if !cont {
			break
		}
		handle(err)
		if imgs == nil {
			continue
		}

		// Add to ppt
		for _, im := range imgs {
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
	}

	err = bar.Finish()
	handle(err)

	fmt.Println("Saving...")
	err = out.Close()
	handle(err)
	fmt.Println("Done!")
}
