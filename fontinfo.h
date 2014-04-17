#ifndef OPENVG_FONTINFO_H
#define OPENVG_FONTINFO_H
typedef struct {
	const short *CharacterMap;
	const int *GlyphAdvances;
	int Count;
	VGPath Glyphs[256];
} Fontinfo;

Fontinfo SansTypeface, SerifTypeface, MonoTypeface;
#endif				// OPENVG_FONTINFO_H
