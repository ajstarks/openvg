// first OpenVG program
package main

import (
	"bufio"
	"github.com/ajstarks/openvg"
	"os"
)

func main() {
	width, height := openvg.Init() // OpenGL, etc initialization

	w2 := float64(width / 2)
	h2 := float64(height / 2)
	w := float64(width)

	openvg.Start(width, height)                               // Start the picture
	openvg.BackgroundColor("black")                           // Black background
	openvg.FillRGB(44, 77, 232, 1)                               // Big blue marble
	openvg.Circle(w2, 0, w)                                   // The "world"
	openvg.FillColor("white")                                 // White text
	openvg.TextMid(w2, h2, "hello, world", "serif", width/10) // Greetings 
	openvg.End()                                              // End the picture
	bufio.NewReader(os.Stdin).ReadBytes('\n')                 // Pause until [RETURN]
	openvg.Finish()                                           // Graphics cleanup
}
