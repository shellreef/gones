// Created:20101128
// By Jeff Connelly

// 6502 

// http://nesdev.parodius.com/NESDoc.pdf

package cpu6502

import "fmt"

//import . "./dis6502"
import . "./nesfile"

type CPU struct {
    Memory [0xffff]uint8
}

// Load a game cartridge
func (cpu *CPU) Load(cart *Cartridge) {
    // http://nesdev.parodius.com/NESDoc.pdf
    if (len(cart.Prg) == 1) {
        // One PRG bank loads into $8000 and $C000 (mirrored)
        copy(cpu.Memory[0x8000:], cart.Prg[0])
        copy(cpu.Memory[0xC000:], cart.Prg[0])
    } else if (len(cart.Prg) == 2) {
        // Loads first into $8000 and second into $C000
        copy(cpu.Memory[0x8000:], cart.Prg[0])
        copy(cpu.Memory[0xC000:], cart.Prg[1])
    } else {
        // TODO: mappers
        panic(fmt.Sprintf("Load: PRG banks not yet supported: %d", len(cart.Prg)))
    }
}
