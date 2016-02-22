// twh: time, weather, headlines
package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ajstarks/openvg"
)

// Feed is the Atom feed from the Associated Press
type Feed struct {
	Title   string  `xml:"title"`
	Updated string  `xml:"updated"`
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Title string `xml:"title"`
}

// Forecast is weather information from forecast.io
type Forecast struct {
	Lat       float64 `json:"latitude"`
	Long      float64 `json:"longitude"`
	Currently Current `json:"currently"`
}

type Current struct {
	Summary     string  `json:"summary"`
	Icon        string  `json:"icon"`
	PrecipProb  float64 `json:"precipProbability"`
	Temperature float64 `json:"temperature"`
	FeelsLike   float64 `json:"apparentTemperature"`
}

const (
	weatherURL = "https://api.forecast.io/forecast"
	APIkey     = "/-api-key-here/"
	options    = "?exclude=hourly,daily,minutely,flags"
)

// Mapping of headline types to AP API URLs
var feeds = map[string]string{
	"top":       "http://hosted2.ap.org/atom/APDEFAULT/3d281c11a96b4ad082fe88aa0db04305",
	"us":        "http://hosted2.ap.org/atom/APDEFAULT/386c25518f464186bf7a2ac026580ce7",
	"world":     "http://hosted2.ap.org/atom/APDEFAULT/cae69a7523db45408eeb2b3a98c0c9c5",
	"politics":  "http://hosted2.ap.org/atom/APDEFAULT/89ae8247abe8493fae24405546e9a1aa",
	"business":  "http://hosted2.ap.org/atom/APDEFAULT/f70471f764144b2fab526d39972d37b3",
	"sports":    "http://hosted2.ap.org/atom/APDEFAULT/347875155d53465d95cec892aeb06419",
	"tech":      "http://hosted2.ap.org/atom/APDEFAULT/495d344a0d10421e9baa8ee77029cfbd",
	"entertain": "http://hosted2.ap.org/atom/APDEFAULT/4e67281c3f754d0696fbfdee0f3f1469",
	"health":    "http://hosted2.ap.org/atom/APDEFAULT/bbd825583c8542898e6fa7d440b9febc",
	"science":   "http://hosted2.ap.org/atom/APDEFAULT/b2f0ca3a594644ee9e50a8ec4ce2d6de",
	"strange":   "http://hosted2.ap.org/atom/APDEFAULT/aa9398e6757a46fa93ed5dea7bd3729e",
}

// netread derefernces a URL, returning the Reader, with an error
func netread(url string) (io.ReadCloser, error) {
	conn := &http.Client{Timeout: 1 * time.Minute}
	resp, err := conn.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to retreive network data for %s (%s)", url, resp.Status)
	}
	return resp.Body, nil
}

func countdown(w, h openvg.VGfloat) {
	tick := time.NewTicker(1 * time.Second)
	for delay := 60 - time.Now().Second(); delay > 0; delay-- {
		select {
		case <-tick.C:
			openvg.FillColor("white")
			openvg.BackgroundColor("gray")
			openvg.TextMid(w/2, h/2, fmt.Sprintf("starting in %2d seconds", delay), "sans", int(w/20))
			openvg.End()
		}
	}
	openvg.BackgroundColor("gray")
}

// regionFill colors a rectangular region, and sets the fill color for subsequent text
func regionFill(x, y, w, h openvg.VGfloat, color string) {
	openvg.FillColor(color)
	openvg.Rect(x, y, w, h)
	openvg.FillColor("white")
}

