// Created:20101212
// By Jeff Connelly

// Leggo my Allegro: an Allegro 5 wrapper for Go

package leggo

// #include "leggo.h"
// #define ALLEGRO_NO_MAGIC_MAIN
// #include <allegro5/allegro.h>
import "C"

// TODO: export something real
//export Init
func Init() (bool) {
    x := C.al_init_wrapper()

    return x != 0
}

func CreateDisplay(width int, height int) (*C.ALLEGRO_DISPLAY) {
    return C.al_create_display(C.int(width), C.int(height))
}

func AllegroMain() {
    // TODO: pass channel
    C.al_run_main_wrapper()
}

