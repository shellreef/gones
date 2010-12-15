// Created:20101212
// By Jeff Connelly

// Leggo my Allegro: an Allegro 5 wrapper for Go

package leggo

// #include "leggo.h"
// #define ALLEGRO_NO_MAGIC_MAIN
// #include <allegro5/allegro.h>
import "C"

/*
func CreateDisplay(width int, height int) (*C.ALLEGRO_DISPLAY) {
    return C.al_create_display(C.int(width), C.int(height))
}*/

// Get things going
//export Main
func Main() {
    C.al_run_main_wrapper()
}

