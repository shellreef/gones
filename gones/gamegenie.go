// Created:20101205
// By Jeff Connelly

// Game Genie codes

// TODO: optionally allow using Codemaster's official Game Genie ROM to enter codes: http://nesdev.parodius.com/genie.zip

// http://nesdev.com/nesgg.txt

package GameGenie

import (
    "fmt"
    "strings"
)

type GameGenieCode struct {
    Address uint16
    Value uint8
    HasKey bool
    WantsKey bool
    Key uint8
}

// Position in this string corresponds to numerical value of code
const LETTERS = "APZLGITYEOXUKSVN"


// Get Game Genie numerical hex digit of a letter, 0-15
func letterToDigit(letter string) (digit int) {
    digit = strings.Index(LETTERS, strings.ToUpper(letter))
    if digit == -1 {
        panic(fmt.Sprintf("letterToValue(%s): invalid", letter))
        // TODO: error handling
    }
    return digit
}

// Decode a game Genie Code into its consistutent fields
func Decode(s string) (c GameGenieCode) {
    var digits [8]int
    for i, characterCode := range s {
        letter := string(characterCode)
        digits[i] = letterToDigit(letter)
    }

    // Unscramble digits
    c.Value = uint8(((digits[0]&8)<<4)+((digits[1]&7)<<4)+(digits[0]&7))
    c.Address = uint16(((digits[3]&7)<<12)+((digits[4]&8)<<8)+((digits[5]&7)<<8)+
            ((digits[1]&8)<<4)+((digits[2]&7)<<4)+(digits[3]&8)+(digits[4]&7))
    if len(s) == 8 {
        c.Value += uint8(digits[7] & 8)
        c.Key = uint8(((digits[6]&8)<<4)+((digits[7]&7)<<4)+(digits[5]&8)+(digits[6]&7))
        c.HasKey = true
    } else {
        c.Value += uint8(digits[5] & 8)
    }

    // Codes like this don't automagically go to the next line when the sixth 
    // letter is typed. They should be 8 letters, but if you only type 6, they'll
    // still take affect. This only matters for code entry.
    // NOTE: When patching, use HasKey to tell whether to apply Key! 
    // .. not WantsKey, because the code doesn't always get what it wants.
    c.WantsKey = digits[2] >> 3 != 0

    return c
}

// String representation of hex addresses and value
func (c GameGenieCode) String() (string) {
    if c.HasKey {
        return fmt.Sprintf("%.4X:%.2X?%.2X", c.Address, c.Value, c.Key)
    } 
    return fmt.Sprintf("%.4X:%.2X", c.Address, c.Value)
}

// Encode to Game Genie
func (c GameGenieCode) Encode() (s string) {
    var digits [8]uint8

    digits[0]=uint8((c.Value&7)+((c.Value>>4)&8))
    digits[1]=uint8(((c.Value>>4)&7)+uint8((c.Address>>4)&8))
    digits[2]=uint8(((c.Address>>4)&7))
    digits[3]=uint8((c.Address>>12)+(c.Address&8))
    digits[4]=uint8((c.Address&7)+((c.Address>>8)&8))
    digits[5]=uint8(((c.Address>>8)&7))

    var length int
    if c.HasKey {
        digits[2]+=8
        digits[5]+=uint8(c.Key&8)
        digits[6]=uint8((c.Key&7)+((c.Key>>4)&8))
        digits[7]=uint8(((c.Key>>4)&7)+(c.Value&8))
        length = 8
    }  else {
        digits[5]+=uint8(c.Value&8)
        length = 6
    }
    // TODO: encoding WantsKey without HasKey?
    // Currently, this function will only encode canonicalized
    // codes, so that for example SLXPLO which has WantsKey set
    // but is only 6 digits will be, if decoded and re-encoded,
    // changed to SLZPLO which is more correct, and will in
    // a real Game Genie automatically advance to the next code entry line.

    for i := 0; i < length; i += 1 {
        s += string(LETTERS[digits[i]])
    }

    return s
}
