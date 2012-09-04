package openvg

/*
#cgo CFLAGS:   -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads 
#cgo LDFLAGS:  -L/opt/vc/lib -lGLESv2 -ljpeg
#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include <string.h>
#include <time.h>
#include "VG/openvg.h"
#include "VG/vgu.h"
#include "EGL/egl.h"
#include "GLES/gl.h"
#include "DejaVuSans.inc"               // font data
#include "DejaVuSerif.inc"
#include "eglstate.h"                   // data structures for graphics state
#include "fontinfo.h"                   // font data structure
#include "shapes.h"
*/
import "C"

func Init() (int, int) {
	var rh, rw C.int
	C.init(&rw, &rh)
	return int(rw), int(rh)
}

func Finish() {
	C.finish()
}

func Background(r, g, b uint8) {
	C.Background(C.uint(r), C.uint(g), C.uint(b))
}

func Fill(r, g, b uint8, alpha float64) {
	C.Fill(C.uint(r), C.uint(g), C.uint(b), C.VGfloat(alpha))
}

func Stroke(r, g, b uint8, alpha float64) {
	C.Stroke(C.uint(r), C.uint(g), C.uint(b), C.VGfloat(alpha))
}

func StrokeWidth(w float64) {
	C.StrokeWidth(C.VGfloat(w))
}

func Start(w, h int, color ...uint8) {
	C.Start(C.int(w), C.int(h))
}

func End() {
	C.End()
}

func Line(x1, y1, x2, y2 float64, style ...string) {
	C.Line(C.VGfloat(x1), C.VGfloat(y1), C.VGfloat(x2), C.VGfloat(y2))
}

func Rect(x, y, w, h float64, style ...string) {
	C.Rect(C.VGfloat(x), C.VGfloat(y), C.VGfloat(w), C.VGfloat(h))
}

func Roundrect(x, y, w, h, rw, rh float64, style ...string) {
	C.Roundrect(C.VGfloat(x), C.VGfloat(y), C.VGfloat(w), C.VGfloat(h), C.VGfloat(rw), C.VGfloat(rh))
}

func Ellipse(x, y, w, h float64, style ...string) {
	C.Ellipse(C.VGfloat(x), C.VGfloat(y), C.VGfloat(w), C.VGfloat(h))
}

func Circle(x, y, r float64, style ...string) {
	C.Circle(C.VGfloat(x), C.VGfloat(y), C.VGfloat(r))
}

func Qbezier(sx, sy, cx, cy, ex, ey float64, style ...string) {
	C.Qbezier(C.VGfloat(sx), C.VGfloat(sy), C.VGfloat(cx), C.VGfloat(cy), C.VGfloat(ex), C.VGfloat(ey))
}

func Cbezier(sx, sy, cx, cy, px, py, ex, ey float64, style ...string) {
	C.Cbezier(C.VGfloat(sx), C.VGfloat(sy), C.VGfloat(cx), C.VGfloat(cy), C.VGfloat(px), C.VGfloat(py), C.VGfloat(ex), C.VGfloat(ey))
}

func Arc(x, y, w, h, sa, aext float64, style ...string) {
	C.Arc(C.VGfloat(x), C.VGfloat(y), C.VGfloat(w), C.VGfloat(h), C.VGfloat(sa), C.VGfloat(aext))
}

type VGF *C.VGfloat

func Polyline(x, y []float64, n int, style ...string) {
		/*
	if len(x) != len(y) {
			return
	}
	C.Polyline(&x[0], &y[0], C.VGint(len(x)))
	*/
}

func Polygon(x, y []float64, n int, style ...string) {
		/*
	if len(x) != len(y) {
			return
	}
	C.Polygon(&x[0], &y[0], C.VGint(len(x)))
	*/
}

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

func Text(x, y float64, s string, font string, size int, style ...string) {
	C.Text(C.VGfloat(x), C.VGfloat(y), C.CString(s), selectfont(font), C.int(size))
}

func TextMid(x, y float64, s string, font string, size int, style ...string) {
	C.TextMiddle(C.VGfloat(x), C.VGfloat(y), C.CString(s), selectfont(font), C.int(size))
}

func TextEnd(x, y float64, s string, font string, size int, style ...string) {
	C.TextEnd(C.VGfloat(x), C.VGfloat(y), C.CString(s), selectfont(font), C.int(size))
}

func TextWidth(s string, font string, size float64) float64 {
	return float64(C.textwidth(C.CString(s), selectfont(font), C.VGfloat(size)))
}
func Translate(x, y float64) {
	C.Translate(C.VGfloat(x), C.VGfloat(y))
}

func Rotate(r float64) {
	C.Rotate(C.VGfloat(r))
}

func Shear(x, y float64) {
	C.Shear(C.VGfloat(x), C.VGfloat(y))
}

func Scale(x, y float64) {
	C.Scale(C.VGfloat(x), C.VGfloat(y))
}
