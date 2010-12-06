// Created:20101205
// By Jeff Connelly

// Test

package gamegenie

import (
    "testing"
)

func TestDecodeKey(t *testing.T) {
    c := Decode("SLXPLOVS")
    if c.Address != 0x1123 {
        t.Errorf("DecodeKey: bad address: %.4x\n", c.Address)
    }
    if c.Value != 0xbd {
        t.Errorf("DecodeKey: bad value: %.2x\n", c.Value)
    }
    if c.Key != 0xde {
        t.Errorf("DecodeKey: bad key: %.2x\n", c.Key)
    }
}

func TestDecodeMissingKey(t *testing.T) {
    c := Decode("SLXPLO")
    if c.Address != 0x1123 {
        t.Errorf("DecodeMissingKey: bad address: %.4x\n", c.Address)
    }
    if c.Value != 0xbd {
        t.Errorf("DecodeMissingKey: bad value: %.2x\n", c.Value)
    }
}

func TestEncode(t *testing.T) {
}
