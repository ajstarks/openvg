LIBFLAGS=-I/opt/vc/include -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads
FONTLIB=/usr/share/fonts/truetype/ttf-dejavu/
all:	libshapes.o oglinit.o gopenvg

libshapes.o:	libshapes.c shapes.h fontinfo.h oglinit.o
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
