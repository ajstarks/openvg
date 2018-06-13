
/**
 * High-level OpenVG API.
 *
 * @author  Anthony Starks <ajstarks@gmail.com>
 * @author  George Thomas <m.g.thomas99@gmail.com>
 * @author  <github.com/paeryn>
 */

#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <termios.h>
#include <jpeglib.h>
#include "EGL/egl.h"
#include "VG/openvg.h"
#include "VG/vgu.h"
#include "bcm_host.h"
#include "fontinfo.h"
#include "oglinit.h"
#include "./../lib/DejaVuSans.inc"

// createImageFromJpeg decompresses a JPEG image to the standard image format
// source: https://github.com/ileben/ShivaVG/blob/master/examples/test_image.c
VGImage createImageFromJpeg(const char *filename) {
    FILE *infile;
    struct jpeg_decompress_struct jdc;
    struct jpeg_error_mgr jerr;
    JSAMPARRAY buffer;
    unsigned int bstride;
    unsigned int bbpp;

    VGImage img;
    VGubyte *data;
    unsigned int width;
    unsigned int height;
    unsigned int dstride;
    unsigned int dbpp;

    VGubyte *brow;
    VGubyte *drow;
    unsigned int x;
    unsigned int lilEndianTest = 1;
    VGImageFormat rgbaFormat;

    // Check for endianness
    if (((unsigned char *)&lilEndianTest)[0] == 1)
        rgbaFormat = VG_sABGR_8888;
    else
        rgbaFormat = VG_sRGBA_8888;

    // Try to open image file
    infile = fopen(filename, "rb");
    if (infile == NULL) {
        printf("Failed opening '%s' for reading!\n", filename);
        return VG_INVALID_HANDLE;
    }
    // Setup default error handling
    jdc.err = jpeg_std_error(&jerr);
    jpeg_create_decompress(&jdc);

    // Set input file
    jpeg_stdio_src(&jdc, infile);

    // Read header and start
    jpeg_read_header(&jdc, TRUE);
    jpeg_start_decompress(&jdc);
    width = jdc.output_width;
    height = jdc.output_height;

    // Allocate buffer using jpeg allocator
    bbpp = jdc.output_components;
    bstride = width * bbpp;
    buffer = (*jdc.mem->alloc_sarray)
        ((j_common_ptr) & jdc, JPOOL_IMAGE, bstride, 1);

    // Allocate image data buffer
    dbpp = 4;
    dstride = width * dbpp;
    data = (VGubyte *) malloc(dstride * height);

    // Iterate until all scanlines processed
    while (jdc.output_scanline < height) {

        // Read scanline into buffer
        jpeg_read_scanlines(&jdc, buffer, 1);
        drow = data + (height - jdc.output_scanline) * dstride;
        brow = buffer[0];
        // Expand to RGBA
        for (x = 0; x < width; ++x, drow += dbpp, brow += bbpp) {
            switch (bbpp) {
            case 4:
                drow[0] = brow[0];
                drow[1] = brow[1];
                drow[2] = brow[2];
                drow[3] = brow[3];
                break;
            case 3:
                drow[0] = brow[0];
                drow[1] = brow[1];
                drow[2] = brow[2];
                drow[3] = 255;
                break;
            }
        }
    }

    // Create VG image
    img = vgCreateImage(rgbaFormat, width, height, VG_IMAGE_QUALITY_BETTER);
    vgImageSubData(img, data, dstride, rgbaFormat, 0, 0, width, height);

    // Cleanup
    jpeg_destroy_decompress(&jdc);
    fclose(infile);
    free(data);

    return img;
}

// makeimage makes an image from a raw raster of red, green, blue, alpha values
void makeimage(VGfloat x, VGfloat y, int w, int h, VGubyte * data) {
    unsigned int dstride = w * 4;
    VGImageFormat rgbaFormat = VG_sABGR_8888;
    VGImage img = vgCreateImage(rgbaFormat, w, h, VG_IMAGE_QUALITY_BETTER);
    vgImageSubData(img, (void *)data, dstride, rgbaFormat, 0, 0, w, h);
    vgSetPixels(x, y, img, 0, 0, w, h);
    vgDestroyImage(img);
}


