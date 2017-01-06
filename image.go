package main

import (
	"os"
	"image/color"
	"image/gif"
)

type GifImg struct {
	gif.GIF
}
// Load image
func (gifImg *GifImg) Load(filename string) (*gif.GIF, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, err := gif.DecodeAll(file)
	if err != nil {
		panic(err)
	}
	return img, err
}

// Calculates the average RGB color within the given
// rectangle, and returns the [0,1] range of each component.
func (gifImg *GifImg) CellAvgRGB(img *gif.GIF, dominantColor color.RGBA, startX, startY, endX, endY, index int) uint16 {
	var total = [3]uint32{}
	var count uint32

	for x := startX; x < endX; x++ {
		for y := startY; y < endY; y++ {
			gf := img.Image[index]
			r,g,b,_ := gf.At(x,y).RGBA()
			rd,gd,bd,_ := dominantColor.RGBA()
			// remove background color
			if rd == r && gd == g && bd == b {
				r, g, b = 0x00, 0x00, 0x00
			}
			// reduce color range to fit in range [0,15]
			total[0] += r >> 8
			total[1] += g >> 8
			total[2] += b >> 8
			count++
		}
	}
	r := total[0] / count
	g := total[1] / count
	b := total[2] / count

	// Converts a 32-bit RGB color into a term256 compatible approximation.
	rTerm := (((uint16(r) * 5) + 127) / 255) * 36
	gTerm := (((uint16(g) * 5) + 127) / 255) * 6
	bTerm := (((uint16(b) * 5) + 127) / 255)

	return rTerm + gTerm + bTerm + 16 + 1
}

// Get the most dominant color in the image
func (gifImg *GifImg) GetDominantColor(img *gif.GIF) color.RGBA {
	imgWidth, imgHeight := img.Config.Width, img.Config.Height
	firstFrame := img.Image[0]
	histogram := make(map[uint32][]color.RGBA)

	for x := 0; x < imgWidth; x++ {
		for y := 0; y < imgHeight; y++ {
			r,g,b,a := firstFrame.At(x,y).RGBA()
			// get the value from the RGBA
			r /= 0xff
			g /= 0xff
			b /= 0xff
			a /= 0xff
			pixVal := uint32(r)<<24 | uint32(g)<<16 | uint32(b)<<8 | uint32(a)
			// Add the pixel color from the color range to the histogram map, which index is the pixel color converted to uint32.
			// This way we will store all the identical pixels to the same indexed entry.
			histogram[pixVal] = append(histogram[pixVal], color.RGBA{uint8(r),uint8(g),uint8(b),uint8(a)})
		}
	}

	var maxVal uint32
	var dominantColor color.RGBA
	// Find which uint32 converted color occurs mostly in the color range
	// We lookup for the length of histogram map indexes
	for pix, _ := range histogram {
		colorRange := len(histogram[pix])
		if uint32(colorRange) > maxVal {
			maxVal = uint32(colorRange)
			// get the first color from the color range
			dominantColor = histogram[pix][0]
		}
	}
	return dominantColor
}

// Scale the generated image to fit between terminal width & height
func (gifImg *GifImg) Scale(imgWidth, imgHeight, termWidth, termHeight int, ratio float64) (float64, float64) {
	width := float64(imgWidth) / float64(termWidth)
	height := float64(imgHeight) / (float64(termHeight) * ratio)

	/*fmt.Println(width)
	fmt.Println(height)
	fmt.Println(termWidth)
	fmt.Println(termHeight)
	fmt.Println(imgWidth)
	fmt.Println(imgHeight)
	os.Exit(1)*/
	maxValX := maxValue(1, width, height)
	maxValY := maxValue(2, width, height)

	if float64(imgWidth) / float64(imgHeight) > width + height {
		maxValY = maxValue(0, width, height)
	}
	return maxValX, maxValY
}

// Set terminal cell's dimension
func (gifImg *GifImg) CellSize(x, y int, scaleX, scaleY, ratio float64) (int, int, int, int){
	startX, startY := float64(x) * scaleX, float64(y) * scaleY
	endX, endY := startX + scaleX, startY + scaleY * ratio
	return int(startX), int(startY), int(endX), int(endY)
}

func maxValue(values ...float64) float64 {
	var max float64
	for _, val := range values {
		if val > max {
			max = val
		}
	}
	return max
}