// Created:20101226
// By Jeff Connelly
//
// Read NES cartridge database
// XML from http://bootgod.dyndns.org:7777/

package main

import ("fmt"
        "os"
        "xml")

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

func main() {
    filename := "../NesCarts (2010-11-29).xml"
    r, err := os.Open(filename, os.O_RDONLY, 0)
    if r == nil {
        panic(fmt.Sprintf("cannot open %s: %s", filename, err))
    }

    database := Database{}

    xml.Unmarshal(r, &database)

    fmt.Printf("Loaded %d games\n", len(database.Game))

    for i, game := range database.Game {
        fmt.Printf("#%d. [%s] %s - %s\n", i, game.Region, game.Developer, game.Name)
        for _, cart := range game.Cartridge {
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
