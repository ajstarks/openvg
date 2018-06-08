#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include "VG/openvg.h"
#include "VG/vgu.h"
#include "./../src/fontinfo.h"
#include "./../src/libshapes.h"

int main() {
    const char msg[] = { 'H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd', 0 };
    int width, height;
    char s[3];

    /* We have to initialise the OpenVG canvas and renderer by using
     * `init(&w, &h)`.
     */
    evgInit(&width, &height);
    evgBegin(width, height);
    {
        evgBackground(0, 0, 0);
        evgFill(44, 77, 232, 1);
        evgCircle(width / 2, 0, width);
        evgFill(255, 255, 255, 1);
        evgTextMid(width / 2,
                height * 0.7,
                msg,
                SansTypeface,
                width / 15);
    }
    evgEnd();

    // We can use `fgets` to wait for user input before ending the program. This
    // is a neat way of keeping the rendered assets on screen until the user
    // presses a key.
    fgets(s, 2, stdin);
    evgFinish();
    exit(0);
}
