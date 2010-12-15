// Created:20101212
// By Jeff Connelly

// Test wrapper

package main

import (
    "fmt"

    "leggo")

func main() {
    // never returns
    leggo.AllegroMain()

    fmt.Printf("about to call init\n")
    fmt.Printf("init = %t\n", leggo.Init())

    //allego5.CreateDisplay(640, 480)
}
