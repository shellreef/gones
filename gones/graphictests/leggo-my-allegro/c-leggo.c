// Created:20101214
// By Jeff Connelly

// Leggo my Allegro: an Allegro 5 wrapper for Go

// Turn off magic re#defining of main(), so we can wrap it manually
#define ALLEGRO_NO_MAGIC_MAIN
#include <allegro5/allegro5.h>

#include <stdio.h>
#include <stdlib.h>

#include "leggo.h"
#include "_cgo_export.h"

// The "user-defined" callback function given to al_run_main()
// Allegro will call this as part of its UI wrapper, so when it 
// returns, the program will exit. All code is invoked from here.
int leggo_user_main(int argc, char **argv) {
    printf("leggo_user_main: starting up\n");

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

    // Main event loop
    ALLEGRO_EVENT_QUEUE *queue;
    ALLEGRO_EVENT event;

    queue = al_create_event_queue();
    al_register_event_source(queue, (ALLEGRO_EVENT_SOURCE *)al_get_keyboard_event_source());
    while(1) {
        al_wait_for_event(queue, &event);

        if (event.type == ALLEGRO_EVENT_KEY_DOWN && event.keyboard.keycode == ALLEGRO_KEY_ESCAPE) {
            printf("leggo_user_main: calling LeggoExit()\n");
            LeggoExit();
            break;
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

