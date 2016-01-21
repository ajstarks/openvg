INCLUDEFLAGS=-I/opt/vc/include -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads -I..
LIBFLAGS=-L/opt/vc/lib -lEGL -lGLESv2 -lbcm_host -lpthread  -ljpeg

all: shapedemo hellovg mouse-hellovg particles screensize clip

shapedemo:	shapedemo.c ../libshapes.o ../oglinit.o
	gcc -Wall $(INCLUDEFLAGS) -o shapedemo shapedemo.c ../libshapes.o ../oglinit.o $(LIBFLAGS)

test:	shapedemo
	./shapedemo demo 5

hellovg:	hellovg.c ../libshapes.o ../oglinit.o
	gcc -Wall $(INCLUDEFLAGS) -o  hellovg hellovg.c ../libshapes.o ../oglinit.o $(LIBFLAGS)

mouse-hellovg:	mouse-hellovg.c ../libshapes.o ../oglinit.o
	gcc -Wall $(INCLUDEFLAGS) -o  mouse-hellovg mouse-hellovg.c ../libshapes.o ../oglinit.o $(LIBFLAGS)

particles:	particles.c ../libshapes.o ../oglinit.o
	gcc -Wall $(INCLUDEFLAGS) -o  particles particles.c ../libshapes.o ../oglinit.o $(LIBFLAGS)

screensize:	screensize.c ../libshapes.o ../oglinit.o
	gcc -Wall $(INCLUDEFLAGS) -o  screensize screensize.c ../libshapes.o ../oglinit.o $(LIBFLAGS)

clip:	clip.c ../libshapes.o ../oglinit.o
	gcc -Wall $(INCLUDEFLAGS) -o  clip clip.c ../libshapes.o ../oglinit.o $(LIBFLAGS)

indent:
	indent -linux -c 60 -brf -l 132 shapedemo.c hellovg.c mouse-hellovg.c particles.c screensize.c
