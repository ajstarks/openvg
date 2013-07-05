// deck: make slide decks
package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/ajstarks/openvg"
)

// Deck defines the structure of a presentation deck
// The size of the canvas, and series of slides
type Deck struct {
	Canvas canvas  `xml:"canvas"`
	Slide  []slide `xml:"slide"`
}

type canvas struct {
	Width  int `xml:"width,attr"`
	Height int `xml:"height,attr"`
}

type slide struct {
	Bg    string  `xml:"bg,attr"`
	Fg    string  `xml:"fg,attr"`
	List  []list  `xml:"list"`
	Text  []text  `xml:"text"`
	Image []image `xml:"image"`
}

type CommonAttr struct {
	Xp    float64 `xml:"xp,attr"`
	Yp    float64 `xml:"yp,attr"`
	Sp    float64 `xml:"sp,attr"`
	Type  string  `xml:"type,attr"`
	Align string  `xml:"align,attr"`
	Color string  `xml:"color,attr"`
	Font  string  `xml:"font,attr"`
}

type list struct {
	CommonAttr
	Li []string `xml:"li"`
}

type text struct {
	CommonAttr
	Wp    float64 `xml:"wp,attr"`
	Tdata string  `xml:",chardata"`
}

type image struct {
	Xp     float64 `xml:"xp,attr"`
	Yp     float64 `xml:"yp,attr"`
	Width  int     `xml:"width,attr"`
	Height int     `xml:"height,attr"`
	Name   string  `xml:"name,attr"`
}

// dodeck sets up the graphics environment and kicks off the interaction
func dodeck(filename string) {
	w, h := openvg.Init()
	openvg.Background(0, 0, 0)
	interact(filename, w, h)
	openvg.Finish()
}

// readdeck reads the deck description file
func readdeck(filename string, w, h int) (Deck, error) {
	var d Deck
	r, err := os.Open(filename)
	if err != nil {
		return d, err
	}
	err = xml.NewDecoder(r).Decode(&d)
	if d.Canvas.Width == 0 {
		d.Canvas.Width = w
	}
	if d.Canvas.Height == 0 {
		d.Canvas.Height = h
	}
	r.Close()
	return d, err
}

// interact controls the display of the deck
func interact(filename string, w, h int) {
	openvg.SaveTerm()
	defer openvg.RestoreTerm()
	var d Deck
	var err error
	d, err = readdeck(filename, w, h)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	openvg.RawTerm()
	r := bufio.NewReader(os.Stdin)
	firstslide := 0
	lastslide := len(d.Slide) - 1
	n := firstslide

	// respond to keyboard commands, 'q' to exit
	for cmd := byte('0'); cmd != 'q'; cmd = readcmd(r) {
		switch cmd {
		// read/reload
		case 'r', 18: // r, Ctrl-R
			d, err = readdeck(filename, w, h)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return
			}
			showslide(d, n)

		// save slide
		case 's':
			openvg.SaveEnd(fmt.Sprintf("slide-%04d", n))

		// first slide
		case '0', '1', 1, '^': // 0,1,Ctrl-A,^
			n = firstslide
			showslide(d, n)

		// last slide
		case '*', 5, '$': // *, Crtl-E, $
			n = lastslide
			showslide(d, n)

		// next slide
		case '+', 'n', '\n', ' ', '\t', 14: // +,n,newline,space,tab,Crtl-N
			n++
			if n > lastslide {
				n = firstslide
			}
			showslide(d, n)

		// previous slide
		case '-', 'p', 8, 16, 127: // -,p,Backspace,Ctrl-P,Del
			n--
			if n < firstslide {
				n = lastslide
			}
			showslide(d, n)

		// grid
		case 'g':
			showgrid(d, n)

		// search
		case '/', 6: // slash, Ctrl-F
			openvg.RestoreTerm()
			searchterm, serr := r.ReadString('\n')
			openvg.RawTerm()
			if serr != nil {
				continue
			}
			if len(searchterm) > 2 {
				ns := searchdeck(d, searchterm[0:len(searchterm)-1])
				if ns >= 0 {
					showslide(d, ns)
					n = ns
				}
			}
		}
	}
}

// showgrid xrays a slide
func showgrid(d Deck, n int) {
	w := float64(d.Canvas.Width)
	h := float64(d.Canvas.Height)
	fs := w * 0.01 // labels are 1% of the width
	pct := 10.0    // grid at 10% intervals
	xpct := (pct / 100.0) * w
	ypct := (pct / 100.0) * h

	openvg.StrokeColor("lightgray", 0.5)
	openvg.StrokeWidth(2)

	// vertical gridlines
	xl := 0.0
	for x := xl; x <= w; x += xpct {
		openvg.Line(x, 0, x, h)
		openvg.Text(x, 0, fmt.Sprintf("%.0f%%", xl), "sans", int(fs))
		xl += pct
	}

	// horizontal gridlines
	yl := 0.0
	for y := yl; y <= h; y += ypct {
		openvg.Line(0, y, w, y)
		openvg.Text(0, y, fmt.Sprintf("%.0f%%", yl), "sans", int(fs))
		yl += pct
	}

	// show boundary and location of images
	if n < 0 || n > len(d.Slide) {
		return
	}
	for _, im := range d.Slide[n].Image {
		x := (im.Xp / 100) * w
		y := (im.Yp / 100) * h
		iw := float64(im.Width)
		ih := float64(im.Height)
		openvg.FillRGB(127, 0, 0, 0.3)
		openvg.Circle(x, y, fs)
		openvg.FillRGB(255, 0, 0, 0.1)
		openvg.Rect(x-iw/2, y-ih/2, iw, ih)
	}
	openvg.End()
}

