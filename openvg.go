// High-level 2D vector graphics library built on OpenVG
package openvg

/*
#cgo CFLAGS:   -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads 
#cgo LDFLAGS:  -L/opt/vc/lib -lGLESv2 -ljpeg
#include "VG/openvg.h"
#include "VG/vgu.h"
#include "EGL/egl.h"
#include "GLES/gl.h"
#include "fontinfo.h" // font information
#include "shapes.h"   // C API
*/
import "C"
import "unsafe"

// RGB triple
type RGB struct {
	Red, Green, Blue uint8
}

// colornames maps SVG color names to RGB triples
var colornames = map[string]RGB{
	"aliceblue":            {240, 248, 255},
	"antiquewhite":         {250, 235, 215},
	"aqua":                 {0, 255, 255},
	"aquamarine":           {127, 255, 212},
	"azure":                {240, 255, 255},
	"beige":                {245, 245, 220},
	"bisque":               {255, 228, 196},
	"black":                {0, 0, 0},
	"blanchedalmond":       {255, 235, 205},
	"blue":                 {0, 0, 255},
	"blueviolet":           {138, 43, 226},
	"brown":                {165, 42, 42},
	"burlywood":            {222, 184, 135},
	"cadetblue":            {95, 158, 160},
	"chartreuse":           {127, 255, 0},
	"chocolate":            {210, 105, 30},
	"coral":                {255, 127, 80},
	"cornflowerblue":       {100, 149, 237},
	"cornsilk":             {255, 248, 220},
	"crimson":              {220, 20, 60},
	"cyan":                 {0, 255, 255},
	"darkblue":             {0, 0, 139},
	"darkcyan":             {0, 139, 139},
	"darkgoldenrod":        {184, 134, 11},
	"darkgray":             {169, 169, 169},
	"darkgreen":            {0, 100, 0},
	"darkgrey":             {169, 169, 169},
	"darkkhaki":            {189, 183, 107},
	"darkmagenta":          {139, 0, 139},
	"darkolivegreen":       {85, 107, 47},
	"darkorange":           {255, 140, 0},
	"darkorchid":           {153, 50, 204},
	"darkred":              {139, 0, 0},
	"darksalmon":           {233, 150, 122},
	"darkseagreen":         {143, 188, 143},
	"darkslateblue":        {72, 61, 139},
	"darkslategray":        {47, 79, 79},
	"darkslategrey":        {47, 79, 79},
	"darkturquoise":        {0, 206, 209},
	"darkviolet":           {148, 0, 211},
	"deeppink":             {255, 20, 147},
	"deepskyblue":          {0, 191, 255},
	"dimgray":              {105, 105, 105},
	"dimgrey":              {105, 105, 105},
	"dodgerblue":           {30, 144, 255},
	"firebrick":            {178, 34, 34},
	"floralwhite":          {255, 250, 240},
	"forestgreen":          {34, 139, 34},
	"fuchsia":              {255, 0, 255},
	"gainsboro":            {220, 220, 220},
	"ghostwhite":           {248, 248, 255},
	"gold":                 {255, 215, 0},
	"goldenrod":            {218, 165, 32},
	"gray":                 {128, 128, 128},
	"green":                {0, 128, 0},
	"greenyellow":          {173, 255, 47},
	"grey":                 {128, 128, 128},
	"honeydew":             {240, 255, 240},
	"hotpink":              {255, 105, 180},
	"indianred":            {205, 92, 92},
	"indigo":               {75, 0, 130},
	"ivory":                {255, 255, 240},
	"khaki":                {240, 230, 140},
	"lavender":             {230, 230, 250},
	"lavenderblush":        {255, 240, 245},
	"lawngreen":            {124, 252, 0},
	"lemonchiffon":         {255, 250, 205},
	"lightblue":            {173, 216, 230},
	"lightcoral":           {240, 128, 128},
	"lightcyan":            {224, 255, 255},
	"lightgoldenrodyellow": {250, 250, 210},
	"lightgray":            {211, 211, 211},
	"lightgreen":           {144, 238, 144},
	"lightgrey":            {211, 211, 211},
	"lightpink":            {255, 182, 193},
	"lightsalmon":          {255, 160, 122},
	"lightseagreen":        {32, 178, 170},
	"lightskyblue":         {135, 206, 250},
	"lightslategray":       {119, 136, 153},
	"lightslategrey":       {119, 136, 153},
	"lightsteelblue":       {176, 196, 222},
	"lightyellow":          {255, 255, 224},
	"lime":                 {0, 255, 0},
	"limegreen":            {50, 205, 50},
	"linen":                {250, 240, 230},
	"magenta":              {255, 0, 255},
	"maroon":               {128, 0, 0},
	"mediumaquamarine":     {102, 205, 170},
	"mediumblue":           {0, 0, 205},
	"mediumorchid":         {186, 85, 211},
	"mediumpurple":         {147, 112, 219},
	"mediumseagreen":       {60, 179, 113},
	"mediumslateblue":      {123, 104, 238},
	"mediumspringgreen":    {0, 250, 154},
	"mediumturquoise":      {72, 209, 204},
	"mediumvioletred":      {199, 21, 133},
	"midnightblue":         {25, 25, 112},
	"mintcream":            {245, 255, 250},
	"mistyrose":            {255, 228, 225},
	"moccasin":             {255, 228, 181},
	"navajowhite":          {255, 222, 173},
	"navy":                 {0, 0, 128},
	"oldlace":              {253, 245, 230},
	"olive":                {128, 128, 0},
	"olivedrab":            {107, 142, 35},
	"orange":               {255, 165, 0},
	"orangered":            {255, 69, 0},
	"orchid":               {218, 112, 214},
	"palegoldenrod":        {238, 232, 170},
	"palegreen":            {152, 251, 152},
	"paleturquoise":        {175, 238, 238},
	"palevioletred":        {219, 112, 147},
	"papayawhip":           {255, 239, 213},
	"peachpuff":            {255, 218, 185},
	"peru":                 {205, 133, 63},
	"pink":                 {255, 192, 203},
	"plum":                 {221, 160, 221},
	"powderblue":           {176, 224, 230},
	"purple":               {128, 0, 128},
	"red":                  {255, 0, 0},
	"rosybrown":            {188, 143, 143},
	"royalblue":            {65, 105, 225},
	"saddlebrown":          {139, 69, 19},
	"salmon":               {250, 128, 114},
	"sandybrown":           {244, 164, 96},
	"seagreen":             {46, 139, 87},
	"seashell":             {255, 245, 238},
	"sienna":               {160, 82, 45},
	"silver":               {192, 192, 192},
	"skyblue":              {135, 206, 235},
	"slateblue":            {106, 90, 205},
	"slategray":            {112, 128, 144},
	"slategrey":            {112, 128, 144},
	"snow":                 {255, 250, 250},
	"springgreen":          {0, 255, 127},
	"steelblue":            {70, 130, 180},
	"tan":                  {210, 180, 140},
	"teal":                 {0, 128, 128},
	"thistle":              {216, 191, 216},
	"tomato":               {255, 99, 71},
	"turquoise":            {64, 224, 208},
	"violet":               {238, 130, 238},
	"wheat":                {245, 222, 179},
	"white":                {255, 255, 255},
	"whitesmoke":           {245, 245, 245},
	"yellow":               {255, 255, 0},
	"yellowgreen":          {154, 205, 50},
}

