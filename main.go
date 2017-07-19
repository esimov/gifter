package main

import (
	"os"
	"flag"
	"math/rand"
	"time"
	"os/signal"
	"syscall"
	"fmt"
	"sync"
	"image/color"
	"image/gif"
	"unicode/utf8"
	"github.com/nsf/termbox-go"
	"math"
)

var (
	wg		sync.WaitGroup
	gifImg		*GifImg
	terminal	*Terminal
	termWidth  	int 	= Window.Width
	termHeight 	int 	= Window.Height
	ratio		float64 = Window.Ratio

	// Flags
	outputFile	string
	blockCode 	string
	removeBg 	bool
	frameRates 	int
	count		uint64

	commands 	flag.FlagSet
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	signalChan := make(chan os.Signal, 2)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	commands = *flag.NewFlagSet("commands", flag.ExitOnError)
	commands.BoolVar(&removeBg, "removeBg", false, "Remove GIF file background color")
	commands.StringVar(&outputFile, "output", "output.gif", "Create new GIF file with the background color removed")
	commands.StringVar(&blockCode, "block", "▄", "Select unicode character as cell block")
	commands.IntVar(&frameRates, "frame", 120, "Frame rates")
	commands.Uint64Var(&count, "count", math.MaxUint64, "Loop count")

	if len(os.Args) <= 1 {
		fmt.Println("Please provide a GIF image, or type --help for the supported command line arguments\n")
		fmt.Println("Terminate GIF animation by pressing <ESC> or 'q'.\n")
		os.Exit(1)
	}

	if (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Println(`
Usage of commands:
  -block string
    	Select unicode character as cell block (default "▄")
  -count uint
    	Loop count (default 18446744073709551615)
  -frame int
    	Frame rates (default 120)
  -output string
    	Create new GIF file with the background color removed (default "output.gif")
  -removeBg
    	Remove GIF file background color
		`)
		os.Exit(1)
	}

	commands.Parse(os.Args[2:])
	terminal = &Terminal{}
	//terminal.Flush()

	img := loadGif(os.Args[1])
	gifImg = &GifImg{}

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetOutputMode(termbox.Output256)

	if commands.Parsed() && removeBg {
		dominantColor := gifImg.GetDominantColor(img)
		for idx := 0; idx < len(img.Image); idx++ {
			for x := 0; x < img.Config.Width; x++ {
				for y := 0; y < img.Config.Height; y++ {
					gf := img.Image[idx]
					r,g,b,a := gf.At(x,y).RGBA()
					rd,gd,bd,_ := dominantColor.RGBA()
					// remove background color
					if rd == r && gd == g && bd == b {
						r, g, b = 0x00, 0x00, 0x00
					}
					gf.Set(x,y, color.NRGBA{uint8(r),uint8(g),uint8(b),uint8(a)})
				}
			}
		}
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()

		// Write out the data into the new GIF file
		err = gif.EncodeAll(file, img)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		img = loadGif(outputFile)
	}
	// Render the gif image
	draw(img)
}

func loadGif(fileName string) *gif.GIF {
	img, err := gifImg.Load(fileName)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	return img
}

// Render gif on terminal window
func draw(img *gif.GIF) {
	var startX, startY, endX, endY int
	var loopCount uint64

	ticker := time.Tick(time.Millisecond * time.Duration(frameRates))
	imgWidth, imgHeight := img.Config.Width, img.Config.Height
	scaleX, scaleY := gifImg.Scale(imgWidth, imgHeight, termWidth, termHeight, ratio)
	dominantColor := gifImg.GetDominantColor(img)

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	// This where the magic happens
	loop:
	for {
		if loopCount >= count {
			os.Remove(outputFile)
			break loop
		}
		for idx := 0; idx < len(img.Image); idx++ {
			select {
			case ev := <-eventQueue:
				switch ev.Type {
				case termbox.EventKey :
					if ev.Ch == 'q' || ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC || ev.Key == termbox.KeyCtrlD {
						os.Remove(outputFile)
						break loop
					}
				case termbox.EventResize :
					//break loop
				}
			default:
				wg.Add(1)
				<-ticker
				go func(idx int) {
					defer wg.Done()
					for x := 0; x < img.Config.Width; x++ {
						for y := 0; y < img.Config.Height; y++ {
							if img.Config.Width <= img.Config.Height {
								startX, startY, endX, endY = gifImg.CellSize(x, y, scaleX, scaleY*ratio, ratio)
							} else {
								startX, startY, endX, endY = gifImg.CellSize(x, y, scaleX, scaleY, ratio)
							}
							col := gifImg.CellAvgRGB(img, dominantColor, startX, startY, endX, (startY+endY)/2, idx)
							colorUp := termbox.Attribute(col)

							col = gifImg.CellAvgRGB(img, dominantColor, startX, (startY+endY)/2, endX, endY, idx)
							colorDown := termbox.Attribute(col)

							r, _ := utf8.DecodeRuneInString(blockCode)
							if commands.Parsed() {
								termbox.SetCell(x, y, r, colorDown, colorUp)
							} else {
								termbox.SetCell(x, y, '▄', colorDown, colorUp)
							}
						}
					}
					termbox.Flush()
				}(idx)
			}
			wg.Wait()
			time.Sleep(10 * time.Millisecond)
		}
		loopCount++
	}
}