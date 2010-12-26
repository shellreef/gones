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

const SOCKET_FILE = "/tmp/leggo.sock"   // TODO: stop using insecure temporary directory

//export GoLeggoExit
func GoLeggoExit() {
    // Note: don't call this GoExit! It will cause dyld to fail to find anything here.
    fmt.Printf("in LeggoExit!\n");
}

/*
func CreateDisplay(width int, height int) (*C.ALLEGRO_DISPLAY) {
    return C.al_create_display(C.int(width), C.int(height))
}*/

// Setup a Unix domain socket to communicate with the Allegro side
// Calling an //export'd Go function from C in Allegro fails with a
// null pointer exception, which is why I'm doing this. Another alternative
// may be using pthread condition variables as described on 
// http://blog.labix.org/2010/12/10/integrating-go-with-c-the-zookeeper-binding-experience
// and https://github.com/0xe2-0x9a-0x9b/Go-SDL/tree/master/sdl/audio/
// and http://bazaar.launchpad.net/%7Eensemble/gozk/trunk/annotate/head%3A/helpers.c
func LeggoServer() {
    os.Remove(SOCKET_FILE)

    listener, err := net.Listen("unix", SOCKET_FILE)
    if err != nil {
        fmt.Printf("LeggoServer: failed to listen on socket %s: %s\n", SOCKET_FILE, err)
        return
    }
    fmt.Printf("listener = %s\n", listener)

    for {
        fmt.Printf("LeggoServer: waiting for connection\n")
        conn, err := listener.Accept()
        if err != nil {
            fmt.Printf("LeggoServer: Accept() failed: %s", err)
            return
        }

        fmt.Printf("Established connection: %s\n", conn)
   
        for {
            var buffer [1024]byte
            bytesRead, err := conn.Read(buffer[:])
            if err != nil {
                fmt.Printf("Error reading from client: %s\n", err)
                continue
            }
            fmt.Printf(">>> Read %d bytes from client: %s\n", bytesRead, buffer)

            // TODO: dispatch events to Go channel
            if buffer[0] == 'x' {
                fmt.Printf("Exiting\n")
                os.Exit(0)
            } else if buffer[0] == ' ' {
                fmt.Printf("TODO: do something\n")
            }
        }
    }
}

// Get things going
//export LeggoMain
func LeggoMain() {
    // TODO: mmap anonymous for pixel communication

    go LeggoServer()

    fmt.Printf("LeggoMain: about to call al_run_main_wrapper\n")
    C.al_run_main_wrapper()
}

func init() {
    runtime.GOMAXPROCS(2)
}


