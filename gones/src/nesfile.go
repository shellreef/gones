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
    PrgPageCount, ChrPageCount, Flags6, Flags7 uint8
    // Defined by NES 2.0
    Flags8, Flags9, Flags10, Flags11, Flags12, Flags13 uint8
    Reserved [2]byte
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

type DisplayType string
const (
    DisplayPAL="DisplayPAL";
    DisplayNTSC="DisplayNTSC";
    DisplayBoth="DisplayBoth";
)

type PPUType string
const (
    PPU_NES="PPU_NES";                  // all NES games use this
    // Vs. PPUs http://wiki.nesdev.com/w/index.php/NES_2.0#Vs._hardware
    PPU_RP2C03B="PPU_RP2C03B";          // "bog standard RGB palette"
    PPU_RP2C03G="PPU_RP2C03G";          // "similar pallete to above, might have 1 changed colour"
    PPU_RP2C04_0001="PPU_RP2C04_0001";  // "scrambled palette + new colours"
    PPU_RP2C04_0002="PPU_RP2C04_0002";  // "same as above, different scrambling, diff new colours"
    PPU_RP2C04_0003="PPU_RP2C04_0003";
    PPU_RP2C04_0004="PPU_RP2C04_0004";
    PPU_RC2C03B="PPU_RC2C03B";          // "bog standard palette, seems identical to RP2C03B"
    PPU_RC2C03C="PPU_RC2C03C";          // "similar to above, but with 1 changed colour or so"
    PPU_RC2C05_01="PPU_RC2C05_01";      // "all five of these have the normal palette..."
    PPU_RC2C05_02="PPU_RC2C05_02";      // "...but with different bits returned on 2002"
    PPU_RC2C05_03="PPU_RC2C05_03";      
    PPU_RC2C05_04="PPU_RC2C05_04";
    PPU_RC2C05_05="PPU_RC2C05_05";
)


// Represents a game cartridge
type Cartridge struct {
    Prg []([]byte)                  // Program data
    Chr []([]byte)                  // Character data
    Platform Platform               // Platform to run on
    DisplayType DisplayType
    MapperCode, SubmapperCode int   // Extra hardware inside the cart
    Mirroring Mirroring
    SRAMBacked bool                 // "SRAM in CPU $6000-$7FFF, if present, is battery backed"
    NoSRAM bool                     // "SRAM in CPU $6000-$7FFF is 0: present; 1: not present"
    BusConflicts bool
    Trainer []byte
    HintScreen []byte
    PPUType PPUType

    // PRG-RAM and CHR-RAM, battery-backed and not
    RamPrgBacked, RamPrgUnback int
    RamChrBacked, RamChrUnback int
}

