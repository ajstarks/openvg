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

type list struct {
	Xp    float64  `xml:"xp,attr"`
	Yp    float64  `xml:"yp,attr"`
	Sp    float64  `xml:"sp,attr"`
	Type  string   `xml:"type,attr"`
	Align string   `xml:"align,attr"`
	Color string   `xml:"color,attr"`
	Li    []string `xml:"li"`
}

type text struct {
	Xp    float64 `xml:"xp,attr"`
	Yp    float64 `xml:"yp,attr"`
	Sp    float64 `xml:"sp,attr"`
	Type  string  `xml:"type,attr"`
	Align string  `xml:"align,attr"`
	Color string  `xml:"color,attr"`
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
	r, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	defer r.Close()

	var d Deck
	err = xml.NewDecoder(r).Decode(&d)
	if err != nil {
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

	// every image in the slide
	for _, im := range s.Image {
		x = (im.Xp / 100) * float64(d.Canvas.Width)
		y = (im.Yp / 100) * float64(d.Canvas.Height)
		openvg.Image(x-float64(im.Width/2), y-float64(im.Height/2), im.Width, im.Height, im.Name)
	}
	// every list in the slide
	var offset float64
	for _, l := range s.List {
		x, y, fontsize = dimen(d.Canvas, l.Xp, l.Yp, l.Sp)
		if l.Type == "bullet" {
			offset = 1.2 * float64(fontsize)
		} else {
			offset = 0
		}
		if l.Color != "" {
			openvg.FillColor(l.Color)
		} else {
			openvg.FillColor(s.Fg)
		}
		for _, li := range l.Li {
			if l.Type == "bullet" {
				boffset := float64(fontsize) / 2
				openvg.Circle(x, y+boffset, boffset)
			}
			showtext(x+offset, y, li, l.Align, "sans", fontsize)
			y -= float64(fontsize) * 2.0
		}
	}
	openvg.FillColor(s.Fg)

	// every text in the slide
	var font string
	for _, t := range s.Text {
		x, y, fontsize = dimen(d.Canvas, t.Xp, t.Yp, t.Sp)
		if t.Type == "code" {
			font = "mono"
		} else {
			font = "sans"
		}
		if t.Color != "" {
			openvg.FillColor(t.Color)
		} else {
			openvg.FillColor(s.Fg)
		}
		td := strings.Split(t.Tdata, "\n")
		for _, txt := range td {
			showtext(x, y, txt, t.Align, font, fontsize)
			y -= float64(fontsize) * 1.8
		}
	}
	openvg.FillColor(s.Fg)

	openvg.End()
}

// interact controls the display of the deck
func interact(d Deck) {
	firstslide := 0
	lastslide := len(d.Slide) - 1
	n := firstslide
	openvg.SaveTerm()
	openvg.RawTerm()
	defer openvg.RestoreTerm()
	r := bufio.NewReader(os.Stdin)

	var cmd byte = '0'

	for ; cmd != 'q'; cmd = readcmd(r) {
		switch cmd {
		case '0', '1':
			n = firstslide
			showslide(d, n)

		case '*':
			n = lastslide
			showslide(d, n)

		case '+', 'n', '\n', ' ':
			n++
			if n > lastslide {
				n = firstslide
			}
			showslide(d, n)

		case '-', 'p', 8, 127:
			n--
			if n < firstslide {
				n = lastslide
			}
			showslide(d, n)

		case '/':
			openvg.RestoreTerm()
			searchterm, err := r.ReadString('\n')
			openvg.RawTerm()
			if err != nil {
				continue
			}
			if len(searchterm) > 2 {
				ns := searchdeck(d, searchterm[0:len(searchterm)-1])
				if ns >= 0 {
					showslide(d, ns)
				}
			}
		}
	}
}

// searchdeck searches the deck for the specified text, returning the slide number if found
func searchdeck(d Deck, s string) int {
	// for every slide...
	for i := 0; i < len(d.Slide); i++ {
		// search lists
		for _, l := range d.Slide[i].List {
			for _, ll := range l.Li {
				if strings.Contains(ll, s) {
					return i
				}
			}
		}
		// search text
		for _, t := range d.Slide[i].Text {
			if strings.Contains(t.Tdata, s) {
				return i
			}
		}
	}
	return -1
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
