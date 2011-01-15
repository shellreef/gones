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
    "io/ioutil"
    
    "readline"  // http://bitbucket.org/taruti/go-readline/
    "leggo" // install from ../leggo-my-allegro

    "cartridge"
    "cartdb"
    "cpu6502"
    "ppu2c02"
    "io2a03"
    "gamegenie"
    "cheatdb"
)

type ControlDeck struct {
    InternalRAM [0x800]uint8      // 2 K

    CPU *cpu6502.CPU
    PPU *ppu2c02.PPU
    IO *io2a03.IO

    ShowGui bool
}

func New() (*ControlDeck) {
    deck := new(ControlDeck)
    deck.CPU = new(cpu6502.CPU)
    deck.PPU = new(ppu2c02.PPU)
    deck.IO = new (io2a03.IO)
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


    // $4018-FFFF are available to cartridge
  
    // TODO: properly emulate open bus (such as for missing SRAM), reads last value on bus
    // (currently it will crash on nil pointer panic)
    deck.CPU.Map(0x4000, 0x7fff,
        func(address uint16)(value uint8) { /*fmt.Printf("read unmapped: %.4x\n", address); */ return 0 },
        func(address uint16, value uint8) { /*fmt.Printf("write unmapped: %.4x:%.2x\n", address, value)*/ },
        "Open bus")

    // APU & I/O registers, $4000-4017
    // These *are* completely decoded, but for simplicity, map all 4K
    deck.CPU.Map(0x4000, 0x4fff, 
        func(address uint16)(value uint8) { return deck.IO.ReadRegister(address) },
        func(address uint16, value uint8) { deck.IO.WriteRegister(address, value) },
        "2A03 I/O")

    deck.CPU.Map(0x8000, 0xffff,
        func(address uint16)(value uint8) { /*fmt.Printf("PRG read unmapped: %.4x\n", address);*/ return 0 },
        func(address uint16, value uint8) { /*fmt.Printf("PRG write unmapped: %.4x:%.2x\n", address, value)*/ },
        "Missing Game Pak")


    deck.CPU.PowerUp()

    return deck
}

// Load a game
func (deck *ControlDeck) Load(filename string) {
    cart := cartridge.LoadFile(filename)

    // Check ROM against known hashes
    // TODO: maybe this should be in nesfile, or in controldeck?
    /*
    matches := cartdb.Identify(cartdb.Load(), cart)
    if len(matches) != 0 {
        fmt.Printf("cartdb: found %d matching cartridges:\n", len(matches))
    }
    for i, match := range matches {
        fmt.Printf("%d. ", i)
        cartdb.DumpMatch(match)
    }*/
    // TODO: use info here, to display title/image, or use mapper


    cart.LoadPRG(deck.CPU)
    cart.LoadCHR(deck.PPU)

    // nestest.nes (TODO: better way to have local additions to cartdb!)
    hash := cartdb.CartHash(cart)
    if hash == "4131307F0F69F2A5C54B7D438328C5B2A5ED0820" {
        // TODO: only do this when running in automated mode, still want GUI mode unaffected
        fmt.Printf("Identified nestest; mapping for automation\n")

        deck.CPU.MapOver(0xc66e, 
                // The known-correct log http://nickmass.com/images/nestest.log ends at $C66E, on RTS
                func(address uint16)(value uint8) { 
                    // http://nesdev.com/bbs/viewtopic.php?t=7130
                    result := deck.CPU.ReadFrom(2) << 8 | deck.CPU.ReadFrom(3)

                    if result == 0 {
                        fmt.Printf("Nestest automation: Pass\n")
                    } else {
                        fmt.Printf("Nestest automation: FAIL with code %.4x\n", result)
                    }
                    // Signal to stop CPU (KIL instruction)
                    // TODO: find a better way to quit the program. Could quit here, but
                    // then last instruction trace wouldn't match expected log.
                    deck.InternalRAM[0x0001] = 0x02

                    return 0x60 },
                func(address uint16, value uint8) { },
                "nestest-automation")
    }
}

// Read a bunch of files from a directory
func (deck *ControlDeck) LoadDir(dirname string) {
    filenames, err := ioutil.ReadDir(dirname)
    if err != nil {
        panic(fmt.Sprintf("LoadDir(%s) failed: %s", dirname, err))
    }

    mapperUsage := make(map[string]int)

    // Show information about each file
    for _, fileInfo := range filenames {
        path := fmt.Sprintf("%s/%s", dirname, fileInfo.Name)
        fmt.Printf("%s\n", path)
        cart := cartridge.LoadFile(path)
        fmt.Printf("%s\t%s\n", cart.MapperName, fileInfo.Name)
        mapperUsage[cart.MapperName] += 1
    }

    // Show mapper usage frequency
    fmt.Printf("\n")
    for mapper, count := range mapperUsage {
        fmt.Printf("%d\t%s\n", count, mapper)
    }
}


