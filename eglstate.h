typedef struct {
	uint32_t screen_width;
	uint32_t screen_height;
	// OpenGL|ES objects
	EGLDisplay display;

	EGLSurface surface;
	EGLContext context;
} STATE_T;

extern void init_ogl(STATE_T *);
