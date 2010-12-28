// Created:20101212
// By Jeff Connelly

// Test wrapper

package main

import (
    "fmt"
    "unsafe"
    "runtime"
    "rand"
    //"time"

    "leggo")

func something(screen unsafe.Pointer) {
    fmt.Printf("... doing something\n")

    for {
        w, h := 256, 240

        for y := 0; y < h; y += 1 {
            for x := 0; x < w; x += 1 {
                offset := 4*(x + y*w)
                if offset > w*h*4 {
                    panic(fmt.Sprintf("out of range for (%d,%d): %d > %d",
                                x,y,offset,w*h*4))
                }
                value := byte(rand.Uint32())
                leggo.WriteByte(offset + 0, value)
                leggo.WriteByte(offset + 1, value)
                leggo.WriteByte(offset + 2, value)
                leggo.WriteByte(offset + 3, 0)
            }
        }
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
