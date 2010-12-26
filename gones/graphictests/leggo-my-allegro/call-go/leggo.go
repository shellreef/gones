// Created:20101212
// By Jeff Connelly

// Leggo my Allegro: an Allegro 5 wrapper for Go

package leggo

import "fmt"

// extern int user_main(int argc, char **argv);
// extern void wrapper();
// #define ALLEGRO_NO_MAGIC_MAIN
// #include <allegro5/allegro.h>
import "C"

//export GoFoo
func GoFoo() {
    fmt.Printf("*** It worked! In GoFoo()\n")
}

// Get things going
//export LeggoMain
func LeggoMain() {
    fmt.Printf("LeggoMain: about to call wrapper\n")
    C.wrapper()
}
