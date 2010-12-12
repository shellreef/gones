// Created:20101207
// By Jeff Connelly

// Ricoh 2C02 PPU (Picture Processing Unit) emulation

// http://nesdev.parodius.com/2C02%20technical%20reference.TXT
// http://wiki.nesdev.com/w/index.php/PPU_registers

package ppu2c02

import (
    "fmt"

    "cpu6502"
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

const SCANLINES_INIT = 20
const SCANLINES_DUMMY = 1
const SCANLINES_DATA = 240
const SCANLINES_PER_FRAME = SCANLINES_INIT + SCANLINES_DUMMY + SCANLINES_DATA + SCANLINES_DUMMY // 262 total

const PPU_MASTER_CYCLES = 5     // 5 "master cycles" per PPU cycle

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

type PPU struct {
    Pixel int
    Scanline int

    CycleChannel chan int       // Synchronize with CPU
    CPU *cpu6502.CPU

    // Set by PPU_CTRL
    nametableBase uint16        // Base nametable address
    vramIncrement int           // "VRAM address increment per CPU read/write of PPUDATA"
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

    Verbose bool

    // http://wiki.nesdev.com/w/index.php/PPU_memory_map
    Memory [0x4000]uint8        // $0000-3FFF
    vramAddress uint16          // Address latch used by PPU_SCROLL and PPU_ADDRESS
    partialAddress bool      // true if in middle of reading/writing second byte of address
    

    // Object Attribute Memory, information on up to 64 sprites
    OAM [0x100]uint8            // $00-FF, see http://wiki.nesdev.com/w/index.php/PPU_OAM
    oamAddress uint8
} 

// Continuously run
func (ppu *PPU) Run() {
    for {
        // One pixel per clock cycle
        ppu.Pixel += 1

        if ppu.Pixel > PIXELS_PER_SCANLINE {
            ppu.Pixel = 0
            ppu.Scanline += 1
        }

        if ppu.Scanline > SCANLINES_PER_FRAME {
            ppu.VBlank()
            ppu.Scanline = -1
        }

        // TODO: render

        // Tick
        ppu.CycleChannel <- PPU_MASTER_CYCLES
        //fmt.Printf("PPU: %d,%d\n", ppu.Pixel, ppu.Scanline)
    }
}

// Vertical blank
func (ppu *PPU) VBlank() {
    // TODO: only run if VBlank flag is not disabled

    if ppu.Verbose {
        fmt.Printf("** VBLANK **\n")
    }
    ppu.vblankStarted = true
    if ppu.nmiEnabled {
        ppu.CPU.NMI() 
    }
}


// Read registers
func (ppu *PPU) ReadMapper(operAddr uint16) (wants bool, ret uint8) {
        if operAddr < 0x2000 || operAddr > 0x3fff {
            fmt.Printf("mapper: PPU doesn't care about read %.4X\n", operAddr)
            return false, ret
        }

        // $2000-2007 is mirrored every 8 bytes
        operAddr &^= 0x1ff8

        //fmt.Printf("mapper: before %X\n", operAddr) 

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

            fmt.Printf("read PPUSTATUS = %.2X\n", ret)

        case PPU_OAM_DATA:
            // Note: reading OAM is unreliable in real hardware
            ret = ppu.OAM[ppu.oamAddress]

        default:
            wants = false
        }

        return wants, ret
}

// Write registers
func (ppu *PPU) WriteMapper(operAddr uint16, b uint8) (bool) {
        if operAddr > 0x3fff {
            // Not our territory
            fmt.Printf("mapper: PPU doesn't care about write %.4X -> %.2X\n", operAddr, b)
            return false
        }

        // $2000-2007 is mirrored every 8 bytes up to $3FFF
        operAddr &^= 0x1ff8 

        fmt.Printf("mapper: write %.4X -> %X (%.8b)\n", operAddr, b, b)

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

            fmt.Printf("PPU_CTRL: nametable=%.4X, vramIncrement=%d, spritePatternBase8x8=%.4X, backgroundBase=%.4X, spriteSize=%t, nmiEnabled=%t\n",
                ppu.nametableBase, ppu.vramIncrement, ppu.spritePatternBase8x8, ppu.backgroundBase, ppu.spriteSize, ppu.nmiEnabled)

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

            fmt.Printf("PPU_MASK = %.8b\n", b)

        case PPU_OAM_ADDR:
            // $2003, sets address of OAM to write to
            ppu.oamAddress = b
            fmt.Printf("PPU_OAM_ADDR = %.2X\n", b)

        case PPU_OAM_DATA:
            // $2004, to access OAM, but most games use DMA ($4014) instead
            fmt.Printf("PPU_OAM_DATA write to %.2X: %.2X\n", ppu.oamAddress, b)
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

                fmt.Printf("PPU_ADDRESS latched %.4X\n", ppu.vramAddress)
            }
                
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

    for tile := 0; tile < 255; tile += 1 {
        ppu.ShowPattern(0, tile)
    }
}

// table: 0 or 1
// tile: 0 to 255
func (ppu *PPU) ShowPattern(table int, tile int) {
    // There are two pattern tables (of 16 bytes), with 256 titles each
    base := table << 12 | tile << 4
    fmt.Printf("\n#%d\n", tile)
    for row := 0; row < 8; row += 1 {
        // Palette data is split into two planes, for bit 0 and bit 1
        plane0Row := ppu.Memory[base + row]
        plane1Row := ppu.Memory[base + row + 8]

        // Combine bits of both planes to get color of bitmapped pattern
        for column := uint(0); column < 8; column += 1 {
            bit0 := (plane0Row & (1 << column)) >> column
            bit1 := (plane1Row & (1 << column)) >> column
            color := bit0 | (bit1 << 1)

            fmt.Printf("%d ", color)
        }
        fmt.Printf("\n")
    }
}
