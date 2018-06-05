#include <stdio.h>
#include <stdlib.h>
#include <termios.h>
#include <jpeglib.h>
#include "EGL/egl.h"
#include "VG/openvg.h"
#include "VG/vgu.h"
#include "./../lib/DejaVuSans.inc"

// createImageFromJpeg decompresses a JPEG image to the standard image format
// source: https://github.com/ileben/ShivaVG/blob/master/examples/test_image.c
VGImage createImageFromJpeg(const char *filename) {
    FILE *infile;
    struct jpeg_decompress_struct jdc;
    struct jpeg_error_mgr jerr;
    JSAMPARRAY buffer;
    unsigned int bstride;
    unsigned int bbpp;

    VGImage img;
    VGubyte *data;
    unsigned int width;
    unsigned int height;
    unsigned int dstride;
    unsigned int dbpp;

    VGubyte *brow;
    VGubyte *drow;
    unsigned int x;
    unsigned int lilEndianTest = 1;
    VGImageFormat rgbaFormat;

    // Check for endianness
    if (((unsigned char *)&lilEndianTest)[0] == 1)
        rgbaFormat = VG_sABGR_8888;
    else
        rgbaFormat = VG_sRGBA_8888;

    // Try to open image file
    infile = fopen(filename, "rb");
    if (infile == NULL) {
        printf("Failed opening '%s' for reading!\n", filename);
        return VG_INVALID_HANDLE;
    }
    // Setup default error handling
    jdc.err = jpeg_std_error(&jerr);
    jpeg_create_decompress(&jdc);

    // Set input file
    jpeg_stdio_src(&jdc, infile);

    // Read header and start
    jpeg_read_header(&jdc, TRUE);
    jpeg_start_decompress(&jdc);
    width = jdc.output_width;
    height = jdc.output_height;

    // Allocate buffer using jpeg allocator
    bbpp = jdc.output_components;
    bstride = width * bbpp;
    buffer = (*jdc.mem->alloc_sarray)
        ((j_common_ptr) & jdc, JPOOL_IMAGE, bstride, 1);

    // Allocate image data buffer
    dbpp = 4;
    dstride = width * dbpp;
    data = (VGubyte *) malloc(dstride * height);

    // Iterate until all scanlines processed
    while (jdc.output_scanline < height) {

        // Read scanline into buffer
        jpeg_read_scanlines(&jdc, buffer, 1);
        drow = data + (height - jdc.output_scanline) * dstride;
        brow = buffer[0];
        // Expand to RGBA
        for (x = 0; x < width; ++x, drow += dbpp, brow += bbpp) {
            switch (bbpp) {
            case 4:
                drow[0] = brow[0];
                drow[1] = brow[1];
                drow[2] = brow[2];
                drow[3] = brow[3];
                break;
            case 3:
                drow[0] = brow[0];
                drow[1] = brow[1];
                drow[2] = brow[2];
                drow[3] = 255;
                break;
            }
        }
    }

    // Create VG image
    img = vgCreateImage(rgbaFormat, width, height, VG_IMAGE_QUALITY_BETTER);
    vgImageSubData(img, data, dstride, rgbaFormat, 0, 0, width, height);

    // Cleanup
    jpeg_destroy_decompress(&jdc);
    fclose(infile);
    free(data);

    return img;
}

// makeimage makes an image from a raw raster of red, green, blue, alpha values
void makeimage(VGfloat x, VGfloat y, int w, int h, VGubyte * data) {
    unsigned int dstride = w * 4;
    VGImageFormat rgbaFormat = VG_sABGR_8888;
    VGImage img = vgCreateImage(rgbaFormat, w, h, VG_IMAGE_QUALITY_BETTER);
    vgImageSubData(img, (void *)data, dstride, rgbaFormat, 0, 0, w, h);
    vgSetPixels(x, y, img, 0, 0, w, h);
    vgDestroyImage(img);
}