// headlines retrieves data from the Associated Press API, decodes and displays it.
func headlines(w, h openvg.VGfloat, name string) {
	var data Feed
	url, ok := feeds[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "%s headlines not supported\n", name)
		return
	}
	r, err := netread(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	defer r.Close()
	err = xml.NewDecoder(r).Decode(&data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	regionFill(0, 0, w, h*.50, "gray")
	headsize := w / 50
	x := w * 0.10
	y := h * 0.45
	spacing := headsize * 2.0
	for _, e := range data.Entries {
		openvg.Text(x, y, e.Title, "sans", int(headsize))
		y = y - spacing
	}
	openvg.End()
}

func fog(x, y, w, h openvg.VGfloat, color string) {
	radius := w / 3
	r2 := radius * 1.8
	openvg.FillColor(color, 0.5)
	openvg.Circle(x+w*0.25, y+h*0.25, radius)
	openvg.Circle(x+w*0.30, y+h*0.45, radius)
	openvg.Circle(x+w*0.60, y+h*0.40, r2)
}

func cloud(x, y, w, h openvg.VGfloat, color string) {
	radius := w / 3
	r2 := radius * 1.8
	openvg.FillColor(color)
	openvg.Circle(x+w*0.25, y+h*0.25, radius)
	openvg.Circle(x+w*0.30, y+h*0.45, radius)
	openvg.Circle(x+w*0.60, y+h*0.40, r2)
}

func flake(x, y, w, h openvg.VGfloat, color string) {
	cx := x + (w / 2)
	cy := y + (h / 2)
	r := w * 0.30
	openvg.StrokeColor(color)
	openvg.StrokeWidth(w / 20)
	for t := 0.0; t < 2*math.Pi; t += math.Pi / 4 {
		c := openvg.VGfloat(math.Cos(t))
		s := openvg.VGfloat(math.Sin(t))
		x1 := (r * c) + cx
		y1 := (r * s) + cy
		openvg.Line(cx, cy, x1, y1)
	}
	openvg.StrokeWidth(0)
}

func drop(x, y, w, h openvg.VGfloat, color string) {
	openvg.FillColor(color)
	openvg.Ellipse(x+(w/2), y+(h*0.40), w*0.52, h*0.65)
	xp := []openvg.VGfloat{x + (w / 2), x + (w * 0.25), x + (w * 0.75)}
	yp := []openvg.VGfloat{y + h, y + (h / 2), y + (h / 2)}
	openvg.Polygon(xp, yp)
}

func rain(x, y, w, h openvg.VGfloat, color string) {
	for i := 0; i < 20; i++ {
		rx := openvg.VGfloat(rand.Float64())
		ry := openvg.VGfloat(rand.Float64())
		drop(x+w*rx, y+h*ry, w/6, h/6, color)
	}
}

func snow(x, y, w, h openvg.VGfloat, color string) {
	for i := 0; i < 20; i++ {
		rx := openvg.VGfloat(rand.Float64())
		ry := openvg.VGfloat(rand.Float64())
		flake(x+w*rx, y+h*ry, w/6, h/6, color)
	}
}

func sun(x, y, w, h openvg.VGfloat, color string) {
	cx := x + (w / 2)
	cy := y + (h / 2)
	r0 := w * 0.50
	r1 := w * 0.45
	r2 := w * 0.30
	openvg.FillColor(color)
	openvg.Circle(cx, cy, r0)
	openvg.StrokeColor(color)
	openvg.StrokeWidth(w / 30)
	for t := 0.0; t < 2*math.Pi; t += math.Pi / 6 {
		c := openvg.VGfloat(math.Cos(t))
		s := openvg.VGfloat(math.Sin(t))
		x1 := (r1 * c) + cx
		y1 := (r1 * s) + cy
		x2 := (r2 * c) + cx
		y2 := (r2 * s) + cy
		openvg.Line(x1, y1, x2, y2)
	}
	openvg.StrokeWidth(0)
}

func moon(x, y, w, h openvg.VGfloat, bg, fg string) {
	cx := x + w/2
	cy := y + h/2
	w2 := w / 2
	openvg.FillColor(fg)
	openvg.Circle(cx, cy, w2)
	openvg.FillColor(bg)
	openvg.Circle(x+w*0.65, cy, w2)
}

func pcloud(x, y, w, h openvg.VGfloat, color string) {
	sun(x+w*.2, y+h*.35, w*.7, h*.7, "orange")
	cloud(x, y, w, h, color)
}

func npcloud(x, y, w, h openvg.VGfloat, ccolor, mcolor string) {
	cloud(x, y, w, h, ccolor)
	moon(x+w*0.2, y+h*0.05, w*.7, h*.7, ccolor, mcolor)
}

func wind(x, y, w, h openvg.VGfloat, bg, color string) {
	openvg.FillColor(bg, 0)
	openvg.StrokeWidth(w / 25)
	openvg.StrokeColor(color)
	openvg.Qbezier(x+w*0.10, y+h*0.8, x+w*0.50, y+h*0.60, x+w*0.9, y+h*0.85)
	openvg.Qbezier(x+w*0.10, y+h*0.5, x+w*0.55, y+h*0.30, x+w*0.9, y+h*0.55)
	openvg.Qbezier(x+w*0.10, y+h*0.2, x+w*0.60, y+h*0.10, x+w*0.9, y+h*0.35)
	openvg.StrokeWidth(0)
}

// weather retrieves data from the forecast.io API, decodes and displays it.
func weather(w, h openvg.VGfloat, latlong string) {
	var data Forecast
	r, err := netread(weatherURL + APIkey + latlong + options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Weather read error (%v)\n", err)
		return
	}
	defer r.Close()
	err = json.NewDecoder(r).Decode(&data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	x := w * 0.05
	y := h * 0.70
	wsize := w / 20
	spacing := wsize * 2.0
	w1 := int(wsize)
	w2 := w1 / 2
	w3 := w1 / 4
	temp := fmt.Sprintf("%0.f°", data.Currently.Temperature)
	tw := openvg.TextWidth(temp, "sans", w1)

	regionFill(0, h*0.50, w*.50, h*.50, "gray")
	openvg.Text(x, y, temp, "sans", w1)
	if data.Currently.Temperature-data.Currently.FeelsLike > 1 {
		openvg.Text(x, y-(spacing/3), fmt.Sprintf("(feels like %0.f°)", data.Currently.FeelsLike), "sans", w3)
	}
	openvg.Text(x, y+spacing, data.Currently.Summary, "sans", w2)
	if data.Currently.PrecipProb > 0 {
		openvg.Text(x, y-(spacing*.6), fmt.Sprintf("%0.f%% Chance of precipitation", data.Currently.PrecipProb*100), "sans", w3)
	}
	iw := w / 10
	ih := iw // h / 10
	ix := x + tw + w*0.010

	switch data.Currently.Icon {
	case "clear-day":
		sun(ix, y, iw, ih, "orange")
	case "clear-night":
		moon(ix, y, iw, ih, "gray", "white")
	case "rain":
		rain(ix, y, iw, ih, "skyblue")
	case "snow":
		snow(ix, y, iw, ih, "white")
	case "wind":
		wind(ix, y, iw, ih, "gray", "white")
	case "fog":
		fog(ix, y, iw, ih, "white")
	case "cloudy":
		cloud(ix, y, iw, ih, "white")
	case "partly-cloudy-day":
		pcloud(ix, y, iw, ih, "white")
	case "partly-cloudy-night":
		npcloud(ix, y, iw, ih, "darkgray", "white")
	}

	openvg.End()
}

// clock displays the current time
func clock(w, h openvg.VGfloat) {
	regionFill(w*0.50, h*0.50, w*.50, h*.50, "gray")
	clocksize := w / 20
	cs := int(clocksize)
	x := w * 0.95
	y := h * 0.70

	now := time.Now()
	openvg.TextEnd(x, y, now.Format("3:04 pm"), "sans", cs)
	openvg.TextEnd(x, y+(clocksize*2), now.Format("Monday January _2"), "sans", cs/2)
	openvg.End()
}

// show the current time, weather and headlines
func main() {
	var headtype = flag.String("h", "us", "headline type (business, entertain, health, politics, sports, science, strange, tech, top, us, world)")
	var location = flag.String("loc", "40.6213,-74.4395", "lat,long for weather")
	flag.Parse()

	// initial display
	dw, dh := openvg.Init()
	openvg.Start(dw, dh)
	w, h := openvg.VGfloat(dw), openvg.VGfloat(dh)
	countdown(w, h)
	openvg.End()
	clock(w, h)
	weather(w, h, *location)
	headlines(w, h, *headtype)

	// update on specific intervals, shutdown on interrupt
	dateticker := time.NewTicker(1 * time.Minute)
	weatherticker := time.NewTicker(5 * time.Minute)
	headticker := time.NewTicker(30 * time.Minute)
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	for {
		select {
		case <-dateticker.C:
			clock(w, h)
		case <-weatherticker.C:
			weather(w, h, *location)
		case <-headticker.C:
			headlines(w, h, *headtype)
		case <-sigint:
			openvg.Finish()
			os.Exit(0)
		}
	}
}
