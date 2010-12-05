// Created:20101126
// By Jeff Connelly

//

package main

import (
    "flag"
    "os"
    "fmt"
    //"strconv"

    "./nesfile"
    "./cpu6502"
)

func main() {
    var start string

    flag.StringVar(&start, "start", "RESET", "Initial value for Program Counter in hex, or reset vector")
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage: %s [flags] game.nes\n", os.Args[0])
        flag.PrintDefaults()
        os.Exit(-1)
    }

    flag.Parse()

    if flag.NArg() != 1 { 
        flag.Usage()
    }
    filename := flag.Arg(0)

    cart := nesfile.Open(filename)
    cpu := new(cpu6502.CPU)
    cpu.Load(cart)
    cpu.PowerUp()

    if start != "RESET" {
        // By default, execution starts at reset vector, but can override
        var startInt uint16
        items, err := fmt.Sscanf(start, "%x", &startInt)
        if items != 1 {
            fmt.Fprintf(os.Stderr, "Unable to read start address: %s: %s\n", start, err)
            os.Exit(-1)
        }

        cpu.PC = startInt
    }

    cpu.Run()
}

