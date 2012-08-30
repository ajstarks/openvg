//
// first OpenVG program
// Anthony Starks (ajstarks@gmail.com)
//
// compile and run on the Raspberry Pi:
// cc -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -c libshapes.c
// cc -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -c oglinit.c
// cc -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -o hellovg hellovg.c libshapes.o oglinit.o -L/opt/vc/lib -lGLESv2 -ljpeg && ./hellovg
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