// Init initializes the graphics subsystem
func Init() (int, int) {
	var rh, rw C.int
	C.init(&rw, &rh)
	return int(rw), int(rh)
}

// Finish shuts down the graphics subsystem
func Finish() {
	C.finish()
}

// Background clears the screen with the specified solid background color using RGB triples
func Background(r, g, b uint8) {
	C.Background(C.uint(r), C.uint(g), C.uint(b))
}

// BackgroundRGB clears the screen with the specified background color using a RGBA quad
func BackgroundRGB(r, g, b uint8, alpha float64) {
	C.BackgroundRGB(C.uint(r), C.uint(g), C.uint(b), C.VGfloat(alpha))
}

// BackgroundColor sets the background color
func BackgroundColor(s string, alpha ...float64) {
	c := colorlookup(s)
	if len(alpha) == 0 {
		BackgroundRGB(c.Red, c.Green, c.Blue, 1)
	} else {
		BackgroundRGB(c.Red, c.Green, c.Blue, alpha[0])
	}
}

// FillRGB sets the fill color, using RGB triples
func FillRGB(r, g, b uint8, alpha float64) {
	C.Fill(C.uint(r), C.uint(g), C.uint(b), C.VGfloat(alpha))
}

// StrokeRGB sets the stroke color, using RGB triples
func StrokeRGB(r, g, b uint8, alpha float64) {
	C.Stroke(C.uint(r), C.uint(g), C.uint(b), C.VGfloat(alpha))
}

