GCC_INCLUDEFLAGS=-I/opt/vc/include/interface/vcos/pthreads -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include -fPIC
GCC_LIBFLAGS=-L/opt/vc/lib -lbrcmEGL -lbrcmGLESv2 -ljpeg

all:	lib	src
src:	libshapes	oglinit

install:	lib	src
	install -m 755 -p ./lib/font2openvg /usr/bin/
	install -m 755 -p ./build/libshapes.so /usr/lib/libshapes.so.1.0.0
	strip --strip-unneeded /usr/lib/libshapes.so.1.0.0
	ln -f -s /usr/lib/libshapes.so.1.0.0 /usr/lib/libshapes.so
	ln -f -s /usr/lib/libshapes.so.1.0.0 /usr/lib/libshapes.so.1
	ln -f -s /usr/lib/libshapes.so.1.0.0 /usr/lib/libshapes.so.1.0
	install -m 644 -p ./src/libshapes.h /usr/include/
	install -m 644 -p ./src/fontinfo.h /usr/include/

uninstall:
	rm -f /usr/bin/font2openvg
	rm -f /usr/include/fontinfo.h
	rm -f /usr/lib/libshapes.so.1.0.0
	rm -f /usr/lib/libshapes.so.1.0
	rm -f /usr/lib/libshapes.so.1
	rm -f /usr/lib/libshapes.so
	rm -f /usr/include/libshapes.h

build-dir:
	mkdir "./build/"

libshapes:	./src/libshapes.c	./src/fontinfo.h	./src/libshapes.h	build-dir	fonts
	gcc -O2 -Wall $(GCC_INCLUDEFLAGS) -c ./src/libshapes.c -o ./build/libshapes.o

oglinit:	./src/oglinit.c	build-dir
	gcc -O2 -Wall $(GCC_INCLUDEFLAGS) -c ./src/oglinit.c -o ./build/oglinit.o

libs:	lib
lib:	font2openvg	font

font2openvg:	./lib/font2openvg.cpp
	g++ -I/usr/include/freetype2 lib/font2openvg.cpp -o lib/font2openvg -lfreetype

fonts:	font
font:	/usr/share/fonts/truetype/ttf-dejavu/DejaVuSans.ttf	font2openvg
	./lib/font2openvg /usr/share/fonts/truetype/ttf-dejavu/DejaVuSans.ttf ./lib/DejaVuSans.inc DejaVuSans
