// Created:20101207
// By Jeff Connelly

// Ricoh 2C02 PPU (Picture Processing Unit) emulation

// http://nesdev.parodius.com/2C02%20technical%20reference.TXT

package ppu2c02

//import ("fmt")

const PPU_CTRL      = 0x2000
const PPU_MASK      = 0x2001
const PPU_STATUS    = 0x2002
const PPU_SPR_ADDR  = 0x2003
const PPU_SPR_DATA  = 0x2004
const PPU_SCROLL    = 0x2005
const PPU_ADDRESS   = 0x2006
const PPU_DATA      = 0x2007

const PIXELS_PER_SCANLINE = 341
const SCANLINES_PER_FRAME = 262

type PPU struct {
    CycleChannel chan int
} 

func (ppu *PPU) Run() {
    for {
        ppu.CycleChannel <- 1
        //fmt.Printf("PPU cycle\n")
    }
}

func Before(operAddr uint16) (wants bool, ptr *uint8) {
        if operAddr < 0x2000 || operAddr > 0x3fff {
            //fmt.Printf("mapper: PPU doesn't care about %.4x\n", operAddr)
            return false, ptr
        }

        // $2000-2007 is mirrored every 8 bytes
        operAddr &^= 0x1ff8

        //fmt.Printf("mapper: before %x\n", operAddr) 

        switch operAddr {
        case PPU_STATUS: 
            a := uint8(0x80)    // VBLANK
            return true, &a
        }

        return false, ptr
}

func After(operAddr uint16, ptr *uint8) {
        // $2000-2007 is mirrored every 8 bytes
        operAddr &^= 0x1ff8    

        //fmt.Printf("mapper: after %.4x -> %x (%.8b)\n", operAddr, *ptr, *ptr)
}
