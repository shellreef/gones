// Created:20101212
// By Jeff Connelly

// Test yasdl

package main

import "fmt"

import "yasdl"

func main() {
    x := yasdl.Init()
    fmt.Printf("%d\n", x)
}

