// deck: make slide decks
package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
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

type list struct {
	Xp    float64  `xml:"xp,attr"`
	Yp    float64  `xml:"yp,attr"`
	Sp    float64  `xml:"sp,attr"`
	Type  string   `xml:"type,attr"`
	Align string   `xml:"align,attr"`
	Li    []string `xml:"li"`
}

type text struct {
	Xp    float64 `xml:"xp,attr"`
	Yp    float64 `xml:"yp,attr"`
	Sp    float64 `xml:"sp,attr"`
	Align string  `xml:"align,attr"`
	Type  string  `xml:"type,attr"`
	Tdata string  `xml:",chardata"`
}

type image struct {
	Xp     float64 `xml:"xp,attr"`
	Yp     float64 `xml:"yp,attr"`
	Width  int     `xml:"width,attr"`
	Height int     `xml:"height,attr"`
	Name   string  `xml:"name,attr"`
}

type slide struct {
	Bg    string  `xml:"bg,attr"`
	Fg    string  `xml:"fg,attr"`
	List  []list  `xml:"list"`
	Text  []text  `xml:"text"`
	Image []image `xml:"image"`
}

// dodeck reads and decodes slide files
func dodeck(filename string) {
	var r io.ReadCloser
	var err error
	if len(filename) > 0 {
		r, err = os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
	} else {
		r = os.Stdin
	}
	defer r.Close()

	var d Deck
	derr := xml.NewDecoder(r).Decode(&d)
	if derr != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	w, h := openvg.Init()
	if d.Canvas.Width == 0 {
		d.Canvas.Width = w
	}
	if d.Canvas.Height == 0 {
		d.Canvas.Height = h
	}
	interact(d)
	openvg.Finish()
	//dumpdeck(d)
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

// dimen determines the coordinates and size of an object
func dimen(c canvas, xp, yp, sp float64) (x, y float64, s int) {
	x = (xp / 100) * float64(c.Width)
	y = (yp / 100) * float64(c.Height)
	s = int((sp / 100) * float64(c.Width))
	return
}

// showlide displays slides
func showslide(d Deck, n int) {
	if n < 0 || n > len(d.Slide)-1 {
		return
	}
	s := d.Slide[n]
	openvg.Start(d.Canvas.Width, d.Canvas.Height)
	if s.Bg == "" {
		s.Bg = "white"
	}
	if s.Fg == "" {
		s.Fg = "black"
	}
	openvg.BackgroundColor(s.Bg)
	openvg.FillColor(s.Fg)

	var x, y float64
	var fontsize int

	// every list in the slide
	var offset float64
	for _, l := range s.List {
		x, y, fontsize = dimen(d.Canvas, l.Xp, l.Yp, l.Sp)
		if l.Type == "bullet" {
			offset = 1.2 * float64(fontsize)
		} else {
			offset = 0
		}
		for _, li := range l.Li {
			if l.Type == "bullet" {
				boffset := float64(fontsize)/2
				openvg.Circle(x, y+boffset, boffset)
			}
			switch l.Align {
			case "center", "middle", "mid":
				openvg.TextMid(x+offset, y, li, "sans", fontsize)
			case "right", "end":
				openvg.TextEnd(x+offset, y, li, "sans", fontsize)
			default:
				openvg.Text(x+offset, y, li, "sans", fontsize)
			}
			y -= float64(fontsize) * 2.0
		}
	}

	// every text in the slide
	var font string
	for _, t := range s.Text {
		x, y, fontsize = dimen(d.Canvas, t.Xp, t.Yp, t.Sp)
		if t.Type == "code" {
			font = "mono"
		} else {
			font = "sans"
		}
		td := strings.Split(t.Tdata, "\n")
		for _, txt := range td {
			switch t.Align {
			case "center", "middle", "mid":
				openvg.TextMid(x, y, txt, font, fontsize)
			case "right", "end":
				openvg.TextEnd(x, y, txt, font, fontsize)
			default:
				openvg.Text(x, y, txt, font, fontsize)
			}
			y -= float64(fontsize) * 1.8 
		}
	}

	// every image in the slide
	for _, im := range s.Image {
		x = (im.Xp / 100) * float64(d.Canvas.Width)
		y = (im.Yp / 100) * float64(d.Canvas.Height)
		openvg.Image(x-float64(im.Width/2), y-float64(im.Height/2), im.Width, im.Height, im.Name)
	}
	openvg.End()
}

// interact controls the display of the deck 
func interact(d Deck) {
	firstslide := 0
	lastslide := len(d.Slide) - 1
	n := firstslide
	r := bufio.NewReader(os.Stdin)

	var cmd = "0"

	for ; rune(cmd[0]) != 'q'; cmd = readcmd(r) {
		switch rune(cmd[0]) {
		case '0', '1':
			n = firstslide
			showslide(d, n)

		case '*':
			n = lastslide
			showslide(d, n)

		case '+', 'n', '\n', '\t':
			n++
			if n > lastslide {
				n = firstslide
			}
			showslide(d, n)

		case '-', 'p':
			n--
			if n < firstslide {
				n = lastslide
			}
			showslide(d, n)

		case '/':
			if len(cmd) > 2 {
				ns := searchdeck(d, n, cmd[1:len(cmd)-1])
				println("searchdeck", ns)
				if ns >= 0 {
					showslide(d, ns)
				}
			}
		}
	}
}

// searchdeck searches the deck for the specified text, returning the slide number if found
func searchdeck(d Deck, n int, s string) int {
	if n < 0 || n > len(d.Slide)-1 {
		return -1
	}
	// for every slide...
	for i := 0; i < len(d.Slide); i++ {
		// search lists
		for _, l := range d.Slide[i].List {
			for _, ll := range l.Li {
				// fmt.Printf("searching %s in list %v\n", s, ll)
				if strings.Contains(ll, s) {
					return i
				}
			}
		}
		// search text
		for _, t := range d.Slide[i].Text {
			// fmt.Printf("searching %s in text %v\n", s, t.Tdata)
			if strings.Contains(t.Tdata, s) {
				return i
			}
		}
	}
	return -1
}
// readcmd reads interaction commands
func readcmd(r *bufio.Reader) string {
	s, err := r.ReadBytes('\n')
	if err != nil {
		return "q"
	}
	return string(s)
}

// for every file, make a deck
func main() {
	if len(os.Args) > 1 {
		for _, f := range os.Args[1:] {
			dodeck(f)
		}
	}
}
