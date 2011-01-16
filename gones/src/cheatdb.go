// Created:20110115
// By Jeff Connelly

// Cheat code database

package cheatdb

import ("fmt"
        "xml"
        "os"
        "path"

        "github.com/kless/go-sqlite/sqlite" // https://github.com/kless/go-sqlite/tree/master/sqlite
        //"gosqlite.googlecode.com/hg/sqlite"    // http://code.google.com/p/gosqlite/
            // there are also other implementations, like https://github.com/feyeleanor/gosqlite3/blob/master/database.go 
        
        "gamegenie"
        )

type Cheats struct {
    XMLName xml.Name "cheat-codes"

    Game []Game
}

type Game struct {
    Name string "attr"
    GaloobID string "attr"
    GaloobName string "attr"
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


type Database struct {
    handle *sqlite.Connection
}

func Open() (db *Database) {
    db = new(Database)

    root, _ := path.Split(os.Args[0])
    filename := path.Join(root, "data", "cheats.sqlite")
    _, err := os.Stat(filename)
    if err != nil {
        // Does not exist, so initialize it before returning
        defer func() {
            db.CreateTables()
        }()
    }
    
    handle, err := sqlite.Connect(filename)
    if err != nil {
        panic(fmt.Sprintf("failed to open: %s", err))
    }

    db.handle = handle

    return db
}

// Execute a SQL statement, checking for error
func (db *Database) exec(sql string, args ...interface{}) { 
    err := db.handle.Exec(sql, args...)
    if err != nil {
        panic(fmt.Sprintf("failed: %s %s\ncheatdb exec failed: %s", sql, args, err))
    }
}

func (db *Database) AllCodes() {
    query, _ := db.handle.Prepare("SELECT game.name,effect.title,code.cart_id,cpu_address,value,compare FROM game,effect,code WHERE effect.game_id=game.id AND code.effect_id=effect.id")
    err := query.Exec()
    if err != nil {
        panic(fmt.Sprintf("AllCodes() failed: %s", err))
    }

    var gameName string
    var effectTitle string
    var cartID int
    var cpuAddress uint16
    var value uint8
    var compare uint8

    for query.Next() {
        query.Scan(&gameName, &effectTitle, &cartID, &cpuAddress, &value, &compare)
        // uh..NULL compare?
        fmt.Printf("%s: %s: %.4X?%.2X:%.2X\n", gameName, effectTitle, cpuAddress, compare, value)
    }
}

// Import GoNES XML cheat database 
// TODO: read Nestopia's NstCheat files, like those found on http://www.mightymo.net/downloads.html
// TODO: read Nesticle .pat files, [+]<code> [<name>] per line, raw format. http://www.zophar.net/faq/nitrofaq/nesticle.readme.txt and some at http://jeff.tk:81/consoles/
func (db *Database) ImportXML(filename string) {
    fmt.Printf("Importing %s\n", filename)
    r, err := os.Open(filename, os.O_RDONLY, 0)
    if r == nil {
        panic(fmt.Sprintf("cannot open %s: %s", filename, err))
    }

    cheats := Cheats{}

    xml.Unmarshal(r, &cheats)

    db.exec("INSERT INTO user(fullname) VALUES(?)", "Galoob")   // TODO: pass as argument, search for existing, and stop hardcoding
    userID := db.handle.LastInsertRowID()

    gamesInserted := 0
    cartsInserted := 0
    effectsInserted := 0
    codesInserted := 0

    // Insert into database
    for _, game := range cheats.Game {
        db.exec("INSERT INTO game(name,galoob_id,galoob_name) VALUES(?,?,?)", game.Name, game.GaloobID, game.GaloobName)
        // LastInsertedRowID() requires gosqlite patch at http://code.google.com/p/gosqlite/issues/detail?id=7
        // or for https://github.com/kless/go-sqlite, add to end of connection.go:
        /*
func (c *Connection) LastInsertRowID() (int64) {
    return int64(C.sqlite3_last_insert_rowid(c.db))
}
*/
        gameID := db.handle.LastInsertRowID()
        gamesInserted += 1

        cartName2ID := make(map[string]int64)

        for _, cart := range game.Cartridge {
            db.exec("INSERT INTO cart(game_id,sha1,name,filename) VALUES(?,?,?,?)", gameID, cart.SHA1, cart.Name, cart.Filename)
            cartName2ID[cart.Name] = db.handle.LastInsertRowID()
            cartsInserted += 1
        }

        for _, effect := range game.Effect {
            db.exec("INSERT INTO effect(game_id,number,source,title,hacker_id,create_date) VALUES(?,?,?,?,?,?)", 
                gameID, effect.Number, effect.Source, effect.Title, userID, "1994")
            effectID := db.handle.LastInsertRowID()
            effectsInserted += 1

            for _, code := range effect.Code {
                func() {
                    defer func() {
                        r := recover(); if r != nil {
                            fmt.Printf("Bad code for %s: %s: %s\n", game.Name, code.Genie, r)
                        }
                    }()

                    decoded := gamegenie.Decode(code.Genie)
                    db.exec("INSERT INTO code(effect_id,cpu_address,value,compare) VALUES(?,?,?,?)",
                        effectID, 
                        decoded.CPUAddress(),
                        decoded.Value,
                        decoded.Key)
                    codeID := db.handle.LastInsertRowID()
                    codesInserted += 1

                    // If code only applies to a specific cart, set it; otherwise, leave NULL to apply to all
                    if code.Applies != "" {
                        db.exec("UPDATE code SET cart_id=? WHERE id=?", cartName2ID[code.Applies], codeID)
                    }
                }()
            }
        }
    }
    fmt.Printf("Imported %d games, %d carts, %d effects, and %d total codes\n",
        gamesInserted, cartsInserted, effectsInserted, codesInserted)
}

// Database schema
func (db *Database) CreateTables() {
    // Code info
    db.exec(`CREATE TABLE game(     -- an abstract "game", has ≥1 carts
        id INTEGER PRIMARY KEY, 
        name TEXT NOT NULL, 
        galoob_id TEXT NULL, 
        galoob_name TEXT NULL
        )`)
    db.exec(`CREATE TABLE cart(     -- a physical cartridge, possibly different versions
        id INTEGER PRIMARY KEY, 
        game_id INTEGER NOT NULL, 
        sha1 TEXT NOT NULL, 
        filename TEXT NOT NULL,
        name TEXT NULL,             -- PRG0, PRG1, etc. if multiple versions, otherwise NULL

        FOREIGN KEY(game_id) REFERENCES game(id)
        )`)
    db.exec(`CREATE TABLE effect(
        id INTEGER PRIMARY KEY,
        game_id INTEGER NOT NULL,
        title TEXT NOT NULL,
        hacker_id INTEGER NOT NULL,
        create_date DATETIME NOT NULL,
        
        number INTEGER,
        source TEXT,

        FOREIGN KEY(game_id) REFERENCES game(id),
        FOREIGN KEY(hacker_id) REFERENCES user(id)
        )`)
    db.exec(`CREATE TABLE code(     -- a decoded Game Genie code, has ≥0 patches
        id INTEGER PRIMARY KEY, 
        cart_id INTEGER NULL,       -- set if code only applies to one cartridge version; otherwise NULL for all
        effect_id INTEGER NOT NULL,
        cpu_address INTEGER NOT NULL, -- the address referred to, $8000-FFFF
        value INTEGER NOT NULL, 
        compare INTEGER NULL,       -- compare byte if 8-letter code, or NULL for 6-letter

        FOREIGN KEY(cart_id) REFERENCES cart(id),
        FOREIGN KEY(effect_id) REFERENCES effect(id)
        )`)
    db.exec(`CREATE TABLE patch(    -- an actual known modification to ROM
        id INTEGER PRIMARY KEY,
        code_id INTEGER NOT NULL,
        rom_address INTEGER NOT NULL,
        rom_before BLOB NOT NULL,   -- bytes in the ROM before the affected address, for context
        rom_after BLOB NOT NULL,    -- bytes _including_ and after the affected address

        FOREIGN KEY(code_id) REFERENCES code(id)
        )`)

    // Extra info
    db.exec(`CREATE TABLE user(     -- someone that created codes or commented on them
        id INTEGER PRIMARY KEY,
        fullname TEXT NOT NULL UNIQUE
        )`)
    db.exec(`CREATE TABLE comment(  -- user commentary on an effect
        id INTEGER PRIMARY KEY,
        effect_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        message TEXT NOT NULL,
        create_date DATETIME NOT NULL,

        FOREIGN KEY(effect_id) REFERENCES effect(id),
        FOREIGN KEY(user_id) REFERENCES user(id)
        )`)
}


// TODO: web interface
