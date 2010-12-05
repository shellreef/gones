// Created:20101128
// By Jeff Connelly

// .nes file reader

package nesfile

import (
    "fmt"
    "os"
    "bytes"
    "encoding/binary"
)

// iNES (.nes) file header
// http://wiki.nesdev.com/w/index.php/INES
type NesfileHeader struct {
    Magic uint32
    PrgPageCount, ChrPageCount, Flags6, Flags7, RamPageCount, Flags9, Flags10 uint8
    Reserved [5]byte
}

// Represents a game cartridge
type Cartridge struct {
    Prg []([]byte)
    Chr []([]byte)
    Mapper int // TODO: string constant not flag
    // TODO: flags
}

const NESFILE_MAGIC = 0x1a53454e        // NES^Z
const PRG_PAGE_SIZE = 16384
const CHR_PAGE_SIZE = 8192 

// Load a .nes file
func Open(filename string) (*Cartridge) {
    data := slurp(filename)
    // Convert data to an object compatible with http://golang.org/pkg/io/
    buffer := bytes.NewBuffer(data)

    header := new(NesfileHeader)
    cart := new(Cartridge)

    binary.Read(buffer, binary.LittleEndian, header)

    if header.Magic != NESFILE_MAGIC {
        panic(fmt.Sprintf("invalid nesfile signature: %x != %x", header.Magic, NESFILE_MAGIC))
    }

    cart.Prg = readPages(buffer, PRG_PAGE_SIZE, int(header.PrgPageCount))
    cart.Chr = readPages(buffer, CHR_PAGE_SIZE, int(header.ChrPageCount))
   
    // This is weird, but the mapper nibbles are spread across two bytes
    cart.Mapper = int(header.Flags7 & 0xf0 | header.Flags6 & 0xf0 >> 4)

    // TODO: read and set flags

    fmt.Printf("ROM: %d, VROM: %d, Mapper %d\n", header.PrgPageCount, header.ChrPageCount, cart.Mapper)

    return cart
}

// Read chunks of data from a buffer
func readPages(buffer *bytes.Buffer, size int, pageCount int) ([]([]byte)) {
    pages := make([]([]byte), pageCount)

    for i := 0; i < pageCount; i++ {
        page := make([]byte, size)
        readLength, err := buffer.Read(page)
        //fmt.Printf("read page %d size=%d\n", i, readLength)
        if err != nil {
            panic(fmt.Sprintf("readPages(%d, %d) #%d failed: %d %s", size, pageCount, i, readLength, err))
        }

        pages[i] = page
    }

    return pages
}

// Read all the bytes from a file, terminating if an error occurs
func slurp(filename string) []byte {
    f, err := os.Open(filename, os.O_RDONLY, 0)
    if f == nil {
        panic(fmt.Sprintf("cannot open %s: %s", filename, err))
    }
    stat, err := f.Stat()
    expectedLength := stat.Size
    data := make([]byte, expectedLength)
    readLength, err := f.Read(data)
    if int64(readLength) != expectedLength {
        panic(fmt.Sprintf("failed to read all %d bytes (only %d) from %s: %s", expectedLength, readLength, filename, err))
    }

    f.Close()

    //fmt.Printf("Read %d bytes from %s\n", len(data), filename)

    return data
}

