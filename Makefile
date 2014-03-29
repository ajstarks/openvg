LIBFLAGS=-I/opt/vc/include -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads
FONTLIB=/usr/share/fonts/truetype/ttf-dejavu/
all:	libshapes.o oglinit.o gopenvg

libshapes.o:	libshapes.c shapes.h fontinfo.h fonts
	gcc -O2  -Wall $(LIBFLAGS) -c libshapes.c

gopenvg:	openvg.go
	go install .

oglinit.o:	oglinit.c
	gcc  -O2  -Wall $(LIBFLAGS) -c oglinit.c

font2openvg:	fontutil/font2openvg.cpp
	g++ -I /usr/include/freetype2 fontutil/font2openvg.cpp -o font2openvg -lfreetype

fonts: font2openvg
	for f in $(FONTLIB)/*.ttf; do fn=`basename $$f .ttf`; ./font2openvg $$f $$fn.inc $$fn; done

indent:
	indent -linux -brf -l 132  libshapes.c oglinit.c shapes.h fontinfo.h

clean:
	rm -f *.o *.inc *.so font2openvg

library: oglinit.o libshapes.o indent
	gcc -L/opt/vc/lib -lGLESv2 -lEGL -ljpeg -shared -o libshapes.so oglinit.o libshapes.o

install:
	install -m 755 -p font2openvg /usr/bin/
	install -m 755 -p libshapes.so /usr/lib/libshapes.so.1.0.0
	strip --strip-unneeded /usr/lib/libshapes.so.1.0.0
	ln -f -s /usr/lib/libshapes.so.1.0.0 /usr/lib/libshapes.so
	ln -f -s /usr/lib/libshapes.so.1.0.0 /usr/lib/libshapes.so.1
	ln -f -s /usr/lib/libshapes.so.1.0.0 /usr/lib/libshapes.so.1.0
	install -m 644 -p shapes.h /usr/include/
	install -m 644 -p fontinfo.h /usr/include/

uninstall:
	rm -f /usr/bin/font2openvg
	rm -f /usr/lib/libshapes.so.1.0.0 /usr/lib/libshapes.so.1.0 /usr/lib/libshapes.so.1 /usr/lib/libshapes.so
	rm -f /usr/include/shapes.h /usr/include/fontinfo.h