static EVG_STATE_T _state, *state = &_state;	// global graphics state
static const int MAXFONTPATH = 500;
static int init_x = 0;		// Initial window position and size
static int init_y = 0;
static unsigned int init_w = 0;
static unsigned int init_h = 0;
//
// Terminal settings
//

// terminal settings structures
struct termios new_term_attr;
struct termios orig_term_attr;

// saveterm saves the current terminal settings
void saveterm() {
    tcgetattr(fileno(stdin), &orig_term_attr);
}

// rawterm sets the terminal to raw mode
void rawterm() {
    memcpy(&new_term_attr, &orig_term_attr, sizeof(struct termios));
    new_term_attr.c_lflag &= ~(ICANON | ECHO | ECHOE | ECHOK | ECHONL | ECHOPRT | ECHOKE | ICRNL);
    new_term_attr.c_cc[VTIME] = 0;
    new_term_attr.c_cc[VMIN] = 0;
    tcsetattr(fileno(stdin), TCSANOW, &new_term_attr);
}

// restore resets the terminal to the previously saved setting
void restoreterm() {
    tcsetattr(fileno(stdin), TCSANOW, &orig_term_attr);
}

//
// Font functions
//

// loadfont loads font path data
// derived from http://web.archive.org/web/20070808195131/http://developer.hybrid.fi/font2openvg/renderFont.cpp.txt
Fontinfo loadfont(const int *Points,
        const int *PointIndices,
        const unsigned char *Instructions,
        const int *InstructionIndices,
        const int *InstructionCounts,
        const int *adv,
        const short *cmap,
        int ng) {

    Fontinfo f;
    int i;

    memset(f.Glyphs, 0, MAXFONTPATH * sizeof(VGPath));
    if (ng > MAXFONTPATH) {
        return f;
    }
    for (i = 0; i < ng; i++) {
        const int *p = &Points[PointIndices[i] * 2];
        const unsigned char *instructions = &Instructions[InstructionIndices[i]];
        int ic = InstructionCounts[i];
        VGPath path = vgCreatePath(VG_PATH_FORMAT_STANDARD, VG_PATH_DATATYPE_S_32,
                       1.0f / 65536.0f, 0.0f, 0, 0,
                       VG_PATH_CAPABILITY_ALL);
        f.Glyphs[i] = path;
        if (ic) {
            vgAppendPathData(path, ic, instructions, p);
        }
    }
    f.CharacterMap = cmap;
    f.GlyphAdvances = adv;
    f.Count = ng;
    f.descender_height = 0;
    f.font_height = 0;
    return f;
}

// unloadfont frees font path data
void unloadfont(VGPath * glyphs, int n) {
    int i;
    for (i = 0; i < n; i++) {
        vgDestroyPath(glyphs[i]);
    }
}

// Image places an image at the specifed location
void Image(VGfloat x, VGfloat y, int w, int h, const char *filename) {
    VGImage img = createImageFromJpeg(filename);
    vgSetPixels(x, y, img, 0, 0, w, h);
    vgDestroyImage(img);
}

// dumpscreen writes the raster
void dumpscreen(int w, int h, FILE * fp) {
    void *ScreenBuffer = malloc(w * h * 4);
    vgReadPixels(ScreenBuffer, (w * 4), VG_sABGR_8888, 0, 0, w, h);
    fwrite(ScreenBuffer, 1, w * h * 4, fp);
    free(ScreenBuffer);
}

Fontinfo SansTypeface;

// initWindowSize requests a specific window size & position, if not called
// then init() will open a full screen window.
// Done this way to preserve the original init() behaviour.
void initWindowSize(int x, int y, unsigned int w, unsigned int h) {
    init_x = x;
    init_y = y;
    init_w = w;
    init_h = h;
}

// init sets the system to its initial state
void evgInit(int *w, int *h) {
    bcm_host_init();
    memset(state, 0, sizeof(*state));
    state->window_x = init_x;
    state->window_y = init_y;
    state->window_width = init_w;
    state->window_height = init_h;
    oglinit(state);
    SansTypeface = loadfont(DejaVuSans_glyphPoints,
                DejaVuSans_glyphPointIndices,
                DejaVuSans_glyphInstructions,
                DejaVuSans_glyphInstructionIndices,
                DejaVuSans_glyphInstructionCounts,
                DejaVuSans_glyphAdvances,
                DejaVuSans_characterMap,
                DejaVuSans_glyphCount);
    SansTypeface.descender_height = DejaVuSans_descender_height;
    SansTypeface.font_height = DejaVuSans_font_height;

    *w = state->window_width;
    *h = state->window_height;
}

