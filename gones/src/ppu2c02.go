// Created:20101207
// By Jeff Connelly

// Ricoh 2C02 PPU (Picture Processing Unit) emulation

// http://nesdev.parodius.com/2C02%20technical%20reference.TXT

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

type PPU struct {
    Pixel int
    Scanline int

    CycleChannel chan int       // Synchronize with CPU
    CPU *cpu6502.CPU

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
    ppu.CPU.NMI() 
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

        return true
}
