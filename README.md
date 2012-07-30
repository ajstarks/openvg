=Tools for exploring OpenVG on the Raspberry Pi.=

For now there is a single program "shapes". If you run the program with no command line arguments, you get a "reference card" that 
demonstrates the calls in the library (Text, Arc, Circle, Ellipse, Rect, Roundrect, Line, Polyline, Polygon, Cubic Bezier, 
and Quadratic Bezier). Running with a numeric command line argument shows the specified number of random shapes.

Also included is the font2openvg program, which will turn font information into C source that 
you can embed in your program. The Makefile makes font code from files found in /usr/share/fonts/truetype/ttf-dejavu/. 
If you want to use other fonts, adjust the Makefile accordingly, or generate them on your own once, the font2openvg program is built.


	pi@raspberrypi ~/vg $ make fonts shapes
	g++ -I /usr/include/freetype2 font2openvg.cpp   -o font2openvg -lfreetype
	for f in /usr/share/fonts/truetype/ttf-dejavu/*.ttf; do fn=`basename $f .ttf`; ./font2openvg $f $fn.inc $fn; done
	224 glyphs written
	224 glyphs written
	224 glyphs written
	224 glyphs written
	224 glyphs written
	224 glyphs written
	cc -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -o shapes shapes.c -L/opt/vc/lib -lGLESv2)
	pi@raspberrypi ~/vg $ ./shapes # hit return when you are done looking at the awesomness

	pi@raspberrypi ~/vg $ ./shapes 100 # show 100 random shapes