[header-image]: https://raw.githubusercontent.com/mgthomas99/easy-vg/develop/.github/assets/raspi-spiral.png
[git-repository-url]: https://github.com/mgthomas99/openvg
[license-shield-url]: https://img.shields.io/github/license/mgthomas99/easy-vg.svg?style=flat-square
[license-url]: https://github.com/mgthomas99/openvg/blob/master/LICENSE
[raspberry-pi-shield-url]: https://img.shields.io/badge/raspberry%20pi-compatible-C5194A.svg?style=flat-square
[raspberry-pi-url]: https://www.raspberrypi.org/

# EasyVG

[![One][header-image]][git-repository-url]
[![LICENSE][license-shield-url]][license-url]
[![RaspberryPi][raspberry-pi-shield-url]][raspberry-pi-url]

EasyVG provides an abstraction layer above native OpenVG, for quickly and easily
creating OpenVG contexts and drawing shapes, images, and text.

EasyVG follows standard C and EGL API naming conventions.

EasyVG is fully compatible with the Raspberry Pi.

*This project is a fork of [OpenVG by Anthony Starks <ajstarks>](https://github.com/ajstarks/openvg)*

### Text & Fonts

EasyVG is capable of rendering text.

To use custom TrueType fonts, developers should convert the font into C code use
the [font2openvg](https://github.com/mgthomas99/font2openvg) library, and then
load it the font using `loadfont()`.

*There are plans to support loading a font directly from a `.ttf` file*.

#### Using font2openvg

The `font2openvg` repository contains build instructions for building the
`font2openvg` source. Once the library is built, compile a TrueType font file.
For the below demonstrations, it will be assumed that you are using a source
font file named `DejaVuSans.ttf` and a compiled output named `DejaVuSans.inc`.

Once a font is compiled, it can be included in your code like so:

```c
    #include "DejaVuSans.inc"
    Fontinfo DejaFont

    loadfont(DejaVuSans_glyphPoints,
            DejaVuSans_glyphPointIndices,
            DejaVuSans_glyphInstructions,
            DejaVuSans_glyphInstructionIndices,
            DejaVuSans_glyphInstructionCounts,
            DejaVuSans_glyphAdvances,
            DejaVuSans_characterMap,
            DejaVuSans_glyphCount);

    // Unload the font when done
    unloadfont(DejaFont.Glyphs, DejaFont.Count);
```

## Build and Run Examples

*Note that you will need at least 64 MB of GPU RAM*. You will also need the
DejaVu fonts, and the jpeg and freetype libraries.

```shell
pi@raspberrypi ~ $ sudo apt-get install libfreetype6-dev libjpeg8-dev ttf-dejavu-core
```

Next, build the library:

```shell
pi@raspberrypi ~ $ git clone https://github.com/mgthomas99/easy-vg
pi@raspberrypi ~ $ cd easy-vg
pi@raspberrypi ~/easy-vg $ make
 g++ -I/usr/include/freetype2 fontutil/font2openvg.cpp -o font2openvg -lfreetype
 ./font2openvg /usr/share/fonts/truetype/ttf-dejavu/DejaVuSans.ttf DejaVuSans.inc DejaVuSans
 224 glyphs written
 gcc -O2 -Wall -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads -c libshapes.c
 gcc -O2 -Wall -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads -c oglinit.c
```

Next, build the examples:

```shell
pi@raspberrypi ~/easy-vg/client $ cd example
pi@raspberrypi ~/easy-vg/client $ make hello
 cc -Wall -I/opt/vc/include -I/opt/vc/include/interface/vcms_host/linux -I/opt/vc/include/interface/vcos/pthreads -I.. hello.c -o hello ../build/libshapes.o ../build/oglinit.o -L/opt/vc/lib -lEGL -lGLESv2 -lbcm_host -ljpeg -lpthread
pi@raspberrypi ~/easy-vg/client $ ./hello
```

### Installing as a Global Library

To install the library as a system-wide shared library:

```shell
pi@raspberrypi ~/easy-vg $ make library
pi@raspberrypi ~/easy-vg $ sudo make install
```

The EasyVG shapes library can now be used in C code by including `shapes.h` and
`fontinfo.h`, and linking with `libshapes.so`:

```c
#include <shapes.h>
#include <fontinfo.h>
```

```shell
pi@raspberrypi ~ $ gcc -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads anysource.c -o anysource -lshapes
pi@raspberrypi ~ $ ./anysource
```

## License

See the `LICENSE` file for license information.
