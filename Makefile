all: shapedemo

libshapes.o:	libshapes.c shapes.h fontinfo.h
	cc -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -c libshapes.c

shapedemo:	shapedemo.c libshapes.o oglinit.o
	cc -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -o shapedemo shapedemo.c libshapes.o oglinit.o -L/opt/vc/lib -lGLESv2 -ljpeg

oglinit.o:	oglinit.c
	cc -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -c oglinit.c

font2openvg:	font2openvg.cpp
	g++ -I /usr/include/freetype2 font2openvg.cpp -o font2openvg -lfreetype

fonts: font2openvg
	for f in /usr/share/fonts/truetype/ttf-dejavu/*.ttf; do fn=`basename $$f .ttf`; ./font2openvg $$f $$fn.inc $$fn; done

test:	shapedemo
		./shapedemo ; ./shapedemo rand 50 ; ./shapedemo test hello ; ./shapedemo image ; ./shapedemo text ; ./shapedemo rotate 20 hello; ./shapedemo play

raw2png:	raw2png.go
	go build raw2png.go
indent:
	indent -linux -brf -l 132 shapedemo.c libshapes.c oglinit.c shapes.h fontinfo.h

