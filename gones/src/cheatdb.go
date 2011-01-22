// Created:20110115
// By Jeff Connelly

// Cheat code database

package cheatdb

import ("fmt"
        "xml"
        "os"
        "path"
        "io"
        "json"

        "github.com/kless/go-sqlite/sqlite" // https://github.com/kless/go-sqlite/tree/master/sqlite
        //"gosqlite.googlecode.com/hg/sqlite"    // http://code.google.com/p/gosqlite/
        // there are also other incompatible SQLite wrappers, like https://github.com/feyeleanor/gosqlite3/blob/master/database.go 
        // probably the right direction is: https://github.com/thomaslee/go-dbi - but it only supports MySQL

        "web" // http://www.getwebgo.com/

        
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
    // TODO: have ROM patches
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


// Prepare a SQL statement, then execute it and return the query, checking for error
// Like exec(), but you can receive the results of the query
func (db *Database) query(sql string, args ...interface{}) (*sqlite.Stmt) {
    query, err := db.handle.Prepare(sql)
    if err != nil {
        panic(fmt.Sprintf("failed to prepare: %s", sql))
    }
    err = query.Exec(args...)
    if err != nil {
        panic(fmt.Sprintf("failed to exec %s %s", sql, args))
    }
    return query
}


// Get all cartridges recognized in the cheat database
func (db *Database) AllCarts(analyzeCart func(cartPath string)(bool), analyzeCode func(code gamegenie.GameGenieCode)(bool,uint32,string,string,string)) {
    query, _ := db.handle.Prepare("SELECT game.name,game.id,cart.filename,cart.name,cart.sha1,cart.id FROM game,cart WHERE cart.game_id=game.id")
    err := query.Exec()
    if err != nil {
        panic(fmt.Sprintf("AllCarts() failed: %s", err))
    }

    for query.Next() {
        var gameName string
        var gameID int
        var cartFilename string
        var cartName string
        var cartSHA1 string
        var cartID int

        err := query.Scan(&gameName, &gameID, &cartFilename, &cartName, &cartSHA1, &cartID)
        if err != nil {
            panic(fmt.Sprintf("AllCarts() scan error: %s\n", err))
        }

        fmt.Printf("\n%s(%s) = %s\n", gameName, cartName, cartFilename)
        success := analyzeCart(path.Join("../roms/best/", gameName, cartFilename))
        if !success {
            continue
        }


        db.CodesFor(gameID, analyzeCode)
    }
}

//TODO: func (db *Database) CodesFor(cart *cartridge.Cartridge) {
func (db *Database) CodesFor(gameID int, analyzeCode func(code gamegenie.GameGenieCode)(bool,uint32,string,string,string)) {
    query, err := db.handle.Prepare("SELECT effect.title,code.cart_id,code.id,cpu_address,value,compare FROM effect,code WHERE code.effect_id=effect.id AND effect.game_id=?")
    if err != nil {
        panic(fmt.Sprintf("AllCodes() prepare failed: %s", err))
    }
    err = query.Exec(gameID)
    if err != nil {
        panic(fmt.Sprintf("AllCodes() exec failed: %s", err))
    }


    for query.Next() {
        var effectTitle string
        var cartID string    // TODO: why can't this be an int? fails with: scan error: arg 2 as int: parsing "": invalid argument
        var codeID int
        var cpuAddress int
        var value int
        var compare *int = new(int)

        // Requires patch to support **int in Scan() argument, in order to read NULLs:
/*
--- a/sqlite/statement.go
+++ b/sqlite/statement.go
@@ -70,6 +70,7 @@ func (s *Stmt) Scan(args ...interface{}) os.Error {
        for i, v := range args {
                n := C.sqlite3_column_bytes(s.stmt, C.int(i))
                p := C.sqlite3_column_blob(s.stmt, C.int(i))
+                t := C.sqlite3_column_type(s.stmt, C.int(i))
                if p == nil && n > 0 {
                        return os.NewError("got nil blob")
                }
@@ -91,6 +92,17 @@ func (s *Stmt) Scan(args ...interface{}) os.Error {
                                        err.String())
                        }
                        *v = x
+                case **int:
+                       if t == C.SQLITE_NULL {
+                               *v = nil
+                        } else {
+                               x, err := strconv.Atoi(string(data))
+                               if err != nil {
+                                       return os.NewError("arg " + strconv.Itoa(i) + " as int: " +
+                                               err.String())
+                               }
+                               **v = x
+                       }
*/

        err := query.Scan(&effectTitle, &cartID, &codeID, &cpuAddress, &value, &compare)
        if err != nil {
            panic(fmt.Sprintf("AllCodes() scan error: %s\n", err))
        }
        fmt.Printf("%s: ", effectTitle)

        var code gamegenie.GameGenieCode
        code.Address = uint16(cpuAddress) & 0x7fff
        code.Value = uint8(value)
        if compare != nil {
            fmt.Printf("%.4X?%.2X:%.2X\n", cpuAddress, *compare, value)
            code.HasKey = true
            code.Key = uint8(*compare)
        } else {
            fmt.Printf("%.4X:%.2X\n", cpuAddress, value)
            code.HasKey = false
        }

        found, romAddress, romChip, romBefore, romAfter := analyzeCode(code)
        // Add patch in database
        if found {
            fmt.Printf("Code %s: ROM Address: %.6X (%s)\n", code, romAddress, romChip)
            fmt.Printf("ROM context: %s | %s\n", romBefore, romAfter)

            db.exec("INSERT INTO patch(code_id,rom_address,rom_before,rom_after) VALUES(?,?,?,?)", codeID, romAddress, romBefore, romAfter)
            patchID := db.handle.LastInsertRowID()
            fmt.Printf("%d\n", patchID)
        }
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
                    db.exec("INSERT INTO code(effect_id,cpu_address,value) VALUES(?,?,?)",
                        effectID, 
                        decoded.CPUAddress(),
                        decoded.Value)
                    codeID := db.handle.LastInsertRowID()
                    codesInserted += 1

                    // Only set if has a compare value, otherwise leave NULL
                    if decoded.HasKey {
                        db.exec("UPDATE code SET compare=? WHERE id=?", decoded.Key, codeID)
                    }

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

// Web interface
func (db *Database) Serve() {
    web.Get("/patches.js", func(w *web.Context) { // TODO: accept arguments to filter
        query := db.query("SELECT rom_address,rom_before,rom_after, cpu_address,value,compare, title,game.name FROM patch,code,effect,game WHERE patch.code_id=code.id AND code.effect_id=effect.id AND effect.game_id=game.id;") // TODO: get effect,game,cart

        w.SetHeader("Content-Type", "text/javascript", false)
        w.WriteString("var PATCHES = [\n")

        for query.Next() {
            var r struct {
                RomAddress int "romAddress"
                RomBefore string "romBefore"
                RomAfter string "romAfter"
                CpuAddress int "cpuAddress"
                Value int "value"
                Compare *int "compare"
                Title string "title"
                GameName string "gameName"
            }

            r.Compare = new(int)

            err := query.Scan(&r.RomAddress, &r.RomBefore, &r.RomAfter, &r.CpuAddress, &r.Value, &r.Compare, &r.Title, &r.GameName)
            if err != nil {
                panic(fmt.Sprintf("failed to Scan: %s", err))
            }
            // TODO: use JSON module, avoid XSS
            html, err2 := json.MarshalForHTML(r)

            if err2 != nil {
                panic(fmt.Sprintf("failed to marshal: %s", err2))
            }

            w.Write(html)
            w.WriteString(",\n")
        }
        w.WriteString("]\n")
    })

    web.Get("/games.js", func(w *web.Context) {
        w.SetHeader("Content-Type", "text/javascript", false)
        w.WriteString("[\n")

        //query := db.query("SELECT game.id,game.name,game.galoob_name,game.galoob_id, COUNT(effect.id) FROM game,effect WHERE effect.game_id=game.id GROUP BY game.id ORDER BY COUNT(effect.id) DESC;")
        query := db.query("SELECT game.id,game.name,game.galoob_name,game.galoob_id, COUNT(effect.id) FROM game,effect WHERE effect.game_id=game.id GROUP BY game.id ORDER BY game.name;")
        for query.Next() {
            var r struct {
                GameID int "id"
                GameName string "name"
                GameGaloob string "galoob"
                GameGaloobID string "galoobID"
                EffectCount int "effectCount"
            }

            err := query.Scan(&r.GameID, &r.GameName, &r.GameGaloob, &r.GameGaloobID, &r.EffectCount)
            if err != nil {
                panic(fmt.Sprintf("failed to Scan: %s", err))
            }

            html, err2 := json.MarshalForHTML(r)
            if err2 != nil {
                panic(fmt.Sprintf("failed to marshal: %s", err2))
            }
            w.Write(html)
            w.WriteString(",\n")
        }
        w.WriteString("]\n")
    })

    web.Get("/", func(w *web.Context) {
        root, _ := path.Split(os.Args[0])
        filename := path.Join(root, "..", "genie", "jsdis.html")
        f, err := os.Open(filename, os.O_RDONLY, 0)
        if err != nil {
            panic(fmt.Sprintf("failed to open %s: %s", filename, err))
        }

        io.Copy(w, f)
    })


    web.Get("/favicon.ico", func(w *web.Context) {
        // TODO: real favicon
        w.SetHeader("Content-Type", "image/png", false)
    })


    web.Run("0.0.0.0:9999")
}
