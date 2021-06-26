package main

import (
	"embed"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"runtime/pprof"

	"github.com/Nv7-Github/pptx"
	"github.com/Nv7-Github/vidego"
	arg "github.com/alexflint/go-arg"
	"github.com/nfnt/resize"
)

type Args struct {
	Input  string `help:"input video file" arg:"-i"`
	Output string `help:"output video file" arg:"-o"`
	Fast   bool   `help:"use no compression, leads to bigger output but faster encoding" arg:"-f"`

	Cpuprof string `help:"cpu benchmark profile"`
	Memprof string `help:"memory benchmark profile"`
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

	// Profiling
	if args.Cpuprof != "" {
		cpufile, err := os.Create(args.Cpuprof)
		handle(err)
		err = pprof.StartCPUProfile(cpufile)
		handle(err)
		defer pprof.StopCPUProfile()
	}

	// Check if file exists
	_, err := os.Stat(args.Input)
	if os.IsNotExist(err) {
		handle(fmt.Errorf("file does not exist"))
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
	compLevel := png.DefaultCompression
	if args.Fast {
		compLevel = png.NoCompression
	}
	out, err := pptx.Open(args.Output, compLevel)
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
			im = resize.Resize(width, height, im, resize.Bilinear)

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

	// Memprof
	if args.Memprof != "" {
		memprof, err := os.Create(args.Memprof)
		handle(err)
		err = pprof.WriteHeapProfile(memprof)
		handle(err)
	}

	fmt.Println("Saving...")
	err = out.Close()
	handle(err)
	fmt.Println("Done!")
}