// AreaClear clears a given rectangle in window coordinates (not affected by
// transformations)
void evgClearRect(unsigned int x, unsigned int y, unsigned int w, unsigned int h) {
    vgClear(x, y, w, h);
}

// WindowClear clears the window to previously set background colour
void evgClear() {
    evgClearRect(0, 0, state->window_width, state->window_height);
}

// finish cleans up
void evgFinish() {
    unloadfont(SansTypeface.Glyphs, SansTypeface.Count);
    eglSwapBuffers(state->display, state->surface);
    eglMakeCurrent(state->display, EGL_NO_SURFACE, EGL_NO_SURFACE, EGL_NO_CONTEXT);
    eglDestroySurface(state->display, state->surface);
    eglDestroyContext(state->display, state->context);
    eglTerminate(state->display);
}

//
// Transformations
//

// Translate the coordinate system to x,y
void evgTranslate(VGfloat x, VGfloat y) {
    vgTranslate(x, y);
}

// Rotate around angle r
void evgRotate(VGfloat r) {
    vgRotate(r);
}

// Shear shears the x coordinate by x degrees, the y coordinate by y degrees
void evgShear(VGfloat x, VGfloat y) {
    vgShear(x, y);
}

// Scale scales by  x, y
void evgScale(VGfloat x, VGfloat y) {
    vgScale(x, y);
}

//
// Style functions
//

// setfill sets the fill color
void evgSetFill(VGfloat color[4]) {
    VGPaint fillPaint = vgCreatePaint();
    vgSetParameteri(fillPaint, VG_PAINT_TYPE, VG_PAINT_TYPE_COLOR);
    vgSetParameterfv(fillPaint, VG_PAINT_COLOR, 4, color);
    vgSetPaint(fillPaint, VG_FILL_PATH);
    vgDestroyPaint(fillPaint);
}

// setstroke sets the stroke color
void evgSetStroke(VGfloat color[4]) {
    VGPaint strokePaint = vgCreatePaint();
    vgSetParameteri(strokePaint, VG_PAINT_TYPE, VG_PAINT_TYPE_COLOR);
    vgSetParameterfv(strokePaint, VG_PAINT_COLOR, 4, color);
    vgSetPaint(strokePaint, VG_STROKE_PATH);
    vgDestroyPaint(strokePaint);
}

// StrokeWidth sets the stroke width
void evgStrokeWidth(VGfloat width) {
    vgSetf(VG_STROKE_LINE_WIDTH, width);
    vgSeti(VG_STROKE_CAP_STYLE, VG_CAP_BUTT);
    vgSeti(VG_STROKE_JOIN_STYLE, VG_JOIN_MITER);
}

//
// Color functions
//
//

// RGBA fills a color vectors from a RGBA quad.
void RGBA(unsigned int r, unsigned int g, unsigned int b, VGfloat a, VGfloat color[4]) {
    if (r > 255) {
        r = 0;
    }
    if (g > 255) {
        g = 0;
    }
    if (b > 255) {
        b = 0;
    }
    if (a < 0.0 || a > 1.0) {
        a = 1.0;
    }
    color[0] = (VGfloat) r / 255.0f;
    color[1] = (VGfloat) g / 255.0f;
    color[2] = (VGfloat) b / 255.0f;
    color[3] = a;
}

// RGB returns a solid color from a RGB triple
void RGB(unsigned int r, unsigned int g, unsigned int b, VGfloat color[4]) {
    RGBA(r, g, b, 1.0f, color);
}

// Stroke sets the stroke color, defined as a RGB triple.
void evgStroke(unsigned int r, unsigned int g, unsigned int b, VGfloat a) {
    VGfloat color[4];
    RGBA(r, g, b, a, color);
    evgSetStroke(color);
}

// Fill sets the fillcolor, defined as a RGBA quad.
void evgFill(unsigned int r, unsigned int g, unsigned int b, VGfloat a) {
    VGfloat color[4];
    RGBA(r, g, b, a, color);
    evgSetFill(color);
}

