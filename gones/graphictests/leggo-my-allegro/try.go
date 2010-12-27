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
        base := (*uintptr)(screen)
       
        x := uintptr(screen) + uintptr(1)
        p := unsafe.Pointer(x)
        *(*uintptr)(p) = 128
        fmt.Printf("base = %x, x = %x\n", base, x)
        if rand.Uint32() % 2 == 0 {
            *base = 255 
        } else {
            *base = 0 
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
