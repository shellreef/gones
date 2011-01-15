// Created:20110115
// By Jeff Connelly

// Cheat code database

package cheatdb

import ("fmt"
        "xml"
        "os"
        )

type Games struct {
    XMLName xml.Name "cheat-codes"

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
    Code []Code
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

    for _, game := range games.Game {
        fmt.Printf("\n%s\n", game.Name)
        for _, effect := range game.Effect {
            for _, code := range effect.Code {
                fmt.Printf("%s ", code.Genie)
            }
            fmt.Printf("%s\n", effect.Title)
        }
    }
}

// TODO: read Nestopia's NstCheat files, like those found on http://www.mightymo.net/downloads.html
// TODO: read Nesticle .pat files, [+]<code> [<name>] per line, raw format. http://www.zophar.net/faq/nitrofaq/nesticle.readme.txt and some at http://jeff.tk:81/consoles/
