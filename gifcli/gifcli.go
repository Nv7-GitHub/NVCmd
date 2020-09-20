package main

import (
	"fmt"     // Printing
	"math"    // Rounding
	"os"      // For reading from stdin
	"os/exec" // Clearing console

	// Clearing
	"strconv" // Parsing args
	"strings" // Parsing args

	"image"     // For dimensions in resizing
	"image/gif" // Decoding gifs

	"github.com/alexflint/go-arg"       // Arguments
	"github.com/schollz/progressbar/v3" // Progress bar

	"runtime" // Call GC, clearing console
	"time"    // FPS

	"github.com/nfnt/resize" // Image resizing
)

var frames []*string
var w int
var h int
var p *arg.Parser
var framesRendered []chan bool
var images []*image.Paletted
var delays []int

func main() {
	var Args struct {
		Input  string  `arg:"required,positional" help:"Input video file"`
		Scale  float64 `arg:"-s" help:"The amount of scaling (0.5 for half, etc.)"`
		Width  int     `arg:"-w" help:"The new width, keeping aspect ratio"`
		Resize string  `arg:"-r" help:"The new width and height, in the WidthxHeight format"`
	}
	p = arg.MustParse(&Args) // Parse arguments

	// Check if Input file exists
	if !Exists(Args.Input) {
		p.Fail("input file does not exist.")
	}

	// Open video file
	file, err := os.Open(Args.Input)
	if err != nil {
		p.Fail(err.Error())
	}
	defer file.Close()

	// Decode Video File
	gif, err := gif.DecodeAll(file)
	if err != nil {
		p.Fail(err.Error())
	}
	images = gif.Image
	delays = gif.Delay

	// Get video scale

	if Args.Scale != *new(float64) {
		// Scale supplied

		// Get dimensions
		width := images[0].Bounds().Size().X
		height := images[0].Bounds().Size().Y

		// Calculate frame width and height
		w = int(math.Round(float64(width) * Args.Scale))
		h = int(math.Round(float64(height) * Args.Scale))
	} else if Args.Width != *new(int) {
		// New width supplied

		// Get dimensions
		width := images[0].Bounds().Size().X
		height := images[0].Bounds().Size().Y

		// Get new height, based on aspect ratio, get width based on supplied value
		w = Args.Width
		h = height * Args.Width / width
	} else if Args.Resize != *new(string) {
		// Resizing

		// Get new dimensions
		dimens := strings.Split(Args.Resize, "x")

		// If wrong input format, error
		if len(dimens) != 2 {
			p.Fail("Invalid resizing dimensions")
		}

		// Set width and height to dimensions
		var err error

		w, err = strconv.Atoi(dimens[0])
		if err != nil {
			p.Fail(err.Error())
		}

		h, err = strconv.Atoi(dimens[1])
		if err != nil {
			p.Fail(err.Error())
		}
	} else {
		p.Fail("no scale specified, use vidcli --scale 1 to play at original resolution")
	}

	// Render to an array of strings
	Render()

	// Play rendered text
	Play()
}

// Exists is for checking if file exists
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Render renders all the frames to text
func Render() {
	fmt.Println("Rendering to text...")

	framesRendered = make([]chan bool, 0)

	bar := progressbar.Default(int64(len(images))) // Create Progress Bar
	for _, img := range images {
		size := img.Bounds().Size() // Get image size
		var frame image.Image
		if (size.X != w) || (size.Y != h) { // Check if frame is right size
			frame = resize.Resize(uint(w), uint(h), img, resize.NearestNeighbor) // Resize with Nearest Neighbor if not right size
		} else {
			frame = img
		}
		frametext := new(string)                                             // Create pointer to string
		framesRendered = append(framesRendered, make(chan bool))             // Add Boolean channel for checking if done
		frames = append(frames, frametext)                                   // Add to frames
		RenderFrame(frame, frametext, framesRendered[len(framesRendered)-1]) // If right, render frame to string
		bar.Add(1)                                                           // Increment Loading bar
	}

	bar = progressbar.Default(int64(len(framesRendered))) // Create Progress Bar
	// Wait for all goroutines to finish
	for _, channel := range framesRendered {
		<-channel  // Wait for channel to finish
		bar.Add(1) // Increment Loading Bar
	}

	fmt.Println("Rendered!")
}

// RenderFrame renders a gocv.Mat to text
func RenderFrame(img image.Image, output *string, finished chan bool) {
	go func() {
		size := img.Bounds().Size()

		(*output) = ""

		for y := 0; y < size.Y; y++ {
			for x := 0; x < size.X; x++ {
				color := img.At(x, y)
				r, g, b, _ := color.RGBA()
				(*output) += fmt.Sprintf("\033[38;2;%d;%d;%dm█\033[0m", r/257, g/257, b/257)
			}
			(*output) += "\n"
		}

		runtime.GC()

		finished <- true
	}()
}

// Play displays the frames
func Play() {
	clears := make(map[string]func(), 0)
	clears["darwin"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clears["linux"] = clears["darwin"]
	clears["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	for i, frame := range frames {
		clears[runtime.GOOS]()
		fmt.Println(*frame)
		time.Sleep((time.Second / 100) * time.Duration(delays[i]))
	}
}
