//
// shapes: testbed for OpenVG APIs
// Anthony Starks (ajstarks@gmail.com)
//
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

static STATE_T _state, *state=&_state;  // global graphics state
static const int MAXFONTPATH=256;
Fontinfo DejaFont, DejaSerif;

//
// Font functions
//

// loadfont loads font path data
// derived from http://web.archive.org/web/20070808195131/http://developer.hybrid.fi/font2openvg/renderFont.cpp.txt
Fontinfo loadfont(	const int *Points, 
					const int *PointIndices, 
					const unsigned char *Instructions, 
					const int *InstructionIndices, 
					const int *InstructionCounts, 
					const int *adv, 
					const short *cmap, 
					int ng) {

	Fontinfo f;
	int i;

 	memset(f.Glyphs, 0, MAXFONTPATH*sizeof(VGPath));
	if (ng > MAXFONTPATH) {
		return f;
	}
	for(i=0; i < ng; i++) {
		const int* p = &Points[PointIndices[i]*2];
		const unsigned char* instructions = &Instructions[InstructionIndices[i]];
		int ic = InstructionCounts[i];

		VGPath path = vgCreatePath(VG_PATH_FORMAT_STANDARD, VG_PATH_DATATYPE_S_32, 1.0f/65536.0f, 0.0f, 0, 0, VG_PATH_CAPABILITY_ALL);
		f.Glyphs[i] = path;
		if(ic) {
			vgAppendPathData(path, ic, instructions, p);
		}
	}
	f.Points = Points;
	f.PointIndices = PointIndices;
	f.Instructions = Instructions;
	f.CharacterMap = cmap;
	f.InstructionIndices = InstructionIndices;
	f.InstructionCounts = InstructionCounts;
	f.GlyphAdvances = adv;
	f.Count = ng;
	return f;
}

// unloadfont frees font path data
void unloadfont(VGPath *glyphs, int n) {
	int i;
 	for(i=0; i<n; i++) {
		vgDestroyPath(glyphs[i]);
	}
}

void dumpscreen(int w, int h) {
	void *ScreenBuffer = malloc(w*h*4);
	vgReadPixels(ScreenBuffer, (w*4), VG_sABGR_8888, 0, 0, w, h);
	fwrite(ScreenBuffer, 1, w*h*4, stdout); 
	free(ScreenBuffer);
}

// init sets the system to its initial state
void init(int *w, int *h) {
	bcm_host_init();
	memset( state, 0, sizeof( *state ) );
	oglinit(state);
	DejaFont = loadfont(DejaVuSans_glyphPoints, 
						DejaVuSans_glyphPointIndices, 
						DejaVuSans_glyphInstructions, 
						DejaVuSans_glyphInstructionIndices, 
						DejaVuSans_glyphInstructionCounts, 
						DejaVuSans_glyphAdvances,
						DejaVuSans_characterMap, 
						DejaVuSans_glyphCount);

	DejaSerif = loadfont(DejaVuSerif_glyphPoints,
                        DejaVuSerif_glyphPointIndices,
                        DejaVuSerif_glyphInstructions,
                        DejaVuSerif_glyphInstructionIndices,
                        DejaVuSerif_glyphInstructionCounts,
                        DejaVuSerif_glyphAdvances,
                        DejaVuSerif_characterMap,
                        DejaVuSerif_glyphCount);

	*w = state->screen_width;
	*h = state->screen_height;
}

// finish cleans up
static void finish(int w, int h) {
	// Release font data
	unloadfont(DejaFont.Glyphs, DejaFont.Count);
	unloadfont(DejaSerif.Glyphs, DejaSerif.Count);
	// clear screen
	glClear( GL_COLOR_BUFFER_BIT );
	eglSwapBuffers(state->display, state->surface);
	// Release OpenGL resources
	eglMakeCurrent( state->display, EGL_NO_SURFACE, EGL_NO_SURFACE, EGL_NO_CONTEXT );
	eglDestroySurface( state->display, state->surface );
	eglDestroyContext( state->display, state->context );
	eglTerminate( state->display );
}

// 
// Transformations
//

// Translate the coordinate system to x,y
void Translate(VGfloat x, VGfloat y) {
	vgTranslate(x, y);
}

// Rotate around angle r
void Rotate(VGfloat r) {
	vgRotate(r);
}

// Shear shears the x coordinate by x degrees, the y coordinate by y degrees
void Shear(VGfloat x, VGfloat y) {
	vgShear(x,y);
}

// Scale scales by  x, y
void Scale(VGfloat x, VGfloat y) {
	vgScale(x,y);
}

