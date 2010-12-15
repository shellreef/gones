// Created:20101212
// By Jeff Connelly

// Test wrapper

package main

import (
    "fmt"

    "leggo")

func main() {
    // never returns
    leggo.Main()

    fmt.Printf("returned?!\n")
}
