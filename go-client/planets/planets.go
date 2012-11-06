// planets: an exploration of scale
package main

import (
	"bufio"
	"github.com/ajstarks/openvg"
	"os"
)

var ssDist = []float64{ // Astronomical Units
	0.00,  // Sun
	0.34,  // Mercury
	0.72,  // Venus
	1.00,  // Earth
	1.54,  // Mars
	5.02,  // Jupiter
	9.46,  // Saturn
	20.11, // Uranus
	30.08} // Netpune

var ssRad = []float64{ // Planet radius in miles
	423200.0, // Sun
	1516.0,   // Mercury
	3760.0,   // Venus
	3957.0,   // Earth
	2104.0,   // Mars
	42980.0,  // Jupiter
	35610.0,  // Saturn
	15700.0,  // Uranus
	15260.0}  // Neptune

var ssColor = []openvg.RGB{ // Planet colors
	{247, 115, 12},  // Sun
	{250, 248, 242}, // Mercury
	{255, 255, 242}, // Venus
	{11, 92, 227},   // Earth
	{240, 198, 29},  // Mars
	{253, 199, 145}, // Jupiter
	{224, 196, 34},  // Saturn
	{220, 241, 245}, // Uranus
	{57, 182, 247},  // Neptune
}

func vmap(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

func light(x, y, r float64, c openvg.RGB) {
	stops := []openvg.Offcolor{
		{0.0, c, 1},
		{0.50, openvg.RGB{c.Red / 2, c.Green / 2, c.Blue / 2}, 1},
	}
	openvg.FillRadialGradient(x, y, (x-r)*.75, y, r, stops)
}

func main() {

	width, height := openvg.Init()
	nobj := len(ssDist)
	y := float64(height) / 2.0
	margin := 100.0
	minsize := 7.0
	labeloc := 100.0
	bgcolor := "black"
	labelcolor := "white"
	maxh := (float64(height) / 2) * 0.05
	openvg.Start(width, height)
	openvg.BackgroundColor(bgcolor)

	for i := 0; i < nobj; i++ {
		x := vmap(ssDist[i], ssDist[0], ssDist[nobj-1], margin, float64(width)-margin)
		r := vmap(ssRad[i], ssRad[1], ssRad[nobj-1], minsize, maxh)

		if ssDist[i] == 0 { // Sun
			openvg.FillRGB(ssColor[0].Red, ssColor[0].Green, ssColor[0].Blue, 1)
			openvg.Circle(margin-(r/2), y, r)
			continue
		}
		if ssDist[i] == 1.0 { // earth
			openvg.StrokeColor(labelcolor)
			openvg.StrokeWidth(1)
			openvg.Line(x, y+(r/2), x, y+labeloc)
			openvg.StrokeWidth(0)
			openvg.FillColor(labelcolor)
			openvg.TextMid(x, y+labeloc+10, "You are here", "sans", 12)
		}

		light(x, y, r, ssColor[i])
		openvg.Circle(x, y, r)
	}
	openvg.End()
	bufio.NewReader(os.Stdin).ReadByte()
	openvg.Finish()
}
