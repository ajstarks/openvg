#ifndef OPENVG_FONTINFO_H
#define OPENVG_FONTINFO_H

#if defined(__cplusplus)
extern "C" {
#endif
	typedef struct {
		const short *CharacterMap;
		const int *GlyphAdvances;
		int Count;
		VGPath Glyphs[500];
	} Fontinfo;

	extern Fontinfo SansTypeface, SerifTypeface, MonoTypeface;

#if defined(__cplusplus)
}
#endif				// OPENVG_FONTINFO_H
#endif
