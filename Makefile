
GCC_INCLUDEFLAGS=-I/opt/vc/include -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads -fPIC
GCC_LIBFLAGS=-L/opt/vc/lib -lbrcmEGL -lbrcmGLESv2 -ljpeg
FONTFILES=DejaVuSans.inc

all:	lib	src	library
src:	libshapes	oglinit

clean:
	rm -f *.o *.inc *.so font2openvg *.c~ *.h~
	indent -linux -c 60 -brf -l 132  libshapes.c oglinit.c shapes.h fontinfo.h

install:
	install -m 755 -p font2openvg/ /usr/bin/
	install -m 755 -p libshapes.so /usr/lib/libshapes.so.1.0.0
	strip --strip-unneeded /usr/lib/libshapes.so.1.0.0
	ln -f -s /usr/lib/libshapes.so.1.0.0 /usr/lib/libshapes.so
	ln -f -s /usr/lib/libshapes.so.1.0.0 /usr/lib/libshapes.so.1
	ln -f -s /usr/lib/libshapes.so.1.0.0 /usr/lib/libshapes.so.1.0
	install -m 644 -p shapes.h /usr/include/
	install -m 644 -p fontinfo.h /usr/include/

uninstall:
	rm -f /usr/bin/font2openvg/
	rm -f /usr/lib/libshapes.so.1.0.0 /usr/lib/libshapes.so.1.0 /usr/lib/libshapes.so.1 /usr/lib/libshapes.so
	rm -f /usr/include/shapes.h /usr/include/fontinfo.h


libshapes:	libshapes.c	fontinfo.h	shapes.h	fonts
	gcc -O2 -Wall $(GCC_INCLUDEFLAGS) -c libshapes.c

oglinit:	oglinit.c
	gcc -O2 -Wall $(GCC_INCLUDEFLAGS) -c oglinit.c

fonts:	font
libs:	lib
lib:	font2openvg	font

font2openvg:	lib/font2openvg.cpp
	g++ -I/usr/include/freetype2 lib/font2openvg.cpp -o font2openvg -lfreetype

font:	/usr/share/fonts/truetype/ttf-dejavu/DejaVuSans.ttf	font2openvg
	font2openvg /usr/share/fonts/truetype/ttf-dejavu/DejaVuSans.ttf
	./font2openvg /usr/share/fonts/truetype/ttf-dejavu/DejaVuSans.ttf DejaVuSans.inc DejaVuSans


library:	libshapes	oglinit
	gcc $(GCC_LIBFLAGS) -shared -o libshapes.so oglinit.o libshapes.o
