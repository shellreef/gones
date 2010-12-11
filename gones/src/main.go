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
    "dis6502"
    "ppu2c02"
    "gamegenie"
)

func Start(cpu *cpu6502.CPU) {
    ppu := new(ppu2c02.PPU)
    ppu.CycleChannel = make(chan int)
    ppu.CPU = cpu
    cpu.CycleChannel = make(chan int)
    go cpu.Run()
    go ppu.Run()


    masterCycles := 0

    for {
        // For every CPU cycle...
        masterCycles += <-cpu.CycleChannel

        // ..run PPU appropriate number of cycles
        // This is 3 for NTSC, but keep the left over for PAL (1 master cycle)
        for masterCycles > ppu2c02.PPU_MASTER_CYCLES {
            masterCycles -= <-ppu.CycleChannel
        }
    }
}

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
            Start(cpu)
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
        Start(cpu)
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
    // enter
    case "e":
        if len(args) < 2 {
            fmt.Printf("usage: e <address> <byte0> [<byte1> [<byteN...]]\n")
            return
        }

        baseAddress, err := strconv.Btoui64(args[0], 16)
        if err != nil {
            fmt.Printf("bad address: %s: %s\n", args[0], err)
            return
        }

        for offset, arg := range args[1:] {
            argInt, err := strconv.Btoui64(arg, 16)
            if err == nil {
                writeAddress := uint16(baseAddress) + uint16(offset)
                writeByte := uint8(argInt)

                fmt.Printf("%.4X = %X\n", writeAddress, writeByte)

                cpu.Memory[writeAddress] = writeByte
            } else {
                fmt.Printf("bad value: %s: %s\n", arg, err)
            }
        }
    // assemble
    case "a":
        if len(args) != 2 {
            fmt.Printf("usage: a <opcode> <addrmode>\n")
            return
        }
        opcode := dis6502.Opcode(args[0])
        addrMode := dis6502.AddrMode(args[1])
        // TODO: show all matches
        fmt.Printf("%.2x\n", dis6502.OpcodeByteFor(opcode, addrMode))
        // TODO: read operands and actually assemble and write
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

    cpu.MappersBeforeExecute[0] = ppu2c02.Before
    cpu.MappersAfterExecute[0] = ppu2c02.After
 
    for _, cmd := range(os.Args[1:]) {
        RunCommand(cpu, cmd)
    }

    if len(os.Args) < 2 {
        Shell(cpu)
    }
}

