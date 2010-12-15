// Created:20101214
// By Jeff Connelly

// Leggo my Allegro: an Allegro 5 wrapper for Go

#include <allegro5/allegro5.h>
#include <stdio.h>
#include <stdlib.h>

#include "_cgo_export.h"

// We have to wrap this since it is a C macro
bool al_init_wrapper() {
    return al_init();
}

int user_main(int argc, char **argv) {
    ALLEGRO_DISPLAY *display;
    
    al_init();

    display = al_create_display(640, 480);
    if (!display) {
        fprintf(stderr, "failed to create display\n");
        exit(EXIT_FAILURE);
    }

    if (!al_install_keyboard()) {
        fprintf(stderr, "failed to install keyboard\n");
        exit(EXIT_FAILURE);
    }

    al_set_target_bitmap(al_get_backbuffer(display));
    al_clear_to_color(al_map_rgb(128, 255, 128));
    al_flip_display();


    ALLEGRO_EVENT_QUEUE *queue;
    ALLEGRO_EVENT event;

    queue = al_create_event_queue();
    al_register_event_source(queue, (ALLEGRO_EVENT_SOURCE *)al_get_keyboard_event_source());
    while(1) {
        al_wait_for_event(queue, &event);

        if (event.type == ALLEGRO_EVENT_KEY_DOWN && event.keyboard.keycode == ALLEGRO_KEY_ESCAPE) {
            break;
        }
    }

    return 0;
}

void al_run_main_wrapper() {
    al_run_main(0, 0, user_main);
}

