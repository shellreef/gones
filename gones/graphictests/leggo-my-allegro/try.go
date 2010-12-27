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
        offset := 100
        ptr := unsafe.Pointer(uintptr(screen) + uintptr(offset))
        pixel := (*uintptr)(ptr) 
        if rand.Uint32() % 2 == 0 {
            *pixel = 255 
        } else {
            *pixel= 0 
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
