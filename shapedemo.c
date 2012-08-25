//
// shapedemo: testbed for OpenVG APIs
// Anthony Starks (ajstarks@gmail.com)
//
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <fcntl.h>
#include "VG/openvg.h"
#include "VG/vgu.h"
#include "fontinfo.h"
#include "shapes.h"

// randcolor returns a random number 0..255
unsigned int randcolor() {
	return (unsigned int)(drand48() * 255.0);
}

// randf returns a floating point number bounded by n
VGfloat randf(int n) {
	return drand48() * n;
}

// coordpoint marks a coordinate, preserving a previous color
void coordpoint(VGfloat x, VGfloat y, VGfloat size, VGfloat pcolor[4]) {
	Fill(128, 0, 0, 0.3);
	Circle(x, y, size);
	setfill(pcolor);
}

// grid draws a grid
void grid(VGfloat x, VGfloat y, int n, int w, int h) {
	VGfloat ix, iy;
	Stroke(128, 128, 128, 0.5);
	StrokeWidth(2);
	for (ix = x; ix <= x + w; ix += n) {
		Line(ix, y, ix, y + h);
	}

	for (iy = y; iy <= y + h; iy += n) {
		Line(x, iy, x + w, iy);
	}
}

typedef struct {
	Fontinfo font;
	VGfloat tw;
	int fontsize;
} FW;

// adjust the font to fit within a width
void fitwidth(int width, int adj, char *s, FW * f) {
	f->tw = textwidth(s, f->font, f->fontsize);
	while (f->tw > width) {
		f->fontsize -= adj;
		f->tw = textwidth(s, f->font, f->fontsize);
	}
}

// testpattern shows a test pattern 
void testpattern(int w, int h, char *s) {
	VGfloat midx, midy1, midy2, midy3;
	int fontsize = 256, h2 = h / 2;
	FW tw1 = { MonoTypeface, 0, fontsize };
	FW tw2 = { SansTypeface, 0, fontsize };
	FW tw3 = { SerifTypeface, 0, fontsize };

	Start(w, h);

	// colored squares in the corners
	Fill(255, 0, 0, 1);
	Rect(0, 0, 100, 100);
	Fill(0, 255, 0, 1);
	Rect(0, h - 100, 100, 100);
	Fill(0, 0, 255, 1);
	Rect(w - 100, 0, 100, 100);
	Fill(128, 128, 128, 1);
	Rect(w - 100, h - 100, 100, 100);

	// for each font, (Sans, Serif, Mono), adjust the string to the w
	fitwidth(w, 20, s, &tw1);
	fitwidth(w, 20, s, &tw2);
	fitwidth(w, 20, s, &tw3);

	midx = w / 2;

	// Adjust the baselines to be medial
	midy1 = h2 + 20 + (tw1.fontsize) / 2;
	midy2 = h2 - (tw2.fontsize) / 2;
	midy3 = h2 - 20 - tw2.fontsize - (tw3.fontsize) / 2;

	Fill(128, 128, 128, 1);
	TextMiddle(midx, midy1, s, tw1.font, tw1.fontsize);
	Fill(128, 0, 0, 1);
	TextMiddle(midx, midy2, s, tw2.font, tw2.fontsize);
	Fill(0, 0, 128, 1);
	TextMiddle(midx, midy3, s, tw3.font, tw3.fontsize);
	End();
}

// textlines writes lines of text
void textlines(VGfloat x, VGfloat y, char *s[], Fontinfo f, int fontsize, VGfloat leading) {
	int i;
	for (i = 0;; i++) {
		if (s[i] == NULL) {
			break;
		}
		Text(x, y, s[i], f, fontsize);
		y -= leading;
	}
}

// tb draws a block of text
void tb(int w, int h) {
	char *para[] = {
		"For lo, the winter is past,",
		"the rain is over and gone",
		"the flowers appear on the earth;",
		"the time for the singing of birds is come,",
		"and the voice of the turtle is heard in our land",
		NULL
	};

	VGfloat tmargin = w * 0.25, lmargin = w * 0.10, top = h * .9, mid = h * .6, bot = h * .3;

	int fontsize = 24, leading = 40, lfontsize = fontsize * 2, midb = ((leading * 2) + (leading / 2)) - (lfontsize / 2);

	Start(w, h);
	Fill(49, 79, 79, 1);
	textlines(tmargin, top, para, SerifTypeface, fontsize, leading);
	textlines(tmargin, mid, para, SansTypeface, fontsize, leading);
	textlines(tmargin, bot, para, MonoTypeface, fontsize, leading);
	Text(lmargin, top - midb, "Serif", SansTypeface, lfontsize);
	Text(lmargin, mid - midb, "Sans", SansTypeface, lfontsize);
	Text(lmargin, bot - midb, "Mono", SansTypeface, lfontsize);
	End();
}

