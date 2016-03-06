// twh: time, weather, headlines
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ajstarks/openvg"
)

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

// NYTHeadlines is the headline info from the New York Times
type NYTHeadlines struct {
	Status     string   `json:"status"`
	Copyright  string   `json:"copyright"`
	NumResults int      `json:"num_results"`
	Results    []result `json:"results"`
}

type result struct {
	Section    string `json:"section"`
	Subsection string `json:"subsection"`
	Title      string `json:"title"`
	Abstract   string `json:"abstract"`
}

const (
	weatherfmt    = "https://api.forecast.io/forecast/%s/%s/?exclude=hourly,daily,minutely,flags"
	NYTfmt        = "http://api.nytimes.com/svc/news/v3/content/all/%s/.json?api-key=%s&limit=5"
	bgcolor       = "gray"
	textcolor     = "white"
	weatherAPIkey = "-api-key-"
	NYTAPIkey     = "-api-key-"
)

var fromHTML = strings.NewReplacer(
	"‘", "'",
	"’", "'",
	"—", "--",
	"&#8216;", "'",
	"&#8217;", "'",
	"&#8220;", `"`,
	"&#8221;", `"`,
	"&lsquo;", "'",
	"&rsquo;", "'",
	"&ndash;", "-",
	"&mdash;", "--",
	"&amp;", "&")

// netread derefernces a URL, returning the Reader, with an error
func netread(url string) (io.ReadCloser, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to retreive network data for %s (%s)", url, resp.Status)
	}
	return resp.Body, nil
}

// countdown shows a countdown display to the top of minute
func countdown(w, h openvg.VGfloat) {
	tick := time.NewTicker(1 * time.Second)
	ty := h / 2
	th := h / 20
	size := w / 70
	for delay := 60 - time.Now().Second(); delay > 0; delay-- {
		select {
		case <-tick.C:
			tx := w * (openvg.VGfloat(60-delay) / 60)
			openvg.BackgroundColor(bgcolor)
			openvg.FillColor("black")
			openvg.TextMid(tx, ty+(size*2.5), fmt.Sprintf("%d", delay), "sans", int(size))
			openvg.Rect(0, ty, w, th)
			openvg.FillColor("white")
			openvg.Rect(tx, ty, w-tx, th)
			openvg.End()
		}
	}
	openvg.BackgroundColor(bgcolor)
}

// regionFill colors a rectangular region, and sets the fill color for subsequent text
func regionFill(x, y, w, h openvg.VGfloat, color string) {
	openvg.FillColor(color)
	openvg.Rect(x, y, w, h)
	openvg.FillColor(textcolor)
}

// gerror makes a graphical error message
func gerror(x, y, w, h openvg.VGfloat, s string) {
	regionFill(x, y, w, h, bgcolor)
	openvg.TextMid(x+w/2, y+h/2, s, "sans", int(w/20))
	openvg.End()
}

// headlines retrieves data from the New York Times API, decodes and displays it.
func headlines(w, h openvg.VGfloat, section string) {
	r, err := netread(fmt.Sprintf(NYTfmt, section, NYTAPIkey))
	if err != nil {
		fmt.Fprintf(os.Stderr, "headline read error: %v\n", err)
		gerror(0, 0, w, h*0.5, "no headlines")
		return
	}
	defer r.Close()
	var data NYTHeadlines
	err = json.NewDecoder(r).Decode(&data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decode: %v\n", err)
		gerror(0, 0, w, h*0.5, "no headlines")
		return
	}
	x := w * 0.50
	y := h * 0.10
	regionFill(0, 0, w, h*.50, bgcolor)
	headsize := w / 70
	spacing := headsize * 2.0
	for i := len(data.Results) - 1; i >= 0; i-- {
		openvg.TextMid(x, y, fromHTML.Replace(data.Results[i].Title), "sans", int(headsize))
		y = y + spacing
	}
	openvg.Image(w*0.05, 15, 30, 30, "poweredby_nytimes_30a.png")
	openvg.End()
}

// fog shows the fog icon
func fog(x, y, w, h openvg.VGfloat, color string) {
	radius := w / 3
	r2 := radius * 1.8
	openvg.FillColor(color, 0.5)
	openvg.Circle(x+w*0.25, y+h*0.25, radius)
	openvg.Circle(x+w*0.30, y+h*0.45, radius)
	openvg.Circle(x+w*0.60, y+h*0.40, r2)
}

// cloud shows the cloudy icon
func cloud(x, y, w, h openvg.VGfloat, color string) {
	radius := w / 3
	r2 := radius * 1.8
	openvg.FillColor(color)
	openvg.Circle(x+w*0.25, y+h*0.25, radius)
	openvg.Circle(x+w*0.30, y+h*0.45, radius)
	openvg.Circle(x+w*0.60, y+h*0.40, r2)
}

