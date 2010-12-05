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

// Get Game Genie value of a letter, 0-15
func letterToValue(letter string) (value int) {
    value = strings.Index(LETTERS, strings.ToUpper(letter))
    if value == -1 {
        panic(fmt.Sprintf("letterToValue(%s): invalid", letter))
        // TODO: error handling
    }
    return value
}

func Decode(s string) (c GameGenieCode) {
    for _, characterCode := range s {
        letter := string(characterCode)

        fmt.Printf("letter=%s, value=%d\n", letter, letterToValue(letter))
    }

    return c
}
