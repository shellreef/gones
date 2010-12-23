// Created:20101212
// By Jeff Connelly

// Leggo my Allegro: an Allegro 5 wrapper for Go

package leggo

import ("fmt"
    "net"
    "runtime")

// #include "leggo.h"
// #define ALLEGRO_NO_MAGIC_MAIN
// #include <allegro5/allegro.h>
import "C"

//export GoLeggoExit
func GoLeggoExit() {
    // Note: don't call this GoExit! It will cause dyld to fail to find anything here.
    fmt.Printf("in LeggoExit!\n");
}

/*
func CreateDisplay(width int, height int) (*C.ALLEGRO_DISPLAY) {
    return C.al_create_display(C.int(width), C.int(height))
}*/

// Get things going
//export LeggoMain
func LeggoMain() {
    conn, err := net.Dial("unix", "", "sock")
    fmt.Printf("conn = %s\n", conn)
    fmt.Printf("err = %s\n", err)

    fmt.Printf("LeggoMain: about to call al_run_main_wrapper\n")
    C.al_run_main_wrapper()
}

func init() {
    runtime.GOMAXPROCS(2)
}


