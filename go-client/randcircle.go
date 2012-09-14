package main

import (
	"bufio"
	"flag"
	"github.com/ajstarks/openvg"
	"math/rand"
	"os"
	"time"
)

func main() {
	var nr = flag.Int("n", 500, "number of objects")
	var message = flag.String("m", "Go/OpenVG", "message")
	var bgcolor = flag.String("bg", "white", "background color")
	var fgcolor = flag.String("fg", "maroon", "text color")

	flag.Parse()
	rand.Seed(int64(time.Now().Nanosecond()) % 1e9)

	width, height := openvg.Init()
	fw := float64(width)
	fh := float64(height)

	openvg.Start(width, height)
	openvg.BackgroundColor(*bgcolor)
	for i := 0; i < *nr; i++ {

		red := uint8(rand.Intn(255))
		green := uint8(rand.Intn(255))
		blue := uint8(rand.Intn(255))
		alpha := rand.Float64()

		x := rand.Float64() * fw
		y := rand.Float64() * fh
		radius := rand.Float64() * fw / 10

		openvg.FillRGB(red, green, blue, alpha)
		openvg.Circle(x, y, radius)
	}
	openvg.FillColor(*fgcolor)
	openvg.TextMid(fw/2, fh/2, *message, "sans", width/25)
	openvg.End()

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	openvg.Finish()
}
