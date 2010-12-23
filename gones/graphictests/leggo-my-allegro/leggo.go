// Created:20101212
// By Jeff Connelly

// Leggo my Allegro: an Allegro 5 wrapper for Go

package leggo

import ("fmt"
    "net"
    "os"
    "runtime")

// #include "leggo.h"
// #define ALLEGRO_NO_MAGIC_MAIN
// #include <allegro5/allegro.h>
import "C"

const SOCKET_FILE = "sock"

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
    // Setup a Unix domain socket to communicate with other thread 
    // TODO: use pthread conditions? Some other IPC?
    os.Remove(SOCKET_FILE)

    listener, err := net.Listen("unix", SOCKET_FILE)
    if err != nil {
        fmt.Printf("Failed to listen on socket %s: %s\n", SOCKET_FILE, err)
        return
    }
    fmt.Printf("listener = %s\n", listener)

    fmt.Printf("LeggoMain: about to call al_run_main_wrapper\n")
    C.al_run_main_wrapper()
}

func init() {
    runtime.GOMAXPROCS(2)
}


