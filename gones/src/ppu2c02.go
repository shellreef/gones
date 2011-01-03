// Created:20101207
// By Jeff Connelly

// Ricoh 2C02 PPU (Picture Processing Unit) emulation

// http://nesdev.parodius.com/2C02%20technical%20reference.TXT
// http://wiki.nesdev.com/w/index.php/PPU_registers

package ppu2c02

import (
    "fmt"
    "time"

    "cpu6502"
    "leggo"
)

import . "nesfile"

// Registers
const PPU_CTRL      = 0x2000
const PPU_MASK      = 0x2001
const PPU_STATUS    = 0x2002
const PPU_OAM_ADDR  = 0x2003
const PPU_OAM_DATA  = 0x2004
const PPU_SCROLL    = 0x2005
const PPU_ADDRESS   = 0x2006
const PPU_DATA      = 0x2007

const PIXELS_PER_SCANLINE = 341
const PIXELS_VISIBLE = 256

const SCANLINES_INIT = 20
const SCANLINES_DUMMY = 1
const SCANLINES_DATA = 240
const SCANLINES_PER_FRAME = SCANLINES_INIT + SCANLINES_DUMMY + SCANLINES_DATA + SCANLINES_DUMMY // 262 total
const SCANLINES_VISIBLE = SCANLINES_DATA

const SPRITE_SIZE_8x8 = false
const SPRITE_SIZE_8x16 = true

// Memory map http://wiki.nesdev.com/w/index.php/PPU_memory_map
const PATTERN_TABLE_0 = 0x0000      // to 0x0fff, aka "left http://wiki.nesdev.com/w/index.php/PPU_pattern_tables
const PATTERN_TABLE_1 = 0x1000      // to 0x1fff, aka "right"
const NAME_TABLE_0 = 0x2000         // to 0x03c0
const ATTR_TABLE_0 = 0x23c0         // to 0x2400
// Name and attribute tables 1 to 3 follow, possibly mirrored
const IMAGE_PALETTE_1 = 0x3f00      // to 0x3f10
const SPRITE_PALETTE_1 = 0x3f10     // to 0x3f20

// Master cycles per CPU and PPU cycles. This determines their
// synchronization.. NTSC and Dendy is 3:1, PAL is 3.2:1.
// See http://nesdev.parodius.com/bbs/viewtopic.php?p=61467
const NTSC_CPU_CYCLES, NTSC_PPU_CYCLES = 12, 4
const DENDY_CPU_CYCLES, DENDY_PPU_CYCLES = 15, 5
const PAL_CPU_CYCLES, PAL_PPU_CYCLES = 15, 6
type VideoMode string
const (
        NTSC="NTSC";
        Dendy="Dendy";
        PAL="PAL")


type PPU struct {
    Pixel int
    Scanline int

    CycleCount uint             // PPU cycle count
    masterCycles uint           // Current "master cycle" count for CPU synchronization
    ppuMasterCycles uint        // Master cycles per PPU cycle
    cpuMasterCycles uint        // Master cycles per CPU cycle
    CPU *cpu6502.CPU

    // Set by PPU_CTRL
    nametableBase uint16        // Base nametable address
    vramIncrement uint16        // "VRAM address increment per CPU read/write of PPUDATA"
    spritePatternBase8x8 uint16 // Sprite pattern table address for 8x8 sprites
    backgroundBase uint16       // Background pattern table address
    spriteSize bool             // SPRITE_SIZE_8x8 or 8x16
    nmiEnabled bool             // "Generate an NMI at the start of the vertical blanking interval"

    // Set by PPU_MASK
    grayscale bool              // If enabled, AND all palette entries with 0x30 for monochrome display
    backgroundLeftmost bool     // "Enable backgrounds in leftmost 8 pixels of screen"
    spritesLeftmost bool        // "Enable sprites in leftmost 8 pixels of screen"
    backgroundEnabled bool     
    spritesEnabled bool
    intensifyRed bool
    intensifyGreen bool
    intensifyBlue bool

    // Returned by PPU_STATUS
    spriteOverflow bool
    spriteZeroHit bool
    vblankStarted bool

    vblankStartedAt int64       // For framerate

    Verbose bool

    // http://wiki.nesdev.com/w/index.php/PPU_memory_map
    Memory [0x4000]uint8        // $0000-3FFF
    vramAddress uint16          // Address latch used by PPU_SCROLL and PPU_ADDRESS
    partialAddress bool      // true if in middle of reading/writing second byte of address
    

    // Object Attribute Memory, information on up to 64 sprites
    OAM [0x100]uint8            // $00-FF, see http://wiki.nesdev.com/w/index.php/PPU_OAM
    oamAddress uint8
} 

