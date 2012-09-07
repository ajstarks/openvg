// package openvg is a high-level 2D Vector graphics library built on OpenVG
package openvg

/*
#cgo CFLAGS:   -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads 
#cgo LDFLAGS:  -L/opt/vc/lib -lGLESv2 -ljpeg
#include "VG/openvg.h"
#include "VG/vgu.h"
#include "EGL/egl.h"
#include "GLES/gl.h"
#include "fontinfo.h"                   // font data structure
#include "shapes.h"
*/
import "C"

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

// Background clears the screen with the specified background color
func Background(r, g, b uint8) {
	C.Background(C.uint(r), C.uint(g), C.uint(b))
}

// Fill sets the fill color
func Fill(r, g, b uint8, alpha float64) {
	C.Fill(C.uint(r), C.uint(g), C.uint(b), C.VGfloat(alpha))
}

// Stroke sets the stroke color
func Stroke(r, g, b uint8, alpha float64) {
	C.Stroke(C.uint(r), C.uint(g), C.uint(b), C.VGfloat(alpha))
}

// StrokeWidth sets the stroke width
func StrokeWidth(w float64) {
	C.StrokeWidth(C.VGfloat(w))
}

// Start begins a picture
func Start(w, h int, color ...uint8) {
	C.Start(C.int(w), C.int(h))
}

// End ends the picture
func End() {
	C.End()
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

type VGF *C.VGfloat

// Polyline draws a polyline with coordinates in x, y
func Polyline(x, y []float64, n int, style ...string) {
	/*
		if len(x) != len(y) {
				return
		}
		C.Polyline(&x[0], &y[0], C.VGint(len(x)))
	*/
}

// Polyline draws a polygon with coordinates in x, y
func Polygon(x, y []float64, n int, style ...string) {
	/*
		if len(x) != len(y) {
				return
		}
		C.Polygon(&x[0], &y[0], C.VGint(len(x)))
	*/
}

// select font specifies the font by generic name
func selectfont(s string) C.Fontinfo {
	if len(s) == 0 {
		return C.SansTypeface
	}
	switch s {
	case "sans":
		return C.SansTypeface
	case "serif":
		return C.SerifTypeface
	case "mono":
		return C.MonoTypeface
	default:
		return C.SansTypeface
	}
	return C.SansTypeface
}

// Text draws text whose aligment begins (x,y)
func Text(x, y float64, s string, font string, size int, style ...string) {
	C.Text(C.VGfloat(x), C.VGfloat(y), C.CString(s), selectfont(font), C.int(size))
}

// TextMid draws text centered at (x,y)
func TextMid(x, y float64, s string, font string, size int, style ...string) {
	C.TextMid(C.VGfloat(x), C.VGfloat(y), C.CString(s), selectfont(font), C.int(size))
}

// TextEnd draws text end-aligned at (x,y)
func TextEnd(x, y float64, s string, font string, size int, style ...string) {
	C.TextEnd(C.VGfloat(x), C.VGfloat(y), C.CString(s), selectfont(font), C.int(size))
}

// TextWidth returns the length of text at a specified font and size
func TextWidth(s string, font string, size float64) float64 {
	return float64(C.textwidth(C.CString(s), selectfont(font), C.VGfloat(size)))
}

// Translate translates the coordinate system to (x,y)
func Translate(x, y float64) {
	C.Translate(C.VGfloat(x), C.VGfloat(y))
}

// Rotate rotates the coordinate system around the specifed angle
func Rotate(r float64) {
	C.Rotate(C.VGfloat(r))
}

// Shear transforms the coordinate system by (x,y)
func Shear(x, y float64) {
	C.Shear(C.VGfloat(x), C.VGfloat(y))
}

// Scale scales the coordinate system by (x,y)
func Scale(x, y float64) {
	C.Scale(C.VGfloat(x), C.VGfloat(y))
}
