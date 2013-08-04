// colortab -- make a color/code placemat
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ajstarks/openvg"
)

func main() {
	var (
		filename           = flag.String("f", "svgcolors.txt", "input file")
		fontname           = flag.String("font", "sans", "fontname")
		neg                = flag.Bool("neg", false, "negative")
		showrgb            = flag.Bool("rgb", false, "show RGB")
		showcode           = flag.Bool("showcode", true, "show colors and codes")
		circsw             = flag.Bool("circle", true, "circle swatch")
		outline            = flag.Bool("outline", false, "outline swatch")
		fontsize           = flag.Int("fs", 12, "fontsize")
		rowsize            = flag.Int("r", 32, "rowsize")
		colw               = flag.Float64("c", 340, "column size")
		swatch             = flag.Float64("s", 16, "swatch size")
		gutter             = flag.Float64("g", 12, "gutter")
		err          error
		tcolor, line string
	)

	flag.Parse()
	f, oerr := os.Open(*filename)
	if oerr != nil {
		fmt.Fprintf(os.Stderr, "%v\n", oerr)
		return
	}
	width, height := openvg.Init()

	openvg.Start(width, height)
	fw := float32(width)
	fh := float32(height)
	if *neg {
		openvg.FillColor("black")
		openvg.Rect(0, 0, fw, fh)
		tcolor = "white"
	} else {
		openvg.FillColor("white")
		openvg.Rect(0, 0, fw, fh)
		tcolor = "black"
	}
	top := fh - 32.0
	left := float32(32.0)
	in := bufio.NewReader(f)

	for x, y, nr := left, top, 0; err == nil; nr++ {
		line, err = in.ReadString('\n')
		fields := strings.Split(strings.TrimSpace(line), "\t")
		if nr%*rowsize == 0 && nr > 0 {
			x += float32(*colw)
			y = top
		}
		if len(fields) == 3 {
			var red, green, blue uint8
			fmt.Sscanf(fields[2], "%d,%d,%d", &red, &green, &blue)
			openvg.FillRGB(red, green, blue, 1)
			if *outline {
				openvg.StrokeColor("black")
				openvg.StrokeWidth(1)
			}
			if *circsw {
				openvg.Circle(x+float32(*swatch)/2.0, y+float32(*swatch)/2.0, float32(*swatch))
			} else {
				openvg.Rect(x, y, float32(*swatch), float32(*swatch))
			}
			openvg.StrokeWidth(0)
			openvg.FillColor(tcolor)
			openvg.Text(x+float32(*swatch)+float32(*fontsize/2), y, fields[0], *fontname, *fontsize)
			var label string
			if *showcode {
				if *showrgb {
					label = fields[1]
				} else {
					label = fields[2]
				}
				openvg.FillColor("gray")
				openvg.TextEnd(x+float32(*colw)-(float32(*swatch)+float32(*gutter)), y, label, *fontname, *fontsize)
			}
		}
		y -= float32(*swatch) + float32(*gutter)
	}
	openvg.End()
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	openvg.Finish()
}
