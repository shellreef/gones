// Created:20101126
// By Jeff Connelly

//

package main

import (
    "flag"
    "os"
    "fmt"
    "strconv"

    "./nesfile"
    "./cpu6502"
)

func main() {
    var start string

    flag.StringVar(&start, "start", "RESET", "Initial value for Program Counter")
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
        startInt, err := strconv.Atoi(start)
        if err == nil {
            cpu.PC = uint16(startInt)
        } else {
            // TODO: need to accept hex, 0xc000.. fmt.Sscanf?
            fmt.Fprintf(os.Stderr, "Unable to read start address: %s\n", start)
            os.Exit(-1)
        }
    }

    cpu.Run()
}

