// Created:20101209
// By Jeff Connelly

// Test 6502 CPU

package cpu6502

import (
    "testing"
    "fmt"

    "dis6502"
    )

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

}

func TestTiming(t *testing.T) {
    for _, pair := range Timings {
        opcodeByte := pair.OpcodeByte
        expected := pair.Cycles

        actual := cyclesForOp(opcodeByte)

        disasm := dis6502.Opcodes[opcodeByte]

        if actual != expected {
            t.Errorf("FAIL: %s %d != %d (%d)", disasm, actual, expected, actual - expected)
        } else {
            fmt.Printf("Pass: %s %d == %d\n", disasm, actual, expected);
        }
    }
}
