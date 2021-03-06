// Created:20101128
// By Jeff Connelly

// Emulate the NES CPU, a 6502-based processor
// Specifically, the NES uses a 6502 in the Ricoh RP2A03, which:
// - is NMOS-based 6502 (as opposed to the CMOS 650C2)
// - has no decimal mode (D flag in processor status is ignored)
// This package only handles the CPU portion of the 2A03

package cpu6502

import . "dis6502"
import (
    "fmt"
    "os"
)

type Interrupt int
const (INT_NONE=iota;
        INT_NMI;
        INT_RESET;
        INT_BRK;
        INT_IRQ)

type CPU struct {
    // Functions to read and write to memory, indexed by each 4 KB,
    // using the upper 4 bits of the address. Most mappers use 
    // 8 KB, 16 KB, or 32 KB banks (NSF uses 4 KB) so they'll map multiple.
    MemRead  [16](func(address uint16)(value uint8))
    MemWrite [16](func(address uint16, value uint8))
    MemName  [16]string         // informational, name of memory segment
    MemROMOffset [16]uint32     // informational, ROM offset mapped to, if applicable (useful for ROM patch/Game Genie code decoding)
    MemHasROM [16]bool

    PC uint16   // Program counter
    S uint8     // Stack pointer, offset from $0100
    A uint8     // Accumulator
    X, Y uint8  // Index registers
    P uint8     // Processor Status (7-0 = N V - B D I Z C) (TODO: consider using separate flags, pack/unpack as needed)

    CycleCount uint         // CPU cycle count
    CycleCallback func()    // Called on every CPU cycle
    PendingInt Interrupt    // Run interrupt after next instr finishes

    Verbose bool
    InstrTrace bool

    Instruction *Instruction  // Current instruction
    OperPtr *uint8  // Pointer to operand of current instruction
    HaveCalculatedAddress bool
    CalculatedAddress uint16
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

// Low-level read/write functions, which do not increment cycles.
// Note: instructions use *Operand methods instead
// Read from CPU address
func (cpu *CPU) ReadFrom(address uint16) (value uint8) {
    value = cpu.MemRead[(address & 0xf000) >> 12](address)
    // TODO: breakpoints on read (bpr), execute (bpx), logging 
    //fmt.Printf("read from %.4x -> %.2x\n", address, value)
    return value
}

// Write to CPU address
func (cpu *CPU) WriteTo(address uint16, value uint8) {
    // TODO: breakpoints on write (bpw)
    //fmt.Printf("WRITE %.4x:%.2x\n", address, value)
    cpu.MemWrite[(address & 0xf000) >> 12](address, value)
}

// Read byte from memory, advancing program counter
func (cpu *CPU) NextUInt8() (b uint8) {
    b = cpu.ReadFrom(cpu.PC)
    cpu.PC += 1

    cpu.Tick("fetch byte, increment PC")

    return b
}

// Get the address an operand refers to
func (cpu *CPU) AddressOperand() (address uint16) {
   
    if cpu.HaveCalculatedAddress {
        return cpu.CalculatedAddress
    }

    switch cpu.Instruction.AddrMode {
    case Rel: address = (cpu.PC) + uint16(cpu.Instruction.Operand)
    case Ind: address = cpu.ReadUInt16Wraparound(uint16(cpu.Instruction.Operand))
    case Abs: address = uint16(cpu.Instruction.Operand)
    case Abx: 
        // Loads account for page crossing
        // TODO: this is painfully ugly, can it not be special-cased please?
        if cpu.Instruction.Opcode == LDA || cpu.Instruction.Opcode == LDY || cpu.Instruction.Opcode == LDX || cpu.Instruction.Opcode == TOP {
            address = cpu.AddIndex(uint16(cpu.Instruction.Operand), uint16(cpu.X))
        } else {
            address = uint16(cpu.Instruction.Operand) + uint16(cpu.X)
        }
    case Aby: 
        if cpu.Instruction.Opcode == LDA || cpu.Instruction.Opcode == LDY || cpu.Instruction.Opcode == LDX || cpu.Instruction.Opcode == TOP {
            address = cpu.AddIndex(uint16(cpu.Instruction.Operand), uint16(cpu.Y))
        } else {
            address = uint16(cpu.Instruction.Operand) + uint16(cpu.Y)
        }
    case Zpg: address = uint16(cpu.Instruction.Operand)
    case Zpx: 
        cpu.Tick("add index register")
        address = uint16(uint8(cpu.Instruction.Operand) + cpu.X)
    case Zpy: 
        cpu.Tick("add index register")
        address = uint16(uint8(cpu.Instruction.Operand) + cpu.Y)
    case Ndx: pointer := uint8(cpu.Instruction.Operand) + cpu.X         // ($%.2X,X)
        cpu.Tick("read from address, and X to it")
        address = cpu.ReadUInt16ZeroPage(pointer)
    case Ndy: pointer := uint8(cpu.Instruction.Operand)                 // ($%.2X),Y
        address = cpu.ReadUInt16ZeroPage(pointer) 
        address = cpu.AddIndex(address, uint16(cpu.Y))
        //address += uint16(cpu.Y)

    // Accumulator, implied, immediate have no address
    default: panic(fmt.Sprintf("Address() on invalid mode: %s\n", cpu.Instruction.AddrMode))
    }

    // TODO: find a better way to avoid recalculating addres (causing extra cycles) for Modify()
    cpu.HaveCalculatedAddress = true
    cpu.CalculatedAddress = address

    return address
}

// Write to instruction operand
func (cpu *CPU) WriteOperand(b uint8) {
    // TODO: can this be not hardcoded? 
    switch cpu.Instruction.Opcode {
    case STA, AXA, SXA, SYA: 
        switch cpu.Instruction.AddrMode {
        case Abx, Aby, Ndy: 
            // http://nesdev.parodius.com/6502_cpu.txt
            // Absolute indexed addressing
            //      Write instructions (STA, STX[sic], STY[sic], SHA[AXA], SHX[SXA], SHY[SYA])
            //         4  address+I* R  read from effective address,
            /*
                  * The high byte of the effective address may be invalid
                    at this time, i.e. it may be smaller by $100. Because
                    the processor cannot undo a write to an invalid
                    address, it always reads from the address first.
            */
            _ = cpu.ReadOperand()
        }
    }

    switch cpu.Instruction.AddrMode {
    case Acc: cpu.A = b
    case Zpg, Zpx, Zpy, Ndx, Ndy:
        cpu.Tick("write to effective address")
        // TODO: only can access zero page, so write directly from memory
        cpu.WriteTo(cpu.AddressOperand(), b)

    case Abs, Abx, Aby:
        cpu.Tick("write to effective address")

        address := cpu.AddressOperand()
        cpu.WriteTo(address, b)

    // Can't write to implied, immediate, etc.
    default: panic(fmt.Sprintf("WriteOperand() bad mode: %s\n", cpu.Instruction.AddrMode))
    }
}

// Read from operand for instruction
func (cpu *CPU) ReadOperand() (b uint8) {
    switch cpu.Instruction.AddrMode {
    case Acc: return cpu.A
    case Imd: return uint8(cpu.Instruction.Operand)

    case Zpg, Zpx, Zpy, Ndx, Ndy:
        // TODO: only can access zero page, so read directly from memory
        cpu.Tick("read from effective address")
        return cpu.ReadFrom(cpu.AddressOperand())

    case Abs, Abx, Aby:
        cpu.Tick("read from effective address")

        address := cpu.AddressOperand()

        return cpu.ReadFrom(address)

    // Can't read from implied (nothing to read from), indirect or relative (no instruction does)
    default: panic(fmt.Sprintf("ReadOperand() bad mode: %s\n", cpu.Instruction.AddrMode))
    }

    return 0
}


// Read something, modify it with the given modifier function, and write it back out (for read-modify-write instructions)
func (cpu *CPU) Modify(modify func(in uint8) (out uint8)) (out uint8) {
    in := cpu.ReadOperand()

    if cpu.Instruction.AddrMode == Abx || cpu.Instruction.AddrMode == Aby {
        // Read-modify-write to absolute indexed has an extra cycle:
        // http://nesdev.parodius.com/6502_cpu.txt
        // "5  address+X  R  re-read from effective address"
        in = cpu.ReadOperand()
    }

    // read-modify-write operations write unmodified value back first
    cpu.WriteOperand(in)

    out = modify(in)
    cpu.WriteOperand(out)
    return out
}

// Read unsigned 16-bits at given address, not advancing PC
func (cpu *CPU) ReadUInt16(address uint16) (w uint16) {
    low := cpu.ReadFrom(address)
    cpu.Tick("fetch ReadUInt16 low")
    high := cpu.ReadFrom(uint16(address + 1))
    cpu.Tick("fetch ReadUInt16 high")
    return uint16(high) << 8 + uint16(low)
}

// Read unsigned 16-bits from zero page
func (cpu *CPU) ReadUInt16ZeroPage(address uint8) (w uint16) {
    low := cpu.ReadFrom(uint16(address))
    cpu.Tick("fetch effective address low (zero page)")
    high := cpu.ReadFrom(uint16(uint8(address + 1)))    // 0xff + 1 wraps around
    cpu.Tick("fetch effective address high (zero page)")
    return uint16(high) << 8 + uint16(low)
}

// An an index to a base address obtaining an effective address
// This is the same as baseAddress + offset, except accounts for extra cycles
func (cpu *CPU) AddIndex(baseAddress uint16, offset uint16) (effectiveAddress uint16) {
    effectiveAddress = baseAddress + offset

    // An extra cycle is needed to increment the high
    if effectiveAddress & 0xff00 != baseAddress & 0xff00 {
        cpu.Tick("page crossing")
    }

    return effectiveAddress
}

// Read unsigned 16-bits, but wraparound lower byte
// JMP ($20FF) jumps the address at $20FF low and $2000 high (not $2100)
// TODO: replace ReadUInt16ZeroPage with this, and maybe ReadUInt16 entirely?
func (cpu *CPU) ReadUInt16Wraparound(address uint16) (w uint16) {
    low := cpu.ReadFrom(address)
    cpu.Tick("fetch low address ReadUInt16Wraparound")
    high := cpu.ReadFrom((address & 0xff00) | (address + 1) & 0xff)
    cpu.Tick("fetch high address ReadUInt16Wraparound")
    return uint16(high) << 8 + uint16(low)
}

// Read an operand for a given addressing mode
// Byte count read is retrievable by addrMode.OperandSize()
func (cpu *CPU) NextOperand(addrMode AddrMode) (int) {
    switch addrMode {
    case Imd, Zpg, Zpx, Zpy, Ndx, Ndy: // read 8 bits
        return int(cpu.NextUInt8())
    case Abs, Abx, Aby, Ind:           // read 16 bits
        low := cpu.NextUInt8()
        high := cpu.NextUInt8()
        address := uint16(high) << 8 + uint16(low)
        return int(address)
    case Imp, Acc: 
        // all others fetch next byte and use it
        cpu.Tick("read next instruction byte (and throw it away, Imp/Acc)")
        return 0
    case Rel:                          // read 8 bits
        // Note this is signed! int8 not uint8
        offset := int8(cpu.NextUInt8())    // TODO: calculate from PC
        return int(offset)
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
    // TODO: show stack
    // TODO: show disassembly of next/current instruction
}

// Show current memory map for debugging purposes
func (cpu *CPU) DumpMemoryMap() {
    fmt.Printf("Start\tEnd\tMapped to\n")
    // TODO: coalesce 4 KB chunks?
    for i, name := range cpu.MemName {
        fmt.Printf("%.4x\t%.4x\t%-20s",
                i << 12,
                (i << 12) + 0xfff,
                name)
 
        if cpu.MemHasROM[i] {
            fmt.Printf("@ %.8x", cpu.MemROMOffset[i])
        }
    
        fmt.Print("\n")
    }
}

// Show 256-byte stack for debugging purposes
func (cpu *CPU) DumpStack() {
    var i uint16
    for i = 0; i < 0x100; i += 1 {
        if uint8(i) == cpu.S {
            fmt.Printf(">%.2x<", cpu.ReadFrom(0x100 + i)) // indicate stack pointer
        } else {
            fmt.Printf(" %.2x ", cpu.ReadFrom(0x100 + i))
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
        cpu.Tick("taking branch")
        if cpu.PC & 0xff00 != address & 0xff00 {
            cpu.Tick("branch new page")
        }
        cpu.PC = address
    }
}

// Pull 8 bits from stack
func (cpu *CPU) Pull() (b uint8) {
    cpu.S += 1
    cpu.Tick("increment S")
    b = cpu.ReadFrom(0x100 + uint16(cpu.S))
    cpu.Tick("pull stack")
    return b
}

// Pull 16 bits from stack
func (cpu *CPU) Pull16() (w uint16) {
    // pipelining causes this to take one fewer cycles than two 8-bit Pull()s
    cpu.S += 1
    cpu.Tick("increment S")

    low := cpu.ReadFrom(0x100 + uint16(cpu.S))
    cpu.S += 1
    cpu.Tick("pull low from stack, increment S")

    high := cpu.ReadFrom(0x100 + uint16(cpu.S))
    cpu.Tick("pull high from stack")

    return uint16(high) * 0x100 + uint16(low)
}

// Push 8 bits onto stack
func (cpu *CPU) Push(b uint8) {
    cpu.WriteTo(0x100 + uint16(cpu.S), b)
    cpu.S -= 1
    cpu.Tick("push stack, decrement S")
}

// Push 16 bits onto stack
func (cpu *CPU) Push16(w uint16) {
    cpu.Push(uint8(w >> 8))
    cpu.Push(uint8(w & 0xff))
}

// Increment clock cycle
func (cpu *CPU) Tick(reason string) {
    if cpu.Verbose {
        fmt.Printf("tick: %s\n", reason)
    }
    cpu.CycleCount += 1

    // TODO: remove check, but test suites don't setup channel
    if cpu.CycleCallback != nil {
        cpu.CycleCallback()
    }
}


// Initialize CPU to power-up state
// http://wiki.nesdev.com/w/index.php/CPU_power_up_state
func (cpu *CPU) PowerUp() {
    cpu.P = FLAG_R | FLAG_I   // Reserved bit set, see http://wiki.nesdev.com/w/index.php/CPU_status_flag_behavior
    cpu.S = 0xfd              // Top of stack
}

// Some instructions worth having in their own functions

func (cpu *CPU) OpADC(operVal uint8) {
    var carryIn, tmp int
    if cpu.P & FLAG_C == 0 {
        carryIn = 0
    } else {
        carryIn = 1
    }
    tmp = int(operVal) + int(cpu.A) + carryIn
    cpu.SetSZ(uint8(tmp))
    cpu.SetOverflow(((int(cpu.A) ^ int(operVal)) & 0x80 == 0) && ((int(cpu.A) ^ tmp) & 0x80 != 0))
    cpu.SetCarry(tmp > 0xff)
    cpu.A = uint8(tmp)
}

func (cpu *CPU) OpSBC(operVal uint8) {
    var carryIn, tmp uint
    if cpu.P & FLAG_C == 0 {
        carryIn = 1
    } else {
        carryIn = 0
    }
    tmp = uint(cpu.A) - uint(operVal) - carryIn
    cpu.SetSZ(uint8(tmp))
    cpu.SetOverflow(((uint(cpu.A) ^ uint(tmp)) & 0x80 != 0) && ((uint(cpu.A) ^ uint(operVal)) & 0x80 != 0))
    cpu.SetCarry(tmp < 0x100)
    cpu.A = uint8(tmp)
}

func (cpu *CPU) OpROR() (ret uint8) {
    ret = cpu.Modify(func(x uint8) (y uint8) {
        tmp := uint(x)
        if cpu.P & FLAG_C != 0 {
            tmp |= 0x100
        }
        cpu.SetCarry(tmp & 0x01 != 0)
        return uint8(tmp >> 1)
    })

    cpu.SetSZ(ret)

    return ret
}

func (cpu *CPU) OpROL() (ret uint8) {
    ret = cpu.Modify(func(x uint8) (y uint8) {
        tmp := int(x)    // larger than uint8 so can store carry bit
        tmp <<= 1

        if cpu.P & FLAG_C != 0 {
            tmp |= 1
        }
        cpu.SetCarry(tmp > 0xff)

        return uint8(tmp)
    })

    cpu.SetSZ(ret)

    return ret
}


// Execute one instruction
func (cpu *CPU) ExecuteInstruction() {
    start := cpu.PC
    startCyc := cpu.CycleCount
    cpu.HaveCalculatedAddress = false
    cpu.Instruction = cpu.NextInstruction()

    // Instruction trace
    if cpu.InstrTrace || cpu.Verbose {
        fmt.Printf("%.4X  ", start)
        fmt.Printf("%.2X ", cpu.Instruction.OpcodeByte)
        if cpu.Instruction.AddrMode.OperandSize() >= 1 {
           fmt.Printf("%.2X ", cpu.ReadFrom(start + 1))
        } else {
           fmt.Printf("   ")
        }

        if cpu.Instruction.AddrMode.OperandSize() >= 2 {
           fmt.Printf("%.2X ", cpu.ReadFrom(start + 2))
        } else {
           fmt.Printf("   ")
        }


        // Trace *before* the instruction executes
        fmt.Printf(" %-31s A:%.2X X:%.2X Y:%.2X P:%.2X SP:%.2X CYC:%3d SL:%d\n",
           cpu.Instruction, cpu.A, cpu.X, cpu.Y, cpu.P, cpu.S, 
           (startCyc * 3) % 341,   // 3 PPU cycles per 1 CPU cycle, wrap around at 341 (TODO: refactor)
           0, // TODO: scanline
           )
    }

    switch cpu.Instruction.Opcode {
    // http://nesdev.parodius.com/6502.txt
    // http://www.obelisk.demon.co.uk/6502/reference.html#ADC

    case NOP:
    case DOP: _ = cpu.ReadOperand()     // zero page, so two byte "double no-op"
    case TOP: _ = cpu.ReadOperand()     // absolute, so three byte "triple no-op"

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
    case LDA: cpu.A = cpu.ReadOperand(); cpu.SetSZ(cpu.A)
    case LDX: cpu.X = cpu.ReadOperand(); cpu.SetSZ(cpu.X)
    case LDY: cpu.Y = cpu.ReadOperand(); cpu.SetSZ(cpu.Y)
    case LAX: cpu.A = cpu.ReadOperand(); cpu.X = cpu.A; cpu.SetSZ(cpu.X)
    // Store register to memory
    case STA: cpu.WriteOperand(cpu.A)
    case STX: cpu.WriteOperand(cpu.X)
    case STY: cpu.WriteOperand(cpu.Y)

    // TODO: AXA(SHA), SXA(SHX), SYA(SHY), XAS(SHA), they & high byte of address + 1

    // Transfers
    case TAX: cpu.X = cpu.A; cpu.SetSZ(cpu.A)   // would like to do cpu.SetSZ((cpu.X=cpu.A)) like in C, but can't in Go
    case TAY: cpu.Y = cpu.A; cpu.SetSZ(cpu.A)
    case TSX: cpu.X = cpu.S; cpu.SetSZ(cpu.S)
    case TXA: cpu.A = cpu.X; cpu.SetSZ(cpu.X)
    case TXS: cpu.S = cpu.X
    case TYA: cpu.A = cpu.Y; cpu.SetSZ(cpu.Y)

    // Bitwise operations
    case AND: cpu.A &= cpu.ReadOperand(); cpu.SetSZ(cpu.A)
    case EOR: cpu.A ^= cpu.ReadOperand(); cpu.SetSZ(cpu.A)
    case ORA: cpu.A |= cpu.ReadOperand(); cpu.SetSZ(cpu.A)
    case ASL: cpu.SetSZ(cpu.Modify(func(x uint8) (uint8) {
        cpu.SetCarry(x & 0x80 != 0)
        return x << 1
        }))
    case LSR: cpu.SetSZ(cpu.Modify(func(x uint8) (uint8) {
        cpu.SetCarry(x & 0x01 != 0)
        return x >> 1
        }))

    case ROL: cpu.OpROL()
    case ROR: cpu.OpROR()
    case RLA: cpu.A &= cpu.OpROL(); cpu.SetSZ(cpu.A)
    case BIT: tmp := cpu.ReadOperand()
        cpu.SetSign(tmp)
        cpu.SetOverflow(0x40 & tmp != 0)
        cpu.SetZero(tmp & cpu.A)
    case SAX: cpu.WriteOperand(cpu.X & cpu.A)
    case SLO: cpu.A |= cpu.Modify(func(x uint8) (uint8) {
            cpu.SetCarry(x & 0x80 != 0)
            return x << 1
        })
        cpu.SetSZ(cpu.A)
    case SRE: cpu.A ^= cpu.Modify(func(x uint8) (uint8) {
            cpu.SetCarry(x & 0x01 != 0)
            return x >> 1
        })
        cpu.SetSZ(cpu.A)
    case AAC: cpu.A &= cpu.ReadOperand()
        // Despite <http://wiki.nesdev.com/w/index.php/Programming_with_unofficial_opcodes>, the zero
        // flag is actually set too, according to blargg's instr-test_v3/02-immediate.nes.
        cpu.SetSZ(cpu.A)
        cpu.SetCarry(cpu.P & FLAG_N != 0)   // Set C to N
    case ASR: cpu.A &= cpu.ReadOperand()
        cpu.SetCarry(cpu.A & 0x01 != 0)
        cpu.A >>= 1
        cpu.SetSZ(cpu.A)
    case ARR: cpu.A &= cpu.ReadOperand()
        cpu.Instruction.AddrMode = Acc
        cpu.OpROR()
        // C is bit 6, V is bit 6 XOR bit 5. http://wiki.nesdev.com/w/index.php/Programming_with_unofficial_opcodes
        cpu.SetCarry(cpu.A & 0x40 != 0)
        cpu.SetOverflow((cpu.A & 0x40 != 0) != (cpu.A & 0x20 != 0))
    case AXS: cpu.X &= cpu.A
        var tmp uint
        tmp = uint(cpu.X) - uint(cpu.ReadOperand())
        cpu.SetCarry(tmp < 0x100)
        cpu.X = uint8(tmp)
        cpu.SetSZ(cpu.X)

    // Unreliable/unknown unofficial opcodes
    // http://nesdev.parodius.com/undocumented_opcodes.txt
    // http://nesdev.parodius.com/extra_instructions.txt
    // http://nesdev.parodius.com/6502_cpu.txt
    // ATX ($AB) fails Blarggs intr_test-v3 in 02-immediate
    case ATX: 
        // This is unreliable/unknown (opcode $AB)
        // http://nesdev.parodius.com/bbs/viewtopic.php?t=3831&highlight=atx
        // says A is ORed with $EE, $EF, $FE, or $FF depending on internal state
        cpu.A |= 0xee
        cpu.A &= cpu.X
        cpu.A &= cpu.ReadOperand()
        cpu.X = cpu.A
        cpu.SetSZ(cpu.A)
    // $9C and $9E (SYA and SXA) fail instruction tests also
    // They fail Blargg instr_test-v3, but that is not the newest version..
    // http://nesdev.parodius.com/bbs/viewtopic.php?t=3831 mentions blargg_nes_cpu_test5, but
    // the link is broken, although Blargg suggests these opcodes are inconsistent so their
    // tests will be removed. Nestopia, FWIW, fails v3 too (06-abs_xy). And on v5, too..
    // Some discussion of how to make pass tests on http://nesdev.parodius.com/bbs/viewtopic.php?t=3831
    case SYA: cpu.WriteOperand(cpu.Y & uint8(cpu.AddressOperand() >> 8) + 1)
    case SXA: cpu.WriteOperand(cpu.X & uint8(cpu.AddressOperand() >> 8) + 1)

    // Arithmetic
    case DEC: cpu.SetSZ(cpu.Modify(func(x uint8) (uint8) { return x - 1}))
    case DEX: cpu.X -= 1; cpu.SetSZ(cpu.X)
    case DEY: cpu.Y -= 1; cpu.SetSZ(cpu.Y)
    case INC: cpu.SetSZ(cpu.Modify(func(x uint8) (uint8) { return x + 1 }))
    case INX: cpu.X += 1; cpu.SetSZ(cpu.X)
    case INY: cpu.Y += 1; cpu.SetSZ(cpu.Y)
    case ADC: cpu.OpADC(cpu.ReadOperand())
    case SBC: cpu.OpSBC(cpu.ReadOperand())
    case RRA: cpu.OpADC(cpu.OpROR())
    case DCP: tmp := uint(cpu.A) - uint(cpu.Modify(func(x uint8) (uint8) { return x - 1 }))
        cpu.SetCarry(tmp < 0x100)
        cpu.SetSZ(uint8(tmp))
    case ISB: cpu.OpSBC(cpu.Modify(func(x uint8) (uint8) { return x + 1 }))
    case CMP, CPX, CPY:
        var tmp uint
        switch cpu.Instruction.Opcode {
        case CMP: tmp = uint(cpu.A) - uint(cpu.ReadOperand())
        case CPX: tmp = uint(cpu.X) - uint(cpu.ReadOperand())
        case CPY: tmp = uint(cpu.Y) - uint(cpu.ReadOperand())
        }
        cpu.SetCarry(tmp < 0x100)
        cpu.SetSZ(uint8(tmp))

    // Stack
    case PHA: cpu.Push(cpu.A)
    case PHP: cpu.Push(cpu.P | FLAG_B)                         // no actual "B" flag, but 4th P bit is set on PHP (and BRK)
    case PLA: cpu.A = cpu.Pull(); cpu.SetSZ(cpu.A)
    case PLP: cpu.P = cpu.Pull() | FLAG_R; cpu.P &^= FLAG_B    // on pull, R "flag" is always set and B "flag" always clear (same with RTI). See http://wiki.nesdev.com/w/index.php/CPU_status_flag_behavior

    // Branches
    // Instructions that access cpu.AddressOperand() directly, always referring to ROM, can skip mapper checks
    case BCC: cpu.BranchIf(cpu.AddressOperand(), cpu.P & FLAG_C == 0)
    case BCS: cpu.BranchIf(cpu.AddressOperand(), cpu.P & FLAG_C != 0)
    case BNE: cpu.BranchIf(cpu.AddressOperand(), cpu.P & FLAG_Z == 0)
    case BEQ: cpu.BranchIf(cpu.AddressOperand(), cpu.P & FLAG_Z != 0)
    case BPL: cpu.BranchIf(cpu.AddressOperand(), cpu.P & FLAG_N == 0)
    case BMI: cpu.BranchIf(cpu.AddressOperand(), cpu.P & FLAG_N != 0)
    case BVC: cpu.BranchIf(cpu.AddressOperand(), cpu.P & FLAG_V == 0)
    case BVS: cpu.BranchIf(cpu.AddressOperand(), cpu.P & FLAG_V != 0)
 
    // Jumps
    case JMP: 
        /*if cpu.PC - 3 == cpu.AddressOperand() {
            fmt.Printf("*** Infinite loop detected - halting\n") // TODO: on branches, too; TODO: continue waiting for NMI, just halt CPU
            os.Exit(0)
        }*/
        cpu.PC = cpu.AddressOperand()
    case JSR: cpu.PC -= 1
        cpu.Push16(cpu.PC)
        cpu.Tick("JSR copy PC")
        cpu.PC = cpu.AddressOperand()
    case RTI: 
        cpu.S += 1
        cpu.Tick("RTI increment S")

        cpu.P = cpu.ReadFrom(0x100 + uint16(cpu.S)) | FLAG_R
        cpu.P &^= FLAG_B
        cpu.S += 1
        cpu.Tick("RTI pull P from stack, increment S")
        
        // For cycle accuracy, since pulling PC is pipelined with
        // incrementing for pulling P, do manually instead of Pull16()
        pcl := cpu.ReadFrom(0x100 + uint16(cpu.S))
        cpu.S += 1
        cpu.Tick("RTI pull PCL from stack, increment S")

        pch := cpu.ReadFrom(0x100 + uint16(cpu.S))
        cpu.Tick("RTI pull PCH from stack")

        cpu.PC = uint16(pch) << 8 + uint16(pcl)

    case RTS: cpu.PC = cpu.Pull16() + 1
        cpu.Tick("RTS increment PC")
    case BRK: cpu.PC += 1
        cpu.Push16(cpu.PC)
        // TODO: allow being interrupted by NMI, see http://nesdev.parodius.com/6502_cpu.txt
        cpu.Push(cpu.P | FLAG_B)
        cpu.SetInterrupt(true)
        cpu.PC = cpu.ReadUInt16(BRK_VECTOR)
        // TODO: refactor into same code as NMI()/IRQ()/RESET()

    case KIL:
        fmt.Printf("Halting on KIL instruction\n")
        // TODO: only halt CPU, APU might still keep playing, and may want to reset, continue seeing graphics
        os.Exit(0)

    // Blargg's CPU tests don't test these because they have inconsistent/unknown behavior
    // TODO: simulate their inconsistency, for accuracy
    case XAA, AXA, XAS, LAR:
        fmt.Printf("warning: ignoring unimplemented/inconsistent opcode: %s\n", cpu.Instruction.Opcode)

    case U__:
        fmt.Printf("halting on undefined opcode\n")  // won't happen anymore now that all are 'defined' even undocumented
        os.Exit(-1)
 
    default:
        fmt.Printf("unrecognized opcode! %s\n", cpu.Instruction.Opcode) // shouldn't happen either because should all be in table
        os.Exit(-1)
    }
    if cpu.Verbose {
        fmt.Printf("%s took %d cycles\n", cpu.Instruction.Opcode, cpu.CycleCount - startCyc)
    }

    // Finished instruction, so now can run interrupt if one is pending
    cpu.CheckInterrupt()

    // Running blargg's emulator tests?
    // TODO: post-execute instruction hook for debugging?
    // TODO: move to an automation-specific mapper, invoked on writing to 0x6000, for instr tests only
    if cpu.ReadFrom(0x6001) == 0xde && cpu.ReadFrom(0x6002) == 0xb0 && cpu.ReadFrom(0x6003) == 0x61 {
        status := cpu.ReadFrom(0x6000)
        if cpu.Verbose || status != 0x80 {
            fmt.Printf("@@ Status Code: %.2X\n", status)
            p := uint16(0x6004)
            fmt.Printf("@@ Status Text: ")
            for cpu.ReadFrom(p) != 0 {
                fmt.Printf("%s", string(cpu.ReadFrom(p)))
                p += 1
            }
            fmt.Printf("\n")
        }
        if status < 0x80  {
            os.Exit(int(status))
        }
    }
}

// Start execution
func (cpu *CPU) Run() {
    for {
        cpu.ExecuteInstruction()
    }
}

// Set interrupt to execute after next instruction finishes
func (cpu *CPU) NMI() {
    cpu.PendingInt = INT_NMI
}

func (cpu *CPU) IRQ() {
    cpu.PendingInt = INT_IRQ
}

func (cpu *CPU) Reset() {
    cpu.PendingInt = INT_RESET
}

// TODO: IRQ, it is same as NMI except uses $FFFE vector
// http://www.pagetable.com/?p=410 "Internals of BRK/IRQ/NMI/RESET on a MOS 6502"

// Check for incoming NMI, and if it is present, run it
func (cpu *CPU) CheckInterrupt() {
    if cpu.PendingInt == INT_NONE {
        return
    }

    // TODO: priorities, RESET and NMI over IRQ/BRK

    if cpu.Verbose {
        fmt.Printf("*** RUNNING INTERRUPT\n")
    }
    cpu.Push16(cpu.PC)
    cpu.Push(cpu.P)

    var vector uint16

    switch cpu.PendingInt {
    case INT_NMI: vector = NMI_VECTOR
    case INT_RESET: vector = RESET_VECTOR
    case INT_BRK, INT_IRQ: vector = BRK_VECTOR
    }

    cpu.PC = cpu.ReadUInt16(vector)
 
    cpu.PendingInt = INT_NONE
}


// Map a range of memory to call given functions on read/write
// Returns first and last of 4 KB banks mapped
func (cpu *CPU) Map(start uint16, end uint16,
        read func(address uint16) (value uint8), 
        write func(address uint16, value uint8),
        name string) (firstBankMapped int, lastBankMapped int) {

    if start & 0xfff != 0 || end & 0xfff != 0xfff {
        panic(fmt.Sprintf("Map(%.4X,%.4X): invalid memory range", start, end))
    }

    size := int(end) - int(start) + 1

    // Check to make sure the bank size makes sense
    switch size {
    case 0x10000:// 64 KB - for mapping entire address space for testing
    case 0x8000: // 32 KB   \
    case 0x4000: // 16 KB    | common in NES mappers
    case 0x2000: // 8 KB    /
    case 0x1000: // 4 KB - for NSF mapper
    default: panic(fmt.Sprintf("Map(%.4X,%.4X): invalid memory range size: %.4x", start, end, size))
    }

    count := uint16(size >> 12)

    // Map in 4 KB sections
    for i := uint16(0); i < count; i++ {
        // Map $i000 to $(i+1)fff
        base := start + (i << 12)
        index := (base & 0xf000) >> 12
        
        if i == 0 {
            firstBankMapped = int(index)
        } else if (i == count - 1) {
            lastBankMapped = int(index)
        }

        cpu.MemRead[index] = read
        cpu.MemWrite[index] = write
        cpu.MemName[index] = name
    }

    return firstBankMapped, lastBankMapped
}

// Map one address, overlaying the ordinary mappers which map at 
// multiples of 4 KB boundaries. 
//
// Potentially useful because the $4000-5FFF range is only occupied
// by the APU and I/O registers in $4000-4016, since the addresses are
// completely decoded (not mirrored); all of $4017-FFFF available for
// usage. Not sure if any cartridge uses it, but they can. Also could
// be useful for debugging or automated testing.
func (cpu *CPU) MapOver(overAddress uint16,
        readOver func(address uint16) (value uint8), 
        writeOver func(address uint16, value uint8),
        name string) {

    index := (overAddress & 0xf000) >> 12

    readOriginal := cpu.MemRead[index]
    writeOriginal := cpu.MemWrite[index]

    cpu.MemRead[index] = func(address uint16) (value uint8) {
        if address == overAddress {
            return readOver(address)
        }

        return readOriginal(address)
    }

    cpu.MemWrite[index] = func(address uint16, value uint8) {
        if address == overAddress {
            writeOver(address, value)
        } else {
            writeOriginal(address, value)
        }
    }

    cpu.MemName[index] += fmt.Sprintf(" + %.4x=%s", overAddress, name)
}

// Map memory from a ROM (at dest) into the CPU address space
// addressLinesOffset: bitmask of address lines wired from the CPU to the ROM
//  This mask is AND'd with the CPU address to get the offset within the ROM bank
// addressLinesBank: pointer to bitmask of additional ROM address lines wired to bank register
//  This mask is OR'd with the previous result to select the ROM bank, and 
//  should point to the mapper register (if any), shifted appropriately
// mapperWrite is an optional function to call when writing to the ROM, if nil, it defaults to nothing
func (cpu *CPU) MapROM(start uint16, end uint16, dest []byte, name string, addressLinesOffset uint32, addressLinesBank *uint32, mapperWrite func(uint16, uint8)) {
    // Default to ignoring wires
    if mapperWrite == nil {
        mapperWrite = func(address uint16, value uint8) { /* ignore write */ }
    }

    // Default to no bank select
    zero := uint32(0)
    if addressLinesBank == nil {
        addressLinesBank = &zero
    }

    firstBank, lastBank := cpu.Map(start, end, 
        func(cpuAddress uint16)(value uint8) { 
            romAddress := uint32(cpuAddress) & addressLinesOffset | *addressLinesBank 
            return dest[romAddress]
        },
        mapperWrite,
        name)

    // Store ROM offset for each bank, for informational purposes
    romOffsetStart := 0 & addressLinesOffset | *addressLinesBank
    for i := firstBank; i <= lastBank; i += 1 {
        cpu.MemROMOffset[i] = romOffsetStart
        cpu.MemHasROM[i] = true     // for printing ROM offset
        romOffsetStart += 0x1000    // 4 KB
    }
}

// Convert a CPU address to a ROM address
func (cpu *CPU) Address2ROM(cpuAddress uint16) (romAddress uint32, romChip string) {
    // Search through CPU addresses for ROMs that are mapped to it, live

    // For emulation purposes we map the CPU address space in 4K chunks
    // Find that chunk index for the given CPU address
    i := cpuAddress >> 12

    if !cpu.MemHasROM[i] {
        // This is a problem
        return ^uint32(0), "no ROM"
    }

    // What part of ROM it maps to, and the offset within 
    romAddrStart := cpu.MemROMOffset[i]
    romBankOffset := uint32(cpuAddress & 0xfff)

    // Absolute ROM address
    romAddress = romAddrStart | romBankOffset
    romChip = cpu.MemName[i]

    return romAddress, romChip
}
