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

type Mirroring string 
const (
    MirrorVertical="MirrorVertical";
    MirrorHorizontal="MirrorHorizontal";
    MirrorFourscreen="MirrorFourscreen";
)

type Platform string
const (
    PlatformNES="PlatformNES";
    PlatformVs="PlatformVS";
    PlatformPlayChoice="PlatformPlayChoice";
)

// Represents a game cartridge
type Cartridge struct {
    Prg []([]byte)      // Program data
    Chr []([]byte)      // Character data
    Platform Platform
    MapperCode int
    Mirroring Mirroring
    BatteryBacked bool
    Trainer []byte
    HintScreen []byte     
}

const NESFILE_MAGIC = 0x1a53454e        // NES^Z
const PRG_PAGE_SIZE = 16384
const CHR_PAGE_SIZE = 8192 
const TRAINER_SIZE = 512
const HINTSCREEN_SIZE = 8192

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

    // This is weird, but the mapper nibbles are spread across two bytes
    cart.MapperCode = int(header.Flags7 & 0xf0 | header.Flags6 & 0xf0 >> 4)

    // http://wiki.nesdev.com/w/index.php/INES
    // Flags 6: mirroring, battery backed, trainer

    if header.Flags6 & 1 == 1 {
        cart.Mirroring = MirrorHorizontal
    } else {
        cart.Mirroring = MirrorVertical
    }
    if header.Flags6 & 8 == 1 {
        cart.Mirroring = MirrorFourscreen
    }

    cart.BatteryBacked = header.Flags6 & 2 == 2   // "SRAM in CPU $6000-$7FFF, if present, is battery backed"

    if header.Flags6 & 4 == 4 {
        // "512-byte trainer at $7000-$71FF (stored before PRG data)"
        cart.Trainer = readPages(buffer, TRAINER_SIZE, 1)[0]
    }
    
    // Flags 7: platform, NES 2.0 format
    switch header.Flags7 & 3 {
    case 0: cart.Platform = PlatformNES
    case 1: cart.Platform = PlatformVs
    case 2: cart.Platform = PlatformPlayChoice
    case 3: panic("nesfile flags7 platform is both Vs and PlayChoice")
    }

    if header.Flags7 & 0x0c == 8 {
        fmt.Printf("NES 2.0 detected\n")
        // TODO: read flags 8-15 http://wiki.nesdev.com/w/index.php/NES_2.0
    }

    // Read banks
    cart.Prg = readPages(buffer, PRG_PAGE_SIZE, int(header.PrgPageCount))
    cart.Chr = readPages(buffer, CHR_PAGE_SIZE, int(header.ChrPageCount))

    if cart.Platform == PlatformPlayChoice {
        // "PlayChoice-10 (8KB of Hint Screen data stored after CHR data)"
        cart.HintScreen = readPages(buffer, HINTSCREEN_SIZE, 1)[0]
    }
   
    fmt.Printf("ROM: %d, VROM: %d, Mapper #%d\n", len(cart.Prg), len(cart.Chr), cart.MapperCode)
    fmt.Printf("Mirroring: %s, BatteryBacked=%t\n", cart.Mirroring, cart.BatteryBacked)
    fmt.Printf("Platform type: %s\n", cart.Platform)

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

