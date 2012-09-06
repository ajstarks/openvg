all:	libshapes.o oglinit.o

libshapes.o:	libshapes.c shapes.h fontinfo.h
	gcc   -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -c libshapes.c
	go install .

oglinit.o:	oglinit.c
	gcc   -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -c oglinit.c

font2openvg:	font2openvg.cpp
	g++ -I /usr/include/freetype2 font2openvg.cpp -o font2openvg -lfreetype

fonts: font2openvg
	for f in /usr/share/fonts/truetype/ttf-dejavu/*.ttf; do fn=`basename $$f .ttf`; ./font2openvg $$f $$fn.inc $$fn; done

indent:
	indent -linux -brf -l 132  libshapes.c oglinit.c shapes.h fontinfo.h