// ScaleX scales the x coordinate by a factor of x
void ScaleX(VGfloat x) {
	vgScale(x, 1);
}

// ScaleY scales the y coordinate by a factor of y
void ScaleY(VGfloat y) {
	vgScale(1, y);
}

// ShearX scales the x coordinate by a factor of x
void ShearX(VGfloat x) {
	vgShear(x,0);
}

// ShearY scales the y coordinate by a factor of y
void ShearY(VGfloat y) {
	vgShear(0,y);
}

//
// Style functions
//

// setfill sets the fill color
void setfill(float color[4]) {
	VGPaint fillPaint = vgCreatePaint();
	vgSetParameteri(fillPaint, VG_PAINT_TYPE, VG_PAINT_TYPE_COLOR);
	vgSetParameterfv(fillPaint, VG_PAINT_COLOR, 4, color);	
	vgSetPaint(fillPaint, VG_FILL_PATH);
	vgDestroyPaint(fillPaint);
}

// setstroke sets the stroke color
void setstroke(float color[4]) {
	VGPaint strokePaint = vgCreatePaint();
	vgSetParameteri(strokePaint, VG_PAINT_TYPE, VG_PAINT_TYPE_COLOR);
	vgSetParameterfv(strokePaint, VG_PAINT_COLOR, 4, color);	
	vgSetPaint(strokePaint, VG_STROKE_PATH);
	vgDestroyPaint(strokePaint);
}

// strokeWidth sets the stroke width
void strokeWidth(float width) {
	vgSetf(VG_STROKE_LINE_WIDTH, width);
	vgSeti(VG_STROKE_CAP_STYLE, VG_CAP_BUTT);
	vgSeti(VG_STROKE_JOIN_STYLE, VG_JOIN_MITER);
}

// Text renders a string of text at a specified location, size, and fill, using the specified font glyphs
// derived from http://web.archive.org/web/20070808195131/http://developer.hybrid.fi/font2openvg/renderFont.cpp.txt
void Text(VGfloat x, VGfloat y, const char* s, Fontinfo f, int pointsize, VGfloat fillcolor[4]) {
	float size = (float)pointsize;
	float xx = x;
	float mm[9];
	int i;
	vgGetMatrix(mm);
	setfill(fillcolor);
	for(i=0; i < (int)strlen(s); i++) {
		unsigned int character = (unsigned int)s[i];
		int glyph = f.CharacterMap[character];
		if( glyph == -1 ) {
			continue;	//glyph is undefined
		}
		VGfloat mat[9] = {
			size, 0.0f, 0.0f,
			0.0f, size, 0.0f,
			xx,	y, 1.0f
		};
		vgLoadMatrix(mm);
		vgMultMatrix(mat);
		vgDrawPath(f.Glyphs[glyph], VG_FILL_PATH);
		xx += size * f.GlyphAdvances[glyph] / 65536.0f;
	}
	vgLoadMatrix(mm);
}

// textwidth returns the width of a text string in a font
VGfloat textwidth(char *s, Fontinfo f, VGfloat size) {
	int i;
	VGfloat tw = 0.0;
	for(i=0; i < (int)strlen(s); i++) {
        unsigned int character = (unsigned int)s[i];
        int glyph = f.CharacterMap[character];
        if( glyph == -1 ) {
            continue;   //glyph is undefined
        }
        tw += size * f.GlyphAdvances[glyph] / 65536.0f;
    }
	return tw;
}
//
// Shape functions
//


// newpath creates path data
VGPath newpath() {
	return vgCreatePath(VG_PATH_FORMAT_STANDARD, VG_PATH_DATATYPE_F, 1.0f, 0.0f, 0, 0, VG_PATH_CAPABILITY_ALL);
}

// makecurve makes path data using specified segments and coordinates
void makecurve(VGubyte *segments, VGfloat *coords) {
	VGPath path = newpath();
	vgAppendPathData( path, 2, segments, coords );
	vgDrawPath(path, VG_FILL_PATH | VG_STROKE_PATH);
	vgDestroyPath(path);
}

// CBezier makes a quadratic bezier curve
void Cbezier( VGfloat sx, VGfloat sy, VGfloat cx, VGfloat cy, VGfloat px, VGfloat py, VGfloat ex, VGfloat ey) {
	VGubyte segments[] = { VG_MOVE_TO_ABS, VG_CUBIC_TO };
	VGfloat coords[]   = { sx, sy, cx, cy, px, py, ex, ey };
	makecurve(segments, coords);
}

