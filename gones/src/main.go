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
    "leggo" // install from ../leggo-my-allegro

    "nesfile"
    "cpu6502"
    //"dis6502"
    "ppu2c02"
    "cartdb"
    "gamegenie"
)

// TODO: move to a NES Control Deck abstraction?

// Handle keystrokes from GUI
func processKeystroke(ch chan leggo.Event) {
    for {
        e := <-ch
        if e.Type == leggo.EVENT_KEY_DOWN {
            switch e.Keycode {
            case leggo.KEY_ESCAPE: os.Exit(0)
            case leggo.KEY_SPACE: fmt.Printf("display FPS = %f\n", leggo.FPS()) // not PPU FPS (TODO)
            }
        }
    }
}

// Start emulation
func Start(cpu *cpu6502.CPU, ppu *ppu2c02.PPU, showGui bool) {
    ppu.CPU = cpu

    if showGui {
        leggo.LeggoMain(func() { ppu.Run() }, processKeystroke)
    } else {
        ppu.Run()
    }
}

// Load a game
func Load(cpu *cpu6502.CPU, ppu *ppu2c02.PPU, filename string) {
    cart := nesfile.Open(filename)
    cpu.Load(cart)
    ppu.Load(cart)

    // Check ROM against known hashes
    // TODO: maybe this should be in nesfile, or in controldeck?
    matches := cartdb.Identify(cartdb.Load(), cart)
    if len(matches) != 0 {
        fmt.Printf("cartdb: found %d matching cartridges:\n", len(matches))
    }
    for i, match := range matches {
        fmt.Printf("%d. ", i)
        cartdb.DumpMatch(match)
    }
    // TODO: use info here, to display title/image, or use mapper
}

var showGui bool = true

// Run a command to do something with the unit
// TODO: in shell module
func RunCommand(cpu *cpu6502.CPU, ppu *ppu2c02.PPU, cmd string) {
    // If given filename (no command), load and run it
    if strings.HasSuffix(cmd, ".nes") || strings.HasSuffix(cmd, ".NES") {
        _, err := os.Stat(cmd)
        if err == nil {
            Load(cpu, ppu, cmd)
            Start(cpu, ppu, showGui)
        } // have to ignore non-existing files, since might be a command
    }

    // TODO: better cmd parsing. real scripting language?
    tokens := strings.Split(cmd, " ", -1)
    name := tokens[0]
    args := tokens[1:]
    switch name {
    // go (not to be confused with the programming language)
    case "g", "go": 
        if len(args) > 0 {
            // optional start address to override RESET vector
            startInt, err := strconv.Btoui64(args[0], 0)
            if err != nil {
                fmt.Fprintf(os.Stderr, "unable to read start address: %s: %s\n", args[0], err)
                os.Exit(-1)
            }

            cpu.PC = uint16(startInt)
        }
        // TODO: number of instructions, cycles, pixels, or frames to execute
        // (2i, 2c, 2p, 2f)
        Start(cpu, ppu, showGui)
    // load
    case "l", "load":
        if len(args) > 0 {
            Load(cpu, ppu, strings.Join(args, " ")) // TODO: use cmd directly instead of rejoining
            // TODO: verbose load
        } else {
            fmt.Printf("usage: l <filename>\n")
        }
    // hide GUI
    case "h", "hide": showGui = false

    // interactive
    case "i", "shell": Shell(cpu, ppu)
    // registers
    case "r", "registers": cpu.DumpRegisters()
    // TODO: search
    // TODO: unassemble

    // pattern table
    case "p", "show-pattern":
        if len(args) < 2 {
            fmt.Printf("usage: p <table-number> <tile-number>\n")
            return
        }
        // TODO: allow omitting to show all
        table, err := strconv.Btoui64(args[0], 0)
        if err == nil {
            table = 0
        }
        tile, err := strconv.Btoui64(args[1], 0) 
        if err == nil {
            tile = 0
        }
        ppu.PrintPattern(ppu.GetPattern(int(table), int(tile)))

    /* TODO: be able to interact with shell when GUI is open, so could do this
    case "draw-patterns":
        ppu.ShowPatterns()
    */

    // enter
    case "e", "enter":
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
    // TODO: get string -> Opcode/AddrMode working again
    /*
    case "a", "assemble":
        if len(args) != 2 {
            fmt.Printf("usage: a <opcode> <addrmode>\n")
            return
        }
        opcode := dis6502.Opcode(args[0])
        addrMode := dis6502.AddrMode(args[1])
        // TODO: show all matches
        fmt.Printf("%.2x\n", dis6502.OpcodeByteFor(opcode, addrMode))
        // TODO: read operands and actually assemble and write
    */
    // trace
    case "t", "trace":
        // TODO: optional argument of instructions to execute
        cpu.ExecuteInstruction()
        cpu.DumpRegisters()
    // cheat code
    case "c", "code": 
        for _, code := range(args) {
            // TODO: PAR codes (RAM patches), FCEU ("emulator PAR"; 001234 FF) and PAR format (1234 FF)
            // TODO: ROM patches, complete addresses in PRG chip to patch, more precise than
            //  GG codes since can specify entire address (rXXXXX:XX?)
            // TODO: maybe CHR ROM patches too? vXXXXX:XX

            patch := gamegenie.Decode(code)
            fmt.Printf("%s = %s\n", patch, patch.Encode())
            // TODO: apply
        }
        if len(args) == 0 {
            fmt.Printf("usage: c game-genie-code\n")
        }


    // verbose
    case "v", "show-cycles": cpu.Verbose = true; ppu.Verbose = true
    case "b", "show-vblank": ppu.Verbose = true
    case "V", "show-trace": cpu.InstrTrace = true

    // TODO: breakpoints
    // TODO: watch
    default: fmt.Printf("unknown command: %s\n", cmd)
    }
}

// Interactive shell: prompt for commands on stdin and run them
func Shell(cpu *cpu6502.CPU, ppu *ppu2c02.PPU) {
    prompt := "-"
    for {
        line := readline.ReadLine(&prompt)

        if line == nil {
            break
        }

        if *line == "" {
            continue
        }

        RunCommand(cpu, ppu, *line)
        readline.AddHistory(*line)
    }
}

func main() {
    cpu := new(cpu6502.CPU)
    ppu := new(ppu2c02.PPU)

    // Would like to be able to do this, but go says: method ppu.ReadMapper is not an expression, must be called
    //cpu.ReadMappers[0] = ppu.ReadMapper
    //cpu.WriteMappers[0] = ppu.WriteMapper
    // so I have to wrap it in a closure:
    cpu.ReadMappers[0] = func(address uint16) (wants bool, ret uint8) { return ppu.ReadMapper(address) }
    cpu.WriteMappers[0] = func(address uint16, b uint8) (bool) { return ppu.WriteMapper(address, b) }

    cpu.PowerUp()

    for _, cmd := range(os.Args[1:]) {
        RunCommand(cpu, ppu, cmd)
    }

    if len(os.Args) < 2 {
        Shell(cpu, ppu)
    }
}

