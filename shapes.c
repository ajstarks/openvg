// shapes: testbed for OpenVG APIs
// Anthony Starks (ajstarks@gmail.com)
#include <stdio.h>
#include <stdlib.h>
#include <assert.h>
#include <string.h>
#include <time.h>

#include "VG/openvg.h"
#include "VG/vgu.h"
#include "EGL/egl.h"
#include "GLES/gl.h"
#include "DejaVuSans.inc"

typedef struct {
	uint32_t screen_width;
	uint32_t screen_height;
	// OpenGL|ES objects
	EGLDisplay display;

	EGLSurface surface;
	EGLContext context;
} STATE_T;

static void finish(void);
static STATE_T _state, *state=&_state;
static const int MAXFONTPATH=256;
VGPath DejaVuSansPaths[256];

//
//  Font functions
//

// loadfont loads font path data
void  loadfont(const int *Points, const int *PointIndices, const unsigned char *Instructions, const int *InstructionIndices, const int *InstructionCounts, int ng, VGPath *glyphs) {
	int i;
	if (ng > MAXFONTPATH) {
		return;
	}
 	memset(glyphs, 0, MAXFONTPATH*sizeof(VGPath));
	for(i=0; i < ng; i++) {
		const int* p = &Points[PointIndices[i]*2];
		const unsigned char* instructions = &Instructions[InstructionIndices[i]];
		int ic = InstructionCounts[i];

		//printf("glyph %3d = %x, %x %d\n", i, p, instructions, ic);
		VGPath path = vgCreatePath(VG_PATH_FORMAT_STANDARD, VG_PATH_DATATYPE_S_32, 1.0f/65536.0f, 0.0f, 0, 0, VG_PATH_CAPABILITY_ALL);
		glyphs[i] = path;
		if(ic) {
			vgAppendPathData(path, ic, instructions, p);
		}
	}
}

// unloadfont frees font path data
void unloadfont(VGPath *glyphs, int n) {
	int i;
 	for(i=0; i<n; i++) {
		vgDestroyPath(glyphs[i]);
	}
}

// init_ogl sets the display, OpenGL|ES context and screen information
// state holds the OGLES model information
static void init_ogl(STATE_T *state) {
	int32_t success = 0;
	EGLBoolean result;
	EGLint num_config;

	static EGL_DISPMANX_WINDOW_T nativewindow;

	DISPMANX_ELEMENT_HANDLE_T dispman_element;
	DISPMANX_DISPLAY_HANDLE_T dispman_display;
	DISPMANX_UPDATE_HANDLE_T dispman_update;
	VC_RECT_T dst_rect;
	VC_RECT_T src_rect;

	static const EGLint attribute_list[] = {
		EGL_RED_SIZE, 8,
		EGL_GREEN_SIZE, 8,
		EGL_BLUE_SIZE, 8,
		EGL_ALPHA_SIZE, 8,
		EGL_SURFACE_TYPE, EGL_WINDOW_BIT,
		EGL_NONE
	};

	EGLConfig config;

	// get an EGL display connection
	state->display = eglGetDisplay(EGL_DEFAULT_DISPLAY);
	assert(state->display!=EGL_NO_DISPLAY);

	// initialize the EGL display connection
	result = eglInitialize(state->display, NULL, NULL);
	assert(EGL_FALSE != result);

	// bind OpenVG API
	eglBindAPI(EGL_OPENVG_API);

	// get an appropriate EGL frame buffer configuration
	result = eglChooseConfig(state->display, attribute_list, &config, 1, &num_config);
	assert(EGL_FALSE != result);

	// create an EGL rendering context
	state->context = eglCreateContext(state->display, config, EGL_NO_CONTEXT, NULL);
	assert(state->context!=EGL_NO_CONTEXT);

	// create an EGL window surface
	success = graphics_get_display_size(0 /* LCD */, &state->screen_width, &state->screen_height);
	assert( success >= 0 );

	dst_rect.x = 0;
	dst_rect.y = 0;
	dst_rect.width = state->screen_width;
	dst_rect.height = state->screen_height;

	src_rect.x = 0;
	src_rect.y = 0;
	src_rect.width = state->screen_width << 16;
	src_rect.height = state->screen_height << 16;

	dispman_display = vc_dispmanx_display_open( 0 /* LCD */);
	dispman_update = vc_dispmanx_update_start( 0 );

	dispman_element = vc_dispmanx_element_add ( dispman_update, dispman_display,
		0/*layer*/, &dst_rect, 0/*src*/,
		&src_rect, DISPMANX_PROTECTION_NONE, 0 /*alpha*/, 0/*clamp*/, 0/*transform*/);

	nativewindow.element = dispman_element;
	nativewindow.width = state->screen_width;
	nativewindow.height = state->screen_height;
	vc_dispmanx_update_submit_sync( dispman_update );

	state->surface = eglCreateWindowSurface( state->display, config, &nativewindow, NULL );
	assert(state->surface != EGL_NO_SURFACE);

	// connect the context to the surface
	result = eglMakeCurrent(state->display, state->surface, state->surface, state->context);
	assert(EGL_FALSE != result);

	//DAVE - Set up screen ratio
	glViewport(0, 0, (GLsizei)state->screen_width, (GLsizei)state->screen_height);

	glMatrixMode(GL_PROJECTION);
	glLoadIdentity();

	float ratio = (float)state->screen_width / (float)state->screen_height;
	glFrustumf(-ratio, ratio, -1.0f, 1.0f, 1.0f, 10.0f);
}