// QBezier makes a quadratic bezier curve
void Qbezier(VGfloat sx, VGfloat sy, VGfloat cx, VGfloat cy, VGfloat ex, VGfloat ey) {
	VGubyte segments[] = { VG_MOVE_TO_ABS, VG_QUAD_TO };
	VGfloat coords[]   = { sx, sy, cx, cy, ex, ey };
	makecurve(segments, coords);
}

// interleave interleaves arrays of x, y into a single array
void interleave(VGfloat *x, VGfloat *y, int n, VGfloat *points) {
	while (n--) {
		*points++ = *x++;
		*points++ = *y++;
	} 
}

// poly makes either a polygon or polyline
void poly(VGfloat *x, VGfloat *y, VGint n, VGbitfield flag) {
    VGfloat points[n*2];
    VGPath path = newpath();
    interleave(x, y, n, points);
    vguPolygon(path, points, n, VG_FALSE);
    vgDrawPath(path, flag);
    vgDestroyPath(path);
}

// Polygon makes a filled polygon with vertices in x, y arrays
void Polygon(VGfloat *x, VGfloat *y, VGint n) {
	poly(x, y, n, VG_FILL_PATH);
}

// Polyline makes a polyline with vertices at x, y arrays
void Polyline(VGfloat *x, VGfloat *y, VGint n) {
	poly(x, y, n, VG_STROKE_PATH);
}

// Rect makes a rectangle at the specified location and dimensions
void Rect(VGfloat x, VGfloat y, VGfloat w, VGfloat h) {
	VGPath path = newpath();
	vguRect(path, x, y, w, h);
	vgDrawPath(path, VG_FILL_PATH | VG_STROKE_PATH);
	vgDestroyPath(path);
}

// Line makes a line from (x1,y1) to (x2,y2)
void Line(VGfloat x1, VGfloat y1, VGfloat x2, VGfloat y2) {
	VGPath path = newpath();
	vguLine(path, x1, y1, x2, y2);
	vgDrawPath(path, VG_STROKE_PATH);
	vgDestroyPath(path);
}
// Roundrect makes an rounded rectangle at the specified location and dimensions
void Roundrect(VGfloat x, VGfloat y, VGfloat w, VGfloat h, VGfloat rw, VGfloat rh) { 
	VGPath path = newpath();
	vguRoundRect(path, x, y, w, h, rw, rh);
	vgDrawPath(path, VG_FILL_PATH | VG_STROKE_PATH);
	vgDestroyPath(path);
}

// Ellipse makes an ellipse at the specified location and dimensions
void Ellipse(VGfloat x, VGfloat y, VGfloat w, VGfloat h) {
	VGPath path = newpath();
	vguEllipse(path, x, y, w, h);
	vgDrawPath(path, VG_FILL_PATH | VG_STROKE_PATH);
	vgDestroyPath(path);
}

// Circle makes a circle at the specified location and dimensions
void Circle(VGfloat x, VGfloat y, VGfloat r) {
	Ellipse(x, y, r, r);
}

// Arc makes an elliptical arc at the specified location and dimensions
void Arc(VGfloat x, VGfloat y, VGfloat w, VGfloat h, VGfloat sa, VGfloat aext) {
	VGPath path = newpath();
	vguArc(path, x, y, w, h, sa, aext, VGU_ARC_OPEN);
	vgDrawPath(path, VG_FILL_PATH | VG_STROKE_PATH);
	vgDestroyPath(path);
}

// Start begins the picture, clearing a rectangular region with a specified color
void Start(int width, int height, float fill[4]) {
	vgSetfv(VG_CLEAR_COLOR, 4, fill);
	vgClear(0, 0, width, height);
	VGfloat black[4] = {0,0,0,1};
	setfill(black);
	setstroke(black);
	strokeWidth(0);
	vgLoadIdentity();
}

// End checks for errors, and renders to the display
void End() {
	assert(vgGetError() == VG_NO_ERROR);
	eglSwapBuffers(state->display, state->surface);
	assert(eglGetError() == EGL_SUCCESS);
}

void SaveEnd() {
	assert(vgGetError() == VG_NO_ERROR);
    dumpscreen(state->screen_width, state->screen_height);
    eglSwapBuffers(state->display, state->surface);
    assert(eglGetError() == EGL_SUCCESS);
}

// randcolor returns a fraction of 255
VGfloat randcolor() {
	return (VGfloat)(rand() % 256) / 255.0;
}

// randf returns a floating point number bounded by n
VGfloat randf(n) {
	return (VGfloat)(rand() % n);
}

