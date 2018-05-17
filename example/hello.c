#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include "VG/openvg.h"
#include "VG/vgu.h"
#include "fontinfo.h"
#include "shapes.h"

int main() {
	const char msg[] = { 'H', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd', 0 };
	int width, height;
	char s[3];

	init(&width, &height);
	Start(width, height);
    {
        Background(0, 0, 0);
        Fill(44, 77, 232, 1);
        Circle(width / 2, 0, width);
        Fill(255, 255, 255, 1);
        TextMid(width / 2, (height * 0.7), msg, SerifTypeface, width / 15);
    }
	End();

    // We can use `fgets` to wait for user input before ending the program. This
    // is a neat way of keeping the rendered assets on screen until input is
    // detected.
	fgets(s, 2, stdin);
	finish();
	exit(0);
}
