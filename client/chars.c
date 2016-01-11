// first OpenVG program
// Anthony Starks (ajstarks@gmail.com)
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include "VG/openvg.h"
#include "VG/vgu.h"
#include "fontinfo.h"
#include "shapes.h"

int main() {
	int width, height, i;
	char s[3], buf[400];

	init(&width, &height);				   // Graphics initialization

	Start(width, height);				   // Start the picture
	Background(0, 0, 0);				   // Black background
	Fill(44, 77, 232, 1);				   // Big blue marble
	Circle(width / 2, 0, width);			   // The "world"
	Fill(255, 255, 255, 1);				   // White text
	for ( i=0; i < 400; i++ ) {
		if (i < 32) {
			buf[i] = ' ';
		} else {
			buf[i] = (char)i;
		}
	}
	Text(100, 100, buf+200, SerifTypeface, 15);
	End();						   // End the picture

	fgets(s, 2, stdin);				   // look at the pic, end with [RETURN]
	finish();					   // Graphics cleanup
	exit(0);
}
