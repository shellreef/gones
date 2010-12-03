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
    S uint8    // Stack pointer, offset from $0100
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
    if b != 0 {
        return "X"
    }
    return "-"
}

// Show registers for debugging purposes
func (cpu *CPU) DumpRegisters() {
    fmt.Printf("PC   A  X  Y  NV-DBIZC\n")
    fmt.Printf("%.4X %.2X %.2X %.2X %s%s%s%s%s%s%s%s\n", cpu.PC, cpu.A, cpu.X, cpu.Y, 
        bitize(cpu.P & FLAG_N), bitize(cpu.P & FLAG_V), bitize(cpu.P & FLAG_R), bitize(cpu.P & FLAG_D),
        bitize(cpu.P & FLAG_B), bitize(cpu.P & FLAG_I), bitize(cpu.P & FLAG_Z), bitize(cpu.P & FLAG_C))
}

// Set/reset flag given by mask FLAG_* if byte is non-zero/zero
func (cpu *CPU) SetFlag(mask byte, f bool) {
    if f {
        cpu.P |= mask
    } else {
        cpu.P &^= mask
    }
}

// Set/reset sign flag (N, for negative if 1) from 7th bit in given byte
func (cpu *CPU) SetSign(b byte) {
    cpu.SetFlag(FLAG_N, (b & 0x80) == 0x80)
}

// Set/reset zero flag (Z) if byte is zero/non-zero
func (cpu *CPU) SetZero(b byte) {
    cpu.SetFlag(FLAG_Z, b == 0)
}

// Set sign and zero
func (cpu *CPU) SetSZ(b byte) {
    cpu.SetSign(b)
    cpu.SetZero(b)
}

// Convenience functions to set/reset flags depending on non-zero/zero
func (cpu *CPU) SetCarry(b byte) { cpu.SetFlag(FLAG_C, b != 0) }
func (cpu *CPU) SetOverflow(b byte) { cpu.SetFlag(FLAG_V, b != 0) }
func (cpu *CPU) SetInterrupt(b byte) { cpu.SetFlag(FLAG_I, b != 0) }
func (cpu *CPU) SetBreak(b byte) { cpu.SetFlag(FLAG_B, b != 0) }
func (cpu *CPU) SetDecimal(b byte) { cpu.SetFlag(FLAG_D, b != 0) }

// Branch if flag is set
func (cpu* CPU) BranchIf(address uint16 , flag bool) {
    if flag {
        cpu.PC = address
    }
}

