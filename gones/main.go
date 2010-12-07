// Created:20101126
// By Jeff Connelly

//

package main

import (
    "os"
    "fmt"
    "strconv"
    "strings"
    
    "readline"  // http://bitbucket.org/taruti/go-readline/

    "nesfile"
    "cpu6502"
)

// Run a command to do something with the unit
// TODO: in shell module
func RunCommand(cpu *cpu6502.CPU, cmd string) {
    switch {
    // TODO: better cmd parsing
    case strings.HasPrefix(cmd, "g"):
        args := strings.Split(cmd, " ", -1)
        if len(args) > 1 {
            startInt, err := strconv.Btoui64(args[1], 0)
            if err != nil {
                fmt.Fprintf(os.Stderr, "unable to read start address: %s: %s\n", args[1], err)
                os.Exit(-1)
            }

            cpu.PC = uint16(startInt)
        }
        cpu.Run()
    case strings.HasPrefix(cmd, "l"):
        args := strings.Split(cmd, " ", 2)
        if len(args) == 2 {
            cart := nesfile.Open(args[1])
            cpu.Load(cart)
            // TODO: verbose load
        } else {
            fmt.Printf("usage: l <filename>\n")
        }
    case strings.HasPrefix(cmd, "i"):
        Shell(cpu)
    default: fmt.Printf("unknown command: %s\n", cmd)
    }
}

// Interactive shell: prompt for commands on stdin and run them
func Shell(cpu *cpu6502.CPU) {
    prompt := "-"
    for {
        line := readline.ReadLine(&prompt)

        if line == nil {
            break
        }

        if *line == "" {
            continue
        }

        RunCommand(cpu, *line)
        readline.AddHistory(*line)
    }
}

func main() {
    cpu := new(cpu6502.CPU)
    cpu.PowerUp()

    for _, cmd := range(os.Args[1:]) {
        RunCommand(cpu, cmd)
    }

    if len(os.Args) < 2 {
        Shell(cpu)
    }
}

