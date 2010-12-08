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
    "gamegenie"
)

// Run a command to do something with the unit
// TODO: in shell module
func RunCommand(cpu *cpu6502.CPU, cmd string) {
    // If given filename (no command), load and run it
    if strings.HasSuffix(cmd, ".nes") || strings.HasSuffix(cmd, ".NES") {
        _, err := os.Stat(cmd)
        if err == nil {
            // TODO: refactor
            cart := nesfile.Open(cmd)
            cpu.Load(cart)
            cpu.Run()
        } // have to ignore non-existing files, since might be a command
    }

    // TODO: better cmd parsing. real scripting language?
    name := string(cmd[0])
    args := strings.Split(cmd, " ", -1)[1:]
    switch name {
    // go (not to be confused with the programming language)
    case "g": 
        if len(args) > 0 {
            // optional start address to override RESET vector
            startInt, err := strconv.Btoui64(args[0], 0)
            if err != nil {
                fmt.Fprintf(os.Stderr, "unable to read start address: %s: %s\n", args[0], err)
                os.Exit(-1)
            }

            cpu.PC = uint16(startInt)
        }
        cpu.Run()
    // load
    case "l":
        if len(args) > 0 {
            cart := nesfile.Open(args[0])
            cpu.Load(cart)
            // TODO: verbose load
        } else {
            fmt.Printf("usage: l <filename>\n")
        }
    // interactive
    case "i": Shell(cpu)
    // registers
    case "r": cpu.DumpRegisters()
    // TODO: search
    // TODO: unassemble
    // TODO: enter
    // TODO: assemble?
    // trace
    case "t":
        // TODO: optional argument of instructions to execute
        cpu.ExecuteInstruction()
        cpu.DumpRegisters()
    // code
    case "c": 
        for _, code := range(args) {
            patch := gamegenie.Decode(code)
            fmt.Printf("%s = %s\n", patch, patch.Encode())
            // TODO: apply
        }
        if len(args) == 0 {
            fmt.Printf("usage: c game-genie-code\n")
        }

    // TODO: breakpoints
    // TODO: watch
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

    cpu.MappersBeforeExecute[0] = func(addr uint16) (wants bool, ptr *uint8) {
        fmt.Printf("before mapper got %x\n", addr)
        ptr = &cpu.Memory[0]
        return true, ptr
    }
    cpu.MappersAfterExecute[0] = func(addr uint16, ptr *uint8) {
        fmt.Printf("after %.4x -> %x\n", addr, *ptr)
        ptr = &cpu.Memory[0]
    }
 
    for _, cmd := range(os.Args[1:]) {
        RunCommand(cpu, cmd)
    }

    if len(os.Args) < 2 {
        Shell(cpu)
    }
}

