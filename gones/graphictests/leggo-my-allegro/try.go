// Created:20101212
// By Jeff Connelly

// Test wrapper

package main

import (
    "fmt"
    "unsafe"
    "runtime"
    "rand"
    "time"

    "leggo")

func something(screen unsafe.Pointer) {
    fmt.Printf("... doing something\n")

    for {
        /*
        offset := int(rand.Uint32()) % (256 * 240)
        if rand.Uint32() % 2 == 0 {
            leggo.WriteByte(screen, offset, 255)
        } else {
            leggo.WriteByte(screen, offset, 0)
        }*/

        w, h := 256, 240

        for y := 0; y < h; y += 1 {
            for x := 0; x < w; x += 1 {
                offset := 4*(x + y*w)
                if offset > w*h*4 {
                    panic(fmt.Sprintf("out of range for (%d,%d): %d > %d",
                                x,y,offset,w*h*4))
                }
                leggo.WriteByte(offset + 0, 200)
                leggo.WriteByte(offset + 1, byte(rand.Uint32()))
                leggo.WriteByte(offset + 2, 0)
                leggo.WriteByte(offset + 3, 0)
                time.Sleep(100000)

                //leggo.WriteByte(screen, offset+1, 0)
                //leggo.WriteByte(screen, offset+2, 0)
                //leggo.WriteByte(screen, offset+3, 0)
            }
        }
        break
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
