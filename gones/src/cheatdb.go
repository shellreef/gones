// Created:20110115
// By Jeff Connelly

// Cheat code database

package cheatdb

import ("fmt"
        "xml"
        "os"
        )

type Games struct {
    XMLName xml.Name "cheats"

    Game []Game
}

type Game struct {
    Galoob_id string "attr"
    Galoob_name string "attr"
    Name string "attr"
    Cartridge []Cartridge 
    Effect []Effect
}

type Cartridge struct {
    Filename string "attr"
    SHA1 string "attr"
    Name string "attr"
}

type Effect struct {
    Number string "attr"
    Source string "attr"
    Title string "attr"
    Codes []Code
}

type Code struct {
    Genie string "attr"
    Applies string "attr"       // refers to a Cartridge.Name
}

func Load() {
    filename := "../genie/galoob/galoob.xml"
    r, err := os.Open(filename, os.O_RDONLY, 0)
    if r == nil {
        panic(fmt.Sprintf("cannot open %s: %s", filename, err))
    }

    games := Games{}

    xml.Unmarshal(r, &games)

    fmt.Printf("games = %d, %s\n", len(games.Game), games)
}
