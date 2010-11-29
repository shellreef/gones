// Created:20101128
// By Jeff Connelly

// 6502 

// http://nesdev.parodius.com/NESDoc.pdf

package cpu6502

import "fmt"

import . "./nesfile"
import (
    "bytes"
    "./dis6502"
)

type CPU struct {
    Memory [0xffff]uint8
    PC uint16   // Program counter
    SP uint8    // Stack pointer, offset from $0100
    A uint8     // Accumulator
    X, Y uint8  // Index registers
    P uint8     // Processor Status (7-0 = N V - B D I Z C)
}

// Read from memory, advancing program counter

func (cpu *CPU) NextUInt8() (b uint8) {
    b = cpu.Memory[cpu.PC]
    cpu.PC += 1
    return b
}

func (cpu *CPU) NextInt8() (b int8) {
    b = int8(cpu.Memory[cpu.PC])
    cpu.PC += 1
    return b
}

func (cpu *CPU) NextUInt16 (w uint16) {
    low := cpu.NextUInt8()
    high := cpu.NextUInt8()
    return uint16(high) * 0x100 + uint16(low)
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

// Start execution
func (cpu *CPU) Run() {
    // TODO: stop using a buffer
    buffer := bytes.NewBuffer(cpu.Memory[0x8000:])
    for {
         instr, err := dis6502.ReadInstruction(buffer)
         if err != nil {
             break
         }
         // TODO: show address, but how to get file pointer on a Buffer? no Tell()
         fmt.Printf("%s\n", instr)
    }

}

