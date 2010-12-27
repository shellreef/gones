// Created:20101226
// By Jeff Connelly
//
// Read NES cartridge database
// XML from http://bootgod.dyndns.org:7777/

package main

import ("fmt"
        "os"
        "xml"
        "gob"
        "time"
        )

// See schema http://bootgod.dyndns.org:7777/downloads/nesdb.xsd

type PRG struct {
    Name string "attr"
    Size string "attr"  // like "256k", TODO: integer bytes
    CRC string "attr"
    SHA1 string "attr"
}

type CHR struct {
    Name string "attr"
    Size string "attr"  // like "256k", TODO: integer bytes
    CRC string "attr"
    SHA1 string "attr"
}

type VRAM struct {
    Size string "attr"
}

type WRAM struct {
    Size string "attr"
}

type Pad struct {
    H string "attr"
    V string "attr"
}

type Chip struct {
    Type string "attr"
}

type CIC struct {
    Type string "attr"
}

type Board struct {
    Type string "attr"
    PCB string "attr"
    Mapper string "attr"
    PRG []PRG
    CHR []CHR
    VRAM []VRAM
    WRAM []WRAM
    Pad []Pad
    CIC []CIC
    Chip []Chip
}

type Cartridge struct {
    System string "attr"
    CRC string "attr"
    SHA1 string "attr"
    Dump string "attr"
    Dumper string "attr"
    Revision string "attr"
    DateDumped string "attr"  // TODO: date format
    Board []Board
}

type Game struct {
    Name string "attr"
    AltName string "attr"
    Class string "attr"
    Catalog string "attr"
    Publisher string "attr"
    Developer string "attr"
    Region string "attr"
    Players string "attr"
    Date string "attr"       // TODO: date format

    Cartridge []Cartridge
}

type Database struct {
    XMLName xml.Name "database"
    
    Game []Game
}

func ReadXML() (*Database) {
    filename := "../NesCarts (2010-11-29).xml"
    r, err := os.Open(filename, os.O_RDONLY, 0)
    if r == nil {
        panic(fmt.Sprintf("cannot open %s: %s", filename, err))
    }

    database := Database{}

    start := time.Nanoseconds()
    xml.Unmarshal(r, &database)
    took := time.Nanoseconds() - start

    fmt.Printf("Loaded %d games in %.4f s\n", len(database.Game), float(took) / 1e9)

    return &database
}

func Dump(database *Database) {
    for i, game := range database.Game {
        fmt.Printf("#%d. [%s] %s - %s\n", i, game.Region, game.Developer, game.Name)
        for _, cart := range game.Cartridge {
            rev := ""
            if len(cart.Revision) != 0 {
                rev = fmt.Sprintf(" (rev. %s)", cart.Revision)
            } 
            fmt.Printf("\tCartridge%s for %s\n", rev, cart.System)
            for _, board := range cart.Board {
                fmt.Printf("\tBoard: %s (%s)\n", board.PCB, board.Type)
                for _, chip := range board.Chip { 
                    fmt.Printf("\t\tChip: %s\n", chip.Type)
                }
                for _, prg := range board.PRG {
                    fmt.Printf("\t\tPRG (%s): %s\n", prg.Size, prg.SHA1)
                }
                for _, chr := range board.CHR {
                    fmt.Printf("\t\tCHR (%s): %s\n", chr.Size, chr.SHA1)
                }
            }
        }
    }
}

func main() {
    database := ReadXML()

    Dump(database)

    f, err := os.Open("/tmp/j", os.O_WRONLY | os.O_TRUNC | os.O_CREAT, 0x1a4) // 0644
    if err != nil {
        panic(fmt.Sprintf("error saving: %s", err))
    }
    e := gob.NewEncoder(f)
    e.Encode(database)
}
