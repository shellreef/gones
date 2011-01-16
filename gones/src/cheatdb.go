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


    // Code info
    db.Exec(`CREATE TABLE game(     -- an abstract "game", has ≥1 carts
        id INTEGER PRIMARY KEY, 
        name TEXT NOT NULL, 
        galoob_id TEXT NULL, 
        galoob_name TEXT NULL
        )`)
    db.Exec(`CREATE TABLE cart(     -- a physical cartridge, possibly different versions
        id INTEGER PRIMARY KEY, 
        game_id INTEGER NOT NULL, 
        sha1 TEXT NOT NULL, 
        name TEXT NULL,             -- PRG0, PRG1, etc. if multiple versions, otherwise NULL

        FOREIGN KEY(game_id) REFERENCES game(id)
        )`)
    db.Exec(`CREATE TABLE effect(
        id INTEGER PRIMARY KEY,
        cart_id INTEGER NOT NULL,  
        title TEXT NOT NULL,
        hacker_id INTEGER NOT NULL,
        create_date DATETIME NOT NULL,

        FOREIGN KEY(hacker_id) REFERENCES user(id),
        FOREIGN KEY(cart_id) REFERENCES cart(id)
        )`)
    db.Exec(`CREATE TABLE code(     -- a decoded Game Genie code, has ≥0 patches
        id INTEGER PRIMARY KEY, 
        effect_id INTEGER NOT NULL,
        cpu_address INTEGER NOT NULL, -- the address referred to, $8000-FFFF
        value INTEGER NOT NULL, 
        compare INTEGER NULL,       -- compare byte if 8-letter code, or NULL for 6-letter

        FOREIGN KEY(effect_id) REFERENCES effect(id)
        )`)
    db.Exec(`CREATE TABLE patch(    -- an actual known modification to ROM
        id INTEGER PRIMARY KEY,
        code_id INTEGER NOT NULL,
        rom_address INTEGER NOT NULL,
        rom_before BLOB NOT NULL,   -- bytes in the ROM before the affected address, for context
        rom_after BLOB NOT NULL,    -- bytes _including_ and after the affected address

        FOREIGN KEY(code_id) REFERENCES code(id)
        )`)

    // Extra info
    db.Exec(`CREATE TABLE user(     -- someone that created codes or commented on them
        id INTEGER PRIMARY KEY,
        fullname TEXT NOT NULL
        )`)
    db.Exec(`CREATE TABLE comment(  -- user commentary on an effect
        id INTEGER PRIMARY KEY,
        effect_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        message TEXT NOT NULL,
        create_date DATETIME NOT NULL,

        FOREIGN KEY(effect_id) REFERENCES effect(id),
        FOREIGN KEY(user_id) REFERENCES user(id)
        )`)



    os.Exit(0)
}

// TODO: web interface