// Start emulation
func (deck *ControlDeck) Start() {
    if deck.ShowGui {
        leggo.LeggoMain(func() { deck.PPU.Run() },
            // Handle keystrokes from GUI
            func (ch chan leggo.Event) {
                for {
                    e := <-ch
                    pressed := e.Type == leggo.EVENT_KEY_DOWN
                    switch e.Type {
                    case leggo.EVENT_KEY_DOWN, leggo.EVENT_KEY_UP:
                        switch e.Keycode {
                        case leggo.KEY_ESCAPE, leggo.KEY_Q: os.Exit(0)
                        case leggo.KEY_F: 
                            fmt.Printf("display FPS = %f\n", leggo.FPS())
                            fmt.Printf("PPU FPS = %f\n", deck.PPU.FPS)

                        // TODO: custom controller keyboard key mapping; refactor all this
                        case leggo.KEY_ALT, leggo.KEY_ALTGR:
                            deck.IO.SetButtonState(0, io2a03.BUTTON_A, pressed)
                        case leggo.KEY_LSHIFT, leggo.KEY_RSHIFT:
                            deck.IO.SetButtonState(0, io2a03.BUTTON_B, pressed)
                        case leggo.KEY_TAB:
                            deck.IO.SetButtonState(0, io2a03.BUTTON_SELECT, pressed)
                        case leggo.KEY_ENTER:   
                            deck.IO.SetButtonState(0, io2a03.BUTTON_START, pressed)
                        case leggo.KEY_LEFT:
                            deck.IO.SetButtonState(0, io2a03.BUTTON_LEFT, pressed)
                        case leggo.KEY_RIGHT:
                            deck.IO.SetButtonState(0, io2a03.BUTTON_RIGHT, pressed)
                        case leggo.KEY_UP:
                            deck.IO.SetButtonState(0, io2a03.BUTTON_UP, pressed)
                        case leggo.KEY_DOWN:
                            deck.IO.SetButtonState(0, io2a03.BUTTON_DOWN, pressed)
                        }
                    }
                }
            })

    } else {
        deck.PPU.Run()
    }
}

// Run a command to do something with the unit
func (deck *ControlDeck) RunCommand(cmd string) {
    // If given filename (no command), load and run it
    if strings.HasSuffix(cmd, ".nes") || strings.HasSuffix(cmd, ".NES") ||
       strings.HasSuffix(cmd, ".unf") || strings.HasSuffix(cmd, ".UNF") ||
       strings.HasSuffix(cmd, ".unif") || strings.HasSuffix(cmd, ".UNIF") {
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
    case "load-dir":
        if len(args) > 0 {
            deck.LoadDir(strings.Join(args, " ")) // TODO: use cmd directly instead of rejoining
        }

    // hide GUI
    case "h", "hide": deck.ShowGui = false

    // interactive
    case "i", "shell": deck.Shell()
    // registers
    case "r", "registers": deck.CPU.DumpRegisters()


    case "m", "memory": deck.CPU.DumpMemoryMap()

    case "show-carts": cartdb.Dump(cartdb.Load())

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
        if table == 0 {
            table = 0
        } else {
            table = 0x1000
        }
        deck.PPU.PrintPattern(deck.PPU.GetPattern(uint16(table), int(tile)))

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

                deck.CPU.WriteTo(writeAddress, writeByte)
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
    case "load-cheats":
        cheatdb.Load()
    // cheat code
    case "c", "code": 
        for _, code := range(args) {
            // TODO: Pro Action Replay/PAR codes (RAM patches), FCEU ("emulator PAR"; 001234 FF) and PAR format (1234 FF)
            // TODO: ROM patches, complete addresses in PRG chip to patch, more precise than GG
            //  since can specify entire address (XXXXX=XX?)
            // TODO: maybe CHR ROM patches too? (XXXXX=XX, but with higher address? PRG+CHR)
            // TODO: Pro Action Rocky (ROM patcher like Game Genie but for Famicom) - see http://www.disgruntleddesigner.com/chrisc/fcrocky.html - need more info

            patch := gamegenie.Decode(code)
            fmt.Printf("%s = %s\n", patch, patch.Encode())
            fmt.Printf("CPU address: %.4X\n", patch.CPUAddress())
            fmt.Printf("Value: %.2X\n", patch.Value)
            // TODO: apply code

            // Find out what the CPU address currently maps to in ROM
            // TODO: option to attempt to do this statically.. the code below will only find what is currently loaded.
            romAddress, romChip := deck.CPU.Address2ROM(patch.CPUAddress())

            currentValue := deck.CPU.ReadFrom(patch.CPUAddress())
            if !patch.HasKey || currentValue == patch.Key {
                fmt.Printf("ROM address found: %.6X\n", romAddress)
                fmt.Printf("ROM chip: %s\n", romChip)
                fmt.Printf("iNES offset: %.4X\n", romAddress + 0x10)
            } else {
                // The affected part of ROM is not currently loaded
                fmt.Printf("This is NOT the ROM address: %.6X, since it currently contains %.2X, but key is %.2X\n", romAddress, currentValue, patch.Key)
            }

            // TODO: should we also, after finding the ROM address, search for all CPU
            // addresses it affects, since it may be more than one? For example, 16 KB NROM
            // games (like Mario Bros.) have codes like SXTIEG = 5ce0:a5 = CPU address dce0,
            // which affects ROM address 1ce0, since CPU d000-dfff is mapped to ROM 1000-1fff;
            // however, paching ROM 1ce0 would ALSO affect CPU 9ce0, since CPU 8000-8fff is 
            // mapped there, too. Game Genie, since it uses CPU addresses, can affect each of
            // these mirrors independently. The code for the mirrored bank is 1ce0:a5=SXTPEG,
            // but it has no effect because the program in this case runs from c000-ffff not
            // the mirrored 8000-bfff. 
            // Point is, GG codes can not only be LESS specific than ROM patches (affecting multiple
            // addresses in the ROM), they can also be MORE specific, since they use the CPU address.
            // TODO: also should check compare value, to restrict banks

        }
        if len(args) == 0 {
            fmt.Printf("usage: c game-genie-code\n")
        }


    // verbose
    case "v", "show-cycles": deck.CPU.Verbose = true
    case "show-ppu": deck.PPU.Verbose = true
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