// flake shows the snowflake icon
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

// drop shows the raindrop icon
func drop(x, y, w, h openvg.VGfloat, color string) {
	openvg.FillColor(color)
	openvg.Ellipse(x+(w/2), y+(h*0.40), w*0.52, h*0.65)
	xp := []openvg.VGfloat{x + (w / 2), x + (w * 0.25), x + (w * 0.75)}
	yp := []openvg.VGfloat{y + h, y + (h / 2), y + (h / 2)}
	openvg.Polygon(xp, yp)
}

// rain shows a raindrops
func rain(x, y, w, h openvg.VGfloat, color string) {
	for i := 0; i < 20; i++ {
		rx := openvg.VGfloat(rand.Float64())
		ry := openvg.VGfloat(rand.Float64())
		drop(x+w*rx, y+h*ry, w/6, h/6, color)
	}
}

// snow shows the snow icon
func snow(x, y, w, h openvg.VGfloat, color string) {
	for i := 0; i < 20; i++ {
		rx := openvg.VGfloat(rand.Float64())
		ry := openvg.VGfloat(rand.Float64())
		flake(x+w*rx, y+h*ry, w/6, h/6, color)
	}
}

// sun shows the icon for clear weather
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

// moon shows the icon for clear weather at night
func moon(x, y, w, h openvg.VGfloat, bg, fg string) {
	cx := x + w/2
	cy := y + h/2
	w2 := w / 2
	openvg.FillColor(fg)
	openvg.Circle(cx, cy, w2)
	openvg.FillColor(bg)
	openvg.Circle(x+w*0.65, cy, w2)
}

// pcloud shows the icon for partly cloudy
func pcloud(x, y, w, h openvg.VGfloat, color string) {
	sun(x+w*.2, y+h*.33, w*.7, h*.7, "orange")
	cloud(x, y, w, h, color)
}

// npcloud shows the partly cloudy icon at night
func npcloud(x, y, w, h openvg.VGfloat, ccolor, mcolor string) {
	cloud(x, y, w, h, ccolor)
	moon(x+w*0.2, y+h*0.05, w*.7, h*.7, ccolor, mcolor)
}

// wind shows the windy icon
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
	r, err := netread(fmt.Sprintf(weatherfmt, weatherAPIkey, latlong))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Weather read error %v\n", err)
		gerror(0, h*0.5, w*0.5, h*0.5, "no weather")
		return
	}
	defer r.Close()
	var data Forecast
	err = json.NewDecoder(r).Decode(&data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		gerror(0, h*0.5, w*0.5, h*0.5, "no weather")
		return
	}
	x := w * 0.05
	y := h * 0.70
	wsize := w / 20
	spacing := wsize * 2.0
	w1 := int(wsize)
	w2 := w1 / 2
	w3 := w1 / 4
	c := data.Currently
	temp := fmt.Sprintf("%0.f°", c.Temperature)
	tw := openvg.TextWidth(temp, "sans", w1)
	regionFill(0, h*0.50, w*.50, h*.50, bgcolor)
	openvg.Text(x, y, temp, "sans", w1)
	if c.Temperature-c.FeelsLike > 1 {
		openvg.Text(x, y-(spacing/3), fmt.Sprintf("(feels like %0.f°)", c.FeelsLike), "sans", w3)
	}
	openvg.Text(x, y+spacing, c.Summary, "sans", w2)
	if c.PrecipProb > 0 {
		openvg.Text(x, y-(spacing*.6), fmt.Sprintf("%0.f%% Chance of precipitation", c.PrecipProb*100), "sans", w3)
	}
	iw := w / 10
	ih := iw // h / 10
	ix := x + tw + w*0.010

	switch c.Icon {
	case "clear-day":
		sun(ix, y, iw, ih, "orange")
	case "clear-night":
		moon(ix, y, iw, ih, bgcolor, textcolor)
	case "rain":
		rain(ix, y, iw, ih, "skyblue")
	case "snow":
		snow(ix, y, iw, ih, textcolor)
	case "wind":
		wind(ix, y, iw, ih, bgcolor, textcolor)
	case "fog":
		fog(ix, y, iw, ih, textcolor)
	case "cloudy":
		cloud(ix, y, iw, ih, textcolor)
	case "partly-cloudy-day":
		pcloud(ix, y, iw, ih, textcolor)
	case "partly-cloudy-night":
		npcloud(ix, y, iw, ih, "darkgray", textcolor)
	}
	openvg.End()
}

// clock displays the current time
func clock(w, h openvg.VGfloat) {
	regionFill(w*0.50, h*0.50, w*.50, h*.50, bgcolor)
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
	var headtype = flag.String("h", "u.s.", "headline type (arts, health, sports, science, technology, u.s., world)")
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
	headticker := time.NewTicker(10 * time.Minute)
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
