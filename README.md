#Testbed for exploring OpenVG on the Raspberry Pi.


<a href="http://www.flickr.com/photos/ajstarks/7811750326/" title="rotext by ajstarks, on Flickr"><img src="http://farm8.staticflickr.com/7249/7811750326_614ea891ae.jpg" width="500" height="281" alt="rotext"></a>

## First program

Here is the graphics equivalent of "hello, world"

	//
	// first OpenVG program
	// Anthony Starks (ajstarks@gmail.com)
	//
	// compile and run on the Raspberry Pi:
	// cc -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -o hellovg hellovg.c -L/opt/vc/lib -lGLESv2 -ljpeg && ./hellovg
	//
	//
	#include <stdio.h>
	#include <stdlib.h>
	#include <unistd.h>

	#include "VG/openvg.h"
	#include "VG/vgu.h"
	#include "fontinfo.h"
	#include "shapes.h"

	int main() {
		int width, height;
		init(&width, &height);     // OpenGL, etc initialization

		Start(width, height);                     // Start the picture
		Background(0,0,0);                        // Black background
		Fill(44,77,232,1);                        // Big blue marble
		Circle(width/2, 0, width);                // The "world"
		Fill(255,255,255,1);                      // White text
		TextMiddle(width/2, height/2, 
		"hello, world", SerifTypeface, width/10); // Greetings 
		End();                                    // End the picture

		while (getchar() != '\n')  // look at the pic, end with [RETURN]
			;

		finish();                 // Graphics cleanup
		exit(0);
	}

<a href="http://www.flickr.com/photos/ajstarks/7828969180/" title="hellovg by ajstarks, on Flickr"><img src="http://farm9.staticflickr.com/8436/7828969180_b73db3bf19.jpg" width="500" height="281" alt="hellovg"></a>

## API

<a href="http://www.flickr.com/photos/ajstarks/7717370238/" title="OpenVG refcard by ajstarks, on Flickr"><img src="http://farm8.staticflickr.com/7256/7717370238_1d632cb179.jpg" width="500" height="281" alt="OpenVG refcard"></a>

Coordinates are VGfloat values, with the origin at the lower left, with x increasing to the right, and y increasing up.
OpenVG specifies colors as a VGfloat array containing red, green, blue, alpha values ranging from 0.0 to 1.0, but typically colors are specified as RGBA (0-255 for RGB, A from 0.0 to 1.0)

	void Start(int width, int height)
Begin the picture, clear the screen with a default white, set the stroke and fill to black.

	void End()
End the picture, rendering to the screen.

	void SaveEnd(char *filename)
End the picture, rendering to the screen, save the raster to the named file as 4-byte RGBA words, with a stride of
width*4 bytes. The program raw2png converts the "raw" raster to png.

### Attributes

	void setfill(float color[4])
Set the fill color

	void Background(unsigned int r, unsigned int g, unsigned int b)
Fill the screen with the background color defined from RGB values.

	void StrokeWidth(float width)
Set the stroke width.

	void RGBA(unsigned int r, unsigned int g, unsigned int b, VGfloat a, VGfloat color[4])
fill a color vector from RGBA values.

	void RGB(unsigned int r, unsigned int g, unsigned int b, VGfloat color[4])
fill a color vector from RGB values.

	void Stroke(unsigned int r, unsigned int g, unsigned int b, VGfloat a)
Set the Stroke color using RGBA values.

	void Fill(unsigned int r, unsigned int g, unsigned int b, VGfloat a)
Set the Fill color using RGBA values.

### Shapes

	void Line(VGfloat x1, VGfloat y1, VGfloat x2, VGfloat y2)
Draw a line between (x1, y1) and (x2, y2).

	void Rect(VGfloat x, VGfloat y, VGfloat w, VGfloat h)
Draw a rectangle with its origin (lower left) at (x,y), and size is (width,height).

	void Roundrect(VGfloat x, VGfloat y, VGfloat w, VGfloat h, VGfloat rw, VGfloat rh)
Draw a rounded rectangle with its origin (lower left) at (x,y), and size is (width,height).  
The width and height of the corners are specified with (rw,rh).

	void Polygon(VGfloat *x, VGfloat *y, VGint n)
Draw a polygon using the coordinates in arrays pointed to by x and y.  The number of coordinates is n.

	void Polyline(VGfloat *x, VGfloat *y, VGint n)
Draw a polyline using the coordinates in arrays pointed to by x and y.  The number of coordinates is n.

	void Circle(VGfloat x, VGfloat y, VGfloat r)
