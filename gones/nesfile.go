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
type NesfileHeader struct {
    Magic uint32
    PrgPageCount, ChrPageCount, MapperInfo1, MapperInfo2, RamPageCount, PalFlag uint8
    Reserved [6]byte
}

// Represents a game cartridge
type Cartridge struct {
    Prg []([]byte)
    Chr []([]byte)
    // TODO: mappers
}

const NESFILE_MAGIC = 0x1a53454e        // NES^Z
const PRG_PAGE_SIZE = 16384
const CHR_PAGE_SIZE = 8192 

// Load a .nes file
func Load(filename string) (*Cartridge) {
    data := slurp(filename)
    // Convert data to an object compatible with http://golang.org/pkg/io/
    buffer := bytes.NewBuffer(data)

    header := new(NesfileHeader)
    cart := new(Cartridge)

    binary.Read(buffer, binary.LittleEndian, header)

    if header.Magic != NESFILE_MAGIC {
        fmt.Fprintf(os.Stderr, "invalid nesfile signature: %x != %x", header.Magic, NESFILE_MAGIC)
        os.Exit(1)
    }

    fmt.Printf("ROM: %d, VROM: %d\n", header.PrgPageCount, header.ChrPageCount)

    cart.Prg = readPages(buffer, PRG_PAGE_SIZE, int(header.PrgPageCount))
    cart.Chr = readPages(buffer, CHR_PAGE_SIZE, int(header.ChrPageCount))

    return cart
}

// Read chunks of data from a buffer
func readPages(buffer *bytes.Buffer, size int, pageCount int) ([]([]byte)) {
    pages := make([]([]byte), pageCount)

    for i := 0; i < pageCount; i++ {
        page := make([]byte, size)
        readLength, err := buffer.Read(page)
        fmt.Printf("read page %d size=%d\n", i, readLength)
        if err != nil {
            fmt.Fprintf(os.Stderr, "readPages(%d, %d) #%d failed: %d %s", size, pageCount, i, readLength, err)
            os.Exit(2)
        }

        pages[i] = page
    }

    return pages
}

// Read all the bytes from a file, terminating if an error occurs
func slurp(filename string) []byte {
    f, err := os.Open(filename, os.O_RDONLY, 0)
    if f == nil {
        fmt.Fprintf(os.Stderr, "cannot open %s: %s", filename, err)
        os.Exit(1)
    }
    stat, err := f.Stat()
    expectedLength := stat.Size
    data := make([]byte, expectedLength)
    readLength, err := f.Read(data)
    if int64(readLength) != expectedLength {
        fmt.Fprintf(os.Stderr, "failed to read all %d bytes (only %d) from %s: %s", expectedLength, readLength, filename, err)
        os.Exit(1)
    }

    f.Close()

    fmt.Printf("Read %d bytes from %s\n", len(data), filename)

    return data
}