// StrokeWidth sets the stroke width
func StrokeWidth(w float64) {
	C.StrokeWidth(C.VGfloat(w))
}

// colorlookup returns a RGB triple corresponding to the named color,
// or black on error.
func colorlookup(s string) RGB {
	var black = RGB{0, 0, 0}
	c, ok := colornames[s]
	if ok {
		return c
	}
	return black
}

// FillColor sets the fill color using names to specify the color, optionally applying alpha.
func FillColor(s string, alpha ...float64) {
	c := colorlookup(s)
	if len(alpha) == 0 {
		FillRGB(c.Red, c.Green, c.Blue, 1)
	} else {
		FillRGB(c.Red, c.Green, c.Blue, alpha[0])
	}
}

// StrokeColor sets the fill color using names to specify the color, optionally applying alpha.
func StrokeColor(s string, alpha ...float64) {
	c := colorlookup(s)
	if len(alpha) == 0 {
		StrokeRGB(c.Red, c.Green, c.Blue, 1)
	} else {
		StrokeRGB(c.Red, c.Green, c.Blue, alpha[0])
	}
}

// Start begins a picture
func Start(w, h int, color ...uint8) {
	C.Start(C.int(w), C.int(h))
	if len(color) == 3 {
		Background(color[0], color[1], color[2])
	}
}

// Startcolor begins the picture with the specified color background
func StartColor(w, h int, color string, alpha ...float64) {
	C.Start(C.int(w), C.int(h))
	BackgroundColor(color, alpha...)
}

// End ends the picture
func End() {
	C.End()
}

// Image places an image with dimensions (w,h) at (x,y)
func Image(x, y float64, w, h int, s string) {
	fw := float64(w)
	fh := float64(h)
	FillColor("lightgray")
	Rect(x, y, fw, fh)
	StrokeColor("white")
	StrokeWidth(1)
	Line(x, y, x+fw, y+fh)
	Line(x+fw, y, x, y+fh)
	FillColor("maroon")
	StrokeWidth(0)
	TextMid(x+fw/2, y+fh/2, s, "sans", w/25)
}

// Line draws a line between two points
func Line(x1, y1, x2, y2 float64, style ...string) {
	C.Line(C.VGfloat(x1), C.VGfloat(y1), C.VGfloat(x2), C.VGfloat(y2))
}

// Rect draws a rectangle at (x,y) with dimesions (w,h)
func Rect(x, y, w, h float64, style ...string) {
	C.Rect(C.VGfloat(x), C.VGfloat(y), C.VGfloat(w), C.VGfloat(h))
}

// Rect draws a rounded rectangle at (x,y) with dimesions (w,h).
// the corner radii are at (rw, rh)
func Roundrect(x, y, w, h, rw, rh float64, style ...string) {
	C.Roundrect(C.VGfloat(x), C.VGfloat(y), C.VGfloat(w), C.VGfloat(h), C.VGfloat(rw), C.VGfloat(rh))
}

// Ellipse draws an ellipse at (x,y) with dimensions (w,h)
func Ellipse(x, y, w, h float64, style ...string) {
	C.Ellipse(C.VGfloat(x), C.VGfloat(y), C.VGfloat(w), C.VGfloat(h))
}

// Circle draws a circle centered at (x,y), with radius r
func Circle(x, y, r float64, style ...string) {
	C.Circle(C.VGfloat(x), C.VGfloat(y), C.VGfloat(r))
}