// cookie draws a cookie
void cookie(int w, int h) {
	int ew = 200, eh = 60, h2 = h / 2, w2 = w / 2;

	Start(w, h);
	Fill(128, 128, 128, 1);
	Ellipse(w2, h2, ew, eh);
	Translate(0, 10);

	Fill(255, 255, 255, 1);
	Ellipse(w2, h2, ew, eh);
	Translate(0, 20);

	Fill(0, 0, 0, 1);
	Ellipse(w2, h2, ew, eh);
	End();
}

// imagtest displays four JPEG images, centered on the display
void imagetest(int w, int h) {
	int imgw = 400, imgh = 400;
	VGfloat cx = (w / 2) - (imgw / 2), cy = (h / 2) - (imgh / 2);
	VGfloat ulx = 0, uly = h - imgh;
	VGfloat urx = w - imgw, ury = uly;
	VGfloat llx = 0, lly = 0;
	VGfloat lrx = urx, lry = lly;
	Start(w, h);
	Background(0, 0, 0);
	Image(cx, cy, imgw, imgh, "test_img_violin.jpg");
	Image(ulx, uly, imgw, imgh, "test_img_piano.jpg");
	Image(urx, ury, imgw, imgh, "test_img_sax.jpg");
	Image(llx, lly, imgw, imgh, "test_img_guitar.jpg");
	Image(lrx, lry, imgw, imgh, "test_img_flute.jpg");
	End();
}

// fontrange shows a range of fonts
void fontrange(VGfloat x, VGfloat y, int w, int h) {
	int i, sizes[] = { 6, 7, 8, 9, 10, 11, 12, 14, 16, 18, 21, 24, 36, 48, 60, 72, 96, 0 };
	VGfloat spacing = 100, xs = x, s2 = spacing / 2, lx = xs - sizes[0];
	char num[4];

	Start(w, h);
	Background(255, 255, 255);
	for (i = 0; sizes[i] != 0; i++) {
		Fill(128, 0, 0, 1);
		TextMiddle(x, y, "a", SerifTypeface, sizes[i]);
		Fill(128, 128, 128, 1);
		sprintf(num, "%d", sizes[i]);
		TextMiddle(x, y - spacing, num, SansTypeface, 16);
		x += sizes[i] + 50;
	}
	x -= sizes[i - 1];
	Stroke(150, 150, 150, 0.5);
	Fill(255, 255, 255, 1);
	StrokeWidth(2);
	Line(lx, y - s2, x, y - s2);
	Qbezier(lx, y + s2, x, y + s2, x, y + (spacing * 2));
	End();
}

