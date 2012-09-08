package main

import (
	"fmt"
	"github.com/ajstarks/openvg"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Color struct {
	red, green, blue uint8
	alpha            float64
}

// randcolor returns a random number 0..255
func randcolor() uint8 {
	return uint8(rand.Intn(255))
}

// randf returns a floating point number bounded by n
func randf(n int) float64 {
	return (rand.Float64() * float64(n))
}

// coordpoint marks a coordinate, preserving a previous color
func coordpoint(x, y, size float64, c Color) {
	openvg.FillRGB(128, 0, 0, 0.3)
	openvg.Circle(x, y, size)
	openvg.FillRGB(c.red, c.green, c.blue, c.alpha)
}

// makepi draws the Raspberry Pi
func makepi(x, y, w, h float64) {
	// dimensions
	socw := h / 5
	compw := h / 5
	cjw := h / 10
	cjh := h / 8
	audw := h / 5
	aujw := h / 10
	aujh := cjh / 2
	hdw := w / 6
	hdh := w / 10
	gpw := w / 3
	gph := h / 8
	pw := h / 10
	usw := w / 5
	ush := h / 5
	etw := w / 5
	eth := h / 5
	sdw := w / 6
	sdh := w / 4
	offset := (w / 2) / 10
	w34 := (w * 3) / 4
	w2 := w / 2
	h2 := h / 2
	h40 := (h * 2) / 5

	openvg.FillRGB(0, 128, 0, 1)
	openvg.Rect(x, y, w, h) // board

	openvg.FillRGB(255, 255, 0, 1)
	openvg.Rect(x+w2, (y+h)-compw, compw, compw) // composite
	openvg.FillRGB(192, 192, 192, 1)
	openvg.Rect(x+w2+(cjw/2), y+h, cjw, cjh) // composite jack

	openvg.FillRGB(0, 0, 0, 1)
	openvg.Rect(x+w34, y+h-audw, audw, audw)     // audio
	openvg.Rect(x+w34+(aujw/2), y+h, aujw, aujh) // audio jack

	openvg.FillRGB(192, 192, 192, 1)
	openvg.Rect(x+w2, y, hdw, hdh)                 // HDMI
	openvg.Rect((x+w)-etw, y, etw, eth)            // Ethernet
	openvg.Rect((x+w+offset)-usw, y+h40, usw, ush) // USB
	openvg.Rect(x, y, pw, pw)                      // Power

	openvg.FillRGB(0, 0, 0, 1)
	openvg.Rect(x+(w*2)/5, y+h40, socw, socw) // SoC
	openvg.Rect(x, (y+h)-gph, gpw, gph)       // GPIO
	openvg.FillRGB(0, 0, 255, 1)
	openvg.Rect(x-sdw, (y+h2)-(sdh/2), sdw, sdh) // SD card
}

// raspberry pi, scaled to the screen dimensions
func raspi(w, h int, s string) {
	midx := float64(w) / 2
	midy := float64(h) / 2
	rw := midx
	rh := (rw * 2) / 3
	fontsize := w / 25
	openvg.Start(w, h)
	openvg.Background(255, 255, 255)
	makepi(midx-(rw/2), midy-(rh/2), rw, rh)
	openvg.FillRGB(128, 0, 0, 1)
	openvg.TextMid(midx, midy-(rh/2)-float64(fontsize*2), s, "sans", fontsize)
	openvg.End()
}

// grid draws a grid
func grid(x, y float64, n, w, h int) {
	width := float64(w)
	height := float64(h)
	gn := float64(n)
	openvg.StrokeRGB(128, 128, 128, 0.5)
	openvg.StrokeWidth(2)
	for ix := x; ix <= x+width; ix += gn {
		openvg.Line(ix, y, ix, y+height)
	}

	for iy := y; iy <= y+height; iy += gn {
		openvg.Line(x, iy, x+width, iy)
	}
}

type FW struct {
	font     string
	tw       float64
	fontsize int
}

// adjust the font to fit within a width
func (f *FW) fitwidth(width, adj int, s string) {
	f.tw = openvg.TextWidth(s, f.font, float64(f.fontsize))
	for f.tw > float64(width) {
		f.fontsize -= adj
		f.tw = openvg.TextWidth(s, f.font, float64(f.fontsize))
	}
}

// testpattern shows a test pattern 
func testpattern(w, h int, s string) {
	var midx, midy1, midy2, midy3 float64
	fontsize := 256
	h2 := float64(h / 2)
	by := float64(h - 100)
	bx := float64(w - 100)
	tw1 := &FW{"mono", 0, fontsize}
	tw2 := &FW{"sans", 0, fontsize}
	tw3 := &FW{"serif", 0, fontsize}

	openvg.Start(w, h)

	// colored squares in the corners
	openvg.FillRGB(255, 0, 0, 1)
	openvg.Rect(0, 0, 100, 100)
	openvg.FillRGB(0, 255, 0, 1)
	openvg.Rect(0, by, 100, 100)
	openvg.FillRGB(0, 0, 255, 1)
	openvg.Rect(bx, 0, 100, 100)
	openvg.FillRGB(128, 128, 128, 1)
	openvg.Rect(bx, by, 100, 100)

	// for each font, (Sans, Serif, Mono), adjust the string to the w
	tw1.fitwidth(w, 20, s)
	tw2.fitwidth(w, 20, s)
	tw3.fitwidth(w, 20, s)

	midx = float64(w / 2)

	// Adjust the baselines to be medial
	midy1 = h2 + 20 + float64((tw1.fontsize)/2)
	midy2 = h2 - float64((tw2.fontsize)/2)
	midy3 = h2 - 20 - float64(tw2.fontsize) - float64((tw3.fontsize)/2)

	openvg.FillRGB(128, 128, 128, 1)
	openvg.TextMid(midx, midy1, s, tw1.font, tw1.fontsize)
	openvg.FillRGB(128, 0, 0, 1)
	openvg.TextMid(midx, midy2, s, tw2.font, tw2.fontsize)
	openvg.FillRGB(0, 0, 128, 1)
	openvg.TextMid(midx, midy3, s, tw3.font, tw3.fontsize)
	openvg.End()
}

// textlines writes openvg.Lines of text
func textlines(x, y float64, text []string, f string, fontsize int, leading float64) {
	for _, s := range text {
		openvg.Text(x, y, s, f, fontsize)
		y -= leading
	}
}

// tb draws a block of text
func tb(w, h int) {
	para := []string{
		"For lo, the winter is past,",
		"the rain is over and gone",
		"the flowers appear on the earth;",
		"the time for the singing of birds is come,",
		"and the voice of the turtle is heard in our land",
	}

	tmargin := float64(w) * 0.25
	lmargin := float64(w) * 0.10
	top := float64(h) * .9
	mid := float64(h) * .6
	bot := float64(h) * .3

	fontsize := 24
	leading := 40.0
	lfontsize := fontsize * 2
	midb := ((leading * 2) + (leading / 2)) - float64(lfontsize/2)

	openvg.Start(w, h)
	openvg.FillRGB(49, 79, 79, 1)
	textlines(tmargin, top, para, "serif", fontsize, leading)
	textlines(tmargin, mid, para, "sans", fontsize, leading)
	textlines(tmargin, bot, para, "mono", fontsize, leading)
	openvg.Text(lmargin, top-midb, "Serif", "sans", lfontsize)
	openvg.Text(lmargin, mid-midb, "Sans", "sans", lfontsize)
	openvg.Text(lmargin, bot-midb, "Mono", "sans", lfontsize)
	openvg.End()
}

// cookie draws a cookie
func cookie(w, h int) {
	ew := 200.0
	eh := 60.0
	h2 := float64(h) / 2.0
	w2 := float64(w) / 2.0

	openvg.Start(w, h)
	openvg.FillRGB(128, 128, 128, 1)
	openvg.Ellipse(w2, h2, ew, eh)
	openvg.Translate(0, 10)

	openvg.FillRGB(255, 255, 255, 1)
	openvg.Ellipse(w2, h2, ew, eh)
	openvg.Translate(0, 20)

	openvg.FillRGB(0, 0, 0, 1)
	openvg.Ellipse(w2, h2, ew, eh)
	openvg.End()
}

func imagetest(w, h int) {
	imgw := 400
	imgh := 400
	fiw := float64(imgw)
	fih := float64(imgh)
	fw := float64(w)
	fh := float64(h)
	cx := (fw / 2) - (fiw / 2)
	cy := (fh / 2) - (fih / 2)
	ulx := 0.0
	uly := fh - fih
	urx := fw - fiw
	ury := uly
	llx := 0.0
	lly := 0.0
	lrx := urx
	lry := lly
	openvg.Start(w, h)
	openvg.Background(0, 0, 0)
	openvg.Image(cx, cy, imgw, imgh, "test_img_violin.jpg")
	openvg.Image(ulx, uly, imgw, imgh, "test_img_piano.jpg")
	openvg.Image(urx, ury, imgw, imgh, "test_img_sax.jpg")
	openvg.Image(llx, lly, imgw, imgh, "test_img_guitar.jpg")
	openvg.Image(lrx, lry, imgw, imgh, "test_img_flute.jpg")
	openvg.End()
}

// fontrange shows a range of fonts
func fontrange(w, h int) {
	var x, lx, length float64
	y := float64(h) / 2.0
	w2 := float64(w) / 2.0
	spacing := 50.0
	s2 := spacing / 2.0
	sizes := []int{6, 7, 8, 9, 10, 11, 12, 14, 16, 18, 21, 24, 36, 48, 60, 72, 96}

	openvg.Start(w, h)
	openvg.Background(255, 255, 255)

	// compute the length so we can center
	length = 0.0
	for _, s := range sizes {
		length += float64(s) + spacing
	}
	length -= spacing
	lx = w2 - (length / 2) // center point

	// for each size, display a character and label
	x = lx
	for _, s := range sizes {
		openvg.FillRGB(128, 0, 0, 1)
		openvg.TextMid(x, y, "a", "serif", s)
		openvg.FillRGB(128, 128, 128, 1)
		openvg.TextMid(x, y-spacing, fmt.Sprintf("%d", s), "sans", 16)
		x += float64(s) + spacing
	}
	// draw a openvg.Line below the characters, a curve above
	x -= spacing
	openvg.StrokeRGB(150, 150, 150, 0.5)
	openvg.StrokeWidth(2)
	openvg.Line(lx, y-s2, x, y-s2)
	openvg.FillRGB(255, 255, 255, 1)
	openvg.Qbezier(lx, y+s2, x, y+s2, x, y+(spacing*3))
	openvg.End()
}

// refcard shows a reference card of shapes
func refcard(width, height int) {
	shapenames := []string{
		"Circle",
		"openvg.Ellipse",
		"Rectangle",
		"Rounded Rectangle",
		"Line",
		"Polyline",
		"Polygon",
		"Arc",
		"Quadratic Bezier",
		"Cubic Bezier",
		"Image",
	}
	top := float64(height) * .95
	sx := float64(width) * 0.10
	sy := top
	sw := float64(width) * .05
	sh := float64(height) * .045
	dotsize := 7.0
	spacing := 2.0
	fontsize := int(float64(height) * .033)
	shapecolor := Color{202, 225, 255, 1.0}

	openvg.Start(width, height)
	openvg.FillRGB(128, 0, 0, 1)
	openvg.TextEnd(float64(width-20), float64(height/2), "OpenVG on the Raspberry Pi", "sans", fontsize+(fontsize/2))
	openvg.FillRGB(0, 0, 0, 1)
	for _, s := range shapenames {
		openvg.Text(sx+sw+sw/2, sy, s, "sans", fontsize)
		sy -= sh * spacing
	}
	sy = top
	cx := sx + (sw / 2)
	ex := sx + sw
	openvg.FillRGB(shapecolor.red, shapecolor.green, shapecolor.blue, shapecolor.alpha)
	openvg.Circle(cx, sy, sw)
	coordpoint(cx, sy, dotsize, shapecolor)
	sy -= sh * spacing
	openvg.Ellipse(cx, sy, sw, sh)
	coordpoint(cx, sy, dotsize, shapecolor)
	sy -= sh * spacing
	openvg.Rect(sx, sy, sw, sh)
	coordpoint(sx, sy, dotsize, shapecolor)
	sy -= sh * spacing
	openvg.Roundrect(sx, sy, sw, sh, 20, 20)
	coordpoint(sx, sy, dotsize, shapecolor)
	sy -= sh * spacing

	openvg.StrokeWidth(1)
	openvg.StrokeRGB(204, 204, 204, 1)
	openvg.Line(sx, sy, ex, sy)
	coordpoint(sx, sy, dotsize, shapecolor)
	coordpoint(ex, sy, dotsize, shapecolor)
	sy -= sh

	px := []float64{sx, sx + (sw / 4), sx + (sw / 2), sx + ((sw * 3) / 4), sx + sw}
	py := []float64{sy, sy - sh, sy, sy - sh, sy}

	openvg.Polyline(px, py, 5)
	coordpoint(px[0], py[0], dotsize, shapecolor)
	coordpoint(px[1], py[1], dotsize, shapecolor)
	coordpoint(px[2], py[2], dotsize, shapecolor)
	coordpoint(px[3], py[3], dotsize, shapecolor)
	coordpoint(px[4], py[4], dotsize, shapecolor)
	sy -= sh * spacing

	py[0] = sy
	py[1] = sy - sh
	py[2] = sy - (sh / 2)
	py[3] = py[1] - (sh / 4)
	py[4] = sy
	openvg.Polygon(px, py, 5)
	sy -= (sh * spacing) + sh

	openvg.Arc(sx+(sw/2), sy, sw, sh, 0, 180)
	coordpoint(sx+(sw/2), sy, dotsize, shapecolor)
	sy -= sh * spacing

	var cy, ey float64
	cy = sy + (sh / 2)
	ey = sy
	openvg.Qbezier(sx, sy, cx, cy, ex, ey)
	coordpoint(sx, sy, dotsize, shapecolor)
	coordpoint(cx, cy, dotsize, shapecolor)
	coordpoint(ex, ey, dotsize, shapecolor)
	sy -= sh * spacing

	ey = sy
	cy = sy + sh
	openvg.Cbezier(sx, sy, cx, cy, cx, sy, ex, ey)
	coordpoint(sx, sy, dotsize, shapecolor)
	coordpoint(cx, cy, dotsize, shapecolor)
	coordpoint(cx, sy, dotsize, shapecolor)
	coordpoint(ex, ey, dotsize, shapecolor)

	sy -= (sh * spacing * 1.5)
	// Image(sx, sy, 100, 100, "starx.jpg")

	openvg.End()
}

// rotext draws text, rotated around the center of the screen, progressively faded
func rotext(w, h, n int, s string) {
	fade := (100.0 / float64(n)) / 100.0
	deg := 360.0 / float64(n)
	x := float64(w) / 2.0
	y := float64(h) / 2.0
	alpha := 1.0
	size := w / 8

	openvg.Start(w, h)
	openvg.Background(0, 0, 0)
	openvg.Translate(x, y)
	for i := 0; i < n; i++ {
		openvg.FillRGB(255, 255, 255, alpha)
		openvg.Text(0, 0, s, "serif", size)
		alpha -= fade // fade
		size += n     // enlarge
		openvg.Rotate(deg)
	}
	openvg.End()
}

// rseed seeds the random number generator from the random device
func rseed() {
	rand.Seed(int64(time.Now().Nanosecond()) % 1e9)
}

// rshapes draws shapes with random colors, openvg.Strokes, and sizes. 
func rshapes(width, height, n int) {

	var sx, sy, cx, cy, px, py, ex, ey, pox, poy float64

	np := 10
	polyx := make([]float64, np)
	polyy := make([]float64, np)
	openvg.Start(width, height)
	for i := 0; i < n; i++ {
		openvg.FillRGB(randcolor(), randcolor(), randcolor(), rand.Float64())
		openvg.Ellipse(randf(width), randf(height), randf(200), randf(100))
		openvg.Circle(randf(width), randf(height), randf(100))
		openvg.Rect(randf(width), randf(height), randf(200), randf(100))
		openvg.Arc(randf(width), randf(height), randf(200), randf(200), randf(360), randf(360))

		sx = randf(width)
		sy = randf(height)
		openvg.StrokeRGB(randcolor(), randcolor(), randcolor(), 1)
		openvg.StrokeWidth(randf(5))
		openvg.Line(sx, sy, sx+randf(200), sy+randf(100))
		openvg.StrokeWidth(0)

		sx = randf(width)
		sy = randf(height)
		ex = sx + randf(200)
		ey = sy
		cx = sx + ((ex - sx) / 2.0)
		cy = sy + randf(100)
		openvg.Qbezier(sx, sy, cx, cy, ex, ey)

		sx = randf(width)
		sy = randf(height)
		ex = sx + randf(200)
		ey = sy
		cx = sx + ((ex - sx) / 2.0)
		cy = sy + randf(100)
		px = cx
		py = sy - randf(100)
		openvg.Cbezier(sx, sy, cx, cy, px, py, ex, ey)

		pox = randf(width)
		poy = randf(height)
		for j := 0; j < np; j++ {
			polyx[j] = pox + randf(200)
			polyy[j] = poy + randf(100)
		}
		openvg.Polygon(polyx, polyy, np)

		pox = randf(width)
		poy = randf(height)
		for j := 0; j < np; j++ {
			polyx[j] = pox + randf(200)
			polyy[j] = poy + randf(100)
		}
		openvg.Polyline(polyx, polyy, np)
	}
	openvg.FillRGB(128, 0, 0, 1)
	openvg.Text(20, 20, "OpenVG on the Raspberry Pi", "sans", 32)
	openvg.End()
}

// sunearth shows the relative sizes of the sun and the earth
func sunearth(w, h int) {
	var sun, earth, x, y float64

	openvg.Start(w, h)
	openvg.Background(0, 0, 0)
	openvg.FillRGB(255, 255, 255, 1)
	for i := 0; i < w/4; i++ {
		x = randf(w)
		y = randf(h)
		openvg.Circle(x, y, 2)
	}
	earth = float64(w) * 0.010
	sun = earth * 109
	openvg.FillRGB(0, 0, 255, 1)
	openvg.Circle(float64(w/3), float64(h-(h/10)), earth)
	openvg.FillRGB(255, 255, 224, 1)
	openvg.Circle(float64(w), 0, sun)
	openvg.End()
}

// advert is an ad for the package 
func advert(w, h int) {
	y := (6 * float64(h)) / 10
	fontsize := int(float64(w) * 0.04)
	f3 := float64(fontsize) / 3.0
	s := "github.com/ajstarks/openvg"
	a := "ajstarks@gmail.com"

	//imw := 110
	// imh := 110
	midx := float64(w / 2)

	openvg.Start(w, h)
	openvg.FillRGB(128, 0, 0, 1)
	openvg.TextMid(midx, y-float64(fontsize/4), s, "sans", fontsize)
	y -= 150
	openvg.FillRGB(128, 128, 128, 1)
	openvg.TextMid(midx, y, a, "sans", int(f3))
	// Image((w / 2) - (imw / 2), y - (imh * 2), imw, imh, "starx.jpg")
	openvg.End()
}

// demo shows a timed demonstration
func demo(w, h, s int) {
	sec := time.Duration(s) * time.Second
	refcard(w, h)
	time.Sleep(sec)

	rshapes(w, h, 50)
	time.Sleep(sec)

	testpattern(w, h, "OpenVG on RasPi")
	time.Sleep(sec)

	//imagetest(w, h)
	//time.Sleep(sec)

	rotext(w, h, 30, "Raspi")
	time.Sleep(sec)

	tb(w, h)
	time.Sleep(sec)

	fontrange(w, h)
	time.Sleep(sec)

	sunearth(w, h)
	time.Sleep(sec)

	raspi(w, h, "The Raspberry Pi")
	time.Sleep(sec)

	advert(w, h)
}

func WaitEnd() {
	var x int
	fmt.Scanf("%d", &x)
	openvg.Finish()
}

func usage(s string) {
	fmt.Fprintf(os.Stderr,
		"%s [command]\n\tdemo sec\n\tastro\n\ttest ...\n\trand n\n\trotate n ...\n\timage\n\ttext\n\tfontsize\n\traspi\n", s)
}

// main initializes the system and shows the picture. 
// Exit and clean up when you hit [RETURN].
func main() {
	rseed()
	n := 10
	nargs := len(os.Args)
	w, h := openvg.Init()
	progname := os.Args[0]
	if nargs > 1 {
		switch os.Args[1] {
		case "help":
			usage(progname)
			os.Exit(1)
		case "image":
			imagetest(w, h)
		case "text":
			tb(w, h)
		case "astro":
			sunearth(w, h)
		case "fontsize":
			fontrange(w, h)
		case "raspi":
			raspi(w, h, "The Raspberry Pi")
		case "demo":
			if nargs > 2 {
				n, _ = strconv.Atoi(os.Args[2])
				demo(w, h, n)
			} else {
				usage(progname)
				os.Exit(1)
			}
		case "rand":
			rshapes(w, h, n)
		case "test":
			testpattern(w, h, os.Args[2])
		case "rotate":
			if nargs > 3 {
				n, _ = strconv.Atoi(os.Args[2])
				rotext(w, h, n, os.Args[3])
			} else {
				usage(progname)
				os.Exit(1)
			}
		default:
			refcard(w, h)
		}
	} else {
		refcard(w, h)
	}
	WaitEnd()
}
