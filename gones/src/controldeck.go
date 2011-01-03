// Created:20110202
// By Jeff Connelly
//
// NES Control Deck (the console unit) emulator
// This abstraction ties together various other hardware components

package controldeck

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

type ControlDeck struct {
    InternalRAM [0x7ff]uint8

    CPU *cpu6502.CPU
    PPU *ppu2c02.PPU

    ShowGui bool
}

func New() (*ControlDeck) {
    deck := new(ControlDeck)
    deck.CPU = new(cpu6502.CPU)
    deck.PPU = new(ppu2c02.PPU)
    deck.ShowGui = true

    // PPU calls back to CPU to synchronize execution
    deck.PPU.CPU = deck.CPU

    // RAM: 8 KB address space mirrored across 2 KB of memory
    deck.CPU.Map(0x0000, 0x1fff, 
            func(address uint16)(value uint8) { return deck.InternalRAM[address & 0x7ff] },
            func(address uint16, value uint8) { deck.InternalRAM[address & 0x7ff] = value },
            "Internal RAM")

    // PPU registers, mirrored many times
    deck.CPU.Map(0x2000, 0x3fff,
            // TODO: would be nice to write ppu.ReadRegister directly, but Go doesn't let you :(
            func(address uint16)(value uint8) { _, v := deck.PPU.ReadRegister(address); return v }, // TODO: remove wants junk
            func(address uint16, value uint8) { deck.PPU.WriteRegister(address, value) },
            "PPU registers")

    // TODO: APU & I/O registers, $4000-4017

    // $4018-FFFF are available to cartridge
    

    // XXX: Remove this old stuff, use above instead
    deck.CPU.ReadMappers[0] = func(address uint16) (wants bool, ret uint8) { return deck.PPU.ReadRegister(address) }
    deck.CPU.WriteMappers[0] = func(address uint16, b uint8) (bool) { return deck.PPU.WriteRegister(address, b) }

    deck.CPU.PowerUp()

    return deck
}

// Start emulation
func (deck *ControlDeck) Start() {
    if deck.ShowGui {
        leggo.LeggoMain(func() { deck.PPU.Run() },
            // Handle keystrokes from GUI
            func (ch chan leggo.Event) {
                for {
                    e := <-ch
                    if e.Type == leggo.EVENT_KEY_DOWN {
                        switch e.Keycode {
                        case leggo.KEY_ESCAPE: os.Exit(0)
                        case leggo.KEY_SPACE: fmt.Printf("display FPS = %f\n", leggo.FPS()) // not PPU FPS (TODO)
                        }
                    }
                }
            })

    } else {
        deck.PPU.Run()
    }
}

// Load a game
func (deck *ControlDeck) Load(filename string) {
    cart := nesfile.Open(filename)
    deck.CPU.Load(cart)
    deck.PPU.Load(cart)

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

// Run a command to do something with the unit
func (deck *ControlDeck) RunCommand(cmd string) {
    // If given filename (no command), load and run it
    if strings.HasSuffix(cmd, ".nes") || strings.HasSuffix(cmd, ".NES") {
        _, err := os.Stat(cmd)
        if err == nil {
            deck.Load(cmd)
            deck.Start()
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

            deck.CPU.PC = uint16(startInt)
        }
        // TODO: number of instructions, cycles, pixels, or frames to execute
        // (2i, 2c, 2p, 2f)
        deck.Start()
    // load
    case "l", "load":
        if len(args) > 0 {
            deck.Load(strings.Join(args, " ")) // TODO: use cmd directly instead of rejoining
            // TODO: verbose load
        } else {
            fmt.Printf("usage: l <filename>\n")
        }
    // hide GUI
    case "h", "hide": deck.ShowGui = false

    // interactive
    case "i", "shell": deck.Shell()
    // registers
    case "r", "registers": deck.CPU.DumpRegisters()
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
        deck.PPU.PrintPattern(deck.PPU.GetPattern(int(table), int(tile)))

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

                // TODO: use WriteMemory()!
                deck.CPU.Memory[writeAddress] = writeByte
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
        // TODO: Fix, really have to run PPU & CPU in sync..see 'g'
        deck.CPU.ExecuteInstruction()
        deck.CPU.DumpRegisters()
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
    case "v", "show-cycles": deck.CPU.Verbose = true; deck.PPU.Verbose = true
    case "b", "show-vblank": deck.PPU.Verbose = true
    case "V", "show-trace": deck.CPU.InstrTrace = true

    // TODO: breakpoints
    // TODO: watch
    default: fmt.Printf("unknown command: %s\n", cmd)
    }
}

// Interactive shell: prompt for commands on stdin and run them
func (deck *ControlDeck) Shell() {
    prompt := "-"
    for {
        line := readline.ReadLine(&prompt)

        if line == nil {
            break
        }

        if *line == "" {
            continue
        }

        deck.RunCommand(*line)
        readline.AddHistory(*line)
    }
}


