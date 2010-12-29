// Created:20101214
// By Jeff Connelly

// Leggo my Allegro: an Allegro 5 wrapper for Go

// Turn off magic re#defining of main(), so we can wrap it manually
#define ALLEGRO_NO_MAGIC_MAIN
#include <allegro5/allegro5.h>

#include <stdio.h>
#include <stdlib.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <unistd.h>
#include <stddef.h>
#include <string.h>

#include "leggo.h"
#include "_cgo_export.h"

#define SOCKET_FILENAME "/tmp/leggo.sock"

static unsigned char *screen_map;
static size_t screen_map_size;
static float seconds_per_frame;

// Set pointer to mmap'd memory for displaying screen from Go
void set_screen_map(void *p, size_t s) {
    screen_map = p;
    screen_map_size = s;

    // initialize to a pleasant default gray
    memset(screen_map, 120, screen_map_size);
}

void write_byte(off_t offset, uint8_t value) {
    *(screen_map + offset) = value;
}

float get_seconds_per_frame() {
    return seconds_per_frame;
}

// Update the screen with the contents of screen_map
// TODO: is direct access possible, avoiding copying?
void refresh(ALLEGRO_DISPLAY *display) {
    ALLEGRO_LOCKED_REGION *locked;
    ALLEGRO_BITMAP *bitmap;

    static float last_time = 0;

    seconds_per_frame = al_get_time() - last_time;

    //printf("refresh: %.8f s/frame\n", seconds_per_frame);
    // TODO: calculate this.. for some reason, it is horribly broken,
    // 1.0f/0.013222 is printing 19768790745088.000000 on release.2010-12-22
    //printf("\tfps: %f\n", 1.0f / seconds_per_frame);

    bitmap = al_get_backbuffer(display);
    locked = al_lock_bitmap(bitmap, ALLEGRO_PIXEL_FORMAT_ABGR_8888_LE, ALLEGRO_LOCK_READWRITE);

    int x, y;

    // Copy the screen map, which is a WxHx4 bitmap beginning at 0x0, into the 
    // locked region buffer, which can have a different pitch, meaning there is 
    // extra uncopied space. The pitch can even be negative.
    for (x = 0; x < RESOLUTION_W; x += 1) {
        for (y = 0; y < RESOLUTION_H; y += 1) {
            uint8_t *dst = (uint8_t *)locked->data + (locked->pixel_size * x + locked->pitch * y);
            uint8_t *src = screen_map + 4 * (x + RESOLUTION_W * y);
            // RGBA
            *(dst + 0) = *(src + 0);
            *(dst + 1) = *(src + 1);
            *(dst + 2) = *(src + 2);
            *(dst + 3) = *(src + 3);
        }
    }
    // locked->data always points to first visual scanline
    // but, if pitch < 0, then it is the last scanline in memory, and
    // we have to copy into the memory *before* locked->data
    //TODO: memcpy(locked->data - screen_map_size + (-locked->pitch), screen_map, screen_map_size);
    // This is not a simple memcpy, but getting it right may be faster (TODO: for performance)
    
    al_unlock_bitmap(bitmap);
    al_flip_display();

    last_time = al_get_time();
}

// Connect to the Unix domain socket used for communication with Go
int connect_socket() {
    struct sockaddr_un address;
    size_t length;
    int fd;

    //printf("connect_socket(): connecting...\n");

    fd = socket(AF_UNIX, SOCK_STREAM, 0);
    if (fd < 0) {
        perror("socket() failed");
        exit(EXIT_FAILURE);
    }

    address.sun_family = AF_UNIX;
    strncpy(address.sun_path, SOCKET_FILENAME, sizeof(address.sun_path));
    length = sizeof(address.sun_family) + strlen(address.sun_path) + 1;
    
    if (connect(fd, (struct sockaddr *)&address, length) == -1) {
        perror("connect() failed");
        exit(EXIT_FAILURE);
    }

    return fd;
}


// The "user-defined" callback function given to al_run_main()
// Allegro will call this as part of its UI wrapper, so when it 
// returns, the program will exit. All code is invoked from here.
int leggo_user_main(int argc, char **argv) {
    ALLEGRO_DISPLAY *display;
    ALLEGRO_TIMER *timer;
    int fd;

    //printf("leggo_user_main: starting up\n");
    
    fd = connect_socket();

    
    if (!al_init()) {
        fprintf(stderr, "failed to initialize Allegro\n");
        exit(EXIT_FAILURE);
    }

    display = al_create_display(RESOLUTION_W, RESOLUTION_H);
    if (!display) {
        fprintf(stderr, "failed to create display\n");
        exit(EXIT_FAILURE);
    }

    if (!al_install_keyboard()) {
        fprintf(stderr, "failed to install keyboard\n");
        exit(EXIT_FAILURE);
    }

    // Use monitor refresh rate if applicable; otherwise default to 60 Hz
    int hz = al_get_display_refresh_rate(display);
    float speed_secs;
   
    if (hz == 0) { 
        speed_secs = 1.0f / 60;
    } else {
        speed_secs = 1.0f / hz;
    }
    // This is really bizarre. Removing this call corrupts speed_secs.
    printf("Leggo my Allegro: before call: (%f) %g %d - %.2x%.2x%.2x%.2x %.2x%.2x%.2x%.2x\n", speed_secs, speed_secs, speed_secs > 0, 
            *(unsigned char*)(&speed_secs + 0),
            *(unsigned char*)(&speed_secs + 1),
            *(unsigned char*)(&speed_secs + 2),
            *(unsigned char*)(&speed_secs + 3),
            *(unsigned char*)(&speed_secs + 4),
            *(unsigned char*)(&speed_secs + 5),
            *(unsigned char*)(&speed_secs + 6),
            *(unsigned char*)(&speed_secs + 7)
            );
    timer = al_create_timer(speed_secs);
    al_start_timer(timer);

    // Draw a green background as a test
    al_set_target_bitmap(al_get_backbuffer(display));
    al_clear_to_color(al_map_rgb(128, 255, 128));
    al_flip_display();

    // TODO: al_lock_bitmap() http://alleg.sourceforge.net/a5docs/5.0.0/graphics.html#al_lock_bitmap
    // to get at raw pixel data in mmap'd region, then unlock to update

    // Main event loop
    ALLEGRO_EVENT_QUEUE *queue;
    ALLEGRO_EVENT event;

    queue = al_create_event_queue();
    al_register_event_source(queue, (ALLEGRO_EVENT_SOURCE *)al_get_keyboard_event_source());
    al_register_event_source(queue, (ALLEGRO_EVENT_SOURCE *)timer);
    while(1) {
        al_wait_for_event(queue, &event);

        if (event.type == ALLEGRO_EVENT_TIMER) {
            refresh(display);
            // TODO: pass to Go too? only if not our update timer?
        } else {
            char buf[2];

            buf[0] = event.type;
            // TODO: more events; this is only useful for ALLEGRO_EVENT_KEY_{UP,DOWN}
            buf[1] = event.keyboard.keycode;

            // Received by LeggoServer()
            send(fd, buf, 2, 0);
        }
    }

    return 0;
}

// Wrap calling Allegro's al_run_main(), with our own leggo_main. 
// This is a C function so Go can call it using cgo, in leggo.Main().
void al_run_main_wrapper() {
    //printf("al_run_main_wrapper(): about to call al_run_main()\n");
    al_run_main(0, 0, leggo_user_main);
}