const VRAM_ADDRESS_MASK = 0x3fff    // Mask off valid address range

// Run the CPU continuously, running the PPU every cycle
func (ppu *PPU) Run() {
    ppu.masterCycles = 0

    if ppu.cpuMasterCycles == 0 || ppu.ppuMasterCycles == 0 {
        ppu.SetVideoMode(NTSC)
    }

    ppu.CPU.CycleCallback = func() {
        // for every CPU cycle... (TODO: remove argument, only knowledge of NTSC/PAL in PPU!)
        ppu.masterCycles += ppu.cpuMasterCycles

        for ppu.masterCycles > ppu.ppuMasterCycles {
            ppu.RunOne()
            ppu.masterCycles -= ppu.ppuMasterCycles
        }
    }

    ppu.CPU.Run()
}

// Run the PPU for one PPU cycle
func (ppu *PPU) RunOne() {
    ppu.CycleCount += 1
    
    //fmt.Printf("PPU(%d): %d,%d\n", ppu.CycleCount, ppu.Pixel, ppu.Scanline)

    // http://nesdev.parodius.com/NES%20emulator%20development%20guide.txt
    ppu.Pixel = int(ppu.CycleCount) % PIXELS_PER_SCANLINE
    //ppu.HBlank = ppu.Pixel > 255
    ppu.Scanline = int(ppu.CycleCount) / PIXELS_PER_SCANLINE
    // Note first 21 are not displayed (1 dummy, 20 init)
   
    if ppu.CycleCount == SCANLINES_PER_FRAME * PIXELS_PER_SCANLINE { 
        ppu.CycleCount = 0
        ppu.VBlank()

        ppu.ShowNametable()
        //ppu.ShowPatterns()
    }


    /* Naive branch-heavy pixel/scanline counting: 20 frames/second 
    ppu.Pixel += 1
    if ppu.Pixel >= PIXELS_PER_SCANLINE {
        ppu.Scanline += 1
        ppu.Pixel = 0
    }

    if ppu.Scanline >= SCANLINES_PER_FRAME {
        ppu.Scanline = -1
        ppu.VBlank()
    }*/

    // TODO: render
    /*
    if ppu.Pixel < PIXELS_VISIBLE && ppu.Scanline < SCANLINES_VISIBLE {
        leggo.WritePixel(ppu.Pixel, ppu.Scanline, 255,0,0,0)
    }*/
}

// Set video mode for CPU synchronization
func (ppu *PPU) SetVideoMode(mode VideoMode) {
    // TODO: var MasterCycleRatios = map[string]([2]int){?
    // Then could avoid this redundancy
    switch mode {
    case NTSC: 
        ppu.cpuMasterCycles = NTSC_CPU_CYCLES
        ppu.ppuMasterCycles = NTSC_PPU_CYCLES
    case Dendy:
        ppu.cpuMasterCycles = DENDY_CPU_CYCLES
        ppu.ppuMasterCycles = DENDY_PPU_CYCLES
    case PAL:
        ppu.cpuMasterCycles = PAL_CPU_CYCLES
        ppu.ppuMasterCycles = PAL_PPU_CYCLES

    default:
        panic(fmt.Sprintf("SetVideoMode(%d): invalid", mode))
    }
}

// Vertical blank
func (ppu *PPU) VBlank() {
    // TODO: only run if VBlank flag is not disabled

    nsPerFrame := time.Nanoseconds() - ppu.vblankStartedAt
    fps := 1 / (float(nsPerFrame) / 1e9)
    if ppu.Verbose {
        fmt.Printf("VBLANK: %f frames per second\n", fps)
    }

    ppu.vblankStarted = true
    ppu.vblankStartedAt = time.Nanoseconds()
    if ppu.nmiEnabled {
        ppu.CPU.NMI() 
    }
}


