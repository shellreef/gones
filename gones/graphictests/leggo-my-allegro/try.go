// Created:20101212
// By Jeff Connelly

// Test wrapper

package main

import (
    "fmt"
    "runtime"

    "leggo")

func main() {
    runtime.GOMAXPROCS(2)

    // never returns
    leggo.LeggoMain()

    fmt.Printf("returned?!\n")
}