const NESFILE_MAGIC = 0x1a53454e        // NES^Z
const PRG_PAGE_SIZE = 16384
const CHR_PAGE_SIZE = 8192 
const PRG_RAM_PAGE_SIZE = 8192
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

    prgPageCount := int(header.PrgPageCount)
    chrPageCount := int(header.ChrPageCount)

    if prgPageCount == 0 {
        // Special case for 4 MB PRG
        prgPageCount = 256
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

    cart.SRAMBacked = header.Flags6 & 2 == 2   // "SRAM in CPU $6000-$7FFF, if present, is battery backed"

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

    cart.PPUType = PPU_NES

    if header.Flags7 & 0x0c == 8 {
        fmt.Printf("NES 2.0 detected\n")
        // Flags 8-15 http://wiki.nesdev.com/w/index.php/NES_2.0

        // Mapper variant
        cart.MapperCode |= int(header.Flags8) & 0x0f << 11
        cart.SubmapperCode = int(header.Flags8) & 0x0f >> 4

        // Upper bits of ROM size
        prgPageCount |= int(header.Flags9) & 0x0f << 9
        chrPageCount |= int(header.Flags9) & 0x0f >> 4 << 9
    
        // RAM sizes
        cart.RamPrgBacked = 1 << 6 + (int(header.Flags10) & 0x0f)
        cart.RamPrgUnback = 1 << 6 + (int(header.Flags10) & 0xf0 >> 4)
        cart.RamChrBacked = 1 << 6 + (int(header.Flags11) & 0x0f)
        cart.RamChrUnback = 1 << 6 + (int(header.Flags11) & 0xf0 >> 4)

        // Special case: 0, which encodes to 64 bytes, means 0 bytes
        if cart.RamPrgBacked == 64 { 
            cart.RamPrgBacked = 0 
        }
        if cart.RamPrgUnback == 64 { 
            cart.RamPrgUnback = 0
        } 
        if cart.RamChrBacked == 64 {
            cart.RamPrgBacked = 0
        }
        if cart.RamChrUnback == 64 {
            cart.RamChrUnback = 0
        }

        // Display
        if header.Flags12 & 1 == 1 {
            cart.DisplayType = DisplayPAL   // "312 lines, NMI on line 241, 3.2 dots per CPU clock"
        } else {
            cart.DisplayType = DisplayNTSC  // "262 lines, NMI on line 241, 3 dots per CPU clock"
        }
        if header.Flags12 & 2 == 2{
            cart.DisplayType = DisplayBoth
        }

        if cart.Platform == PlatformVs {
            // Vs hardware
            switch header.Flags13 & 0x0f {
            case  0: cart.PPUType = PPU_RP2C03B
            case  1: cart.PPUType = PPU_RP2C03G
            case  2: cart.PPUType = PPU_RP2C04_0001
            case  3: cart.PPUType = PPU_RP2C04_0002
            case  4: cart.PPUType = PPU_RP2C04_0003
            case  5: cart.PPUType = PPU_RP2C04_0004
            case  6: cart.PPUType = PPU_RC2C03B
            case  7: cart.PPUType = PPU_RC2C03C
            case  8: cart.PPUType = PPU_RC2C05_01
            case  9: cart.PPUType = PPU_RC2C05_02
            case 10: cart.PPUType = PPU_RC2C05_03
            case 11: cart.PPUType = PPU_RC2C05_04
            case 12: cart.PPUType = PPU_RC2C05_05
            default: panic(fmt.Sprintf("Vs system: unknown PPU: %d", header.Flags13 & 0x0f))
            }

            // Not read: upper nibble of Flags13, Vs. mode
        }
    } else {
        // http://wiki.nesdev.com/w/index.php/INES says "8: Size of PRG RAM in 8 KB units (Value 0 infers 8 KB for compatibility; see PRG RAM circuit)"
        // but http://wiki.nesdev.com/w/index.php/NES_2.0 says byte 8 is mapper variant
        if header.Flags8 == 0 {
            cart.RamPrgUnback = PRG_RAM_PAGE_SIZE
        } else {
            // TODO: backed or unbacked?
            cart.RamPrgUnback = int(header.Flags8) * PRG_RAM_PAGE_SIZE
        }

        // "Though in the official specification, very few emulators honor this bit as virtually no ROM images in circulation make use of it."
        if header.Flags9 & 1 == 1 {
            cart.DisplayType = DisplayPAL 
        } else {
            cart.DisplayType = DisplayNTSC 
        }

        // "This byte is not part of the official specification, and relatively few emulators honor it."
        // ..and, it conflicts with Flags9
        switch header.Flags10 & 3 {
        case 0: cart.DisplayType = DisplayNTSC
        case 1, 3: cart.DisplayType = DisplayBoth
        case 2: cart.DisplayType = DisplayPAL
        }
        cart.NoSRAM = header.Flags10 & 0x10 == 0x10
        cart.BusConflicts = header.Flags10 & 0x20 == 0x20
    }

    // Read banks
    cart.Prg = readPages(buffer, PRG_PAGE_SIZE, prgPageCount)
    cart.Chr = readPages(buffer, CHR_PAGE_SIZE, chrPageCount)

    if cart.Platform == PlatformPlayChoice {
        // "PlayChoice-10 (8KB of Hint Screen data stored after CHR data)"
        cart.HintScreen = readPages(buffer, HINTSCREEN_SIZE, 1)[0]
    }
   
    fmt.Printf("ROM: %d, VROM: %d, Mapper #%d\n", len(cart.Prg), len(cart.Chr), cart.MapperCode)

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
