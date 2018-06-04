
GCC_INCLUDEFLAGS=-I/opt/vc/include -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads -fPIC
GCC_LIBFLAGS=-L/opt/vc/lib -lbrcmEGL -lbrcmGLESv2 -ljpeg

all:	lib	src
src:	libshapes	oglinit

install:	lib
	install -m 755 -p ./lib/font2openvg /usr/bin/
	install -m 755 -p ./build/libshapes.so /usr/lib/libshapes.so.1.0.0
	strip --strip-unneeded /usr/lib/libshapes.so.1.0.0
	ln -f -s /usr/lib/libshapes.so.1.0.0 /usr/lib/libshapes.so
	ln -f -s /usr/lib/libshapes.so.1.0.0 /usr/lib/libshapes.so.1
	ln -f -s /usr/lib/libshapes.so.1.0.0 /usr/lib/libshapes.so.1.0
	install -m 644 -p ./src/shapes.h /usr/include/
	install -m 644 -p ./src/fontinfo.h /usr/include/

uninstall:
	rm -f /usr/bin/font2openvg
	rm -f /usr/include/fontinfo.h
	rm -f /usr/lib/libshapes.so.1.0.0
	rm -f /usr/lib/libshapes.so.1.0
	rm -f /usr/lib/libshapes.so.1
	rm -f /usr/lib/libshapes.so
	rm -f /usr/include/shapes.h

libshapes:	./src/libshapes.c	./src/fontinfo.h	./src/shapes.h	fonts
	gcc -O2 -Wall $(GCC_INCLUDEFLAGS) -c ./src/libshapes.c -o ./build/libshapes

oglinit:	./src/oglinit.c
	gcc -O2 -Wall $(GCC_INCLUDEFLAGS) -c ./src/oglinit.c -o ./build/oglinit

fonts:	font
libs:	lib
lib:	font2openvg	font

font2openvg:	./lib/font2openvg.cpp
	g++ -I/usr/include/freetype2 lib/font2openvg.cpp -o lib/font2openvg -lfreetype

font:	/usr/share/fonts/truetype/ttf-dejavu/DejaVuSans.ttf	font2openvg
	./lib/font2openvg /usr/share/fonts/truetype/ttf-dejavu/DejaVuSans.ttf ./lib/DejaVuSans.inc DejaVuSans
