// twh: time, weather, headlines
package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
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
	APIkey     = "/-api-key-/"
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

// regionFill colors a rectangular region, and sets the fill color for subsequent text
func regionFill(x, y, w, h openvg.VGfloat, color string) {
	openvg.FillColor(color)
	openvg.Rect(x, y, w, h)
	openvg.FillColor("white")
}

// headlines retrieves data from the AP API, decodes and formats it.
func headlines(w, h openvg.VGfloat, name string) {
	var data Feed
	url, ok := feeds[name]
	if !ok {
		fmt.Fprintf(os.Stderr, "%s headlines not supported\n", name)
		return
	}
	r, err := netread(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer r.Close()
	err = xml.NewDecoder(r).Decode(&data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	regionFill(0, 0, w, h*.50, "steelblue")
	headsize := w / 65
	x := w * 0.20
	y := h * 0.45
	spacing := headsize * 2.0
	for _, e := range data.Entries {
		openvg.Text(x, y, e.Title, "sans", int(headsize))
		y = y - spacing
	}
	openvg.End()
}

// weather retrieves data from the forecast.io API, decodes and formats it.
func weather(w, h openvg.VGfloat, latlong string) {
	var data Forecast
	r, err := netread(weatherURL + APIkey + latlong + options)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Weather read error (%v)\n", err)
		return
	}
	defer r.Close()
	err = json.NewDecoder(r).Decode(&data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	x := w * 0.05
	y := h * 0.70
	wsize := w / 20
	spacing := wsize * 2.0
	w1 := int(wsize)
	w2 := w1 / 2
	w3 := w1 / 3
	temp := fmt.Sprintf("%0.f°", data.Currently.Temperature)
	tw := openvg.TextWidth(temp, "sans", w1)

	regionFill(0, h*0.50, w*.50, h*.50, "steelblue")
	openvg.Text(x, y, temp, "sans", w1)
	if data.Currently.Temperature-data.Currently.FeelsLike > 1 {
		openvg.Text(x+tw+10, y, fmt.Sprintf("(feels like %0.f°)", data.Currently.FeelsLike), "sans", w3)
	}
	openvg.Text(x, y+spacing, data.Currently.Summary, "sans", w2)
	if data.Currently.PrecipProb > 0 {
		openvg.Text(x, y-(spacing/2), fmt.Sprintf("%0.f%% Chance of precipitation", data.Currently.PrecipProb*100), "sans", w3)
	}
	openvg.End()
}

// clock formats the current time
func clock(w, h openvg.VGfloat) {
	regionFill(w*0.50, h*0.50, w*.50, h*.50, "steelblue")
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
	pause := time.Duration(60-time.Now().Second()) * time.Second
	regionFill(0, 0, w, h, "steelblue")
	openvg.TextMid(w/2, h/2, fmt.Sprintf("Syncing the clock (%v)", pause), "mono", dw/30)
	openvg.End()
	<-time.After(pause) // wait until the beginning of the minute
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
