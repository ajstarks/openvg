typedef struct {
	const int *Points; 
	const int *PointIndices;
	const short *CharacterMap;
	const unsigned char *Instructions; 
	const int *InstructionIndices;
	const int *InstructionCounts;
	const int *GlyphAdvances;
	int Count; 
	VGPath Glyphs[256];
} Fontinfo;

Fontinfo SansTypeface, SerifTypeface;
