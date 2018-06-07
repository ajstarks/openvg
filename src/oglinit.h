
/**
 *
 */
typedef struct {

    /** Screen width, in pixels. */
    uint32_t screen_width;

    /** Screen height, in pixels. */
    uint32_t screen_height;

    /** Window `x` coordinate. */
    int32_t window_x;

    /** Window `y` coordinate. */
    int32_t window_y;

    /** Window width, in pixels. */
    uint32_t window_width;

    /** Window height, in pixels. */
    uint32_t window_height;

    DISPMANX_ELEMENT_HANDLE_T element;

    // EGL data
    EGLDisplay display;
    EGLSurface surface;
    EGLContext context;
} STATE_T;

extern void oglinit(STATE_T *);
extern void dispmanMoveWindow(STATE_T *, int, int);
extern void dispmanChangeWindowOpacity(STATE_T *, unsigned int);
