// Created:20101226
// By Jeff Connelly
//
// Load NES cartridge database
// XML from http://bootgod.dyndns.org:7777/

package cartdb

import ("fmt"
        "os"
        "xml"
        "gob"
        "time"
        "crypto/sha1"
        
        "cartridge"
        "nesfile"
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

// Search result matches - matches on cart, but includes game information too
type CartMatches struct {
    Cartridge Cartridge
    Game Game
}

// Load game database from gob if possible; if not, load from XML
// then create gob for faster future loading
func Load() (*Database) {
    fastFile := "data/cartdb.gob"
    slowFile := "data/cartdb.xml"

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

    //start := time.Nanoseconds()
    decoder := gob.NewDecoder(r)
    decoder.Decode(&database)
    //took := time.Nanoseconds() - start

    //fmt.Printf("Gob: Loaded %d games in %.4f s\n", len(database.Game), float(took) / 1e9)

    return &database, nil
}

// Show the salient parts of the database in human-readable format
func Dump(database *Database) {
    for _, game := range database.Game {
        DumpGame(game)
    }
}

// Show information on all carts for game
func DumpGame(game Game) {
    fmt.Printf("[%s] %s - %s\n", game.Region, game.Developer, game.Name)
    for _, cart := range game.Cartridge {
        DumpCart(cart)
    }
}

// Show information on a cartridge
func DumpCart(cart Cartridge) {
    rev := ""
    if len(cart.Revision) != 0 {
        rev = fmt.Sprintf(" (rev. %s)", cart.Revision)
    } 
    fmt.Printf("\tCartridge%s for %s (%s)\n", rev, cart.System, cart.SHA1)
    for _, board := range cart.Board {
        fmt.Printf("\tBoard: %s (%s)\n", board.PCB, board.Type)
        for _, chip := range board.Chip { 
            fmt.Printf("\t\tChip: %s\n", chip.Type)
        }
        for j, prg := range board.PRG {
            fmt.Printf("\t\tPRG %d: %s\n", j, prg.Size)
        }
        for k, chr := range board.CHR {
            fmt.Printf("\t\tCHR %d: %s\n", k, chr.Size)
        }
    }
}

// Print a short summary of the search result match
func DumpMatch(match CartMatches) {
    game := match.Game
    cart := match.Cartridge
    fmt.Printf("%s ", game.Name)
     if len(cart.Revision) != 0 {
        fmt.Printf(" (rev. %s) ", cart.Revision)
    } 

    fmt.Printf(" by %s/%s - %s - ", game.Developer, game.Publisher, game.Region)

    for _, board := range cart.Board {
        fmt.Printf("board %s (%s); ", board.PCB, board.Type)
        for _, chip := range board.Chip { 
            fmt.Printf("chip %s; ", chip.Type)
        }
        for j, prg := range board.PRG {
            fmt.Printf("PRG %d: %s; ", j, prg.Size)
        }
        for k, chr := range board.CHR {
            fmt.Printf("CHR %d: %s; ", k, chr.Size)
        }
    }

    fmt.Printf("\n")
}

// Find games matching hashes
func Identify(database *Database, cart *cartridge.Cartridge) ([]CartMatches) {
    var found []CartMatches

    hash := cartHash(cart)

    fmt.Printf("Cartridge SHA-1 = %s\n", hash)

    for _, game := range database.Game {
        for _, cart := range game.Cartridge {
            if cart.SHA1 == hash {
                // Match on the cart, but return the game too, which has more details
                found = append(found, CartMatches{cart, game})
            }
        }
    }

    return found
}

// Get SHA-1 digest of cartridge PRG and CHR
func cartHash(cart *cartridge.Cartridge) (string) {
    digester := sha1.New()

    // First PRG.. note, some games have multiple PRG chips (such as Action-52,
    // having 3 x 512k), but the nesfile format stores the data in 8k banks,
    // so there is not enough information to calculate individual chip hashes.
    // So we calculate the overall cartridge hash instead, including everything.
    // (Note: it might be interesting to calculate PRG/CHR and match individually
    // to find differences, since most games have only one chip.)
    digester.Write(cart.Prg)

    digester.Write(cart.Chr)

    hex := ""
    for _, octet := range digester.Sum() {
        hex += fmt.Sprintf("%.2X", octet)
    }
    return hex
    // TODO: compute CRC-32 also (hash/crc32)? The database has both,
    // but they should be equivalent.. though checking CRC-32 might help
    // detect database inconsistencies. Maybe CRC-32 only; faster, and
    // no unintentional collisions in NES games (but could be fabricated easily!)
}

func main() {
    db := Load()
    //Dump(db)

    c := nesfile.Open(os.Args[1])

    finds := Identify(db, c)

    for _, found := range finds {
        game, cart := found.Game, found.Cartridge
        fmt.Printf("[%s] %s - %s\n", game.Region, game.Developer, game.Name)
        DumpCart(cart)
    }
}
