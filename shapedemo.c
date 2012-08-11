//
// shapedemo: testbed for OpenVG APIs
// Anthony Starks (ajstarks@gmail.com)
//
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

#include "VG/openvg.h"
#include "VG/vgu.h"
#include "fontinfo.h"
#include "shapes.h"

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

typedef struct {
	Fontinfo font;
	VGfloat tw;
	int fontsize;
} FW;

// adjust the font to fit within a width
void fitwidth(int width, int adj, char *s, FW *f) {
	f->tw = textwidth(s, f->font, f->fontsize);
	while (f->tw > width) {
		f->fontsize -= adj;
		f->tw = textwidth(s, f->font, f->fontsize);
	}
}

// testpattern shows a test pattern 
void testpattern(int width, int height, char *s) {
	VGfloat llc[4] = {1,0,0,1}, ulc[4] = {0,1,0,1}, lrc[4] = {0,0,1,1}, 
			urc[4] = {0.5,0.5,0.5,1}, tc[4]  = {0.5,0,0,1}, bgcolor[4] = {1,1,1,1},
			midx1, midx2, midx3, midy1, midy2, midy3;
	int fontsize = 256, w2 = width/2, h2=height/2;
	FW tw1={MonoTypeface, 0, fontsize}, tw2={SansTypeface, 0, fontsize}, tw3={SerifTypeface,0, fontsize};

	Start(width, height, bgcolor);

	// colored squares in the corners
	setfill(llc); Rect(0,0,100,100);
	setfill(ulc); Rect(0,height-100,100,100);
	setfill(lrc); Rect(width-100,0,100,100);
	setfill(urc); Rect(width-100,height-100,100,100);
	
	// for each font, (Sans, Serif, Mono), adjust the string to the width
	fitwidth(width, 20, s, &tw1);
	fitwidth(width, 20, s, &tw2);
	fitwidth(width, 20, s, &tw3);

	// Determine the midpoint
	midx1 = w2-(tw1.tw/2);
	midx2 = w2-(tw2.tw/2);
	midx3 = w2-(tw3.tw/2);

	// Adjust the baselines to be medial
	midy1 = h2+20+(tw1.fontsize)/2;
	midy2 = h2-(tw2.fontsize)/2;
	midy3 = h2-20-tw2.fontsize-(tw3.fontsize)/2;

	Text(midx1, midy1, s, tw1.font, tw1.fontsize, urc);
	Text(midx2, midy2, s, tw2.font, tw2.fontsize, lrc);
	Text(midx3, midy3, s, tw3.font, tw3.fontsize, tc);
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
	Text(width*.45, height/2, "OpenVG on the Raspberry Pi", SansTypeface, 48, textcolor);
	textcolor[0] = 0;
	for (i=0; i < ns; i++) {
		Text(sx+sw+sw/2, sy, shapenames[i], SansTypeface, fontsize, textcolor);
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

// rotext draws text, rotated around the center of the screen, progressively faded
void rotext(int w, int h, int n, char *s) {
	int i;
	VGfloat textcolor[4] = {1,1,1,1}, bgcolor[4] = {0,0,0,1};
	VGfloat fade = (100.0/(VGfloat)n)/100.0;
	VGfloat deg = 360.0/n;

	VGfloat x = w/2, y = h/2;
	Start(w, h, bgcolor);
	Translate(x,y);
	for (i=0; i < n; i++) {
		Text(0,0, s, SerifTypeface, 256, textcolor);
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
	Text(50, 100, "OpenVG on the Raspberry Pi", SansTypeface, 64, textcolor);
	End();
}

// main initializes the system and shows the picture. 
// Exit and clean up when you hit [RETURN].
int main (int argc, char **argv) {
	int w, h, nr;
	char *usage = "%s test ...\n%s rand n\n%s rotate n ...\n";
	char *progname = argv[0];
	init(&w, &h);
	switch (argc) {
		case 3:
			if (strncmp(argv[1], "test", 4) == 0) {
				testpattern(w,h,argv[2]);
			} else if (strncmp(argv[1], "rand", 4)  == 0) {
				nr = atoi(argv[2]);
				if (nr < 1 || nr > 1000) {
					nr = 100;
				}
				rshapes(w, h, nr);
			} else {
				fprintf(stderr, usage, progname, progname, progname); 
				return 1;
			}
			break;

		case 4:
			if (strncmp(argv[1], "rotate", 6) == 0) {
				rotext(w, h, atoi(argv[2]), argv[3]);
			} else {
				fprintf(stderr, usage, progname, progname, progname);
				return 2;
			}	
			break;

		default:
			refcard(w,h);
	}
	while (getchar() != '\n') {
		;
	}
	finish();
	return 0;
}