// setstops sets color stops for gradients
void evgSetStop(VGPaint paint, VGfloat * stops, int n) {
    VGboolean multmode = VG_FALSE;
    VGColorRampSpreadMode spreadmode = VG_COLOR_RAMP_SPREAD_REPEAT;
    vgSetParameteri(paint, VG_PAINT_COLOR_RAMP_SPREAD_MODE, spreadmode);
    vgSetParameteri(paint, VG_PAINT_COLOR_RAMP_PREMULTIPLIED, multmode);
    vgSetParameterfv(paint, VG_PAINT_COLOR_RAMP_STOPS, 5 * n, stops);
    vgSetPaint(paint, VG_FILL_PATH);
}

// LinearGradient fills with a linear gradient
void evgFillLinearGradient(VGfloat x1, VGfloat y1, VGfloat x2, VGfloat y2, VGfloat * stops, int ns) {
    VGfloat lgcoord[4] = { x1, y1, x2, y2 };
    VGPaint paint = vgCreatePaint();
    vgSetParameteri(paint, VG_PAINT_TYPE, VG_PAINT_TYPE_LINEAR_GRADIENT);
    vgSetParameterfv(paint, VG_PAINT_LINEAR_GRADIENT, 4, lgcoord);
    evgSetStop(paint, stops, ns);
    vgDestroyPaint(paint);
}

// RadialGradient fills with a linear gradient
void evgFillRadialGradient(VGfloat cx, VGfloat cy, VGfloat fx, VGfloat fy, VGfloat radius, VGfloat * stops, int ns) {
    VGfloat radialcoord[5] = { cx, cy, fx, fy, radius };
    VGPaint paint = vgCreatePaint();
    vgSetParameteri(paint, VG_PAINT_TYPE, VG_PAINT_TYPE_RADIAL_GRADIENT);
    vgSetParameterfv(paint, VG_PAINT_RADIAL_GRADIENT, 5, radialcoord);
    evgSetStop(paint, stops, ns);
    vgDestroyPaint(paint);
}

// ClipRect limits the drawing area to specified rectangle
void evgClipRect(VGint x, VGint y, VGint w, VGint h) {
    vgSeti(VG_SCISSORING, VG_TRUE);
    VGint coords[4] = { x, y, w, h };
    vgSetiv(VG_SCISSOR_RECTS, 4, coords);
}

// ClipEnd stops limiting drawing area to specified rectangle
void evgClipEnd() {
    vgSeti(VG_SCISSORING, VG_FALSE);
}

// Text Functions

// next_utf_char handles UTF encoding
unsigned char *next_utf8_char(unsigned char *utf8, int *codepoint) {
    int seqlen;
    int datalen = (int)strlen((const char *)utf8);
    unsigned char *p = utf8;

    if (datalen < 1 || *utf8 == 0) {		   // End of string
        return NULL;
    }
    if (!(utf8[0] & 0x80)) {			   // 0xxxxxxx
        *codepoint = (wchar_t) utf8[0];
        seqlen = 1;
    } else if ((utf8[0] & 0xE0) == 0xC0) {		   // 110xxxxx
        *codepoint = (int)(((utf8[0] & 0x1F) << 6) | (utf8[1] & 0x3F));
        seqlen = 2;
    } else if ((utf8[0] & 0xF0) == 0xE0) {		   // 1110xxxx
        *codepoint = (int)(((utf8[0] & 0x0F) << 12) | ((utf8[1] & 0x3F) << 6) | (utf8[2] & 0x3F));
        seqlen = 3;
    } else {
        return NULL;				   // No code points this high here
    }
    p += seqlen;
    return p;
}

// Text renders a string of text at a specified location, size, using the specified font glyphs
// derived from http://web.archive.org/web/20070808195131/http://developer.hybrid.fi/font2openvg/renderFont.cpp.txt
void evgText(VGfloat x, VGfloat y, const char *s, Fontinfo f, int pointsize) {
    VGfloat size = (VGfloat) pointsize, xx = x, mm[9];
    vgGetMatrix(mm);
    int character;
    unsigned char *ss = (unsigned char *)s;
    while ((ss = next_utf8_char(ss, &character)) != NULL) {
        int glyph = f.CharacterMap[character];
        if (character >= MAXFONTPATH-1) {
            continue;
        }
        if (glyph == -1) {
            continue;			   //glyph is undefined
        }
        VGfloat mat[9] = {
            size, 0.0f, 0.0f,
            0.0f, size, 0.0f,
            xx, y, 1.0f
        };
        vgLoadMatrix(mm);
        vgMultMatrix(mat);
        vgDrawPath(f.Glyphs[glyph], VG_FILL_PATH);
        xx += size * f.GlyphAdvances[glyph] / 65536.0f;
    }
    vgLoadMatrix(mm);
}

