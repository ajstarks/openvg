INCLUDEFLAGS=-I/opt/vc/include -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads -I..
LIBFLAGS=-L/opt/vc/lib -lGLESv2 -lEGL -ljpeg

all: shapedemo hellovg mouse-hellovg

shapedemo:	shapedemo.c ../libshapes.o ../oglinit.o
	gcc   -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -I.. -o shapedemo shapedemo.c ../libshapes.o ../oglinit.o -L/opt/vc/lib -lGLESv2 -ljpeg

test:	shapedemo
	./shapedemo demo 5

hellovg:	hellovg.c ../libshapes.o ../oglinit.o
	gcc   -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -I.. -o  hellovg hellovg.c ../libshapes.o ../oglinit.o -L/opt/vc/lib -lGLESv2 -ljpeg

mouse-hellovg:	mouse-hellovg.c ../libshapes.o ../oglinit.o
	gcc   -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -I.. -o  mouse-hellovg mouse-hellovg.c ../libshapes.o ../oglinit.o -L/opt/vc/lib -lGLESv2 -ljpeg
indent:
	indent -linux -c 60 -brf -l 132 shapedemo.c hellovg.c mouse-hellovg.c
