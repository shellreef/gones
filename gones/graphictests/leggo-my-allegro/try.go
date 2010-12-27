// Created:20101212
// By Jeff Connelly

// Test wrapper

package main

import (
    "fmt"

    //"unsafe"
    "leggo")

func something() {
    fmt.Printf("... doing something\n")
}

func start() {
    fmt.Printf("starting\n")

    go something()
}

func main() {
    //screen := leggo.LeggoSetup()
    //go something(screen)
    
    _ = leggo.LeggoSetup()

    //go something()

    // never returns
    leggo.LeggoMain(start)

    fmt.Printf("returned?!\n")
}
