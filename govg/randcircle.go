package main

import (
	"os"
	"bufio"
	"github.com/ajstarks/openvg"
	"math/rand"
	"time"
)

func WaitEnd() {
	openvg.End()
bufio.NewReader(os.Stdin).ReadLine()	
	openvg.Finish()
}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()) % 1e9)
	width, height := openvg.Init()
	openvg.Start(width, height)
	for i := 0; i < 500; i++ {
		x := float64(rand.Intn(width))
		y := float64(rand.Intn(height))
		r := float64(rand.Intn(width / 10))
		alpha := rand.Float64()
		red := uint8(rand.Intn(255))
		green := uint8(rand.Intn(255))
		blue := uint8(rand.Intn(255))
		openvg.Fill(red, green, blue, alpha)
		openvg.Circle(x, y, r)
	}
	openvg.Fill(128, 0, 0, 1)
	openvg.TextMid(float64(width/2), float64(height/2), "OpenVG via Go on the Raspberry Pi", "sans", width/25)
	WaitEnd()
}
