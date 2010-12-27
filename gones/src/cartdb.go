// Created:20101226
// By Jeff Connelly
//
// Load NES cartridge database
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

// Top-level database tag
type Database struct {
    XMLName xml.Name "database"
    
    Game []Game
}

// Load game database from gob if possible; if not, load from XML
// then create gob for faster future loading
func Load() (*Database) {
    fastFile := "/tmp/gob"
    slowFile := "../NesCarts (2010-11-29).xml"

    database, err := LoadGob(fastFile)
    if database == nil {
        fmt.Printf("Loading from XML %s -> gob %s\n", slowFile, fastFile)
        database = LoadXML(slowFile)
        WriteGob(fastFile, database)

        database, err = LoadGob(fastFile)
    }

    if err != nil {
        panic(fmt.Sprintf("Failed to load %s: %s\n", fastFile, err))
    }

    return database
}


// Load the game database in XML (the original source format, but slow)
func LoadXML(filename string) (*Database) {
    r, err := os.Open(filename, os.O_RDONLY, 0)
    if r == nil {
        panic(fmt.Sprintf("cannot open %s: %s", filename, err))
    }

    database := Database{}

    start := time.Nanoseconds()
    xml.Unmarshal(r, &database)
    took := time.Nanoseconds() - start

    fmt.Printf("XML: Loaded %d games in %.4f s\n", len(database.Game), float(took) / 1e9)

    return &database
}

// Write the game database in a fast binary Go-specific format
func WriteGob(filename string, database *Database) {
    f, err := os.Open(filename, os.O_WRONLY | os.O_TRUNC | os.O_CREAT, 0x1a4) // 0644
    if err != nil {
        panic(fmt.Sprintf("error saving %s: %s", filename, err))
    }

    e := gob.NewEncoder(f)
    e.Encode(database)
}

// Load the gob written by WriteGob. This is 12x as fast as LoadXML(),
// taking less than 1/10 s while LoadXML takes >1 s, so it is preferred.
func LoadGob(filename string) (*Database, os.Error) {
    r, err := os.Open(filename, os.O_RDONLY, 0)
    if r == nil {
        return nil, err
    }

    database := Database{}

    start := time.Nanoseconds()
    decoder := gob.NewDecoder(r)
    decoder.Decode(&database)
    took := time.Nanoseconds() - start

    fmt.Printf("Gob: Loaded %d games in %.4f s\n", len(database.Game), float(took) / 1e9)

    return &database, nil
}

// Show the salient parts of the database in human-readable format
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
    db := Load()

    Dump(db)

}