// init sets the system to its initial state
void init() {
	bcm_host_init();
	memset( state, 0, sizeof( *state ) );
	init_ogl(state);
   	loadfont(DejaVuSans_glyphPoints, DejaVuSans_glyphPointIndices,
            DejaVuSans_glyphInstructions, DejaVuSans_glyphInstructionIndices,
            DejaVuSans_glyphInstructionCounts, DejaVuSans_glyphCount,
            DejaVuSansPaths);
}
// finish cleans up
static void finish(void) {
	// Release font data
	unloadfont(DejaVuSansPaths, DejaVuSans_glyphCount);
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

// setstroke sets the stroke color and width
void setstroke(float color[4], float width) {
	VGPaint strokePaint = vgCreatePaint();
	vgSetParameteri(strokePaint, VG_PAINT_TYPE, VG_PAINT_TYPE_COLOR);
	vgSetParameterfv(strokePaint, VG_PAINT_COLOR, 4, color);	
	vgSetPaint(strokePaint, VG_STROKE_PATH);
	vgSetf(VG_STROKE_LINE_WIDTH, width);
	vgSeti(VG_STROKE_CAP_STYLE, VG_CAP_BUTT);
	vgSeti(VG_STROKE_JOIN_STYLE, VG_JOIN_MITER);
	vgDestroyPaint(strokePaint);
}

// Text renders a string of text at a specified location, using the specified font glyphs
void Text(VGfloat x, VGfloat y, const char* s, int pointsize, VGfloat fillcolor[4], VGPath *gp, const short *characterMap, const int *glyphAdvances, VGbitfield renderFlags) {
	float size = (float)pointsize;
	float xx = x;
	float mm[9];
	int i;
	vgGetMatrix(mm);
	setfill(fillcolor);
	for(i=0; i < (int)strlen(s); i++) {
		unsigned int character = (unsigned int)s[i];
		int glyph = characterMap[character];
		if( glyph == -1 ) {
			continue;	//glyph is undefined
		}

		VGfloat mat[9] = {
			size,	0.0f,	0.0f,
			0.0f,	size,	0.0f,
			xx,		y,		1.0f
		};

		vgLoadMatrix(mm);
		vgMultMatrix(mat);
		vgDrawPath(gp[glyph], renderFlags);
		xx += size * glyphAdvances[glyph] / 65536.0f;
		//printf("char %c advances %d\n", character, glyphAdvances[glyph]);
	}
	vgLoadMatrix(mm);
}

//
// Shape functions
//

VGPath newpath() {
	return vgCreatePath(VG_PATH_FORMAT_STANDARD, VG_PATH_DATATYPE_F, 1.0f, 0.0f, 0, 0, VG_PATH_CAPABILITY_ALL);
}

// CBezier makes a quadratic bezier curve, applying style
void Cbezier(
	VGfloat sx, VGfloat sy, 
	VGfloat cx, VGfloat cy, 
	VGfloat px, VGfloat py, 
	VGfloat ex, VGfloat ey, 
	VGfloat sw, VGfloat fill[4], VGfloat stroke[4]) {
	
	VGubyte segments[] = { VG_MOVE_TO_ABS, VG_CUBIC_TO };
	VGfloat coords[]   = { sx, sy, cx, cy, px, py, ex, ey };
	VGPath path = newpath();
	vgAppendPathData( path, 2, segments, coords );
	setfill(fill);
	setstroke(stroke, sw);
	vgDrawPath(path, VG_FILL_PATH | VG_STROKE_PATH);
	vgDestroyPath(path);
}

// QBezier makes a quadratic bezier curve, applying style
void Qbezier(
	VGfloat sx, VGfloat sy,
	VGfloat cx, VGfloat cy, 
	VGfloat ex, VGfloat ey, 
	VGfloat sw, VGfloat fill[4], VGfloat stroke[4]) {

	VGubyte segments[] = { VG_MOVE_TO_ABS, VG_QUAD_TO };
	VGfloat coords[]   = { sx, sy, cx, cy, ex, ey };
	VGPath path = newpath();
	vgAppendPathData( path, 2, segments, coords );
	setfill(fill);
	setstroke(stroke, sw);
	vgDrawPath(path, VG_FILL_PATH | VG_STROKE_PATH);
	vgDestroyPath(path);
}

// interleave interleaves arrays of x, y into a single array
void interleave(VGfloat *x, VGfloat *y, int n, VGfloat *points) {
	while (n--)  {
	    *points++ = *x++;
	    *points++ = *y++;
	} 
}

// Poly makes a stroked polyline or a stroked and filled polygon
void Poly(VGfloat *x, VGfloat *y, VGint n, VGfloat sw, VGfloat fill[4], VGfloat stroke[4], VGboolean dofill) {
	VGfloat points[n*2];
	VGPath path = newpath();
	VGbitfield pflag;

	interleave(x, y, n, points);
	vguPolygon(path, points, n, VG_FALSE);
	if (dofill) {
		setfill(fill);
		pflag = VG_FILL_PATH | VG_STROKE_PATH;
	} else {
		pflag = VG_STROKE_PATH;
	}
	
	setstroke(stroke, sw);
	vgDrawPath(path, pflag);
	vgDestroyPath(path);
}

// Rect makes a rectangle  at the specified location and dimensions, applying style
void Rect(VGfloat x, VGfloat y, VGfloat w, VGfloat h, VGfloat sw, VGfloat fill[4], VGfloat stroke[4]) {
	VGPath path = newpath();
	vguRect(path, x, y, w, h);
	setfill(fill);
	setstroke(stroke, sw);
	vgDrawPath(path, VG_FILL_PATH | VG_STROKE_PATH);
	vgDestroyPath(path);
}

// Line makes a line, applying style
void Line(VGfloat x1, VGfloat y1, VGfloat x2, VGfloat y2, VGfloat sw, VGfloat stroke[4]) {
	VGPath path = newpath();
	vguLine(path, x1, y1, x2, y2);
	setstroke(stroke, sw);
	vgDrawPath(path, VG_STROKE_PATH);
	vgDestroyPath(path);
}
// Roundrect makes an rounded rectangle at the specified location and dimensions, applying style
void Roundrect(VGfloat x, VGfloat y, VGfloat w, VGfloat h, VGfloat rw, VGfloat rh, VGfloat sw, VGfloat fill[4], VGfloat stroke[4]) {

	VGPath path = newpath();
	vguRoundRect(path, x, y, w, h, rw, rh);
	setfill(fill);
	setstroke(stroke, sw);
	vgDrawPath(path, VG_FILL_PATH | VG_STROKE_PATH);
	vgDestroyPath(path);
}

// Ellipse makes an ellipse at the specified location and dimensions, applying style
void Ellipse(VGfloat x, VGfloat y, VGfloat w, VGfloat h, VGfloat sw, VGfloat fill[4], VGfloat stroke[4]) {
	VGPath path = newpath();
	vguEllipse(path, x, y, w, h);
	setfill(fill);
	setstroke(stroke, sw);
	vgDrawPath(path, VG_FILL_PATH | VG_STROKE_PATH);
	vgDestroyPath(path);
}

// Circle makes a circle at the specified location and dimensions, applying style
void Circle(VGfloat x, VGfloat y, VGfloat r, VGfloat sw, VGfloat fill[4], VGfloat stroke[4]) {
	Ellipse(x, y, r, r, sw, fill, stroke);
}

// Arc makes an elliptical arc at the specified location and dimensions, applying style
void Arc(VGfloat x, VGfloat y, VGfloat w, VGfloat h, VGfloat sa, VGfloat aext, VGfloat sw, VGfloat fill[4], VGfloat stroke[4]) {
	VGPath path = newpath();
	vguArc(path, x, y, w, h, sa, aext, VGU_ARC_OPEN);
	setfill(fill);
	setstroke(stroke, sw);
	vgDrawPath(path, VG_FILL_PATH | VG_STROKE_PATH);
	vgDestroyPath(path);
}

// Start begins the picture, clearing a rectangular region with a specified color
void Start(int width, int height, float fill[4]) {
	vgSetfv(VG_CLEAR_COLOR, 4, fill);
	vgClear(0, 0, width, height);
	vgLoadIdentity();
}

// End checks for errors, and renders to the display
void End() {
	assert(vgGetError() == VG_NO_ERROR);
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

// rlines makes random lines
void rlines(VGfloat x, VGfloat y, int width, int height) {
	int i = (int)x;
	VGfloat rcolor[4] = {0,0,0,.40};
	for (; i < width; i++ ) {
		rcolor[0] = randcolor();
		rcolor[1] = rcolor[0];
		rcolor[2] = rcolor[0];
		Line(i, y, i+randf(width), height, 10, rcolor);
	}
}

// coordpoint marks a coordinate
void coordpoint(VGfloat x, VGfloat y, VGfloat size) {
	VGfloat dotcolor[4] = {0.3, 0.3, 0.3, 1}, white[4] = {1,1,1,1};
	Circle(x, y, size, 0, dotcolor, white);
}


// refcard shows a reference card of shapes
void refcard(int width, int height) {
	char *shapenames[] = { 
		"Circle", "Ellipse", "Rectangle", "Rounded Rectangle", 
		"Line", "Polyline", "Polygon", "Arc", "Quadratic Bezier", "Cubic Bezier"
	};
	VGfloat strokecolor[4] = {0.8,0.8,0.8,1}, shapecolor[4] = {202.0/255.0, 225.0/255.0,1,1}, textcolor[4] = {0,0,0,1}, bgcolor[4] = {1,1,1,1};

	VGfloat strokewidth = 1;
	VGfloat top=height-100, sx = 500, sy = top, sw=100, sh=50, dotsize=7, spacing=2.0;
	int i, ns = sizeof(shapenames)/sizeof(char *), fontsize = 36;
	Start(width, height, bgcolor);
	setfill(textcolor);
	sx = width * 0.10;
	textcolor[0] = 0.5;
	Text(width*.45, height/2, "OpenVG on the Raspberry Pi", 48, textcolor, DejaVuSansPaths, DejaVuSans_characterMap, DejaVuSans_glyphAdvances, VG_FILL_PATH);
	textcolor[0] = 0;
	for (i=0; i < ns; i++) {
		Text(sx+sw+sw/2, sy, shapenames[i], fontsize, textcolor, DejaVuSansPaths, DejaVuSans_characterMap, DejaVuSans_glyphAdvances, VG_FILL_PATH);
		sy -= sh*spacing; 
	}
	sy = top;
	VGfloat cx = sx+(sw/2), ex = sx + sw;
	Circle(cx, sy, sw, strokewidth, shapecolor, strokecolor); coordpoint(cx, sy, dotsize); sy -= sh*spacing;
	Ellipse(cx, sy, sw, sh, strokewidth, shapecolor, strokecolor); coordpoint(cx, sy, dotsize); sy -= sh*spacing;
	Rect(sx, sy, sw, sh, strokewidth, shapecolor, strokecolor); coordpoint(sx, sy, dotsize); sy -= sh*spacing;
	Roundrect(sx, sy, sw, sh, 20, 20, strokewidth, shapecolor, strokecolor); coordpoint(sx, sy, dotsize); sy -= sh*spacing;
	Line(sx, sy, ex, sy, strokewidth, strokecolor); coordpoint(sx, sy, dotsize); coordpoint(ex, sy, dotsize); sy -= sh;

	VGfloat px[5] = {sx, sx+(sw/4), sx+(sw/2), sx+((sw*3)/4), sx+sw}; 
	VGfloat py[5] = {sy, sy-sh, sy, sy-sh, sy}; 
	Poly(px, py, 5, strokewidth, shapecolor, strokecolor, VG_FALSE); 
	coordpoint(px[0], py[0], dotsize); 
	coordpoint(px[1], py[1], dotsize); 
	coordpoint(px[2], py[2], dotsize); 
	coordpoint(px[3], py[3], dotsize);
	coordpoint(px[4], py[4], dotsize);
	sy -= sh*spacing;

	py[0] = sy;
	py[1] = sy-sh;
	py[2] = sy-(sh/2);
	py[3] = py[1] - (sh/4);
	py[4] = sy;
	Poly(px, py, 5, strokewidth, shapecolor, strokecolor, VG_TRUE); sy -= (sh*spacing) + sh;

	Arc(sx+(sw/2), sy, sw, sh, 0, 180, strokewidth, shapecolor, strokecolor); coordpoint(sx+(sw/2), sy, dotsize); sy -= sh*spacing;
	VGfloat cy = sy + (sh/2), ey = sy;
	Qbezier(sx, sy, cx, cy, ex, ey, strokewidth, shapecolor, strokecolor); coordpoint(sx, sy, dotsize); coordpoint(cx, cy, dotsize); coordpoint(ex, ey, dotsize); sy -= sh*spacing;
	ey = sy;
	cy = sy + sh;
	Cbezier(sx, sy, cx, cy, cx, sy, ex, ey, strokewidth, shapecolor, strokecolor); 
	coordpoint(sx, sy, dotsize); 
	coordpoint(cx, cy, dotsize); 
	coordpoint(cx, sy, dotsize);
	coordpoint(ex, ey, dotsize);
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
		
		Ellipse(randf(width), randf(height), randf(200), randf(100), randf(10), rcolor, scolor);
		Circle(randf(width), randf(height), randf(100), randf(10), rcolor, scolor);
		Rect(randf(width), randf(height), randf(200), randf(100), randf(10), rcolor, scolor);
		Line(randf(width), randf(height), randf(width), randf(height), randf(10), scolor);
		Arc(randf(width), randf(height), randf(200), randf(200), randf(360), randf(360), randf(10), rcolor, scolor);

		sx = randf(width);
		sy = randf(height);
		ex = sx + randf(200);
		ey = sy;
		cx = sx + ((ex - sx )/ 2.0);
		cy = sy + randf(100);
		Qbezier(sx, sy, cx, cy, ex, ey, randf(10), rcolor, scolor);

		sx = randf(width);
		sy = randf(height);
		ex = sx + randf(200);
		ey = sy;
		cx = sx + ((ex - sx )/ 2.0);
		cy = sy + randf(100);
		px = cx;
		py = sy - randf(100);
		Cbezier(sx, sy, cx, cy, px, py, ex, ey, randf(10), rcolor, scolor);
		
		pox = randf(width);
		poy = randf(height);
		for (j=0; j < np; j++) {
			polyx[j] = pox + randf(200);
			polyy[j] = poy + randf(100);
		}
		Poly(polyx, polyy, np, 4, rcolor, scolor, VG_TRUE);

		pox = randf(width);
		poy = randf(height);
		for (j=0; j < np; j++) {
			polyx[j] = pox + randf(200);
			polyy[j] = poy + randf(100);
		}
		Poly(polyx, polyy, np, 4, rcolor, scolor, VG_FALSE);
	}
	Text(randf(100), randf(height-100), "OpenVG on the Raspberry Pi", 64, textcolor, DejaVuSansPaths, DejaVuSans_characterMap, DejaVuSans_glyphAdvances, VG_FILL_PATH);
	End();
}


// main initializes the system and shows the picture, 
// exit and clean up when you hit [RETURN].
int main (int argc, char **argv) {
	int w, h;

	init();
	w = state->screen_width;
	h = state->screen_height;
	if (argc > 1) {
		rshapes(w, h, atoi(argv[1]));
	} else {
		refcard(w, h);
	}
	while (getchar() != '\n') {
		;
	}
	finish();
	return 0;
}
