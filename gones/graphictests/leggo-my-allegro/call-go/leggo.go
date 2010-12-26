// Created:20101212
// By Jeff Connelly

// Leggo my Allegro: an Allegro 5 wrapper for Go

package leggo

import ("fmt"
    "runtime")

// #include "leggo.h"
// #define ALLEGRO_NO_MAGIC_MAIN
// #include <allegro5/allegro.h>
import "C"

const SOCKET_FILE = "/tmp/leggo.sock"   // TODO: stop using insecure temporary directory

//export GoLeggoExit
func GoLeggoExit() {
    // Note: don't call this GoExit! It will cause dyld to fail to find anything here.
    fmt.Printf("in LeggoExit!\n");
}

// Get things going
//export LeggoMain
func LeggoMain() {
    fmt.Printf("LeggoMain: about to call al_run_main_wrapper\n")
    C.al_run_main_wrapper()
}

func init() {
    runtime.GOMAXPROCS(2)
}


