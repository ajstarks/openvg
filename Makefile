INCLUDEFLAGS=-I/opt/vc/include -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads
LIBFLAGS=-L/opt/vc/lib -lGLESv2 -lEGL -ljpeg
FONTLIB=/usr/share/fonts/truetype/ttf-dejavu
FONTFILES=DejaVuSans.inc  DejaVuSansMono.inc DejaVuSerif.inc
all:	libshapes.o fontinfo.o oglinit.o

libshapes.o:	libshapes.c shapes.h fontinfo.h fonts
	gcc -O2 -Wall $(INCLUDEFLAGS) -c libshapes.c

fontinfo.o:	fontinfo.c fontinfo.h shapes.h fonts
	gcc -O2 -Wall $(INCLUDEFLAGS) -c fontinfo.c

gopenvg:	openvg.go
	go install .

oglinit.o:	oglinit.c
	gcc -O2 -Wall $(INCLUDEFLAGS) -c oglinit.c

font2openvg:	fontutil/font2openvg.cpp
	g++ -I/usr/include/freetype2 fontutil/font2openvg.cpp -o font2openvg -lfreetype

fonts:	$(FONTFILES)

DejaVuSans.inc: font2openvg $(FONTLIB)/DejaVuSans.ttf
	./font2openvg $(FONTLIB)/DejaVuSans.ttf DejaVuSans.inc DejaVuSans

DejaVuSerif.inc: font2openvg $(FONTLIB)/DejaVuSerif.ttf
	./font2openvg $(FONTLIB)/DejaVuSerif.ttf DejaVuSerif.inc DejaVuSerif

DejaVuSansMono.inc: font2openvg $(FONTLIB)/DejaVuSansMono.ttf
	./font2openvg $(FONTLIB)/DejaVuSansMono.ttf DejaVuSansMono.inc DejaVuSansMono

indent:
	indent -linux -c 60 -brf -l 132  libshapes.c oglinit.c fontinfo.c shapes.h fontinfo.h

clean:
	rm -f *.o *.inc *.so font2openvg *.c~ *.h~

library: oglinit.o libshapes.o fontinfo.o indent
	gcc $(LIBFLAGS) -shared -o libshapes.so oglinit.o libshapes.o fontinfo.o

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