// Read registers
func (ppu *PPU) ReadRegister(operAddr uint16) (wants bool, ret uint8) {
        if operAddr < 0x2000 || operAddr > 0x3fff {
            if ppu.Verbose { 
                fmt.Printf("mapper: PPU doesn't care about read %.4X\n", operAddr)
            }
            return false, ret
        }

        // $2000-2007 is mirrored every 8 bytes
        // TODO: replace with & 7!
        operAddr &^= 0x1ff8

        wants = true
        switch operAddr {
        case PPU_STATUS: 
            // $2002.5 "Sprite overflow. The PPU can handle only eight sprites on one
            // scanline and sets this bit if it starts dropping sprites"
            if ppu.spriteOverflow {
                ret |= 0x20
            }
            // $2002.6 "Sprite 0 Hit.  Set when a nonzero pixel of sprite 0 'hits'
            // a nonzero background pixel.  Used for raster timing."
            if ppu.spriteZeroHit {
                ret |= 0x40
            }
            // $2002.7 Vertical blank has started
            if ppu.vblankStarted {
                ret |= 0x80
            }

            // "Reading the status register will clear D7 mentioned above"
            ppu.vblankStarted = false
            
            // "and also the address latch"
            ppu.vramAddress = 0
            ppu.partialAddress = false

            if ppu.Verbose {
                fmt.Printf("read PPUSTATUS = %.2X\n", ret)
            }

            //ppu.ShowNametable()

        case PPU_OAM_DATA:
            // Note: reading OAM is unreliable in real hardware
            ret = ppu.OAM[ppu.oamAddress]

        case PPU_DATA:
            ret = ppu.Memory[ppu.vramAddress]
            ppu.vramAddress += ppu.vramIncrement
            ppu.vramAddress &= VRAM_ADDRESS_MASK

        default:
            wants = false
        }

        return wants, ret
}

// Write registers
func (ppu *PPU) WriteRegister(operAddr uint16, b uint8) (bool) {
        if operAddr > 0x3fff {
            // Not our territory
            if ppu.Verbose {
                fmt.Printf("mapper: PPU doesn't care about write %.4X -> %.2X\n", operAddr, b)
            }
            return false
        }

        // $2000-2007 is mirrored every 8 bytes up to $3FFF
        // TODO: replace with & 7!
        operAddr &^= 0x1ff8 

        if ppu.Verbose {
            fmt.Printf("mapper: write %.4X -> %X (%.8b)\n", operAddr, b, b)
        }

        switch operAddr {
        case PPU_CTRL: 
            // $2000.0-1: "0 = $2000; 1 = $2400; 2 = $2800; 3 = $2C00"
            ppu.nametableBase = 0x2000 | (uint16(b) & 0x03 << 10) 

            // $2000.2: "VRAM address increment per CPU read/write of PPUDATA (0: increment by 1, going across; 1: increment by 32, going down)"
            if b & 0x04 == 0 {
                ppu.vramIncrement = 1   // across
            } else {
                ppu.vramIncrement = 32  // down
            }

            // $2000.3: "(0: $0000; 1: $1000; ignored in 8x16 mode)"
            ppu.spritePatternBase8x8 = (uint16(b) & 0x08) << 8

            // $2000.4: "(0: $0000; 1: $1000)"
            ppu.backgroundBase = (uint16(b) & 0x10) << 8

            // $2000.5: "Sprite size (0: 8x8; 1: 8x16)"
            ppu.spriteSize = b & 0x20 != 0

            // $2005.6: "PPU master/slave select (has no effect on the NES)"
            // $2005.7: "Generate an NMI at the start of the VBLANK"
            ppu.nmiEnabled = b & 0x80 != 0

            if ppu.Verbose {
                fmt.Printf("PPU_CTRL: nametable=%.4X, vramIncrement=%d, spritePatternBase8x8=%.4X, backgroundBase=%.4X, spriteSize=%t, nmiEnabled=%t\n",
                    ppu.nametableBase, ppu.vramIncrement, ppu.spritePatternBase8x8, ppu.backgroundBase, ppu.spriteSize, ppu.nmiEnabled)
            }

        case PPU_MASK:
            // $2001 flags
            ppu.grayscale           = b & 0x01 != 0
            ppu.backgroundLeftmost  = b & 0x02 != 0
            ppu.spritesLeftmost     = b & 0x04 != 0
            ppu.backgroundEnabled   = b & 0x08 != 0
            ppu.spritesEnabled      = b & 0x10 != 0
            // http://wiki.nesdev.com/w/index.php/NTSC_video
            ppu.intensifyRed        = b & 0x20 != 0
            ppu.intensifyGreen      = b & 0x40 != 0
            ppu.intensifyBlue       = b & 0x80 != 0

            if ppu.Verbose {
                fmt.Printf("PPU_MASK = %.8b\n", b)
            }

        case PPU_OAM_ADDR:
            // $2003, sets address of OAM to write to
            ppu.oamAddress = b
            if ppu.Verbose {
                fmt.Printf("PPU_OAM_ADDR = %.2X\n", b)
            }

        case PPU_OAM_DATA:
            // $2004, to access OAM, but most games use DMA ($4014) instead
            if ppu.Verbose {
                fmt.Printf("PPU_OAM_DATA write to %.2X: %.2X\n", ppu.oamAddress, b)
            }
            ppu.OAM[ppu.oamAddress] = b
            ppu.oamAddress += 1

        case PPU_ADDRESS:
            // Write here twice for a 16-bit address
            if !ppu.partialAddress {
                // First write is to the upper byte
                ppu.vramAddress = uint16(b) << 8
                ppu.partialAddress = true
            } else {
                // Second write to the lower byte
                ppu.vramAddress |= uint16(b)
                ppu.partialAddress = false

                // $4000-FFFF is mirror of $0000-3FFF
                ppu.vramAddress &= VRAM_ADDRESS_MASK

                if ppu.Verbose { 
                    fmt.Printf("PPU_ADDRESS latched %.4X\n", ppu.vramAddress)
                }
            }

        case PPU_DATA:
            ppu.Memory[ppu.vramAddress] = b
            if ppu.Verbose {
                fmt.Printf("PPU_DATA write to %.4X: %.2X\n", ppu.vramAddress, b)
            }
            ppu.vramAddress += ppu.vramIncrement
            ppu.vramAddress &= VRAM_ADDRESS_MASK
        }

        return true
}

