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
    "./cpu6502"
)

func main() {
    flag.Parse()

    fmt.Printf("Hello, world\n")

    cart := nesfile.Open(flag.Arg(0))
    cpu := new(cpu6502.CPU)
    cpu.Load(cart)

    buffer := bytes.NewBuffer(cpu.Memory[0x8000:])
    for {
         instr, err := dis6502.ReadInstruction(buffer)
         if err != nil {
             break
         }
         // TODO: show address, but how to get file pointer on a Buffer? no Tell()
         fmt.Printf("%s\n", instr)
    }
}