// refcard shows a reference card of shapes
void refcard(int width, int height) {
	char *shapenames[] = {
		"Circle",
		"Ellipse",
		"Rectangle",
		"Rounded Rectangle",
		"Line",
		"Polyline",
		"Polygon",
		"Arc",
		"Quadratic Bezier",
		"Cubic Bezier",
		"Image"
	};
	VGfloat shapecolor[4];
	RGB(202, 225, 255, shapecolor);
	VGfloat top = height * .95, sx = 500, sy = top, sw = width * .05, sh = height * .045, dotsize = 7, spacing = 2.0;

	int i, ns = sizeof(shapenames) / sizeof(char *), fontsize = height * .033;
	Start(width, height);
	sx = width * 0.10;

	Fill(128, 0, 0, 1);
	TextEnd(width - 20, height / 2, "OpenVG on the Raspberry Pi", SansTypeface, fontsize + (fontsize / 2));
	Fill(0, 0, 0, 1);
	for (i = 0; i < ns; i++) {
		Text(sx + sw + sw / 2, sy, shapenames[i], SansTypeface, fontsize);
		sy -= sh * spacing;
	}
	sy = top;
	VGfloat cx = sx + (sw / 2), ex = sx + sw;
	setfill(shapecolor);
	Circle(cx, sy, sw);
	coordpoint(cx, sy, dotsize, shapecolor);
	sy -= sh * spacing;
	Ellipse(cx, sy, sw, sh);
	coordpoint(cx, sy, dotsize, shapecolor);
	sy -= sh * spacing;
	Rect(sx, sy, sw, sh);
	coordpoint(sx, sy, dotsize, shapecolor);
	sy -= sh * spacing;
	Roundrect(sx, sy, sw, sh, 20, 20);
	coordpoint(sx, sy, dotsize, shapecolor);
	sy -= sh * spacing;

	StrokeWidth(1);
	Stroke(204, 204, 204, 1);
	Line(sx, sy, ex, sy);
	coordpoint(sx, sy, dotsize, shapecolor);
	coordpoint(ex, sy, dotsize, shapecolor);
	sy -= sh;

	VGfloat px[5] = { sx, sx + (sw / 4), sx + (sw / 2), sx + ((sw * 3) / 4), sx + sw };
	VGfloat py[5] = { sy, sy - sh, sy, sy - sh, sy };

	Polyline(px, py, 5);
	coordpoint(px[0], py[0], dotsize, shapecolor);
	coordpoint(px[1], py[1], dotsize, shapecolor);
	coordpoint(px[2], py[2], dotsize, shapecolor);
	coordpoint(px[3], py[3], dotsize, shapecolor);
	coordpoint(px[4], py[4], dotsize, shapecolor);
	sy -= sh * spacing;

	py[0] = sy;
	py[1] = sy - sh;
	py[2] = sy - (sh / 2);
	py[3] = py[1] - (sh / 4);
	py[4] = sy;
	Polygon(px, py, 5);
	sy -= (sh * spacing) + sh;

	Arc(sx + (sw / 2), sy, sw, sh, 0, 180);
	coordpoint(sx + (sw / 2), sy, dotsize, shapecolor);
	sy -= sh * spacing;

	VGfloat cy = sy + (sh / 2), ey = sy;
	Qbezier(sx, sy, cx, cy, ex, ey);
	coordpoint(sx, sy, dotsize, shapecolor);
	coordpoint(cx, cy, dotsize, shapecolor);
	coordpoint(ex, ey, dotsize, shapecolor);
	sy -= sh * spacing;

	ey = sy;
	cy = sy + sh;
	Cbezier(sx, sy, cx, cy, cx, sy, ex, ey);
	coordpoint(sx, sy, dotsize, shapecolor);
	coordpoint(cx, cy, dotsize, shapecolor);
	coordpoint(cx, sy, dotsize, shapecolor);
	coordpoint(ex, ey, dotsize, shapecolor);

	sy -= (sh * spacing * 1.5);
	Image(sx, sy, 100, 100, "starx.jpg");

	End();
}

// rotext draws text, rotated around the center of the screen, progressively faded
void rotext(int w, int h, int n, char *s) {
	VGfloat fade = (100.0 / (VGfloat) n) / 100.0;
	VGfloat deg = 360.0 / n;
	VGfloat x = w / 2, y = h / 2;
	VGfloat alpha = 1.0;	// start solid
	int i, size = w / 8;

	Start(w, h);
	Background(0, 0, 0);
	Translate(x, y);
	for (i = 0; i < n; i++) {
		Fill(255, 255, 255, alpha);
		Text(0, 0, s, SerifTypeface, size);
		alpha -= fade;	// fade
		size += n;	// enlarge
		Rotate(deg);
	}
	End();
}

// rseed seeds the random number generator from the random device
void rseed(void) {
	unsigned char d[sizeof(long int)];
	long int s;
	int fd;

	// read bytes from the random device,
	// pack them into a long int.
	fd = open("/dev/urandom", O_RDONLY);
	if (fd < 0) {
		srand48(1);
		return;
	}
	read(fd, (void *)d, (size_t) sizeof(long int));
	s = d[3] | (d[2] << 8) | (d[1] << 16) | (d[0] << 24);
	srand48(s);
	close(fd);
}

// rshapes draws shapes with random colors, strokes, and sizes. 
void rshapes(int width, int height, int n) {
	int i, j, np = 10;
	VGfloat sx, sy, cx, cy, px, py, ex, ey, pox, poy;
	VGfloat polyx[np], polyy[np];
	rseed();
	Start(width, height);
	for (i = 0; i < n; i++) {
		Fill(randcolor(), randcolor(), randcolor(), drand48());
		Ellipse(randf(width), randf(height), randf(200), randf(100));
		Circle(randf(width), randf(height), randf(100));
		Rect(randf(width), randf(height), randf(200), randf(100));
		Arc(randf(width), randf(height), randf(200), randf(200), randf(360), randf(360));

		sx = randf(width);
		sy = randf(height);
		Stroke(randcolor(), randcolor(), randcolor(), 1);
		StrokeWidth(randf(5));
		Line(sx, sy, sx + randf(200), sy + randf(100));
		StrokeWidth(0);

		sx = randf(width);
		sy = randf(height);
		ex = sx + randf(200);
		ey = sy;
		cx = sx + ((ex - sx) / 2.0);
		cy = sy + randf(100);
		Qbezier(sx, sy, cx, cy, ex, ey);

		sx = randf(width);
		sy = randf(height);
		ex = sx + randf(200);
		ey = sy;
		cx = sx + ((ex - sx) / 2.0);
		cy = sy + randf(100);
		px = cx;
		py = sy - randf(100);
		Cbezier(sx, sy, cx, cy, px, py, ex, ey);

		pox = randf(width);
		poy = randf(height);
		for (j = 0; j < np; j++) {
			polyx[j] = pox + randf(200);
			polyy[j] = poy + randf(100);
		}
		Polygon(polyx, polyy, np);

		pox = randf(width);
		poy = randf(height);
		for (j = 0; j < np; j++) {
			polyx[j] = pox + randf(200);
			polyy[j] = poy + randf(100);
		}
		Polyline(polyx, polyy, np);
	}
	Fill(128, 0, 0, 1);
	Text(20, 20, "OpenVG on the Raspberry Pi", SansTypeface, 32);
	End();
}