// dimen computes the coordinates and size of an object
func dimen(c canvas, xp, yp, sp float64) (x, y float64, s int) {
	x = (xp / 100) * float64(c.Width)
	y = (yp / 100) * float64(c.Height)
	s = int((sp / 100) * float64(c.Width))
	return
}

//showtext displays text
func showtext(x, y float64, s, align, font string, fontsize int) {
	switch align {
	case "center", "middle", "mid":
		openvg.TextMid(x, y, s, font, fontsize)
	case "right", "end":
		openvg.TextEnd(x, y, s, font, fontsize)
	default:
		openvg.Text(x, y, s, font, fontsize)
	}
}

// showlide displays slides
func showslide(d Deck, n int) {
	if n < 0 || n > len(d.Slide)-1 {
		return
	}
	slide := d.Slide[n]
	if slide.Bg == "" {
		slide.Bg = "white"
	}
	if slide.Fg == "" {
		slide.Fg = "black"
	}
	openvg.Start(d.Canvas.Width, d.Canvas.Height)
	cw := float64(d.Canvas.Width)
	ch := float64(d.Canvas.Height)
	openvg.FillColor(slide.Bg)
	openvg.Rect(0, 0, cw, ch)
	openvg.FillColor(slide.Fg)

	var x, y float64
	var fontsize int

	// every image in the slide
	for _, im := range slide.Image {
		x = (im.Xp / 100) * cw
		y = (im.Yp / 100) * ch
		openvg.Image(x-float64(im.Width/2), y-float64(im.Height/2), im.Width, im.Height, im.Name)
	}

	// every list in the slide
	var offset float64
	const blinespacing = 2.0
	for _, l := range slide.List {
		if l.Font == "" {
			l.Font = "sans"
		}
		x, y, fontsize = dimen(d.Canvas, l.Xp, l.Yp, l.Sp)
		fs := float64(fontsize)
		if l.Type == "bullet" {
			offset = 1.2 * fs
		} else {
			offset = 0
		}
		if l.Color != "" {
			openvg.FillColor(l.Color)
		} else {
			openvg.FillColor(slide.Fg)
		}
		// every list item
		for _, li := range l.Li {
			if l.Type == "bullet" {
				boffset := fs / 2
				openvg.Circle(x, y+boffset, boffset)
			}
			showtext(x+offset, y, li, l.Align, l.Font, fontsize)
			y -= fs * blinespacing
		}
	}
	openvg.FillColor(slide.Fg)

	// every text in the slide
	const linespacing = 1.8
	for _, t := range slide.Text {
		if t.Font == "" {
			t.Font = "sans"
		}
		x, y, fontsize = dimen(d.Canvas, t.Xp, t.Yp, t.Sp)
		fs := float64(fontsize)
		td := strings.Split(t.Tdata, "\n")
		if t.Type == "code" {
			t.Font = "mono"
			tdepth := ((fs * linespacing) * float64(len(td))) + fs
			openvg.FillColor("rgb(240,240,240)")
			openvg.Rect(x-20, y-tdepth+(fs*linespacing), cw-20-x, tdepth)
		}
		if t.Color != "" {
			openvg.FillColor(t.Color)
		} else {
			openvg.FillColor(slide.Fg)
		}
		if t.Type == "block" {
			var tw float64
			if t.Wp == 0 {
				tw = cw / 2
			} else {
				tw = (t.Wp / 100) * cw
			}
			textwrap(x, y, tw, t.Tdata, t.Font, fontsize, fs*linespacing, 0.3)
		} else {
			// every text line
			for _, txt := range td {
				showtext(x, y, txt, t.Align, t.Font, fontsize)
				y -= (fs * linespacing)
			}
		}
	}
	openvg.FillColor(slide.Fg)
	openvg.End()
}

// searchdeck searches the deck for the specified text, returning the slide number if found
func searchdeck(d Deck, s string) int {
	// for every slide...
	for i, slide := range d.Slide {
		// search lists
		for _, l := range slide.List {
			for _, ll := range l.Li {
				if strings.Contains(ll, s) {
					return i
				}
			}
		}
		// search text
		for _, t := range slide.Text {
			if strings.Contains(t.Tdata, s) {
				return i
			}
		}
	}
	return -1
}

// dumpdeck shows the decoded description
func dumpdeck(d Deck) {
	fmt.Printf("Canvas = %v\n", d.Canvas)
	for i, s := range d.Slide {
		fmt.Printf("Slide #%d = %v %v\n", i, s.Bg, s.Fg)
		for j, l := range s.List {
			fmt.Printf("\tList #%d = %#v\n", j, l)
		}
		for k, t := range s.Text {
			fmt.Printf("\tText #%d = %#v\n", k, t)
		}
		for m, im := range s.Image {
			fmt.Printf("\tImage #%d = %#v\n", m, im)
		}
	}
}

// whitespace determines if a rune is whitespace
func whitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

// textwrap draws text at location, wrapping at the specified width
func textwrap(x, y, w float64, s string, font string, size int, leading, factor float64) {
	wordspacing := openvg.TextWidth("m", font, size)
	words := strings.FieldsFunc(s, whitespace)
	xp := x
	yp := y
	edge := x + w
	for _, s := range words {
		tw := openvg.TextWidth(s, font, size)
		openvg.Text(xp, yp, s, font, size)
		xp += tw + (wordspacing * factor)
		if xp > edge {
			xp = x
			yp -= leading
		}
	}
}

// readcmd reads interaction commands
func readcmd(r *bufio.Reader) byte {
	s, err := r.ReadByte()
	if err != nil {
		return 'e'
	}
	return s
}

// for every file, make a deck
func main() {
	if len(os.Args) > 1 {
		for _, f := range os.Args[1:] {
			dodeck(f)
		}
	}
}
