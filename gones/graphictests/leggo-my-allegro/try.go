// Created:20101212
// By Jeff Connelly

// Test wrapper

package main

import (
    "fmt"
    "unsafe"
    "runtime"

    "leggo")

func something(screen unsafe.Pointer) {
    fmt.Printf("... doing something\n")

    for {
        leggo.Fill(screen)
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
