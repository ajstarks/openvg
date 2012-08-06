package main

import (
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
)

func main() {

	var (
		width  = 1920
		height = 1080
	)
	if len(os.Args) > 2 {
		width, _ = strconv.Atoi(os.Args[1])
		height, _ = strconv.Atoi(os.Args[2])
	}
	scanline := make([]color.NRGBA, width)
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := height; y > 0; y-- {
		binary.Read(os.Stdin, binary.LittleEndian, &scanline)
		for x := 0; x < width; x++ {
			img.Set(x, y, scanline[x])
			np++
		}
	}
	png.Encode(os.Stdout, img)
}
