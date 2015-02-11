//
// clip: test rectangular clipping
//
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <fcntl.h>
#include <time.h>
#include "VG/openvg.h"
#include "VG/vgu.h"
#include "shapes.h"

int main(int argc, char **argv) {
	int w, h, fontsize;
	VGfloat x, cx, cy, cw, ch, midy, speed;
	char *message = "Now is the time for all good men to come to the aid of the party";
	char done[3];

	init(&w, &h);
	speed = 15;
	midy = (VGfloat) h / 2;
	fontsize = w / 50;
	cx = 0.0;
	ch = fontsize * 2;
	cw = w;
	cy = midy - (ch / 2);

	// scroll the text, only in the clipping rectangle
	for (x = 0; x < cw+speed; x += speed) {
		Start(w, h);
		Background(255, 255, 255);
		Fill(0,0,0,.2);
		Rect(cx, cy, cw, ch);
		ClipRect(cx, cy, cw, ch);
		Translate(x, cy + (fontsize / 2));
		Fill(0, 0, 0, 1);
		Text(0, 0, message, SansTypeface, fontsize);
		ClipEnd();
		End();
	}
	fgets(done, 2, stdin); // press [Return] when done
	finish();
	exit(0);
}
