// Created:20101128
// By Jeff Connelly

// Emulate NES CPU, a 6502-based processor
// Specifically, the NES uses a Ricoh RP2A03, which:
// - is NMOS-based 6502 (as opposed to the CMOS 650C2)
// - has no decimal mode (D flag in processor status is ignored)


package cpu6502

import . "nesfile"
import . "dis6502"
import (
    "fmt"
    "os"
)

type CPU struct {
    Memory [0x10000]uint8          // $0000-$FFFF
    PC uint16   // Program counter
    S uint8     // Stack pointer, offset from $0100
    A uint8     // Accumulator
    X, Y uint8  // Index registers
    P uint8     // Processor Status (7-0 = N V - B D I Z C)

    Cyc int     // CPU cycle counter
    CycleChannel chan int  
    PendingNMI bool // Run non-maskable interrupt after next instr finishes

    Verbose bool
    InstrTrace bool

    Instruction *Instruction  // Current instruction
    OperPtr *uint8  // Pointer to operand of current instruction
    HaveCalculatedAddress bool
    CalculatedAddress uint16

    ReadMappers [10](func(uint16) (bool, uint8))
    WriteMappers [10](func(uint16, uint8) (bool))
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

// "Master cycles" for each CPU cycle
const CPU_MASTER_CYCLES_NTSC = 15
const CPU_MASTER_CYCLES_PAL = 16


// Read byte from memory, advancing program counter
func (cpu *CPU) NextUInt8() (b uint8) {
    b = cpu.Memory[cpu.PC]
    cpu.PC += 1

    cpu.Tick(fmt.Sprintf("fetch NextUInt8 = %.2x, increment PC", b))

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

// Get the address an operand refers to
func (cpu *CPU) AddressOperand() (address uint16) {
   
    if cpu.HaveCalculatedAddress {
        return cpu.CalculatedAddress
    }

    switch cpu.Instruction.AddrMode {
    case Rel: 
        cpu.Tick("relative: fetch next opcode")
        address = (cpu.PC) + uint16(cpu.Instruction.Operand)
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
        cpu.Memory[cpu.AddressOperand()] = b

    case Abs, Abx, Aby:
        cpu.Tick("write to effective address")

        address := cpu.AddressOperand()

        switch {
        case address < 0x2000:
            // $0000-07ff is mirrored three times up to $2000, and is always RAM
            address &^= 0x1800
            cpu.Memory[address] = b
        
        case address >= 0x6000 && address <= 0x7fff:
            // Save RAM (SRAM)
            cpu.Memory[address] = b

        default:
            // Find a mapper that cares
            handled := false
            for _, mapper := range cpu.WriteMappers {
                if mapper != nil && mapper(address, b) {
                    handled = true
                    break
                }
            }
            if !handled {
                fmt.Printf("No mapper claimed write: %.4X -> %.2X, ignored\n", address, b)
            }
        }

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
        // Only can access zero page, so read directly from memory
        cpu.Tick("read from effective address")
        return cpu.Memory[cpu.AddressOperand()]

    case Abs, Abx, Aby:
        cpu.Tick("read from effective address")

        address := cpu.AddressOperand()

        switch {
        case address < 0x2000:
            address &^= 0x1800
            return cpu.Memory[address] 

        case address >= 0x6000 && address <= 0x7fff:
            return cpu.Memory[address]

        default:
            for _, mapper := range cpu.ReadMappers {
                if mapper != nil {
                    wanted, b := mapper(address)
                    if wanted {
                        return b
                    }
                }
            }
            // If there isn't any mapper, read from ROM
            return cpu.Memory[address]
        }

        return cpu.Memory[address]

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
    low := cpu.Memory[address]
    cpu.Tick("fetch ReadUInt16 low")
    high := cpu.Memory[uint16(address + 1)]
    cpu.Tick("fetch ReadUInt16 high")
    return uint16(high) << 8 + uint16(low)
}

// Read unsigned 16-bits from zero page
func (cpu *CPU) ReadUInt16ZeroPage(address uint8) (w uint16) {
    low := cpu.Memory[address]
    cpu.Tick("fetch effective address low (zero page)")
    high := cpu.Memory[uint8(address + 1)]    // 0xff + 1 wraps around
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
    low := cpu.Memory[address]
    cpu.Tick("fetch low address ReadUInt16Wraparound")
    high := cpu.Memory[(address & 0xff00) | (address + 1) & 0xff]
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
        return int(cpu.NextUInt16())
    case Imp, Acc: 
        // all others fetch next byte and use it
        cpu.Tick("read next instruction byte (and throw it away, Imp/Acc)")
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
    if len(cart.Prg) == 0 {
        panic("No PRG found")
    }
   
    var bank8000, bankC000 int

    // By default, load first PRG into 0x8000 bank and 
    // last PRG into 0xC000. This covers many mappers, including:
    // (#0) NROM-64K: one PRG loaded into both banks (http://nesdev.parodius.com/NESDoc.pdf)
    // (#0) NROM-128K: two PRGs, one loaded into each
    // MMC1, MMC3
    // TODO: Any other mappers with different default banks??
    bank8000 = 0
    bankC000 = len(cart.Prg) - 1

    switch cart.MapperCode {
    case 0: // NROM - nothing else needed
    case 1: // SxROM, MMC-1 
    // TODO


    default:
        panic(fmt.Sprintf("sorry, no support for mapper %d", cart.MapperCode))
    }

    // TODO: pointers instead of copying, so can switch banks easier
    copy(cpu.Memory[0x8000:], cart.Prg[bank8000])
    copy(cpu.Memory[0xC000:], cart.Prg[bankC000])

    // Initialize to reset vector.. note, don't use ReadUInt16 since it adds CPU cycles!
    pcl := cpu.Memory[RESET_VECTOR]
    pch := cpu.Memory[RESET_VECTOR + 1]
    cpu.PC = uint16(pch) << 8 + uint16(pcl)

    cpu.Cyc = 0
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
    b = cpu.Memory[0x100 + uint16(cpu.S)]
    cpu.Tick("pull stack")
    return b
}

// Pull 16 bits from stack
func (cpu *CPU) Pull16() (w uint16) {
    // pipelining causes this to take one fewer cycles than two 8-bit Pull()s
    cpu.S += 1
    cpu.Tick("increment S")

    low := cpu.Memory[0x100 + uint16(cpu.S)]
    cpu.S += 1
    cpu.Tick("pull low from stack, increment S")

    high := cpu.Memory[0x100 + uint16(cpu.S)]
    cpu.Tick("pull high from stack")

    return uint16(high) * 0x100 + uint16(low)
}

// Push 8 bits onto stack
func (cpu *CPU) Push(b uint8) {
    cpu.Memory[0x100 + uint16(cpu.S)] = b 
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
    cpu.Cyc += 1

    // TODO: remove check, but test suites don't setup channel
    if cpu.CycleChannel != nil {
        cpu.CycleChannel <- CPU_MASTER_CYCLES_NTSC   // TODO: or PAL
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
    startCyc := cpu.Cyc
    cpu.HaveCalculatedAddress = false
    cpu.Instruction = cpu.NextInstruction()

    // Instruction trace
    if cpu.InstrTrace || cpu.Verbose {
        fmt.Printf("%.4X  ", start)
        fmt.Printf("%.2X ", cpu.Instruction.OpcodeByte)
        if cpu.Instruction.AddrMode.OperandSize() >= 1 {
           fmt.Printf("%.2X ", cpu.Memory[start + 1])
        } else {
           fmt.Printf("   ")
        }

        if cpu.Instruction.AddrMode.OperandSize() >= 2 {
           fmt.Printf("%.2X ", cpu.Memory[start + 2])
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
    case AAC: cpu.A &= cpu.ReadOperand(); cpu.SetSZ(cpu.A); cpu.SetCarry(cpu.P & FLAG_Z != 0)
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
    case ASR: cpu.A &= cpu.ReadOperand()
        cpu.A >>= 1
        cpu.SetCarry(cpu.A & 0x01 != 0)
        cpu.SetSZ(cpu.A)
    //case ARR: cpu.A &= cpu.ReadOperand(); cpu.OpROR()
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

        cpu.P = cpu.Memory[0x100 + uint16(cpu.S)] | FLAG_R
        cpu.P &^= FLAG_B
        cpu.S += 1
        cpu.Tick("RTI pull P from stack, increment S")
        
        // For cycle accuracy, since pulling PC is pipelined with
        // incrementing for pulling P, do manually instead of Pull16()
        pcl := cpu.Memory[0x100 + uint16(cpu.S)]
        cpu.S += 1
        cpu.Tick("RTI pull PCL from stack, increment S")

        pch := cpu.Memory[0x100 + uint16(cpu.S)]
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

    case KIL:
        // The known-correct log http://nickmass.com/images/nestest.log ends at $C66E, but my
        // emulator interprets the last RTS as returning to 0001, executing a bunch of weird instructions, ending in KIL
        // So interpret this as a graceful termination. Only up to $C66E will be compared by tracediff.pl anyways.
        fmt.Printf("Halting on KIL instruction\n")
        // http://nesdev.com/bbs/viewtopic.php?t=7130
        result := cpu.Memory[2] << 8 | cpu.Memory[3]
        if result == 0 {
            fmt.Printf("Nestest automation: Pass\n")
        } else {
            fmt.Printf("Nestest automation: FAIL with code %.4x\n", result)
        }
        os.Exit(0)


    // UNIMPLEMENTED INSTRUCTIONS
    // nestest.nes PC=$c000 doesn't test these, so I didn't implement them
    case ARR, ATX, AXA, AXS, LAR, SXA, SYA, XAA, XAS:
        fmt.Printf("unimplemented opcode: %s\n", cpu.Instruction.Opcode)
        os.Exit(-1)

    case U__:
        fmt.Printf("halting on undefined opcode\n")  // won't happen anymore now that all are 'defined' even undocumented
        os.Exit(-1)
 
    default:
        fmt.Printf("unrecognized opcode! %s\n", cpu.Instruction.Opcode) // shouldn't happen either because should all be in table
        os.Exit(-1)
    }
    if cpu.Verbose {
        fmt.Printf("%s took %d cycles\n", cpu.Instruction.Opcode, cpu.Cyc - startCyc)
    }

    // Finished instruction, so now can run NMI
    cpu.CheckNMI()

    // Running blargg's emulator tests?
    // TODO: post-execute instruction hook for debugging?
    if cpu.Memory[0x6001] == 0xde && cpu.Memory[0x6002] == 0xb0 && cpu.Memory[0x6003] == 0x61 {
        status := cpu.Memory[0x6000]
        if cpu.Verbose || status != 0x80 {
            fmt.Printf("@@ Status Code: %.2X\n", status)
            p := 0x6004
            fmt.Printf("@@ Status Text: ")
            for cpu.Memory[p] != 0 {
                fmt.Printf("%s", string(cpu.Memory[p]))
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

// Set NMI to execute after next instruction finishes
func (cpu *CPU) NMI() {
    cpu.PendingNMI = true
}

// Check for incoming NMI, and if it is present, run it
func (cpu *CPU) CheckNMI() {
    if !cpu.PendingNMI {
        return
    }

    if cpu.Verbose {
        fmt.Printf("*** RUNNING NMI\n")
    }
    cpu.Push16(cpu.PC)
    cpu.Push(cpu.P)
    cpu.PC = cpu.ReadUInt16(NMI_VECTOR)
 
    cpu.PendingNMI = false
}
