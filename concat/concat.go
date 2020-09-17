package main

import (
	"bufio" // For Y/n
	"os"    // For reading from stdin

	"gocv.io/x/gocv" // GoCV for OpenCV, need to install OpenCV

	"fmt"   // Printing
	"image" // For dimensions in resizing

	"github.com/alexflint/go-arg"       // Arguments
	"github.com/schollz/progressbar/v3" // Progress bar
)

func main() {
	var Args struct {
		Input  []string `arg:"-i" arg:"required" arg:"positional" help:"Input video files"`
		Output string   `arg:"-o" arg:"required" arg:"positional" help:"Output video file"`
	}
	p := arg.MustParse(&Args) // Parse arguments

	if len(Args.Input) < 1 { // Is there input?
		p.Fail("no input files specified")
	}
	if len(Args.Input) < 1 { // Is there output?
		p.Fail("no output files specified")
	}

	args := Args.Input[1:]                            // last item is output file, those are excluded
	vidcap, _ := gocv.VideoCaptureFile(Args.Input[0]) // Read first file, to find out dimensions and FPS, assuming that the video is going to be like the first one

	// Get dimensions
	width := int(vidcap.Get(gocv.VideoCaptureFrameWidth))
	height := int(vidcap.Get(gocv.VideoCaptureFrameHeight))

	// Check if file exists, if it does ask Y/n
	if Exists(Args.Output) {
		fmt.Print("File '" + Args.Output + "' already exists. Overwrite? (Y/n): ")
		reader := bufio.NewReader(os.Stdin)
		ans, _ := reader.ReadString('\n')
		if !((ans == "Y\n") || (ans == "y\n")) {
			return
		}
	}

	// Create Video File Writer
	out, _ := gocv.VideoWriterFile(Args.Output, "mp4v", vidcap.Get(gocv.VideoCaptureFPS), width, height, true) // Name of file, FOURCC Codec, FPS, dimensions, isColor=True

	// Call Writing function on first video, since that was read before and not in loop
	fmt.Println("Reading & Writing '" + Args.Input[0] + "'")
	WriteToVid(vidcap, out, width, height)

	// Loop through args for all other video
	for _, vidname := range args {
		vidcap, _ = gocv.VideoCaptureFile(vidname)
		fmt.Println("Reading & Writing '" + vidname + "'")
		WriteToVid(vidcap, out, width, height) // Call Writing Function on video
	}

	out.Close() // Release output file
}

// WriteToVid is for writing videocapture to videowriterfile
func WriteToVid(invid *gocv.VideoCapture, outvid *gocv.VideoWriter, w int, h int) { // * means pointer
	frame := gocv.NewMat()       // Create empty Mat to write to
	opened := invid.Read(&frame) // Pass in pointer to frame so that picture is read to frame. The frame part is useless in this case though, this is just to initialize the Opened variable
	size := image.Pt(w, h)       // Create a Image.Pt object out of the dimensions, for resizing later

	bar := progressbar.Default(int64(invid.Get(gocv.VideoCaptureFrameCount))) // Crete Progress Bar

	for opened { // For loops can also act as while loops in Go, and in this case it is. WThere are no while loops in Go.
		if (int(frame.Cols()) == w) && (int(frame.Rows()) == h) { // Check if frame is right size
			outvid.Write(frame)         // Write frame if already right size
			opened = invid.Read(&frame) // Read frame, to frame pointer
			bar.Add(1)                  // Increment bar
		} else {
			gocv.Resize(frame, &frame, size, 0, 0, gocv.InterpolationLinear) // Otherwise resize frame this loop. Parameters: input mat, output mat (this one is pointer), size (scaling to exact size, not aspect ratio so using this), scaling amount (since not keeping aspect ratio does not matter, put at zero), Interpolation mode
		}
	}

	invid.Close() // Release input file
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
