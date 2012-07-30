all: shapes font2openvg

shapes:	shapes.c
	cc -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -o shapes shapes.c -L/opt/vc/lib -lGLESv2 

font2openvg:	font2openvg.cpp
	g++ -I /usr/include/freetype2 font2openvg.cpp   -o font2openvg -lfreetype

fonts: font2openvg 
	for f in /usr/share/fonts/truetype/ttf-dejavu/*.ttf; do fn=`basename $$f .ttf`; ./font2openvg $$f $$fn.inc $$fn; done
