// Created:20101128
// By Jeff Connelly

// 6502 

// http://nesdev.parodius.com/NESDoc.pdf

package cpu6502

import . "./nesfile"
import . "./dis6502"
import (
    "fmt"
    "os"
)

type CPU struct {
    Memory [0xffff]uint8
    PC uint16   // Program counter
    SP uint8    // Stack pointer, offset from $0100
    A uint8     // Accumulator
    X, Y uint8  // Index registers
    P uint8     // Processor Status (7-0 = N V - B D I Z C)
}

// Processor status bits
const (
    FLAG_C = 1 << iota  // Carry
    FLAG_Z  // Zero
    FLAG_I  // Interrupt disable
    FLAG_D  // Decimal mode
    FLAG_B  // Break command
    FLAG_R  // (Reserved)
    FLAG_V  // Overflow
    FLAG_N  // Negative
)

// Read byte from memory, advancing program counter
func (cpu *CPU) NextUInt8() (b uint8) {
    b = cpu.Memory[cpu.PC]
    cpu.PC += 1
    return b
}

// Read signed byte, advancing program counter
func (cpu *CPU) NextSInt8() (b int8) {
    b = int8(cpu.Memory[cpu.PC])
    cpu.PC += 1
    return b
}

// Read unsigned 16-bits, advancing program counter
func (cpu *CPU) NextUInt16() (w uint16) {
    low := cpu.NextUInt8()
    high := cpu.NextUInt8()
    return uint16(high) * 0x100 + uint16(low)
}

// Read unsigned 16-bits at given address, not advancing PC
func (cpu *CPU) ReadUInt16(address int) (w uint16) {
    low := cpu.Memory[address]
    high := cpu.Memory[address + 1]
    return uint16(high) * 0x100 + uint16(low)
}

// Read an operand for a given addressing mode
func (cpu *CPU) NextOperand(addrMode AddrMode) (int) {
    switch addrMode {
    case Imd, Zpg, Zpx, Zpy, Ndx, Ndy: // read 8 bits
        return int(cpu.NextUInt8())
    case Abs, Abx, Aby, Ind:           // read 16 bits
        return int(cpu.NextUInt16())
    case Imp, Acc: 
        return 0
    case Rel:                          // read 8 bits
        return int(cpu.NextSInt8())    // TODO: calculate from PC
    }
    panic(fmt.Sprintf("readOperand unknown addressing mode: %s", addrMode))
} 

// Read and decode a CPU instruction from a buffer
func (cpu *CPU) NextInstruction() (*Instruction) {
    opcodeByte := cpu.NextUInt8()
   
    opcode, addrMode := Opcodes[opcodeByte].Opcode, Opcodes[opcodeByte].AddrMode
    operand := cpu.NextOperand(addrMode)
    
    instr:= new(Instruction)
    instr.Opcode = opcode
    instr.OpcodeByte = opcodeByte
    instr.AddrMode = addrMode
    instr.Operand = operand

    return instr
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

// Return string representation of truth value, for bit flags
func bitize(b uint8) (string) {
    if (b != 0) {
        return "X"
    }
    return "-"
}

// Show registers for debugging purposes
func (cpu *CPU) DumpRegisters() {
    fmt.Printf("PC   A  X  Y  CZIBD-VN\n")
    fmt.Printf("%.4X %.2X %.2X %.2X %s%s%s%s%s%s%s%s\n", cpu.PC, cpu.A, cpu.X, cpu.Y, 
        bitize(cpu.P & FLAG_C), bitize(cpu.P & FLAG_Z), bitize(cpu.P & FLAG_I), bitize(cpu.P & FLAG_D),
        bitize(cpu.P & FLAG_B), bitize(cpu.P & FLAG_R), bitize(cpu.P & FLAG_V), bitize(cpu.P & FLAG_N))
}

// Start execution
func (cpu *CPU) Run() {
    cpu.PC = 0x8000
    for {
         start := cpu.PC
         instr := cpu.NextInstruction()
         fmt.Printf("%.4X\t%s\n", start, instr)

         // Setup operPtr for writing to operand, and operVal for reading
         var operPtr *uint8;
         var operVal uint8;

         switch instr.AddrMode {
         case Imd: operPtr = nil; operVal = uint8(instr.Operand)
         case Zpg: operPtr = &cpu.Memory[instr.Operand]; operVal = *operPtr
         case Zpx: operPtr = &cpu.Memory[instr.Operand + int(cpu.X)]; operVal = *operPtr
         case Zpy: operPtr = &cpu.Memory[instr.Operand + int(cpu.Y)]; operVal = *operPtr
         case Abs: operPtr = &cpu.Memory[instr.Operand]; operVal = *operPtr
         case Abx: operPtr = &cpu.Memory[instr.Operand + int(cpu.X)]; operVal = *operPtr
         case Aby: operPtr = &cpu.Memory[instr.Operand + int(cpu.Y)]; operVal = *operPtr
         case Ndx: operPtr = &cpu.Memory[cpu.ReadUInt16(instr.Operand) + uint16(cpu.X)]; operVal = *operPtr
         case Ndy: operPtr = &cpu.Memory[int(cpu.ReadUInt16(int(instr.Operand))) + int(cpu.Y)]; operVal = *operPtr
         case Imp: operPtr = nil; operVal = 0
         case Acc: operPtr = nil; operVal = cpu.A
         case Ind: operPtr = &cpu.Memory[cpu.ReadUInt16(instr.Operand)]; operVal = *operPtr
         case Rel: operPtr = nil; //TODO: needs to be 16 bits: operVal = (cpu.PC) + 1 + instr.Operand 
         }
         fmt.Printf("operPtr=%x, operVal=%x\n", operPtr, operVal)

         switch instr.Opcode {
         // http://nesdev.parodius.com/6502.txt
         case SEI: cpu.P |= FLAG_I
         case SEC: cpu.P |= FLAG_C
         case SED: cpu.P |= FLAG_D
         case CLD: cpu.P &^= FLAG_D       // Note: &^ is bit clear operator
         case CLC: cpu.P &^= FLAG_C
         case CLI: cpu.P &^= FLAG_I
         case CLV: cpu.P &^= FLAG_V
         case LDA: cpu.A = uint8(instr.Operand)  // TODO: repeat operand *value*, not address
         case STA: cpu.Memory[instr.Operand] = cpu.A // TODO: other addressing modes etc.
         case U__:
             fmt.Printf("halting on undefined opcode\n")
             os.Exit(0)
         }
         cpu.DumpRegisters()
         fmt.Printf("\n")
    }

}

