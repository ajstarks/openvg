// clock
package main

import (
	"math"
	"os"
	"os/signal"
	"time"

	"github.com/ajstarks/openvg"
)

const (
	radians = math.Pi / 180
)

var hourangles = [12]float64{
	90, 60, 30, // 12, 1, 2
	0, 330, 300, // 3, 4, 5
	270, 240, 210, // 6, 7, 8
	180, 150, 120, // 9, 10, 11
}

var minangles = [60]float64{
	90, 84, 78, 72, 66, // 0-4
	60, 54, 48, 42, 36, // 5-9
	30, 24, 18, 12, 6, // 10-14
	0, 354, 348, 342, 336, // 15-19
	330, 324, 318, 312, 306, // 20-24
	300, 294, 288, 282, 276, // 25-29
	270, 264, 258, 252, 246, // 30-34
	240, 234, 228, 222, 216, // 35-39
	210, 204, 198, 192, 186, // 40-44
	180, 174, 168, 162, 156, // 45-49
	150, 144, 138, 132, 126, // 50-54
	120, 114, 108, 102, 96, // 55-59
}

func timecoord(x, y, r openvg.VGfloat, hour, min, sec int) (hx, hy, mx, my, sx, sy openvg.VGfloat) {
	rad := float64(r)
	hrad := rad * 0.6 // hour hand is 60% to the edge of the face
	mrad := rad * 0.9 // minute hand is 90% to the edge
	srad := rad       // second hand is at the edge

	// if the hour is > half-elapsed, adjust the hour angle to
	// reflect the fraction between the current and subsequent hour
	t := hourangles[hour%12]
	if min > 30 {
		t = t - (30.0 * (float64(min) / 60)) // deflection difference is 30 degrees
	}
	t = t * radians
	hx = x + openvg.VGfloat(hrad*math.Cos(t))
	hy = y + openvg.VGfloat(hrad*math.Sin(t))

	t = minangles[min] * radians
	mx = x + openvg.VGfloat(mrad*math.Cos(t))
	my = y + openvg.VGfloat(mrad*math.Sin(t))

	t = minangles[sec] * radians
	sx = x + openvg.VGfloat(srad*math.Cos(t))
	sy = y + openvg.VGfloat(srad*math.Sin(t))
	return
}

var hourdigits = [12]string{"12", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}

func face(x, y, r openvg.VGfloat, ts int) {
	var fx, fy, va openvg.VGfloat
	va = openvg.VGfloat(ts) / 2.0
	rad := float64(r)
	openvg.FillColor("black")
	for h := 12; h > 0; h-- {
		t := hourangles[h%12] * radians
		fx = x + openvg.VGfloat(rad*math.Cos(t))
		fy = y + openvg.VGfloat(rad*math.Sin(t))
		openvg.TextMid(fx, fy-va, hourdigits[h%12], "sans", ts)
	}
}

func main() {
	width, height := openvg.Init()
	cx := openvg.VGfloat(width / 2)
	cy := openvg.VGfloat(height / 2)
	facesize := openvg.VGfloat(cy * 0.5)
	textsize := facesize / 10.0
	framesize := facesize * 2.5
	hourstroke := textsize
	minstroke := hourstroke * .6

	// set up the ticker and signal handler
	ticker := time.NewTicker(1 * time.Second)
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	// main loop: for each second, draw the clock components
	openvg.Start(width, height)
	for {
		select {
		case <-ticker.C:
			// get the current time
			hr, min, sec := time.Now().Clock()

			// compute element coordinates
			hx, hy, mx, my, sx, sy := timecoord(cx, cy, facesize, hr, min, sec)

			// frame and clock face
			openvg.FillRGB(127, 127, 127, 1)
			openvg.Roundrect(cx-framesize/2, cy-framesize/2, framesize, framesize, textsize, textsize)
			openvg.FillRGB(255, 255, 255, 1)
			openvg.Ellipse(cx, cy, facesize*2.2, facesize*2.2)
			face(cx, cy, facesize, int(textsize))

			// hour hand
			openvg.StrokeWidth(hourstroke)
			openvg.StrokeRGB(127, 127, 127, 1)
			openvg.Line(cx, cy, hx, hy)
			openvg.StrokeWidth(0)
			openvg.FillRGB(127, 127, 127, 1)
			openvg.Ellipse(hx, hy, hourstroke, hourstroke)

			// minute hand
			openvg.StrokeWidth(minstroke)
			openvg.StrokeRGB(127, 0, 0, 1)
			openvg.Line(cx, cy, mx, my)
			openvg.StrokeWidth(0)
			openvg.FillRGB(127, 0, 0, 1)
			openvg.Ellipse(mx, my, minstroke, minstroke)

			// second indicator
			openvg.FillRGB(0, 0, 255, 0.3)
			openvg.Ellipse(sx, sy, textsize, textsize)

			// center dot
			openvg.FillColor("black")
			openvg.Ellipse(cx, cy, textsize, textsize)
			openvg.End()
		case <-sigint:
			openvg.Finish()
			return
		}
	}
}
