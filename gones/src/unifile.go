// Created:20110103
// By Jeff Connelly
//
// UNIF file format reader
//
// http://web.archive.org/web/20040614053359/www.parodius.com/~veilleux/UNIF_current.txt

package cartridge

import (
    "fmt"
    "os"
    "encoding/binary"
    "strings"
)

const UNIF_MAGIC = 0x46494e55       // 'UNIF'

type UnifFileHeader struct {
    Magic uint32
    Revision uint32
    Expansion [24]byte
}

type UnifChunkHeader struct {
    ID [4]byte
    Length uint32
}

// Open a Universal NES Image Format file
func OpenUNIF(f *os.File) (*Cartridge) {
    cart := new(Cartridge)

    // Header
    header := new(UnifFileHeader)
    binary.Read(f, binary.LittleEndian, header)
    if header.Magic != UNIF_MAGIC {
        fmt.Printf("invalid unif signature: %.4x\n", header.Magic)
        // TODO: negotiate to other loaders?
        return nil
    }

    for {
        id, data := readUnifChunk(f)
        if len(data) == 0 {
            break
        }
        fmt.Printf("chunk: %s: %d bytes\n", id, len(data))

        switch id {
        case "MAPR": mapperName := string(data[:])
            // Null-terminated string
            term := strings.Index(mapperName, "\000")
            if term != -1 {
                cart.MapperName = mapperName[0:term]
            } else {
                cart.MapperName = mapperName
            }
        case "NAME": fmt.Printf("Name: %s\n", string(data[:])) // TODO: use somewhere
        //case "MIRR": // TODO: read mirroring

        // TODO: read more chips
        case "PRG0": cart.Prg = data
        case "CHR0": cart.Chr = data
        }
    }

    f.Close()
    return cart
}

func readUnifChunk(f *os.File) (ID string, data []byte) {
    header := new(UnifChunkHeader)
    binary.Read(f, binary.LittleEndian, header)
    data = make([]byte, header.Length)

    binary.Read(f, binary.LittleEndian, data)

    return string(header.ID[:]), data
}
