#include <VG/openvg.h>
#include <VG/vgu.h>
#include "fontinfo.h"
extern void Translate(VGfloat, VGfloat);
extern void Rotate(VGfloat);
extern void Shear(VGfloat, VGfloat);
extern void Scale(VGfloat, VGfloat);
extern void Text(VGfloat, VGfloat, char *, Fontinfo, int);
extern void TextMid(VGfloat, VGfloat, char *, Fontinfo, int);
extern void TextEnd(VGfloat, VGfloat, char *, Fontinfo, int);
extern VGfloat TextWidth(char *, Fontinfo, int);
extern void Cbezier(VGfloat, VGfloat, VGfloat, VGfloat, VGfloat, VGfloat, VGfloat, VGfloat);
extern void Qbezier(VGfloat, VGfloat, VGfloat, VGfloat, VGfloat, VGfloat);
extern void Polygon(VGfloat *, VGfloat *, VGint);
extern void Polyline(VGfloat *, VGfloat *, VGint);
extern void Rect(VGfloat, VGfloat, VGfloat, VGfloat);
extern void Line(VGfloat, VGfloat, VGfloat, VGfloat);
extern void Roundrect(VGfloat, VGfloat, VGfloat, VGfloat, VGfloat, VGfloat);
extern void Ellipse(VGfloat, VGfloat, VGfloat, VGfloat);
extern void Circle(VGfloat, VGfloat, VGfloat);
extern void Arc(VGfloat, VGfloat, VGfloat, VGfloat, VGfloat, VGfloat);
extern void Image(VGfloat, VGfloat, int, int, char *);
extern void Start(int, int);
extern void End();
extern void SaveEnd(char *);
extern void Background(unsigned int, unsigned int, unsigned int);
extern void BackgroundRGB(unsigned int, unsigned int, unsigned int, VGfloat);
extern void init(int *, int *);
extern void finish();
extern void setfill(VGfloat[4]);
extern void setstroke(VGfloat[4]);
extern void StrokeWidth(VGfloat);
extern void Stroke(unsigned int, unsigned int, unsigned int, VGfloat);
extern void Fill(unsigned int, unsigned int, unsigned int, VGfloat);
extern void RGBA(unsigned int, unsigned int, unsigned int, VGfloat, VGfloat[4]);
extern void RGB(unsigned int, unsigned int, unsigned int, VGfloat[4]);
extern void FillLinearGradient(VGfloat, VGfloat, VGfloat, VGfloat, VGfloat *, int);
extern void FillRadialGradient(VGfloat, VGfloat, VGfloat, VGfloat, VGfloat, VGfloat *, int);
extern Fontinfo loadfont(const int *, const int *, const unsigned char *, const int *, const int *, const int *, const short *,
			 int);
extern void unloadfont(VGPath *, int);
extern void makeimage(VGfloat, VGfloat, int, int, VGubyte *);
extern void saveterm();
extern void restoreterm();
extern void rawterm();