// TextWidth returns the width of a text string at the specified font and size.
VGfloat TextWidth(const char *s, Fontinfo f, int pointsize) {
    VGfloat tw = 0.0;
    VGfloat size = (VGfloat) pointsize;
    int character;
    unsigned char *ss = (unsigned char *)s;
    while ((ss = next_utf8_char(ss, &character)) != NULL) {
        int glyph = f.CharacterMap[character];
        if (character >= MAXFONTPATH-1) {
            continue;
        }
        if (glyph == -1) {
            continue;			   //glyph is undefined
        }
        tw += size * f.GlyphAdvances[glyph] / 65536.0f;
    }
    return tw;
}

// TextHeight reports a font's height
VGfloat TextHeight(Fontinfo f, int pointsize) {
    return (f.font_height * pointsize) / 65536;
}

// TextDepth reports a font's depth (how far under the baseline it goes)
VGfloat TextDepth(Fontinfo f, int pointsize) {
    return (-f.descender_height * pointsize) / 65536;
}

/**
 * Creates a new {VGPath}.
 *
 * @return  {VGPath}
 *          The {VGPath} that was created.
 */
VGPath evgNewPath() {
    return vgCreatePath(VG_PATH_FORMAT_STANDARD, VG_PATH_DATATYPE_F,
            1.0f, 0.0f, 0, 0, VG_PATH_CAPABILITY_APPEND_TO);
}

// TextMid draws text, centered on (x,y)
void evgTextMid(VGfloat x, VGfloat y, const char *s, Fontinfo f, int pointsize) {
    VGfloat tw = TextWidth(s, f, pointsize);
    evgText(x - (tw / 2.0), y, s, f, pointsize);
}

// TextEnd draws text, with its end aligned to (x,y)
void evgTextEnd(VGfloat x, VGfloat y, const char *s, Fontinfo f, int pointsize) {
    VGfloat tw = TextWidth(s, f, pointsize);
    evgText(x - tw, y, s, f, pointsize);
}

VGPath evgMakeCurve(VGubyte * segments, VGfloat * coords, VGbitfield flags) {
    VGPath path = evgNewPath();
    vgAppendPathData(path, 2, segments, coords);
    return path;
}

void evgCurve(VGubyte* segments, VGfloat* coords, VGbitfield flags) {
    VGPath path = evgMakeCurve(segments, coords, flags);
    vgDrawPath(path, flags);
    vgDestroyPath(path);
}

VGPath evgMakeCBezier(VGfloat sx, VGfloat sy, VGfloat cx, VGfloat cy, VGfloat px, VGfloat py, VGfloat ex, VGfloat ey) {
    VGubyte segments[] = { VG_MOVE_TO_ABS, VG_CUBIC_TO };
    VGfloat coords[] = { sx, sy, cx, cy, px, py, ex, ey };
    return evgMakeCurve(segments, coords, VG_FILL_PATH | VG_STROKE_PATH);
}

void evgCBezier(VGfloat sx, VGfloat sy, VGfloat cx, VGfloat cy, VGfloat px, VGfloat py, VGfloat ex, VGfloat ey) {
    VGPath path = evgMakeCBezier(sx, sy, cx, cy, px, py, ex, ey);
    vgDrawPath(path, flags);
    vgDestroyPath(path);
}

VGPath evgMakeQBezier(VGfloat sx, VGfloat sy, VGfloat cx, VGfloat cy, VGfloat ex, VGfloat ey) {
    VGubyte segments[] = { VG_MOVE_TO_ABS, VG_QUAD_TO };
    VGfloat coords[] = { sx, sy, cx, cy, ex, ey };
    return evgMakeCurve(segments, coords, VG_FILL_PATH | VG_STROKE_PATH);
}

