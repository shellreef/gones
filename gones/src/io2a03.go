// Created:20110103
// By Jeff Connelly
//
// NES RP2A03 emulation, except the embedded 6502 (see cpu6502)
// Specifically, only I/O and extra component besides the core
//
// http://nesdev.parodius.com/2A03%20technical%20reference.txt

package io2a03

// I/O ports
// http://wiki.nesdev.com/w/index.php/2A03

// Sound registers not here yet

// Joystick input
const IO_JOY1           = 0x4016
const IO_JOY2           = 0x4017
const LAST_REGISTER     = IO_JOY2

// Bits for state of each button
// http://wiki.nesdev.com/w/index.php/Standard_controller
const (
        BUTTON_A=0 << iota
        BUTTON_B
        BUTTON_SELECT
        BUTTON_START
        BUTTON_UP
        BUTTON_DOWN
        BUTTON_LEFT
        BUTTON_RIGHT
        // SNES controller can be connected to read out additional buttons, too
      )

type IO struct {
    ButtonState uint8
    ControllerShiftReg uint32   // TODO: shift register object
}

func ReadRegister(address uint16) (value uint8) {
    if address > LAST_REGISTER {
        return 0
    }

    // TODO
    switch address {
    case IO_JOY1: 
        // TODO: shift out bits of register
    case IO_JOY2: 
    }

    return value
}

func WriteRegister(address uint16, value uint8) {
    if address > LAST_REGISTER {
        return
    }

    switch address {
    case IO_JOY1:
        if value & 1 == 1 {
            // TODO: store state of each button
        } else {
            // TODO: allow buttons to be read back
        }
    }


    return
}
