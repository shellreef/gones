// Created:20101212
// By Jeff Connelly

// Test wrapper

package main

import (
    "fmt"
    "allego5")

func main() {
    fmt.Printf("about to call init\n")
    fmt.Printf("init = %t\n", allego5.Init())

    allego5.CreateDisplay(640, 480)
}
