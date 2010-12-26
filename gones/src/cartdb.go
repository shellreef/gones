// Created:20101226
// By Jeff Connelly
//
// Read NES cartridge database
// XML from http://bootgod.dyndns.org:7777/

package main

import ("fmt"
        "os"
        "xml")

type Game struct {
    Name string "attr"
    Altname string "attr"
    Class string "attr"
    Catalog string "attr"
    Publisher string "attr"
    Developer string "attr"
    Region string "attr"
    Players string "attr"
    Date string "attr"
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

    for i, game := range database.Game {
        fmt.Printf("#%d. %s\n", i, game.Name)
    }
}
