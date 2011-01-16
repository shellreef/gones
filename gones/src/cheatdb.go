// Created:20110115
// By Jeff Connelly

// Cheat code database

package cheatdb

import ("fmt"
        "xml"
        "os"

        "gosqlite.googlecode.com/hg/sqlite"    // http://code.google.com/p/gosqlite/
        )

type Cheats struct {
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
    ROMAddress string "attr"
    ROMBefore string "attr"
    ROMAfter string "attr"
}

// Stringify an effect to how it appears in the Game Genie manual, i.e.:
// "SLXPLOVS Infinite lives for Mario(tm) and Luigi(tm)"
// "PEEPUZAG + IUEPSZAA + TEEPVZPA Start on World 2"
// "AAVENYZA / AAVEUYZA Weak Birdetta"
func (effect Effect) String() (string) {
    codeText := ""
    lastApplies := ""
    for i, code := range effect.Code {
        if i != 0 {
            if code.Applies != "" && lastApplies != "" && code.Applies != lastApplies {
                // Alternate version
                codeText += " / "
            } else {
                // Multiple codes per effect
                codeText += " + "
            }
        }
        codeText += code.Genie
        if code.Applies != "" {
            codeText += "(" + code.Applies + ")"
        }

        if code.Applies != "" {
            lastApplies = code.Applies
        }
    }

    return codeText + " " + effect.Title
}

func Load() (Cheats) {
    filename := "../genie/galoob/galoob.xml"
    r, err := os.Open(filename, os.O_RDONLY, 0)
    if r == nil {
        panic(fmt.Sprintf("cannot open %s: %s", filename, err))
    }

    // TODO: read Nestopia's NstCheat files, like those found on http://www.mightymo.net/downloads.html
    // TODO: read Nesticle .pat files, [+]<code> [<name>] per line, raw format. http://www.zophar.net/faq/nitrofaq/nesticle.readme.txt and some at http://jeff.tk:81/consoles/
    cheats := Cheats{}

    xml.Unmarshal(r, &cheats)

    return cheats
}

func (cheats Cheats) Save() {
    db, err := sqlite.Open("/tmp/foo.db")
    if err != nil {
        panic(fmt.Sprintf("cheatdb.Save(): failed to save"))
    }

    db.Exec(`CREATE TABLE game( -- an abstract "game", has ≥1 carts
        id INTEGER PRIMARY KEY, 
        name TEXT NOT NULL, 
        galoob_id TEXT NULL, 
        galoob_name TEXT NULL
        )`)

    db.Exec(`CREATE TABLE cart( -- a physical cartridge, possibly different versions
        id INTEGER PRIMARY KEY, 
        game_id INTEGER NOT NULL, 
        sha1 TEXT NOT NULL, 
        name TEXT NULL,         -- PRG0, PRG1, etc. if multiple versions, otherwise NULL

        FOREIGN KEY(game_id) REFERENCES game(id)
        )`)
    db.Exec(`CREATE TABLE code( -- a decoded code
        id INTEGER PRIMARY KEY, 
        cart_id INTEGER NOT NULL,  
        cpu_addr INTEGER NOT NULL, 
        value INTEGER NOT NULL, 
        compare INTEGER NULL,   -- compare byte if 8-letter code, or NULL for 6-letter
        -- TODO: what if this is a 6-letter code affecting multiple ROM address??
        rom_addr INTEGER NULL,  -- corresponding ROM address if known, or NULL
        rom_before BLOB NULL,   -- data before the affected address, if known
        rom_after BLOB NULL,    -- data after _and_including_ the affected address

        FOREIGN KEY(cart_id) REFERENCES cart(id)
        )`)
    os.Exit(0)
}

// TODO: web interface