// QBezier makes a quadratic bezier curve
void evgQbezier(VGfloat sx, VGfloat sy, VGfloat cx, VGfloat cy, VGfloat ex, VGfloat ey) {
    VGPath path = evgMakeQBezier(sx, sy, cx, cy, px, py, ex, ey);
    vgDrawPath(path, flags);
    vgDestroyPath(path);
}

VGPath evgMakePolygon(VGfloat* points, VGint n) {
    VGPath path = evgNewPath();
    vguPolygon(path, points, n, VG_FALSE);
    vgDrawPath(path, flag);
    vgDestroyPath(path);
}

void evgPolygon(VGfloat* points, VGint n, VGbitfield flag) {
    VGPath path = evgMakePolygon(points, n);
    vgDrawPath(path, flag);
    vgDestroyPath(path);
}

VGPath evgMakeRect(VGfloat x, VGfloat y, VGfloat width, VGfloat height) {
    VGPath path = evgNewPath();
    vguRect(path, x, y, width, height);
    return path;
}

void evgRect(VGfloat x, VGfloat y, VGfloat w, VGfloat h) {
    VGPath path = evgMakeRect();
    vgDrawPath(path, VG_FILL_PATH | VG_STROKE_PATH);
    vgDestroyPath(path);
}

VGPath evgMakeLine(VGfloat x1, VGfloat y1, VGfloat x2, VGfloat y2) {
    VGPath path = evgNewPath();
    vguLine(path, x1, y1, x2, y2);
    return path;
}

void evgLine(VGfloat x1, VGfloat y1, VGfloat x2, VGfloat y2) {
    VGPath path = evgNewLine(x1, y1, x2, y2);
    vgDrawPath(path, VG_STROKE_PATH);
    vgDestroyPath(path);
}

VGPath evgMakeRoundRect(VGfloat x, VGfloat y, VGfloat width, VGfloat height, VGfloat rw, VGfloat rh) {
    VGPath path = evgNewPath();
    vguRoundRect(path, x, y, w, h, rw, rh);
    return path;
}

void evgRoundRect(VGfloat x, VGfloat y, VGfloat width, VGfloat height, VGfloat rw, VGfloat rh) {
    VGPath path = evgMakeRoundRect(x, y, width, height, rw, rh);
    vgDrawPath(path, VG_FILL_PATH | VG_STROKE_PATH);
    vgDestroyPath(path);
}

VGPath evgMakeEllipse(VGfloat x, VGfloat y, VGfloat w, VGfloat h) {
    VGPath path = evgNewPath();
    vguEllipse(path, x, y, w, h);
    return path;
}

void evgEllipse(VGfloat x, VGfloat y, VGfloat w, VGfloat h) {
    VGPath path = evgMakeEllipse();
    vgDrawPath(path, VG_FILL_PATH | VG_STROKE_PATH);
    vgDestroyPath(path);
}

void evgMakeCircle(VGfloat x, VGfloat y, VGfloat radius) {
    const path = evgMakeEllipse(x, y, radius, radius);
    return path;
}

void evgCircle(VGfloat x, VGfloat y, VGfloat radius) {
    evgEllipse(x, y, radius);
}

VGPath evgMakeArc(VGfloat x, VGfloat y, VGfloat w, VGfloat h, VGfloat sa, VGfloat aext) {
    VGPath path = evgNewPath();
    vguArc(path, x, y, w, h, sa, aext, VGU_ARC_OPEN);
    return path
}

void evgArc(VGfloat x, VGfloat y, VGfloat w, VGfloat h, VGfloat sa, VGfloat aext) {
    VGPath path = evgMakeArc();
    vgDrawPath(path, VG_FILL_PATH | VG_STROKE_PATH);
    vgDestroyPath(path);
}

// Start begins the picture, clearing a rectangular region with a specified color
void evgBegin() {
    VGfloat color[4] = { 0, 0, 0, 1 };
    evgSetFill(color);
    evgSetStroke(color);
    evgStrokeWidth(0);
    vgLoadIdentity();
}

// End checks for errors, and renders to the display
void evgEnd() {
    assert(vgGetError() == VG_NO_ERROR);
    eglSwapBuffers(state->display, state->surface);
    assert(eglGetError() == EGL_SUCCESS);
}

