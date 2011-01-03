// Created:20110202
// By Jeff Connelly
//
// NES Control Deck (the console unit) emulator
// This abstraction ties together various other hardware components

package controldeck

import ("cpu6502"
        "ppu2c02"
       
        )

type ControlDeck struct {
    InternalRAM [0x7ff]uint8
}

func New() (*ControlDeck) {
    deck := new(ControlDeck)
    cpu := new(cpu6502.CPU)
    ppu := new(ppu2c02.PPU)

    // RAM: 2 KB, mirrored throughout an 8 KB address space
    cpu.Map8K(0x0000, 0x1fff, 
            func(address uint16)(value uint8) { return deck.InternalRAM[address & 0x7ff] },
            func(address uint16, value uint8) { deck.InternalRAM[address & 0x7ff] = value })

    // PPU registers
    cpu.Map8K(0x2000, 0x3fff,
            // TODO: would be nice to write ppu.ReadRegister directly, but Go doesn't let you :(
            func(address uint16)(value uint8) { _, v := ppu.ReadRegister(address); return v }, // TODO: remove wants
            func(address uint16, value uint8) { ppu.WriteRegister(address, value) })

    // TODO: APU & I/O registers, $4000-4017

    // $4018-FFFF are available to cartridge

    cpu.PowerUp()

    return deck
}


