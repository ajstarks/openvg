package main

import (
	"bufio"
	"github.com/ajstarks/openvg"
	"math/rand"
	"os"
	"time"
)

func main() {
	rand.Seed(int64(time.Now().Nanosecond()) % 1e9)
	width, height := openvg.Init()
	openvg.Start(width, height)
	for i := 0; i < 500; i++ {

		red := uint8(rand.Intn(255))
		green := uint8(rand.Intn(255))
		blue := uint8(rand.Intn(255))
		alpha := rand.Float64()

		x := float64(rand.Intn(width))
		y := float64(rand.Intn(height))
		radius := float64(rand.Intn(width / 10))

		openvg.FillRGB(red, green, blue, alpha)
		openvg.Circle(x, y, radius)
	}
	openvg.FillColor("maroon")
	openvg.TextMid(float64(width/2), float64(height/2), "OpenVG via Go on the Raspberry Pi", "sans", width/25)
	openvg.End()
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	openvg.Finish()
}
