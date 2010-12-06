// Created:20101205
// By Jeff Connelly

// Test

package gamegenie

import (
    "testing"
    "fmt"
)

// 6-letter code
func TestDecode6(t *testing.T) {
    c := Decode("SLZPLO")
    if c.Address != 0x1123 {
        t.Errorf("Decode6: bad address: %.4x\n", c.Address)
    }
    if c.Value != 0xbd {
        t.Errorf("Decode6: bad value: %.2x\n", c.Value)
    }
    if c.HasKey {
        t.Errorf("Decode6: has key, even though is 6 letters")
    }
    if c.WantsKey {
        t.Errorf("Decode6: wants key, even though shouldn't")
    }
}

// 8-letter code
func TestDecode8(t *testing.T) {
    c := Decode("SLXPLOVS")
    if c.Address != 0x1123 {
        t.Errorf("Decode8: bad address: %.4x\n", c.Address)
    }
    if c.Value != 0xbd {
        t.Errorf("Decode8: bad value: %.2x\n", c.Value)
    }
    if !c.HasKey {
        t.Errorf("Decode8: doesn't have key, even though is 6 letters")
    }
    if !c.WantsKey {
        t.Errorf("Decode8: doesn't want key, even though should")
    }
    if c.Key != 0xde {
        t.Errorf("Decode8: bad key: %.2x\n", c.Key)
    }
}

// 6-letter code that should be 8-letter
func TestDecodeMissingKey(t *testing.T) {
    c := Decode("SLXPLO")
    if c.HasKey {
        t.Errorf("DecodeMissingKey: has key, even though is 6 letters")
    }
    if !c.WantsKey {
        t.Errorf("DecodeMissingKey: doesn't want key, but bit 3 of byte 3 is set")
    }
}

func TestDecodeExtraKey(t *testing.T) {
    c := Decode("SLZPLOVS")
    fmt.Printf("code=%s\n", c)
}

func TestEncode(t *testing.T) {
}
