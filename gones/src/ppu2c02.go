// Created:20101207
// By Jeff Connelly

// Ricoh 2C02 PPU (Picture Processing Unit) emulation

// http://nesdev.parodius.com/2C02%20technical%20reference.TXT
// http://wiki.nesdev.com/w/index.php/PPU_registers

package ppu2c02

import (
    "cpu6502"
    "fmt"
)

const PPU_CTRL      = 0x2000
const PPU_MASK      = 0x2001
const PPU_STATUS    = 0x2002
const PPU_SPR_ADDR  = 0x2003
const PPU_SPR_DATA  = 0x2004
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

    Verbose bool
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
            // VBLANK flag set
            // TODO: for real
            ret = 0x80
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

            fmt.Printf("nametable=%.4X, vramIncrement=%d, spritePatternBase8x8=%.4X, backgroundBase=%.4X, spriteSize=%t, nmiEnabled=%t\n",
                ppu.nametableBase, ppu.vramIncrement, ppu.spritePatternBase8x8, ppu.backgroundBase, ppu.spriteSize, ppu.nmiEnabled)
        }

        return true
}
