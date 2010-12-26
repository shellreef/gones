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

#include "leggo.h"
#include "_cgo_export.h"

#define SOCKET_FILENAME "/tmp/leggo.sock"

// Connect to the Unix domain socket used for communication with Go
int connect_socket() {
    struct sockaddr_un address;
    size_t length;
    int fd;

    printf("connect_socket(): connecting...\n");

    fd = socket(AF_UNIX, SOCK_STREAM, 0);
    if (fd < 0) {
        perror("socket() failed");
        exit(EXIT_FAILURE);
    }

    address.sun_family = AF_UNIX;
    strlcpy(address.sun_path, SOCKET_FILENAME, sizeof(address.sun_path));
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
    int fd;

    printf("leggo_user_main: starting up\n");
    
    fd = connect_socket();

    ALLEGRO_DISPLAY *display;
    
    if (!al_init()) {
        fprintf(stderr, "failed to initialize Allegro\n");
        exit(EXIT_FAILURE);
    }

    display = al_create_display(640, 480);
    if (!display) {
        fprintf(stderr, "failed to create display\n");
        exit(EXIT_FAILURE);
    }

    if (!al_install_keyboard()) {
        fprintf(stderr, "failed to install keyboard\n");
        exit(EXIT_FAILURE);
    }

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
    while(1) {
        al_wait_for_event(queue, &event);

        // TODO: proper dispatching
        if (event.type == ALLEGRO_EVENT_KEY_DOWN) {
            if (event.keyboard.keycode == ALLEGRO_KEY_ESCAPE) {
                printf("leggo_user_main: sending event\n");
                send(fd, "x", 1, 0);
            } else if (event.keyboard.keycode == ALLEGRO_KEY_SPACE) {
                send(fd, " ", 1, 0);

                al_set_target_bitmap(al_get_backbuffer(display));
                al_clear_to_color(al_map_rgb(0, 0, 255));
                al_flip_display();

            }
        }
    }

    return 0;
}

// Wrap calling Allegro's al_run_main(), with our own leggo_main. 
// This is a C function so Go can call it using cgo, in leggo.Main().
void al_run_main_wrapper() {
    printf("al_run_main_wrapper(): about to call al_run_main()\n");
    al_run_main(0, 0, leggo_user_main);
}