// SaveEnd dumps the raster before rendering to the display
void evgSaveEnd(const char *filename) {
    FILE *fp;
    assert(vgGetError() == VG_NO_ERROR);
    if (strlen(filename) == 0) {
        dumpscreen(state->screen_width, state->screen_height, stdout);
    } else {
        fp = fopen(filename, "wb");
        if (fp != NULL) {
            dumpscreen(state->screen_width, state->screen_height, fp);
            fclose(fp);
        }
    }
    eglSwapBuffers(state->display, state->surface);
    assert(eglGetError() == EGL_SUCCESS);
}

// Backgroud clears the screen to a solid background color
void evgBackground(unsigned int r, unsigned int g, unsigned int b) {
    VGfloat colour[4];
    RGB(r, g, b, colour);
    vgSetfv(VG_CLEAR_COLOR, 4, colour);
    evgClear();
}

// BackgroundRGB clears the screen to a background color with alpha
void evgBackgroundRGB(unsigned int r, unsigned int g, unsigned int b, VGfloat a) {
    VGfloat colour[4];
    RGBA(r, g, b, a, colour);
    vgSetfv(VG_CLEAR_COLOR, 4, colour);
    evgClear();
}

// WindowOpacity sets the  window opacity
void evgWindowOpacity(unsigned int a) {
    dispmanChangeWindowOpacity(state, a);
}

// WindowPosition moves the window to given position
void WindowPosition(int x, int y) {
    dispmanMoveWindow(state, x, y);
}

// Outlined shapes
// Hollow shapes -because filling still happens even with a fill of 0,0,0,0
// unlike where using a strokewidth of 0 disables the stroke.
// Either this or change the original functions to require the VG_x_PATH flags

// CBezier makes a quadratic bezier curve, stroked
void CbezierOutline(VGfloat sx, VGfloat sy, VGfloat cx, VGfloat cy, VGfloat px, VGfloat py, VGfloat ex, VGfloat ey) {
    VGubyte segments[] = { VG_MOVE_TO_ABS, VG_CUBIC_TO };
    VGfloat coords[] = { sx, sy, cx, cy, px, py, ex, ey };
    evgMakeCurve(segments, coords, VG_STROKE_PATH);
}

// QBezierOutline makes a quadratic bezier curve, outlined
void QbezierOutline(VGfloat sx, VGfloat sy, VGfloat cx, VGfloat cy, VGfloat ex, VGfloat ey) {
    VGubyte segments[] = { VG_MOVE_TO_ABS, VG_QUAD_TO };
    VGfloat coords[] = { sx, sy, cx, cy, ex, ey };
    evgMakeCurve(segments, coords, VG_STROKE_PATH);
}

// RectOutline makes a rectangle at the specified location and dimensions, outlined
void RectOutline(VGfloat x, VGfloat y, VGfloat w, VGfloat h) {
    VGPath path = evgNewPath();
    vguRect(path, x, y, w, h);
    vgDrawPath(path, VG_STROKE_PATH);
    vgDestroyPath(path);
}

// RoundrectOutline  makes an rounded rectangle at the specified location and dimensions, outlined
void RoundrectOutline(VGfloat x, VGfloat y, VGfloat w, VGfloat h, VGfloat rw, VGfloat rh) {
    VGPath path = evgNewPath();
    vguRoundRect(path, x, y, w, h, rw, rh);
    vgDrawPath(path, VG_STROKE_PATH);
    vgDestroyPath(path);
}

// EllipseOutline makes an ellipse at the specified location and dimensions, outlined
void EllipseOutline(VGfloat x, VGfloat y, VGfloat w, VGfloat h) {
    VGPath path = evgNewPath();
    vguEllipse(path, x, y, w, h);
    vgDrawPath(path, VG_STROKE_PATH);
    vgDestroyPath(path);
}

// CircleOutline makes a circle at the specified location and dimensions, outlined
void CircleOutline(VGfloat x, VGfloat y, VGfloat r) {
    EllipseOutline(x, y, r, r);
}

// ArcOutline makes an elliptical arc at the specified location and dimensions, outlined
void ArcOutline(VGfloat x, VGfloat y, VGfloat w, VGfloat h, VGfloat sa, VGfloat aext) {
    VGPath path = evgNewPath();
    vguArc(path, x, y, w, h, sa, aext, VGU_ARC_OPEN);
    vgDrawPath(path, VG_STROKE_PATH);
    vgDestroyPath(path);
}
