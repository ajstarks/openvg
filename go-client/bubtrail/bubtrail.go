// bubtrail draws a randmonized trail of bubbles
package main

import (
	"bufio"
	"flag"
	"math/rand"
	"os"
	"time"

	"github.com/ajstarks/openvg"
)

var (
	width, height, niter int
	opacity, size        float32
)

func init() {
	width, height = openvg.Init()
	flag.Float32Var(&size, "s", float32(width)*.05, "bubble size")
	flag.IntVar(&niter, "n", width/6, "number of iterations")
	flag.Float32Var(&opacity, "o", 0.5, "opacity")
	flag.Parse()
	rand.Seed(int64(time.Now().Nanosecond()) % 1e9)
}

func random(howsmall, howbig int) int {
	if howsmall >= howbig {
		return howsmall
	}
	return rand.Intn(howbig-howsmall) + howsmall
}

func main() {
	var color string
	openvg.Start(width, height)
	openvg.Background(200, 200, 200)
	for i := 0; i < niter; i++ {
		x := random(0, width)
		y := random(height/3, (height*2)/3)
		r := random(0, 10000)
		switch {
		case r >= 0 && r <= 2500:
			color = "white"
		case r > 2500 && r <= 5000:
			color = "maroon"
		case r > 5000 && r <= 7500:
			color = "gray"
		case r > 7500 && r <= 10000:
			color = "black"
		}
		openvg.FillColor(color, opacity)
		openvg.Circle(float32(x), float32(y), size)
	}
	openvg.End()
	bufio.NewReader(os.Stdin).ReadByte()
	openvg.Finish()
}