Draw a circle centered at (x,y) with radius r.

	void Ellipse(VGfloat x, VGfloat y, VGfloat w, VGfloat h)
Draw an ellipse centered at (x,y) with radii (w, h).

	void Qbezier(VGfloat sx, VGfloat sy, VGfloat cx, VGfloat cy, VGfloat ex, VGfloat ey)
Draw a quadratic bezier curve beginning at (sx, sy), using control points at (cx, cy), ending at (ex, ey).

	void Cbezier(VGfloat sx, VGfloat sy, VGfloat cx, VGfloat cy, VGfloat px, VGfloat py, VGfloat ex, VGfloat ey)
Draw a cubic bezier curve beginning at (sx, sy), using control points at (cx, cy) and (px, py), ending at (ex, ey).

	void Arc(VGfloat x, VGfloat y, VGfloat w, VGfloat h, VGfloat sa, VGfloat aext)
Draw an elliptical arc centered at (x, y), with width and height at (w, h).  Start angle (degrees) is sa, angle extent is aext.

### Text and Images

	void Text(VGfloat x, VGfloat y, const char* s, Fontinfo f, int pointsize)
Draw a the text srtring (s) at location (x,y), using pointsize.

	void TextMiddle(VGfloat x, VGfloat y, const char* s, Fontinfo f, int pointsize)
Draw a the text srtring (s) at centered at location (x,y), using pointsize.

	void TextEnd(VGfloat x, VGfloat y, const char* s, Fontinfo f, int pointsize)
Draw a the text srtring (s) at with its lend aligned to location (x,y), using pointsize

	void Image(VGfloat x, VGfloat y, int w, int h, char * filename)
place a JPEG image with dimensions (w,h) at (x,y).

	
### Transformations

	void Translate(VGfloat x, VGfloat y)
Translate the coordinate system to (x,y).

	void Rotate(VGfloat r)
Rotate the coordinate system around angle r (degrees).

	void Scale(VGfloat x, VGfloat y)
Scale by x,y.

	void Shear(VGfloat x, VGfloat y)
Shear by the angles x,y.


## Using fonts

Also included is the font2openvg program, which turns font information into C source that 
you can embed in your program. The Makefile makes font code from files found in /usr/share/fonts/truetype/ttf-dejavu/. 
If you want to use other fonts, adjust the Makefile accordingly, or generate the font code on your own once the font2openvg program is built.

font2openvg takes three arguments: the TrueType font file, the output file to be included and the prefix for identifiers.
For example to use the DejaVu Sans font:

	./font2openvg /usr/share/fonts/truetype/ttf-dejavu/DejaVuSans.ttf DejaVuSans.inc DejaVuSans

and include the generated code in your program:

	#include "DejaVuSans.inc"
	Fontinfo DejaFont
	
The loadfont function creates OpenVG paths from the font data:

	loadfont(DejaVuSans_glyphPoints, 
            DejaVuSans_glyphPointIndices, 
        	DejaVuSans_glyphInstructions,                
        	DejaVuSans_glyphInstructionIndices, 
            DejaVuSans_glyphInstructionCounts, 
            DejaVuSans_glyphAdvances,
            DejaVuSans_characterMap, 
        	DejaVuSans_glyphCount);

The unloadfont function releases the path information:
	
	unloadfont(DejaFont.Glyphs, DejaFont.Count);

# Build and run

	pi@raspberrypi ~ $ git clone git://github.com/ajstarks/openvg
	pi@raspberrypi ~ $ cd openvg
	pi@raspberrypi ~/openvg $ make test
	cc -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -c libshapes.c
	cc -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -c oglinit.c
	cc -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -o shapedemo shapedemo.c libshapes.o oglinit.o -L/opt/vc/lib -lGLESv2 -ljpeg
	./shapedemo demo 5


	The program "shapedemo" exercises a high-level API built on OpenVG found in libshapes.c. 

	./shapedemo                      # show a reference card
	./shapedemo image                # show four test images
	./shapedemo text                 # show blocks of text in serif, sans, and mono fonts
	./shapedemo rand 10              # show 10 random shapes
	./shapedemo rotate 10 a          # rotated and faded "a"
	./shapedemo test "hello, world"  # show a test pattern, with "hello, world" at mid-display in sans, serif, and mono.
	./shapedemo demo 10              # run through the demo, pausing 10 seconds between each one; contemplate the awesome.

