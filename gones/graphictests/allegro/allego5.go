// Created:20101212
// By Jeff Connelly

// Allego5: an Allegro 5 wrapper for Go

package allego5

import "unsafe"

/*

#include <allegro5/allegro5.h>
#include <stdio.h>

// We have to wrap this since it is a C macro
bool al_init_wrapper() {
    return al_init();
}


*/
import "C"

func Init() (bool) {
    x := C.al_init_wrapper()

    return bool(x)
}

func CreateDisplay(width int, height int) (*C.ALLEGRO_DISPLAY) {
    return C.al_create_display(C.int(width), C.int(height))
}

func AllegroMain(f func(C.int, *[0]C.char) (int)) {
    C.al_run_main(C.int(0), nil, unsafe.Pointer(f))
}

/*
// Created:20101212
// By Jeff Connelly

// Test for using Allegro 5
// Based on http://wiki.allegro.cc/pub/3/3a/Tut2Hello2.c

#include <stdio.h>
#include <allegro5/allegro5.h>

int main(int argc, char **argv) {
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

    al_set_target_bitmap(al_get_backbuffer(display));
    al_clear_to_color(al_map_rgb(128, 255, 128));
    al_flip_display();


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
*/
