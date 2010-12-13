// Created:20101212
// By Jeff Connelly

// Test for using Allegro 5
// Based on http://wiki.allegro.cc/pub/3/3a/Tut2Hello2.c

#include <stdio.h>
#include <allegro5/allegro5.h>

int main(int argc, char *argv[]) {
    ALLEGRO_DISPLAY *display;
    ALLEGRO_EVENT_QUEUE *queue;
    ALLEGRO_EVENT event;
    ALLEGRO_COLOR color;

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

    color = al_map_rgba_f(0.5, 1.0, 0.5, 1.0);
    al_clear_to_color(color);


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
