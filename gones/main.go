// Created:20101126
// By Jeff Connelly

//

package main

import (
    "flag"

    "./nesfile"
    "./cpu6502"
)

func main() {
    flag.Parse()

    cart := nesfile.Open(flag.Arg(0))
    cpu := new(cpu6502.CPU)
    cpu.Load(cart)
    cpu.Run()
}

