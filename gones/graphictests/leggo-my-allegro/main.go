// Created:20101212
// By Jeff Connelly

// Test wrapper

package main

import (
    "fmt"
    "runtime"
    "rand"
    "os"
    //"time"

    "leggo")

var mode byte = 0

func something() {
    fmt.Printf("... doing something\n")

    for {
        w, h := 256, 240

        for y := 0; y < h; y += 1 {
            for x := 0; x < w; x += 1 {
                r, g, b := byte(0), byte(0), byte(0)
                switch mode % 3 {
                case 0: r = byte(rand.Uint32()); g = r; b = r
                case 1: r = byte(rand.Uint32()); g = byte(rand.Uint32()); b = byte(rand.Uint32())
                case 2: r = byte(x*y)
                }

                leggo.WritePixel(x, y, r,g,b,0)
            }
        }
    }
}

func start() {
    fmt.Printf("starting\n")

    go something()
}

func event(kind int, code int) {
    fmt.Printf("got event: %d,%d\n", kind, code);

    if kind == leggo.EVENT_KEY_DOWN {
        switch code {
        case leggo.KEY_ESCAPE:
            fmt.Printf("Exiting\n")
            os.Exit(0)
        case leggo.KEY_DOWN:
            fmt.Printf("mode = %d\n", mode)
            mode += 1
        case leggo.KEY_UP:
            fmt.Printf("mode = %d\n", mode)
            mode -= 1
        case leggo.KEY_SPACE:
            fmt.Printf("FPS = %f\n", leggo.FPS())
        }
    }
}

func main() {
    // never returns
    leggo.LeggoMain(start, event)

    fmt.Printf("returned?!\n")
}

func init() {
    runtime.GOMAXPROCS(2)
}