// Initialize the PPU with a game cartridge
func (ppu *PPU) Load(cart *Cartridge) {
    if len(cart.Chr) > 0 { 
        // Load CHR ROM pattern tables #0 and #1 from cartridge, if present, into $0000-1FFF
        // Note that not all games have this, some use CHR RAM instead, but 
        // in either case it is wired in the cartridge, not the NES itself
        // TODO: pointers instead of copying, so can switch CHR banks easily
        copy(ppu.Memory[0:0x2000], cart.Chr[0])
        // TODO: what if there are more CHR banks, which one to load first?
    }
}

// Get a pattern from the pattern table
// table: 0 or 1
// tile: 0 to 255
// TODO: combine with OAM to get sprite bitmap
func (ppu *PPU) GetPattern(table int, tile int) (pattern [8][8]uint8) {

    // There are two pattern tables (of 16 bytes), with 256 titles each
    base := table << 12 | tile << 4
    for row := 0; row < 8; row += 1 {
        // Palette data is split into two planes, for bit 0 and bit 1
        plane0Row := ppu.Memory[base + row]
        plane1Row := ppu.Memory[base + row + 8]

        // Combine bits of both planes to get color of bitmapped pattern
        for column := uint(0); column < 8; column += 1 {
            bit0 := (plane0Row & (1 << column)) >> column
            bit1 := (plane1Row & (1 << column)) >> column
            color := bit0 | (bit1 << 1)

            pattern[row][column] = color
        }
    }

    return pattern
}

// Print ASCII diagram of a pattern
func (ppu *PPU) PrintPattern(pattern [8][8]uint8) {
    for row := 0; row < 8; row += 1 {
        for column := 0; column < 8; column += 1 {
            fmt.Printf("%d ", pattern[row][column])
        }
        fmt.Printf("\n")
    }
}

// Show pattern table for debugging purposes
func (ppu *PPU) ShowPatterns() {
    for i := 0; i < 255; i += 1 {
        ppu.DrawPattern(ppu.GetPattern(0, i), (i % 16) * 8, (i / 16) * 8)
        //time.Sleep(1000000000)
    }
}

// Draw pattern to screen for debugging purposes
func (ppu *PPU) DrawPattern(pattern [8][8]uint8, offX int, offY int) {
    for row := 0; row < 8; row += 1 {
        for column := 0; column < 8; column += 1 {
            // TODO: palette
            // This is not real
            var r, g, b byte
            switch pattern[7-column][7-row] {
            case 0: r=255
            case 1: g=255
            case 2: b=255
            case 3: r=255; g=255; b=255
            }

            // Clipping
            if row+offX > PIXELS_VISIBLE || column+offY > SCANLINES_VISIBLE {
                continue
            }

            leggo.WritePixel(row+offX, (7-column)+offY, r,g,b,0)
        }
    }
}



func (ppu *PPU) ShowNametable() {
    // http://wiki.nesdev.com/w/index.php/PPU_nametables 
    // TODO: mirroring, access other nametables (4)
    base := 0x2000
    patternTable := 0
    for row := 0; row < 30; row += 1 {
        for column := 0; column < 32; column += 1 {
            tile := ppu.Memory[base + row*30 + column]
            pattern := ppu.GetPattern(patternTable, int(tile))

            ppu.DrawPattern(pattern, column*8, row*8)
        }
    }
}
