// Created:20101209
// By Jeff Connelly

// Test 6502 CPU

package cpu6502

import (
    "testing"
    "fmt"
    )

import . "dis6502"

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

func TestTimingStack(t *testing.T) {
    op := OpcodeByteFor(BRK, Imp)

    actual := cyclesForOp(op)
    expected := 7

    if actual != expected {
        t.Errorf("BRK took %d, expected %d", actual, expected)
    } else {
        fmt.Printf("BRK %d == %d\n", actual, expected);
    }
    // TODO: all others http://nesdev.parodius.com/6502_cpu.txt
}
