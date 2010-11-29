// Created:20101126
// By Jeff Connelly

//

package main

import (
    "fmt"
    "flag"
    "bytes"

    "./nesfile"
    "./dis6502"
)

func main() {
    flag.Parse()

    fmt.Printf("Hello, world\n")

    cart := nesfile.Load(flag.Arg(0))

    buffer := bytes.NewBuffer(cart.Prg[0])
    for {
         instr, err := dis6502.ReadInstruction(buffer)
         if err != nil {
             break
         }
         fmt.Printf("%s\n", instr)
    }
}

