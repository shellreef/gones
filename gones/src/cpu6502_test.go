// Created:20101209
// By Jeff Connelly

// Test 6502 CPU

package cpu6502_test

import (
    "testing"
    "fmt"

    "dis6502"
    )

import . "cpu6502"

// Measure the cycle count of an operation
func cyclesForOp(op uint8) (int) {
    cpu := new(CPU)
    cpu.PowerUp()
    cpu.PC = 0x8000
    cpu.Memory[0x8000] = op
    // Note: cycle count may vary depending on operands.. they're all zeros here
    startCyc := cpu.Cyc
    cpu.ExecuteInstruction()
    cycleCount := cpu.Cyc - startCyc

    return cycleCount
}

type OpcodeByteCycleCount struct {
    OpcodeByte uint8;
    Cycles int;
}

// Based on http://nesdev.parodius.com/6502_cpu.txt
// except using more descriptive mneumonics for 
// multi-byte NOPs: DOP for zero page, TOP for absolute
// and also http://www.obelisk.demon.co.uk/6502/reference.html
// Note: these only test if the page boundary is not crossed
var Timings = [...]OpcodeByteCycleCount{
    // Instructions accessing the stack
    {0x00, 7},  // BRK
    {0x40, 6},  // RTI
    {0x60, 6},  // RTS
    {0x48, 3},  // PHA
    {0x08, 3},  // PHP
    {0x68, 4},  // PLA
    {0x28, 4},  // PLP
    {0x20, 6},  // JSR

    // Accumulator or implied addressing
    // TODO: all
    {0x80, 2},  // NOP
    {0xa0, 2},  // LDY
    {0xc0, 2},  // CPY
    {0xe0, 2},  // CPX
    {0x88, 2},  // DEY
    {0xa8, 2},  // TAY
    {0xc8, 2},  // INY
    {0xe8, 2},  // INX
    {0x0a, 2},  // ASL A
    {0x2a, 2},  // ROL A
    {0x4a, 2},  // LSR A
    {0x6a, 2},  // ROR A
    {0x8a, 2},  // TXA
    {0xaa, 2},  // TAX
    {0xca, 2},  // DEX
    {0xea, 2},  // NOP
    {0x18, 2},  // CLC
    {0x38, 2},  // SEC
    {0x58, 2},  // CLI
    {0x78, 2},  // SEI
    {0x98, 2},  // TYA
    {0xb8, 2},  // CLV
    {0xd8, 2},  // CLD
    {0xf8, 2},  // SED
    {0x1a, 2},  // NOP
    {0x3a, 2},  // NOP
    {0x5a, 2},  // NOP
    {0x7a, 2},  // NOP
    {0x9a, 2},  // TXS
    {0xba, 2},  // TSX
    {0xda, 2},  // NOP
    {0xfa, 2},  // NOP
    {0x69, 2},  // ADC
   
    // Immediate addressing
    // TODO: all

    // Absolute addressing
    // Read 
    {0x4c, 3},  // JMP
    {0xad, 4},  // LDA
    {0xae, 4},  // LDX
    {0xac, 4},  // LDY
    {0x4d, 4},  // EOR
    {0x2d, 4},  // AND
    {0x0d, 4},  // ORA
    {0x6d, 4},  // ADC
    {0xed, 4},  // SBC
    {0xcd, 4},  // CMP
    {0x2c, 4},  // BIT
    {0xaf, 4},  // LAX
    {0x0c, 4},  // TOP
    // Read-modify-write
    {0x0e, 6},  // ASL
    {0x4e, 6},  // LSR
    {0x2e, 6},  // ROL
    {0x6e, 6},  // ROR
    {0xee, 6},  // INC
    {0xce, 6},  // DEC
    {0x0f, 6},  // SLO
    {0x4f, 6},  // SRE
    {0x2f, 6},  // RLA
    {0x6f, 6},  // RRA
    {0xef, 6},  // ISC
    {0xcf, 6},  // DCP
    // Write
    {0x8d, 4},  // STA
    {0x8e, 4},  // STX
    {0x8c, 4},  // STY
    {0x8f, 4},  // SAX

    // Zero page addressing
    // These are very similar to absolute, but take one fewer cycle because
    // the high-order byte of the operand doesn't have to be fetched (it is zero)
    // Read
    {0xa5, 3},  // LDA
    {0xa6, 3},  // LDX
    {0xa4, 3},  // LDY
    {0x45, 3},  // EOR
    {0x25, 3},  // AND
    {0x05, 3},  // ORA
    {0x65, 3},  // ADC
    {0xe5, 3},  // SBC
    {0xc5, 3},  // CMP
    {0x24, 3},  // BIT
    {0xa7, 3},  // LAX
    {0x04, 3},  // DOP
    // Read-modify-write
    {0x06, 5},  // ASL
    {0x46, 5},  // LSR
    {0x26, 5},  // ROL
    {0x66, 5},  // ROR
    {0xe6, 5},  // INC
    {0xc6, 5},  // DEC
    {0x07, 5},  // SLO
    {0x47, 5},  // SRE
    {0x27, 5},  // RLA
    {0x67, 5},  // RRA
    {0xe7, 5},  // ISC
    {0xc7, 5},  // DCP
    // Write
    {0x85, 3},  // STA
    {0x86, 3},  // STX
    {0x84, 3},  // STY
    {0x87, 3},  // SAX

    // Zero page indexed
    // X        Y 
    // Read
    {0xb5, 4},              // LDA
                {0xb6, 4},  // LDX
    {0xb4, 4},              // LDY
    {0x55, 4},              // EOR
    {0x35, 4},              // AND
    {0x15, 4},              // ORA
    {0x75, 4},              // ADC
    {0xf5, 4},              // SBC
    {0xd5, 4},              // CMP
                {0xb7, 4},  // LAX
    {0x14, 4},              // DOP
    {0x34, 4},              // DOP
    {0x54, 4},              // DOP
    {0x74, 4},              // DOP
    {0xd4, 4},              // DOP
    {0xf4, 4},              // DOP
    // Read-modify-write
    {0x16, 6},              // ASL
    {0x56, 6},              // LSR
    {0x36, 6},              // ROL
    {0x76, 6},              // ROR
    {0xf6, 6},              // INC
    {0xd6, 6},              // DEC
    {0x17, 6},              // SLO
    {0x57, 6},              // SRE
    {0x37, 6},              // RLA
    {0x77, 6},              // RRA
    {0xf7, 6},              // ISB
    {0xd7, 6},              // DCP
    // Write
    {0x95, 4},              // STA
                {0x96, 4},  // STX
    {0x94, 4},              // STY
                {0x97, 4},  // SAX

    // Absolute indexed
    // Read
    {0xbd, 4},  {0xb9, 4},  // LDA
                {0xbe, 4},  // LDX
    {0xbc, 4},              // LDY
    {0x5d, 4}, {0x59, 4},   // EOR
    {0x3d, 4}, {0x39, 4},   // AND
    {0x1d, 4}, {0x19, 4},   // ORA
    {0x7d, 4}, {0x79, 4},   // ADC
    {0xfd, 4}, {0xf9, 4},   // SBC
    {0xdd, 4}, {0xd9, 4},   // CMP
               {0xbf, 4},   // LAX
    // TODO: nothing called LAE?
    // TODO: SHS aka XAS
    {0x1c, 4},              // TOP
    {0x3c, 4},              // TOP
    {0x5c, 4},              // TOP
    {0x7c, 4},              // TOP
    {0xdc, 4},              // TOP
    {0xfc, 4},              // TOP
    // Read-modify-write
    //{0x1e, 7},              // ASL
    // TODO: rest
    //{0x99, 5},              // STA
    // TODO: more
}

func TestTiming(t *testing.T) {
    for _, pair := range Timings {
        opcodeByte := pair.OpcodeByte
        expected := pair.Cycles

        actual := cyclesForOp(opcodeByte)

        disasm := dis6502.Opcodes[opcodeByte]

        if actual != expected {
            t.Errorf("FAIL %.2X: %s %d != %d (%d)", opcodeByte, disasm, actual, expected, actual - expected)
        } else {
            fmt.Printf("Pass %.2X: %s %d == %d\n", opcodeByte, disasm, actual, expected);
        }
    }
}
