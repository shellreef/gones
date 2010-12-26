// Created:20101212
// By Jeff Connelly

// Leggo my Allegro: an Allegro 5 wrapper for Go

package leggo

import "fmt"

// #include "leggo.h"
// #define ALLEGRO_NO_MAGIC_MAIN
// #include <allegro5/allegro.h>
import "C"

//export GoFoo
func GoFoo() {
    fmt.Printf("in GoFoo()!\n")
}

// Get things going
//export LeggoMain
func LeggoMain() {
    fmt.Printf("LeggoMain: about to call al_run_main_wrapper\n")
    C.al_run_main_wrapper()
}
