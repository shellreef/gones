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
    Memory [0x10000]uint8          // $0000-$FFFF
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

// Interrupt vector addresses
const (
    NMI_VECTOR   = 0xfffa
    RESET_VECTOR = 0xfffc
    BRK_VECTOR   = 0xfffe
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
    return uint16(high) << 8 + uint16(low)
}

// Read unsigned 16-bits at given address, not advancing PC
func (cpu *CPU) ReadUInt16(address uint16) (w uint16) {
    low := cpu.Memory[address]
    high := cpu.Memory[uint16(address + 1)]
    return uint16(high) << 8 + uint16(low)
}

// Read unsigned 16-bits from zero page
func (cpu *CPU) ReadUInt16ZeroPage(address uint8) (w uint16) {
    low := cpu.Memory[address]
    high := cpu.Memory[uint8(address + 1)]    // 0xff + 1 wraps around
    return uint16(high) << 8 + uint16(low)
}

// Read unsigned 16-bits, but wraparound lower byte
// JMP ($20FF) jumps the address at $20FF low and $2000 high (not $2100)
// TODO: replace ReadUInt16ZeroPage with this, and maybe ReadUInt16 entirely?
func (cpu *CPU) ReadUInt16Wraparound(address uint16) (w uint16) {
    low := cpu.Memory[address]
    high := cpu.Memory[(address & 0xff00) | (address + 1) & 0xff]
    return uint16(high) << 8 + uint16(low)
}

// Read an operand for a given addressing mode
// Byte count read is retrievable by addrMode.OperandSize()
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
    instr.Official = OfficialOpcodes[opcodeByte].Opcode != U__

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

