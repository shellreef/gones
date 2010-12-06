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
    expect(t, c.String() == "1123:BD")
    expect(t, c.Encode() == "SLZPLO")
}

// 8-letter code
func TestDecode8(t *testing.T) {
    c := Decode("SLXPLOVS")
    expect(t, c.Address == 0x1123)
    expect(t, c.Value == 0xbd)
    expect(t, c.HasKey)
    expect(t, c.WantsKey)
    expect(t, c.Key == 0xde)
    expect(t, c.String() == "1123:BD?DE")
    expect(t, c.Encode() == "SLXPLOVS")
}

// 6-letter code that should be 8-letter. These work as expected, but
// in the Game Genie ROM don't advance to the next entry line.
func TestDecodeMissingKey(t *testing.T) {
    c := Decode("SLXPLO")
    expect(t, !c.HasKey)
    expect(t, c.WantsKey)
    expect(t, c.Encode() == "SLZPLO") // converted to proper 6-letter
}

// 8-letter code that should be 6-letter. Game Genie won't let you
// enter these, but we allow them
func TestDecodeExtraKey(t *testing.T) {
    c := Decode("SLZPLOVS")
    expect(t, c.HasKey)
    expect(t, !c.WantsKey)
    expect(t, c.Encode() == "SLXPLOVS") // converted to proper 8-letter
}

func TestDecodePatch(t *testing.T) {
    expect(t, DecodePatch("1123:BD?DE").Encode() == "SLXPLOVS")
    expect(t, DecodePatch("1123:BD").Encode() == "SLZPLO")
}

// Test encoding from scratch
func TestEncode(t *testing.T) {
    var c GameGenieCode 
    c.Address = 0
    c.Value = 0
    expect(t, c.Encode() == "AAAAAA")
}