// Start execution
func (cpu *CPU) Run() {
    cpu.PC = 0x8000
    for {
         start := cpu.PC
         instr := cpu.NextInstruction()
         fmt.Printf("%.4X\t%s\n", start, instr)

         // Setup operPtr for writing to operand, and operVal for reading
         // Not all addressing modes allow writing to the operand; in that case,
         // operPtr will be nil
         var operPtr *uint8    // pointer to write/read, if applicable
         var operAddr uint16   // address to write/read, if applicable
         var operVal uint8     // value to read

         switch instr.AddrMode {
         case Zpg: operAddr = uint16(instr.Operand);                     operPtr = &cpu.Memory[operAddr]
         case Zpx: operAddr = uint16(instr.Operand) + uint16(cpu.X);     operPtr = &cpu.Memory[operAddr]
         case Zpy: operAddr = uint16(instr.Operand) + uint16(cpu.Y);     operPtr = &cpu.Memory[operAddr]
         case Abs: operAddr = uint16(instr.Operand);                     operPtr = &cpu.Memory[operAddr]
         case Abx: operAddr = uint16(instr.Operand) + uint16(cpu.X);     operPtr = &cpu.Memory[operAddr]
         case Aby: operAddr = uint16(instr.Operand) + uint16(cpu.Y);     operPtr = &cpu.Memory[operAddr]
         case Ndx: operAddr = cpu.ReadUInt16(instr.Operand) + uint16(cpu.X);         operPtr = &cpu.Memory[operAddr]
         case Ndy: operAddr = uint16(cpu.ReadUInt16(instr.Operand)) + uint16(cpu.Y); operPtr = &cpu.Memory[operAddr]
         case Ind: operAddr = cpu.ReadUInt16(instr.Operand);  operPtr = &cpu.Memory[operAddr]
         case Rel: operAddr = (cpu.PC) + 1 + uint16(instr.Operand);      operPtr = &cpu.Memory[operAddr] // TODO: clk += ((PC & 0xFF00) != (REL_ADDR(PC, src) & 0xFF00) ? 2 : 1);
         case Acc: operPtr = &cpu.A  /* no address */
         case Imd: operVal = uint8(instr.Operand)
         case Imp: operVal = 0
         }
         if operPtr != nil {
             operVal = *operPtr
         }
         fmt.Printf("operPtr=%x, operAddr=%.4X, operVal=%.2X\n", operPtr, operAddr, operVal)

         // Shorthand convenience for opcode implementation
         //src := operVal

         switch instr.Opcode {
         // http://nesdev.parodius.com/6502.txt

         // Flag setting       
         case SEI: cpu.P |= FLAG_I
         case SEC: cpu.P |= FLAG_C
         case SED: cpu.P |= FLAG_D
         // Flag clearing. &^ is Go bit clear operator; cpu.P &= ^FLAG_D will not compile because of an overflow
         case CLD: cpu.P &^= FLAG_D
         case CLC: cpu.P &^= FLAG_C
         case CLI: cpu.P &^= FLAG_I
         case CLV: cpu.P &^= FLAG_V

         // Load register from memory
         case LDA: cpu.A = operVal; cpu.SetSZ(cpu.A)
         case LDX: cpu.X = operVal; cpu.SetSZ(cpu.X)
         case LDY: cpu.Y = operVal; cpu.SetSZ(cpu.Y)
         // Store register to memory
         case STA: *operPtr = cpu.A
         case STX: *operPtr = cpu.X
         case STY: *operPtr = cpu.Y

         // Transfers
         case TAX: cpu.X = cpu.A; cpu.SetSZ(cpu.A)   // would like to do cpu.SetSZ((cpu.X=cpu.A)) like in C, but can't in Go
         case TAY: cpu.Y = cpu.A; cpu.SetSZ(cpu.A)
         case TSX: cpu.X = cpu.S; cpu.SetSZ(cpu.S)
         case TXA: cpu.A = cpu.X; cpu.SetSZ(cpu.X)
         case TXS: cpu.S = cpu.X; cpu.SetSZ(cpu.X)
         case TYA: cpu.A = cpu.Y; cpu.SetSZ(cpu.Y)

         // Math operations
         case AND: cpu.A &= operVal; cpu.SetSZ(cpu.A)
         case ASL: cpu.SetCarry(operVal & 0x80)
             *operPtr = *operPtr << 1
             cpu.SetSZ(*operPtr)
         case BIT: cpu.SetSign(operVal)
             cpu.SetOverflow(0x40 & operVal)
             cpu.SetZero(operVal & cpu.A)

         // Decrement
         case DEC: *operPtr -= 1; cpu.SetSign(*operPtr); cpu.SetZero(*operPtr)
         case DEX: cpu.X -= 1; cpu.SetSign(cpu.X); cpu.SetZero(cpu.X)
         case DEY: cpu.Y -= 1; cpu.SetSign(cpu.Y); cpu.SetZero(cpu.Y)
         case INC: *operPtr += 1; cpu.SetSign(*operPtr); cpu.SetZero(*operPtr)

         // Branches
         case BCC: cpu.BranchIf(operAddr, cpu.P & FLAG_C == 0)
         case BCS: cpu.BranchIf(operAddr, cpu.P & FLAG_C != 0)
         case BEQ: cpu.BranchIf(operAddr, cpu.P & FLAG_Z == 0)
         case BNE: cpu.BranchIf(operAddr, cpu.P & FLAG_Z != 0)
         case BPL: cpu.BranchIf(operAddr, cpu.P & FLAG_N == 0)
         case BMI: cpu.BranchIf(operAddr, cpu.P & FLAG_N != 0)
         case BVS: cpu.BranchIf(operAddr, cpu.P & FLAG_V == 0)
         case BVC: cpu.BranchIf(operAddr, cpu.P & FLAG_V != 0)

         case U__:
             fmt.Printf("halting on undefined opcode\n")
             os.Exit(0)
         }

         cpu.DumpRegisters()
         fmt.Printf("\n")
    }

}