void sunearth(int w, int h) {
	VGfloat sun, earth, x, y;
	int i;

	rseed();
	Start(w, h);
	Background(0, 0, 0);
	Fill(255, 255, 255, 1);
	for (i = 0; i < w / 4; i++) {
		x = randf(w);
		y = randf(h);
		Circle(x, y, 2);
	}
	earth = (VGfloat) w *0.010;
	sun = earth * 109;
	Fill(0, 0, 255, 1);
	Circle(w / 3, h - (h / 10), earth);
	Fill(255, 255, 224, 1);
	Circle(w, 0, sun);
	End();
}

// advert is an ad for the package 
void advert(int w, int h) {
	VGfloat y = (6 * h) / 10;
	int fontsize = w * 0.04;
	char *s = "github.com/ajstarks/openvg";
	char *a = "ajstarks@gmail.com";
	int imw = 110, imh = 110;
	VGfloat tw = textwidth(s, SansTypeface, fontsize);

	Start(w, h);
	Fill(128, 0, 0, 1);
	Text(w / 2 - (tw / 2), y - (fontsize / 4), s, SansTypeface, fontsize);
	y -= 150;
	tw = textwidth(a, SansTypeface, fontsize / 3);
	Fill(128, 128, 128, 1);
	Text(w / 2 - (tw / 2), y, a, SansTypeface, fontsize / 3);
	Image((w / 2) - (imw / 2), y - (imh * 2), imw, imh, "starx.jpg");
	End();
}

// main initializes the system and shows the picture. 
// Exit and clean up when you hit [RETURN].
int main(int argc, char **argv) {
	int w, h, nr;
	char *usage = "%s [command]\n\tdemo sec\n\tastro\n\ttest ...\n\trand n\n\trotate n ...\n\timage\n\ttext\n";
	char *progname = argv[0];
	init(&w, &h);
	switch (argc) {
	case 2:
		if (strncmp(argv[1], "image", 5) == 0) {
			imagetest(w, h);
		} else if (strncmp(argv[1], "text", 4) == 0) {
			tb(w, h);
		} else if (strncmp(argv[1], "astro", 5) == 0) {
			sunearth(w, h);
		} else if (strncmp(argv[1], "fr", 2) == 0) {
			fontrange(100, h / 2, w, h);
		} else {
			fprintf(stderr, usage, progname);
			return 1;
		}
		break;
	case 3:
		nr = atoi(argv[2]);
		if (strncmp(argv[1], "demo", 4) == 0) {
			if (nr < 1 || nr > 30) {
				nr = 5;
			}
			refcard(w, h);
			sleep(nr);
			rshapes(w, h, 50);
			sleep(nr);
			testpattern(w, h, "OpenVG on RasPi");
			sleep(nr);
			imagetest(w, h);
			sleep(nr);
			rotext(w, h, 30, "Raspi");
			sleep(nr);
			tb(w, h);
			sleep(nr);
			fontrange(100, h / 2, w, h);
			sleep(nr);
			sunearth(w, h);
			sleep(nr);
			advert(w, h);
		} else if (strncmp(argv[1], "rand", 4) == 0) {
			if (nr < 1 || nr > 1000) {
				nr = 100;
			}
			rshapes(w, h, nr);
		} else if (strncmp(argv[1], "test", 4) == 0) {
			testpattern(w, h, argv[2]);
		} else {
			fprintf(stderr, usage, progname);
			return 1;
		}
		break;

	case 4:
		if (strncmp(argv[1], "rotate", 6) == 0) {
			rotext(w, h, atoi(argv[2]), argv[3]);
		} else {
			fprintf(stderr, usage, progname);
			return 1;
		}
		break;

	default:
		refcard(w, h);
	}
	while (getchar() != '\n') {
		;
	}
	finish();
	return 0;
}