// coordpoint marks a coordinate, preserving a previous color
void coordpoint(VGfloat x, VGfloat y, VGfloat size, VGfloat pcolor[4]) {
	VGfloat dotcolor[4] = {0.3, 0.3, 0.3, 1};
	setfill(dotcolor);
	Circle(x, y, size);
	setfill(pcolor);
}
// testpattern shows a test pattern 
void testpattern(int width, int height) {
	VGfloat llc[4] = {1,0,0,1},
			ulc[4] = {0,1,0,1},
			lrc[4] = {0,0,1,1},
			urc[4] = {0.5,0.5,0.5,1},
			tc[4]  = {0,0,0,1},
			bgcolor[4] = {1,1,1,1},
			tw;
	Start(width, height, bgcolor);
	setfill(llc); Rect(0,0,100,100);
	setfill(ulc); Rect(0,height-100,100,100);
	setfill(lrc); Rect(width-100,0,100,100);
	setfill(urc); Rect(width-100,height-100,100,100);
	tw = textwidth("hello, Pi", DejaFont, 256);
	Text((width/2)-(tw/2),height/2,"hello Pi", DejaFont, 256, tc);
	End();
}
// refcard shows a reference card of shapes
void refcard(int width, int height) {
	char *shapenames[] = { 
		"Circle", "Ellipse", "Rectangle", "Rounded Rectangle", 
		"Line", "Polyline", "Polygon", "Arc", "Quadratic Bezier", "Cubic Bezier"
	};
	VGfloat strokecolor[4] = {0.8,0.8,0.8,1}, 
			shapecolor[4] = {202.0/255.0, 225.0/255.0,1,1}, 
			textcolor[4] = {0,0,0,1}, 
			bgcolor[4] = {1,1,1,1};

	VGfloat linewidth = 1;
	VGfloat top=height-100, sx = 500, sy = top, sw=100, sh=50, dotsize=7, spacing=2.0;
	int i, ns = sizeof(shapenames)/sizeof(char *), fontsize = 36;
	Start(width, height, bgcolor);
	setfill(textcolor);
	sx = width * 0.10;
	textcolor[0] = 0.5;
	Text(width*.45, height/2, "OpenVG on the Raspberry Pi", DejaFont, 48, textcolor);
		
	textcolor[0] = 0;
	for (i=0; i < ns; i++) {
		Text(sx+sw+sw/2, sy, shapenames[i], DejaFont, fontsize, textcolor);
		sy -= sh*spacing; 
	}
	sy = top;
	VGfloat cx = sx+(sw/2), ex = sx + sw;
	setfill(shapecolor);
	Circle(cx, sy, sw); coordpoint(cx, sy, dotsize, shapecolor); sy -= sh*spacing;
	Ellipse(cx, sy, sw, sh); coordpoint(cx, sy, dotsize, shapecolor); sy -= sh*spacing;
	Rect(sx, sy, sw, sh); coordpoint(sx, sy, dotsize, shapecolor); sy -= sh*spacing;
	Roundrect(sx, sy, sw, sh, 20, 20); coordpoint(sx, sy, dotsize, shapecolor); sy -= sh*spacing;

	strokeWidth(linewidth);
	setstroke(strokecolor); 
	Line(sx, sy, ex, sy); coordpoint(sx, sy, dotsize, shapecolor); coordpoint(ex, sy, dotsize, shapecolor); sy -= sh;

	VGfloat px[5] = {sx, sx+(sw/4), sx+(sw/2), sx+((sw*3)/4), sx+sw}; 
	VGfloat py[5] = {sy, sy-sh, sy, sy-sh, sy}; 

	Polyline(px, py, 5); 
	coordpoint(px[0], py[0], dotsize, shapecolor); 
	coordpoint(px[1], py[1], dotsize, shapecolor); 
	coordpoint(px[2], py[2], dotsize, shapecolor); 
	coordpoint(px[3], py[3], dotsize, shapecolor);
	coordpoint(px[4], py[4], dotsize, shapecolor);
	sy -= sh*spacing;

	py[0] = sy;
	py[1] = sy-sh;
	py[2] = sy-(sh/2);
	py[3] = py[1] - (sh/4);
	py[4] = sy;
	Polygon(px, py, 5); sy -= (sh*spacing) + sh;

	Arc(sx+(sw/2), sy, sw, sh, 0, 180); coordpoint(sx+(sw/2), sy, dotsize, shapecolor); sy -= sh*spacing;

	VGfloat cy = sy + (sh/2), ey = sy;
	Qbezier(sx, sy, cx, cy, ex, ey);
	coordpoint(sx, sy, dotsize, shapecolor);
	coordpoint(cx, cy, dotsize, shapecolor); 
	coordpoint(ex, ey, dotsize, shapecolor);
	sy -= sh*spacing;

	ey = sy;
	cy = sy + sh;
	Cbezier(sx, sy, cx, cy, cx, sy, ex, ey); 
	coordpoint(sx, sy, dotsize, shapecolor); 
	coordpoint(cx, cy, dotsize, shapecolor); 
	coordpoint(cx, sy, dotsize, shapecolor);
	coordpoint(ex, ey, dotsize, shapecolor);

	End();
}