// Qbezier draws a quadratic bezier curve with extrema (sx, sy) and (ex, ey)
// Control points are at (cx, cy)
func Qbezier(sx, sy, cx, cy, ex, ey float64, style ...string) {
	C.Qbezier(C.VGfloat(sx), C.VGfloat(sy), C.VGfloat(cx), C.VGfloat(cy), C.VGfloat(ex), C.VGfloat(ey))
}

// Cbezier draws a cubic bezier curve with extrema (sx, sy) and (ex, ey).
// Control points at (cx, cy) and (px, py)
func Cbezier(sx, sy, cx, cy, px, py, ex, ey float64, style ...string) {
	C.Cbezier(C.VGfloat(sx), C.VGfloat(sy), C.VGfloat(cx), C.VGfloat(cy), C.VGfloat(px), C.VGfloat(py), C.VGfloat(ex), C.VGfloat(ey))
}

// Arc draws an arc at (x,y) with dimensions (w,h).
// the arc starts at the angle sa, extended to aext
func Arc(x, y, w, h, sa, aext float64, style ...string) {
	C.Arc(C.VGfloat(x), C.VGfloat(y), C.VGfloat(w), C.VGfloat(h), C.VGfloat(sa), C.VGfloat(aext))
}

// poly converts coordinate slices
func poly(x, y []float64) (*C.VGfloat, *C.VGfloat, C.VGint) {
	size := len(x)
	if size != len(y) {
		return nil, nil, 0
	}
	px := make([]C.VGfloat, size)
	py := make([]C.VGfloat, size)
	for i := 0; i < size; i++ {
		px[i] = C.VGfloat(x[i])
		py[i] = C.VGfloat(y[i])
	}
	return &px[0], &py[0], C.VGint(size)
}

// Polygon draws a polygon with coordinates in x, y
func Polygon(x, y []float64, style ...string) {
	px, py, np := poly(x, y)
	if np > 0 {
		C.Polygon(px, py, np)
	}
}

// Polyline draws a polyline with coordinates in x, y
func Polyline(x, y []float64, style ...string) {
	px, py, np := poly(x, y)
	if np > 0 {
		C.Polyline(px, py, np)
	}
}

// select font specifies the font by generic name
func selectfont(s string) C.Fontinfo {
	switch s {
	case "sans":
		return C.SansTypeface
	case "serif":
		return C.SerifTypeface
	case "mono":
		return C.MonoTypeface
	}
	return C.SerifTypeface
}

// Text draws text whose aligment begins (x,y)
func Text(x, y float64, s string, font string, size int, style ...string) {
	t := C.CString(s)
	defer C.free(unsafe.Pointer(t))
	C.Text(C.VGfloat(x), C.VGfloat(y), t, selectfont(font), C.int(size))
}

// TextMid draws text centeRed at (x,y)
func TextMid(x, y float64, s string, font string, size int, style ...string) {
	t := C.CString(s)
	defer C.free(unsafe.Pointer(t))
	C.TextMid(C.VGfloat(x), C.VGfloat(y), t, selectfont(font), C.int(size))
}

// TextEnd draws text end-aligned at (x,y)
func TextEnd(x, y float64, s string, font string, size int, style ...string) {
	t := C.CString(s)
	defer C.free(unsafe.Pointer(t))
	C.TextEnd(C.VGfloat(x), C.VGfloat(y), t, selectfont(font), C.int(size))
}

// TextWidth returns the length of text at a specified font and size
func TextWidth(s string, font string, size float64) float64 {
	t := C.CString(s)
	defer C.free(unsafe.Pointer(t))
	return float64(C.textwidth(t, selectfont(font), C.VGfloat(size)))
}

// Translate translates the coordinate system to (x,y)
func Translate(x, y float64) {
	C.Translate(C.VGfloat(x), C.VGfloat(y))
}

// Rotate rotates the coordinate system around the specifed angle
func Rotate(r float64) {
	C.Rotate(C.VGfloat(r))
}

// Shear warps the coordinate system by (x,y)
func Shear(x, y float64) {
	C.Shear(C.VGfloat(x), C.VGfloat(y))
}

// Scale scales the coordinate system by (x,y)
func Scale(x, y float64) {
	C.Scale(C.VGfloat(x), C.VGfloat(y))
}
