all: shapes font2openvg

shapes:	shapes.c oglinit.o
	cc -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -o shapes oglinit.o shapes.c  -L/opt/vc/lib -lGLESv2 

oglinit.o:	oglinit.c
	cc -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -c oglinit.c

font2openvg:	font2openvg.cpp
	g++ -I /usr/include/freetype2 font2openvg.cpp   -o font2openvg -lfreetype

fonts: font2openvg 
	for f in /usr/share/fonts/truetype/ttf-dejavu/*.ttf; do fn=`basename $$f .ttf`; ./font2openvg $$f $$fn.inc $$fn; done
