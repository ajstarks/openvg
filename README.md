[header-image]: https://raw.githubusercontent.com/mgthomas99/openvg/develop/.github/assets/raspi-spiral.png
[git-repository-url]: https://github.com/mgthomas99/openvg
[license-shield-url]: https://img.shields.io/github/license/mgthomas99/openvg.svg?style=flat-square
[license-url]: https://github.com/mgthomas99/openvg/blob/master/LICENSE

*The repository is currently not accepting issue submissions as it is a fork of
a different repository that already has issues. Issue submissions will be
enabled shortly.*

# OpenVG on Raspberry Pi

[![One][header-image]][git-repository-url]
[![LICENSE][license-shield-url]][license-url]

This project provides an abstraction layer on top of native OpenVG, for quickly
creating contexts and drawing shapes, images, and text.

*This project is a fork of [OpenVG by Anthony Starks <ajstarks>](https://github.com/ajstarks/openvg)*

## First program

```c
    #include <stdio.h>
    #include <stdlib.h>
    #include <unistd.h>
    #include "VG/openvg.h"
    #include "VG/vgu.h"
    #include "fontinfo.h"
    #include "shapes.h"

    int main() {
        int width, height;
        char s[3];

        init(&width, &height);

        Start(width, height);
        Background(0, 0, 0);
        Fill(44, 77, 232, 1);
        Circle(width / 2, 0, width);
        Fill(255, 255, 255, 1);
        TextMid(width / 2, height / 2, "hello, world", SerifTypeface, width / 10);
        End();

        fgets(s, 2, stdin);
        finish();
        exit(0)
    }
```

## API

![One](http://farm8.staticflickr.com/7256/7717370238_1d632cb179.jpg)

Coordinates are VGfloat values, with the origin at the lower left, with x
increasing to the right, and y increasing up. OpenVG specifies colors as a
`VGfloat` array containing red, green, blue, alpha values ranging from 0.0 to
1.0, but typically colors are specified as RGBA (0-255 for RGB, A from 0.0 to
1.0).

## Using fonts

Also included is the font2openvg program, which turns font information into C source that
you can embed in your program. The Makefile makes font code from files found in /usr/share/fonts/truetype/ttf-dejavu/.
If you want to use other fonts, adjust the Makefile accordingly, or generate the font code on your own once the font2openvg program is built.

font2openvg takes three arguments: the TrueType font file, the output file to be included and the prefix for identifiers.
For example to use the DejaVu Sans font:

    ./font2openvg /usr/share/fonts/truetype/ttf-dejavu/DejaVuSans.ttf DejaVuSans.inc DejaVuSans

and include the generated code in your program:

    #include "DejaVuSans.inc"
    Fontinfo DejaFont

The loadfont function creates OpenVG paths from the font data:

    loadfont(DejaVuSans_glyphPoints,
            DejaVuSans_glyphPointIndices,
            DejaVuSans_glyphInstructions,
            DejaVuSans_glyphInstructionIndices,
            DejaVuSans_glyphInstructionCounts,
            DejaVuSans_glyphAdvances,
            DejaVuSans_characterMap,
            DejaVuSans_glyphCount);

The unloadfont function releases the path information:

    unloadfont(DejaFont.Glyphs, DejaFont.Count);

Note that the location of the font files may differ.  (The current location for Jessie is /usr/share/fonts/truetype/ttf-dejavu)
Use the FONTLIB makefile variable to adjust this location.

# Build and run

*Note that you will need at least 64 Mbytes of GPU RAM:*. You will also need the DejaVu fonts, and the jpeg and freetype libraries.
The indent tool is also useful for code formatting.  Install them via:

    pi@raspberrypi ~ $ sudo apt-get install libjpeg8-dev indent libfreetype6-dev ttf-dejavu-core

Next, build the library and test:

    pi@raspberrypi ~ $ git clone git://github.com/ajstarks/openvg
    pi@raspberrypi ~ $ cd openvg
    pi@raspberrypi ~/openvg $ make
    g++ -I/usr/include/freetype2 fontutil/font2openvg.cpp -o font2openvg -lfreetype
    ./font2openvg /usr/share/fonts/truetype/ttf-dejavu/DejaVuSans.ttf DejaVuSans.inc DejaVuSans
    224 glyphs written
    ./font2openvg /usr/share/fonts/truetype/ttf-dejavu/DejaVuSansMono.ttf DejaVuSansMono.inc DejaVuSansMono
    224 glyphs written
    ./font2openvg /usr/share/fonts/truetype/ttf-dejavu/DejaVuSerif.ttf DejaVuSerif.inc DejaVuSerif
    224 glyphs written
    gcc -O2 -Wall -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads -c libshapes.c
    gcc -O2 -Wall -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads -c oglinit.c
    pi@raspberrypi ~/openvg/client $ cd client
    pi@raspberrypi ~/openvg/client $ make test
    cc -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcos/pthreads -o shapedemo shapedemo.c ../libshapes.o ../oglinit.o -L/opt/vc/lib -lGLESv2 -ljpeg
    ./shapedemo demo 5

The program "shapedemo" exercises a high-level API built on OpenVG found in libshapes.c.

    ./shapedemo                      # show a reference card
    ./shapedemo raspi                # show a self-portrait
    ./shapedemo image                # show four test images
    ./shapedemo astro                # the sun and the earth, to scale
    ./shapedemo text                 # show blocks of text in serif, sans, and mono fonts
    ./shapedemo rand 10              # show 10 random shapes
    ./shapedemo rotate 10 a          # rotated and faded "a"
    ./shapedemo test "hello, world"  # show a test pattern, with "hello, world" at mid-display in sans, serif, and mono.
    ./shapedemo fontsize             # show a range of font sizes (per <https://speakerdeck.com/u/idangazit/p/better-products-through-typography>)
    ./shapedemo demo 10              # run through the demo, pausing 10 seconds between each one; contemplate the awesome.

To install the shapes library as a system-wide shared library

    pi@raspberrypi ~/openvg $ make library
    pi@raspberrypi ~/openvg $ sudo make install

The openvg shapes library can now be used in C code by including shapes.h and fontinfo.h and linking with libshapes.so:

    #include <shapes.h>
    #include <fontinfo.h>

    pi@raspberrypi ~ $ gcc -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads anysource.c -o anysource -lshapes
    pi@raspberrypi ~ $ ./anysource