// Show 256-byte stack for debugging purposes
func (cpu *CPU) DumpStack() {
    for i := 0; i < 0x100; i += 1 {
        if i == int(cpu.S) {
            fmt.Printf(">%.2x<", cpu.Memory[0x100 + i]) // indicate stack pointer
        } else {
            fmt.Printf(" %.2x ", cpu.Memory[0x100 + i])
        }
        // TODO: show addresses, multi-column?
    }
    fmt.Printf("\n")
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
func (cpu *CPU) SetCarry(f bool) { cpu.SetFlag(FLAG_C, f) }
func (cpu *CPU) SetOverflow(f bool) { cpu.SetFlag(FLAG_V, f) }
func (cpu *CPU) SetInterrupt(f bool) { cpu.SetFlag(FLAG_I, f) }
//func (cpu *CPU) SetBreak(f bool) { cpu.SetFlag(FLAG_B, f) } // not an actual flag; set on PHP/BRK
func (cpu *CPU) SetDecimal(f bool) { cpu.SetFlag(FLAG_D, f) }

// Branch if flag is set
func (cpu* CPU) BranchIf(address uint16 , flag bool) {
    if flag {
        cpu.PC = address
    }
}

// Pull 8 bits from stack
func (cpu *CPU) Pull() (b uint8) {
    cpu.S += 1
    b = cpu.Memory[0x100 + uint16(cpu.S)]
    return b
}

// Pull 16 bits from stack
func (cpu *CPU) Pull16() (w uint16) {
    low := cpu.Pull()
    high := cpu.Pull()
    return uint16(high) * 0x100 + uint16(low)
}

// Push 8 bits onto stack
func (cpu *CPU) Push(b uint8) {
    cpu.Memory[0x100 + uint16(cpu.S)] = b 
    cpu.S -= 1
}

// Push 16 bits onto stack
func (cpu *CPU) Push16(w uint16) {
    cpu.Push(uint8(w >> 8))
    cpu.Push(uint8(w & 0xff))
}

// Initialize CPU to power-up state
// http://wiki.nesdev.com/w/index.php/CPU_power_up_state
func (cpu *CPU) PowerUp() {
    cpu.P = FLAG_R | FLAG_I   // Reserved bit set, see http://wiki.nesdev.com/w/index.php/CPU_status_flag_behavior
    cpu.S = 0xfd
    // TODO: set $0000-$07FF to $ff
    cpu.Memory[0x0008] = 0xf7
    cpu.Memory[0x0009] = 0xef
    cpu.Memory[0x000a] = 0xdf
    cpu.Memory[0x000f] = 0xbf

    // Memory mapped
    cpu.Memory[0x4017] = 0x00   // frame IRQ enabled
    cpu.Memory[0x4015] = 0x00   // all channels disabled
    // TODO: 0x4000-$400f set to $00
}

// Some instructions worth having in their own functions

func (cpu *CPU) OpROR(operPtr *uint8) {
    var temp int
    temp = int(*operPtr)
    if cpu.P & FLAG_C != 0 {
        temp |= 0x100
    }
    cpu.SetCarry(temp & 0x01 != 0)
    temp >>= 1
    *operPtr = uint8(temp)
    cpu.SetSZ(*operPtr)
}

func (cpu *CPU) OpADC(operVal uint8) {
    var carryIn, temp int
    if cpu.P & FLAG_C == 0 {
        carryIn = 0
    } else {
        carryIn = 1
    }
    temp = int(operVal) + int(cpu.A) + carryIn
    cpu.SetSZ(uint8(temp))
    cpu.SetOverflow(((int(cpu.A) ^ int(operVal)) & 0x80 == 0) && ((int(cpu.A) ^ temp) & 0x80 != 0))
    cpu.SetCarry(temp > 0xff)
    cpu.A = uint8(temp)
}

func (cpu *CPU) OpSBC(operVal uint8) {
    var carryIn, temp uint
    if cpu.P & FLAG_C == 0 {
        carryIn = 1
    } else {
        carryIn = 0
    }
    temp = uint(cpu.A) - uint(operVal) - carryIn
    cpu.SetSZ(uint8(temp))
    cpu.SetOverflow(((uint(cpu.A) ^ uint(temp)) & 0x80 != 0) && ((uint(cpu.A) ^ uint(operVal)) & 0x80 != 0))
    cpu.SetCarry(temp < 0x100)
    cpu.A = uint8(temp)
}

// Execute one instruction
func (cpu *CPU) ExecuteInstruction() {
    start := cpu.PC
    instr := cpu.NextInstruction()

    // Instruction trace
    cycle := 0 // TODO
    scanline := 0 // TODO
    fmt.Printf("%.4X  ", start)
    fmt.Printf("%.2X ", instr.OpcodeByte)
    if instr.AddrMode.OperandSize() >= 1 {
       fmt.Printf("%.2X ", cpu.Memory[start + 1])
    } else {
       fmt.Printf("   ")
    }

    if instr.AddrMode.OperandSize() >= 2 {
       fmt.Printf("%.2X ", cpu.Memory[start + 2])
    } else {
       fmt.Printf("   ")
    }


    fmt.Printf(" %-31s A:%.2X X:%.2X Y:%.2X P:%.2X SP:%.2X CYC:%3d SL:%3d\n",
       instr, cpu.A, cpu.X, cpu.Y, cpu.P, cpu.S, cycle, scanline)

    // Setup operPtr for writing to operand, and operVal for reading
    // Not all addressing modes allow writing to the operand; in that case,
    // operPtr will be nil
    var operPtr *uint8    // pointer to write/read, if applicable
    var operAddr uint16   // address to write/read, if applicable
    var operVal uint8     // value to read

    // http://wiki.nesdev.com/w/index.php/CPU_addressing_modes
    switch instr.AddrMode {
    case Zpg: operAddr = uint16(instr.Operand);                     operPtr = &cpu.Memory[operAddr]
    case Zpx: operAddr = uint16(uint8(instr.Operand) + cpu.X);      operPtr = &cpu.Memory[operAddr]
    case Zpy: operAddr = uint16(uint8(instr.Operand) + cpu.Y);      operPtr = &cpu.Memory[operAddr]
    case Abs: operAddr = uint16(instr.Operand);                     operPtr = &cpu.Memory[operAddr]
    case Abx: operAddr = uint16(instr.Operand) + uint16(cpu.X);     operPtr = &cpu.Memory[operAddr]
    case Aby: operAddr = uint16(instr.Operand) + uint16(cpu.Y);     operPtr = &cpu.Memory[operAddr]
    case Ndx: operAddr = cpu.ReadUInt16ZeroPage(uint8(instr.Operand) + cpu.X); operPtr = &cpu.Memory[operAddr]  // ($%.2X,X)
    case Ndy: operAddr = cpu.ReadUInt16ZeroPage(uint8(instr.Operand)) + uint16(cpu.Y); operPtr = &cpu.Memory[operAddr]  // ($%.2X),Y
    case Ind: operAddr = cpu.ReadUInt16Wraparound(uint16(instr.Operand));  operPtr = &cpu.Memory[operAddr]
    case Rel: operAddr = (cpu.PC) + uint16(instr.Operand);      operPtr = &cpu.Memory[operAddr] // TODO: clk += ((PC & 0xFF00) != (REL_ADDR(PC, src) & 0xFF00) ? 2 : 1);
    case Acc: operPtr = &cpu.A  /* no address */
    case Imd: operVal = uint8(instr.Operand)
    case Imp: operVal = 0
    }
    if operPtr != nil {
        operVal = *operPtr
    }
    
    // Shorthand convenience for opcode implementation
    //src := operVal

    switch instr.Opcode {
    // http://nesdev.parodius.com/6502.txt
    // http://www.obelisk.demon.co.uk/6502/reference.html#ADC

    case NOP, DOP, TOP:

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
    case LAX: cpu.A = operVal; cpu.X = operVal; cpu.SetSZ(cpu.X)
    // Store register to memory
    case STA: *operPtr = cpu.A
    case STX: *operPtr = cpu.X
    case STY: *operPtr = cpu.Y

    // Transfers
    case TAX: cpu.X = cpu.A; cpu.SetSZ(cpu.A)   // would like to do cpu.SetSZ((cpu.X=cpu.A)) like in C, but can't in Go
    case TAY: cpu.Y = cpu.A; cpu.SetSZ(cpu.A)
    case TSX: cpu.X = cpu.S; cpu.SetSZ(cpu.S)
    case TXA: cpu.A = cpu.X; cpu.SetSZ(cpu.X)
    case TXS: cpu.S = cpu.X
    case TYA: cpu.A = cpu.Y; cpu.SetSZ(cpu.Y)

    // Bitwise operations
    case AND: cpu.A &= operVal; cpu.SetSZ(cpu.A)
    case EOR: cpu.A ^= operVal; cpu.SetSZ(cpu.A)
    case ORA: cpu.A |= operVal; cpu.SetSZ(cpu.A)
    case ASL: cpu.SetCarry(operVal & 0x80 != 0)
        *operPtr <<= 1
        cpu.SetSZ(*operPtr)
    case LSR: cpu.SetCarry(operVal & 0x01 != 0)
        *operPtr >>= 1
        cpu.SetSZ(*operPtr)
    case ROL: 
        var temp int
        temp = int(operVal)    // larger than uint8 so can store carry bit
        temp <<= 1
        if cpu.P & FLAG_C != 0 {
            temp |= 1
        }
        cpu.SetCarry(temp > 0xff)
        *operPtr = uint8(temp)
        cpu.SetSZ(*operPtr)
    case ROR: cpu.OpROR(operPtr)
    case BIT: cpu.SetSign(operVal)
        cpu.SetOverflow(0x40 & operVal != 0)
        cpu.SetZero(operVal & cpu.A)
    case AAX: *operPtr = cpu.X & cpu.A  // no flags affected

    // Arithmetic
    case DEC: *operPtr -= 1; cpu.SetSZ(*operPtr)
    case DEX: cpu.X -= 1; cpu.SetSZ(cpu.X)
    case DEY: cpu.Y -= 1; cpu.SetSZ(cpu.Y)
    case INC: *operPtr += 1; cpu.SetSZ(*operPtr)
    case INX: cpu.X += 1; cpu.SetSZ(cpu.X)
    case INY: cpu.Y += 1; cpu.SetSZ(cpu.Y)
    case ADC: cpu.OpADC(operVal)
    case RRA: cpu.OpROR(operPtr); cpu.OpADC(*operPtr)
    case SBC: cpu.OpSBC(operVal)
    case DCP: *operPtr -= 1
        var temp uint
        temp = uint(cpu.A) - uint(*operPtr)
        cpu.SetCarry(temp < 0x100)
        cpu.SetSZ(uint8(temp))
    case ISC: *operPtr += 1; cpu.OpSBC(*operPtr)

    case CMP, CPX, CPY:
        var temp uint
        switch instr.Opcode {
        case CMP: temp = uint(cpu.A) - uint(operVal)
        case CPX: temp = uint(cpu.X) - uint(operVal)
        case CPY: temp = uint(cpu.Y) - uint(operVal)
        }
        cpu.SetCarry(temp < 0x100)
        cpu.SetSZ(uint8(temp))
 
 
     // Branches
    case BCC: cpu.BranchIf(operAddr, cpu.P & FLAG_C == 0)
    case BCS: cpu.BranchIf(operAddr, cpu.P & FLAG_C != 0)
    case BNE: cpu.BranchIf(operAddr, cpu.P & FLAG_Z == 0)
    case BEQ: cpu.BranchIf(operAddr, cpu.P & FLAG_Z != 0)
    case BPL: cpu.BranchIf(operAddr, cpu.P & FLAG_N == 0)
    case BMI: cpu.BranchIf(operAddr, cpu.P & FLAG_N != 0)
    case BVC: cpu.BranchIf(operAddr, cpu.P & FLAG_V == 0)
    case BVS: cpu.BranchIf(operAddr, cpu.P & FLAG_V != 0)
 
     // Stack
    case PHA: cpu.Push(cpu.A)
    case PHP: cpu.Push(cpu.P | FLAG_B)                         // no actual "B" flag, but 4th P bit is set on PHP (and BRK)
    case PLA: cpu.A = cpu.Pull(); cpu.SetSZ(cpu.A)
    case PLP: cpu.P = cpu.Pull() | FLAG_R; cpu.P &^= FLAG_B    // on pull, R "flag" is always set and B "flag" always clear (same with RTI). See http://wiki.nesdev.com/w/index.php/CPU_status_flag_behavior

     // Jumps
    case JMP: cpu.PC = operAddr
    case JSR: cpu.PC -= 1
        cpu.Push16(cpu.PC)
        cpu.PC = operAddr
    case RTI: cpu.P = cpu.Pull() | FLAG_R; cpu.P &^= FLAG_B; cpu.PC = cpu.Pull16()
    case RTS: cpu.PC = cpu.Pull16() + 1
    case BRK: cpu.PC += 1
        cpu.Push16(cpu.PC)
        cpu.Push(cpu.P | FLAG_B)
        cpu.SetInterrupt(true)
        cpu.PC = cpu.ReadUInt16(BRK_VECTOR)
 
    case U__:
         fmt.Printf("halting on undefined opcode\n")
         os.Exit(0)
 
    default:
         fmt.Printf("++ TODO: implement %s\n", instr.Opcode)
         os.Exit(0)
    }
 
    // Post-instruction execution trace
    //fmt.Printf("%.4X  %s  ", start, instr)
    //fmt.Printf("operPtr=%x, operAddr=%.4X, operVal=%.2X\n", operPtr, operAddr, operVal)
    //cpu.DumpRegisters()
    //fmt.Printf("\n")
}

// Start execution
func (cpu *CPU) Run() {
    cpu.PowerUp()

    cpu.PC = cpu.ReadUInt16(RESET_VECTOR)
    cpu.PC = 0xc000  // for nestest

    // PPU status register (TODO: memory mapped I/O)
    // http://nocash.emubase.de/everynes.htm#memorymaps
    //cpu.Memory[0x2002] = 0x80 // VBLANK=1

    for {
        cpu.ExecuteInstruction()
    }
}

