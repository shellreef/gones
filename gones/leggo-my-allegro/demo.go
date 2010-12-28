// Created:20101212
// By Jeff Connelly

// Test wrapper

package main

import (
    "fmt"
    "rand"
    "os"
    "time"

    "leggo")

var mode byte = 0

// Continuously update the screen with something interesting
func render() {
    for {
        w, h := leggo.Dimensions()

        for y := 0; y < h; y += 1 {
            for x := 0; x < w; x += 1 {
                var r, g, b byte
                // Some random things to look at
                switch mode % 5 {
                case 0: r = byte(rand.Uint32()); g = r; b = r
                case 1: r = byte(rand.Uint32()); g = byte(rand.Uint32()); b = byte(rand.Uint32())
                case 2: r = byte(x*y)
                case 3: r = byte(x*y - int(rand.Uint32()) % 50)
                case 4: g = byte(rand.Uint32()); time.Sleep(1000)
                }

                leggo.WritePixel(x, y, r,g,b,0)
            }
        }
    }
}

// Handle events
func process(ch chan leggo.Event) {
    for {
        e := <-ch

        kind, code := e.Type, e.Keycode

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
}

func main() {
    // never returns
    leggo.LeggoMain(render, process)

    fmt.Printf("returned?!\n")
}