// clear the screen to a background color
void Background(int w, int h, VGfloat fill[4]) {
	vgSetfv(VG_CLEAR_COLOR, 4, fill);
	vgClear(0,0,w,h);
}
// rotext draws text, rotated around the center of the screen, progressively faded
void rotext(VGfloat x, VGfloat y, int w, int h, int n, VGfloat deg, char *s) {
	int i;
	VGfloat textcolor[4] = {1,1,1,1}, bgcolor[4] = {0,0,0,1};
	VGfloat fade = (100.0/(VGfloat)n)/100.0;
	
	Start(w, h, bgcolor);
	Translate(x,y);
	for (i=0; i < n; i++) {
		Text(0,0, s, DejaSerif, 256, textcolor);
		textcolor[3] -= fade;
		Rotate(deg);
	}
	End();
}

// rshapes draws shapes (rect and ellipse) with random colors, strokes, and sizes. 
void rshapes(int width, int height, int n) {
	int np = 10;
	VGfloat rcolor[4], scolor[4], bgcolor[4] = {1,1,1,1}, textcolor[4] = {0.5, 0, 0, 1};
	scolor[3] = 1; // strokes are solid
	VGfloat sx, sy, cx, cy, px, py, ex, ey, pox, poy;
	VGfloat polyx[np], polyy[np];
	int i,j;
	srand ( time(NULL) );
	Start(width, height, bgcolor);
	for (i=0; i < n; i++) {
		rcolor[0] = randcolor();
		rcolor[1] = randcolor();
		rcolor[2] = randcolor();
		rcolor[3] = randcolor();

		scolor[1] = randcolor();
		scolor[2] = randcolor();
		scolor[3] = randcolor();
		setfill(rcolor);
		setstroke(scolor);
		strokeWidth(randf(10));
		Ellipse(randf(width), randf(height), randf(200), randf(100));
		Circle(randf(width), randf(height), randf(100));
		Rect(randf(width), randf(height), randf(200), randf(100));
		Arc(randf(width), randf(height), randf(200), randf(200), randf(360), randf(360));
		
		sx = randf(width);
		sy = randf(height);
		Line(sx, sy, sx+randf(200), sy+randf(100));

		sx = randf(width);
		sy = randf(height);
		ex = sx + randf(200);
		ey = sy;
		cx = sx + ((ex - sx )/ 2.0);
		cy = sy + randf(100);
		Qbezier(sx, sy, cx, cy, ex, ey);

		sx = randf(width);
		sy = randf(height);
		ex = sx + randf(200);
		ey = sy;
		cx = sx + ((ex - sx )/ 2.0);
		cy = sy + randf(100);
		px = cx;
		py = sy - randf(100);
		Cbezier(sx, sy, cx, cy, px, py, ex, ey);
		
		pox = randf(width);
		poy = randf(height);
		for (j=0; j < np; j++) {
			polyx[j] = pox + randf(200);
			polyy[j] = poy + randf(100);
		}
		Polygon(polyx, polyy, np);

		pox = randf(width);
		poy = randf(height);
		for (j=0; j < np; j++) {
			polyx[j] = pox + randf(200);
			polyy[j] = poy + randf(100);
		}
		Polyline(polyx, polyy, np);
	}
	Text(50, 100, "OpenVG on the Raspberry Pi", DejaFont, 64, textcolor);
	End();
}

// main initializes the system and shows the picture. 
// Exit and clean up when you hit [RETURN].
int main (int argc, char **argv) {
	int w, h, nr;

	init(&w, &h);
	switch (argc) {
		case 2:
			rshapes(w, h, atoi(argv[1]));
			break;
		case 3:
			nr = atoi(argv[1]);
			rotext(w/2, h/2, w, h, nr, 360.0/(VGfloat)nr, argv[2]);
			break;
		default:
			refcard(w,h);
	}
	while (getchar() != '\n') {
		;
	}
	finish(w, h);
	return 0;
}
