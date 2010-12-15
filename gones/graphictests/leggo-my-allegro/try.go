// Created:20101212
// By Jeff Connelly

// Test wrapper

package main

import (
    "fmt"

    "leggo")

func main() {
    // never returns
    leggo.LeggoMain()

    fmt.Printf("returned?!\n")
}
