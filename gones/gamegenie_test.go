// Created:20101205
// By Jeff Connelly

// Test

package gamegenie

import (
    "testing"
)

// Fail if condition is false
func expect(t *testing.T, f bool) {
    if !f {
        t.Fail()
    }
}

// 6-letter code
func TestDecode6(t *testing.T) {
    c := Decode("SLZPLO")
    expect(t, c.Address == 0x1123)
    expect(t, c.Value == 0xbd)
    expect(t, !c.HasKey)
    expect(t, !c.WantsKey)
}

// 8-letter code
func TestDecode8(t *testing.T) {
    c := Decode("SLXPLOVS")
    expect(t, c.Address == 0x1123)
    expect(t, c.Value == 0xbd)
    expect(t, c.HasKey)
    expect(t, c.WantsKey)
    expect(t, c.Key == 0xde)
}

// 6-letter code that should be 8-letter. These work as expected, but
// in the Game Genie ROM don't advance to the next entry line.
func TestDecodeMissingKey(t *testing.T) {
    c := Decode("SLXPLO")
    expect(t, !c.HasKey)
    expect(t, c.WantsKey)
}

// 8-letter code that should be 6-letter. Game Genie won't let you
// enter these, but we allow them
func TestDecodeExtraKey(t *testing.T) {
    c := Decode("SLZPLOVS")
    expect(t, c.HasKey)
    expect(t, !c.WantsKey)
}

func TestEncode(t *testing.T) {
}
