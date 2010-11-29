// Created:20101126
// By Jeff Connelly

//

package main

import (
    "fmt"
    "flag"

    "./nesfile"
    "./dis6502"
)

func main() {
    flag.Parse()

    fmt.Printf("Hello, world\n")

    nesfile.Load(flag.Arg(0))
    var op dis6502.Opcode = dis6502.BRK
    fmt.Printf("op=%s\n", op)
}

