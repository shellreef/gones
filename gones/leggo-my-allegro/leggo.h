
extern int leggo_user_main(int, char **);
extern void al_run_main_wrapper();
extern void set_screen_map(void *p, size_t s);
extern void write_byte(off_t offset, uint8_t value);
float get_seconds_per_frame();

// TODO: custom
#define RESOLUTION_W 256
#define RESOLUTION_H 240
#define SCALE_W 2
#define SCALE_Y 2
