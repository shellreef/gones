// Created:20101212
// By Jeff Connelly

// Test wrapper

package main

import (
    "fmt"
    "unsafe"
    "runtime"
    "rand"

    "leggo")

func something(screen unsafe.Pointer) {
    fmt.Printf("... doing something\n")

    for {
        offset := int(rand.Uint32()) % (256 * 240)
        if rand.Uint32() % 2 == 0 {
            leggo.WriteByte(screen, offset, 255)
        } else {
            leggo.WriteByte(screen, offset, 0)
        }

        /*
        for x := 0; x < 256; x += 1 {
            for y := 0; y < 240; y += 1 {
                offset := 4 * (x + y*256)
                if offset > 256*240*4 {
                    panic(fmt.Sprintf("out of range for (%d,%d): %d > %d",
                                x,y,offset,256*240*4))
                }
                leggo.WriteByte(screen, offset, byte(rand.Uint32()))
                //leggo.WriteByte(screen, offset+1, 0)
                //leggo.WriteByte(screen, offset+2, 0)
                //leggo.WriteByte(screen, offset+3, 0)
            }
        }*/
    }
}

func start(screen unsafe.Pointer) {
    fmt.Printf("starting\n")
    fmt.Printf("with %x\n", screen)

    go something(screen)
}

func main() {
    // never returns
    leggo.LeggoMain(start)

    fmt.Printf("returned?!\n")
}

func init() {
    runtime.GOMAXPROCS(2)
}
