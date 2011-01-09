// Created:20101128
// By Jeff Connelly

// .nes file reader

package cartridge

import (
    "fmt"
    "os"
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

const NESFILE_MAGIC = 0x1a53454e        // NES^Z
const PRG_PAGE_SIZE = 16384
const CHR_PAGE_SIZE = 8192 
const PRG_RAM_PAGE_SIZE = 8192
const TRAINER_SIZE = 512
const HINTSCREEN_SIZE = 8192

// See Disch's detailed mappers list linked from http://wiki.nesdev.com/w/index.php/Mapper
// and http://wiki.nesdev.com/w/index.php/Category:INES_Mappers
// and http://web.archive.org/web/20040614053901/www.parodius.com/~veilleux/boardtable.txt
// TODO: match up with UNIF board names
// TODO: populate this mapper name from the code, and use it instead
// TODO: see ines_convert.c from libunif from http://www.codef00.com/projects.php
// TODO: and also MESS src/machine/nes_ines.c
var INesMapperCode2Name = []string{
    // Most common mappers
    /*000*/ "NROM",
    /*001*/ "MMC1", // SxROM
    /*002*/ "UxROM",
    /*003*/ "CNROM",
    /*004*/ "MMC3", // TxROM, (MMC6), (HxROM)

    /*005*/ "MMC5", // ExROM
    /*006*/ "FFE-FK4",
    /*007*/ "AxROM",
    /*008*/ "FFE-F3K",
    /*009*/ "MMC2", // PxROM
    /*010*/ "MMC4",
    /*011*/ "ColorDreams",
    /*012*/ "Magic Series", // http://mess.dorando.at/svn/?rev=8384
    /*013*/ "CPROM",
    /*014*/ "Sachen SA72008", // http://mess.dorando.at/svn/?rev=8354
    /*015*/ "100-in-1",
    /*016*/ "Bandai",
    /*017*/ "FFE-Mapper17",
    /*018*/ "JALECO-SS88006",
    /*019*/ "Namcot106",
    /*020*/ "FDS",  // http://nocash.emubase.de/everynes.htm#mapper20disksystemprgrambiosdiskirqsound
    /*021*/ "VRC4",
    /*022*/ "VRC2",
    /*023*/ "VRC4/2 variant",
    /*024*/ "VRC6",
    /*025*/ "VRC4 variant",
    /*026*/ "VRC6 variant",
    /*027*/ "",
    /*028*/ "",
    /*029*/ "",
    /*030*/ "",
    /*031*/ "",
    /*032*/ "G101",
    /*033*/ "TC0190FMC",
    /*034*/ "BxROM or NINA-001",
    /*035*/ "",
    /*036*/ "",
    /*037*/ "",
    /*038*/ "",
    /*039*/ "",
    /*040*/ "",
    /*041*/ "",
    /*042*/ "",
    /*043*/ "",
    /*044*/ "Super Big 7-in-1",
    /*045*/ "1000000-in-1",
    /*046*/ "15-in-1",
    /*047*/ "",
    /*048*/ "TC0350FMR",
    /*049*/ "4-in-1",
    /*050*/ "",
    /*051*/ "",
    /*052*/ "Mario 7-in-1",
    /*053*/ "",
    /*054*/ "",
    /*055*/ "",
    /*056*/ "",
    /*057*/ "6-in-1",
    /*058*/ "68-in-1",
    /*059*/ "",
    /*060*/ "Reset Based 4-in-1",
    /*061*/ "20-in-1",
    /*062*/ "Super 700-in-1",
    /*063*/ "",
    /*064*/ "Tengen RAMBO-1",
    /*065*/ "",
    /*066*/ "GxROM",
    /*067*/ "Sunsoft-1",
    /*068*/ "Sunsoft-4",
    /*069*/ "Sunsoft-5B",
    /*070*/ "",
    /*071*/ "",
    /*072*/ "",
    /*073*/ "VRC3",
    /*074*/ "",
    /*075*/ "VRC1",
    /*076*/ "",
    /*077*/ "",
    /*078*/ "",
    /*079*/ "",
    /*080*/ "",
    /*081*/ "",
    /*082*/ "",
    /*083*/ "",
    /*084*/ "",
    /*085*/ "VRC7",
    /*086*/ "",
    /*087*/ "",
    /*088*/ "",
    /*089*/ "",
    /*090*/ "",
    /*091*/ "",
    /*092*/ "",
    /*093*/ "",
    /*094*/ "",
    /*095*/ "",
    /*096*/ "",
    /*097*/ "",
    /*098*/ "",
    /*099*/ "",
    /*100*/ "",
    /*101*/ "",
    /*102*/ "",
    /*103*/ "",
    /*104*/ "",
    /*105*/ "NES-EVENT",
    /*106*/ "",
    /*107*/ "",
    /*108*/ "",
    /*109*/ "",
    /*110*/ "",
    /*111*/ "",
    /*112*/ "",
    /*113*/ "",
    /*114*/ "",
    /*115*/ "",
    /*116*/ "",
    /*117*/ "",
    /*118*/ "MMC3 modified",
    /*119*/ "TQROM",
    /*120*/ "",
    /*121*/ "",
    /*122*/ "",
    /*123*/ "",
    /*124*/ "",
    /*125*/ "",
    /*126*/ "",
    /*127*/ "",
    /*128*/ "",
    /*129*/ "",
    /*130*/ "",
    /*131*/ "",
    /*132*/ "",
    /*133*/ "",
    /*134*/ "",
    /*135*/ "",
    /*136*/ "",
    /*137*/ "",
    /*138*/ "",
    /*139*/ "",
    /*140*/ "",
    /*141*/ "",
    /*142*/ "",
    /*143*/ "",
    /*144*/ "",
    /*145*/ "",
    /*146*/ "",
    /*147*/ "",
    /*148*/ "",
    /*149*/ "",
    /*150*/ "",
    /*151*/ "",
    /*152*/ "",
    /*153*/ "",
    /*154*/ "",
    /*155*/ "",
    /*156*/ "",
    /*157*/ "",
    /*158*/ "",
    /*159*/ "",
    /*160*/ "",
    /*161*/ "",
    /*162*/ "",
    /*163*/ "",
    /*164*/ "",
    /*165*/ "",
    /*166*/ "",
    /*167*/ "",
    /*168*/ "",
    /*169*/ "",
    /*170*/ "",
    /*171*/ "",
    /*172*/ "",
    /*173*/ "",
    /*174*/ "",
    /*175*/ "",
    /*176*/ "",
    /*177*/ "",
    /*178*/ "",
    /*179*/ "",
    /*180*/ "",
    /*181*/ "",
    /*182*/ "",
    /*183*/ "",
    /*184*/ "",
    /*185*/ "",
    /*186*/ "",
    /*187*/ "",
    /*188*/ "",
    /*189*/ "",
    /*190*/ "",
    /*191*/ "",
    /*192*/ "",
    /*193*/ "",
    /*194*/ "",
    /*195*/ "",
    /*196*/ "",
    /*197*/ "",
    /*198*/ "",
    /*199*/ "",
    /*200*/ "1200-in-1",
    /*201*/ "21-in-1",
    /*202*/ "",
    /*203*/ "35-in-1",
    /*204*/ "",
    /*205*/ "3-in-1",
    /*206*/ "",
    /*207*/ "",
    /*208*/ "",
    /*209*/ "",
    /*210*/ "",
    /*211*/ "",
    /*212*/ "",
    /*213*/ "",
    /*214*/ "",
    /*215*/ "",
    /*216*/ "",
    /*217*/ "",
    /*218*/ "",
    /*219*/ "",
    /*220*/ "",
    /*221*/ "",
    /*222*/ "",
    /*223*/ "",
    /*224*/ "",
    /*225*/ "58-in-1",
    /*226*/ "76-in-1",
    /*227*/ "1200-in-1",
    /*228*/ "4520",
    /*229*/ "",
    /*230*/ "22-in-1",
    /*231*/ "20-in-1",
    /*232*/ "",
    /*233*/ "42-in-1",
    /*234*/ "",
    /*235*/ "",
    /*236*/ "",
    /*237*/ "BMC-TELETUBBIES", // http://nesdev.parodius.com/bbs/viewtopic.php?t=5977
    /*238*/ "",
    /*239*/ "",
    /*240*/ "",
    /*241*/ "",
    /*242*/ "",
    /*243*/ "",
    /*244*/ "",
    /*245*/ "",
    /*246*/ "",
    /*247*/ "",
    /*248*/ "",
    /*249*/ "",
    /*250*/ "",
    /*251*/ "",
    /*252*/ "",
    /*253*/ "",
    /*254*/ "",
    /*255*/ "",
}

// Load a .nes file
func OpenINES(f *os.File) (*Cartridge) {
    header := new(NesfileHeader)
    cart := new(Cartridge)

    binary.Read(f, binary.LittleEndian, header)

    if header.Magic != NESFILE_MAGIC {
        return nil
        //panic(fmt.Sprintf("invalid nesfile signature: %x != %x", header.Magic, NESFILE_MAGIC))
    }

    prgPageCount := int(header.PrgPageCount)
    chrPageCount := int(header.ChrPageCount)

    if prgPageCount == 0 {
        // Special case for 4 MB PRG
        prgPageCount = 256
    }

    // This is weird, but the mapper nibbles are spread across two bytes
    mapperCode := int(header.Flags7 & 0xf0 | header.Flags6 & 0xf0 >> 4)
    submapperCode := 0

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
        cart.Trainer = readChunk(f, TRAINER_SIZE)
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
        mapperCode |= int(header.Flags8) & 0x0f << 11
        submapperCode = int(header.Flags8) & 0x0f >> 4

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
    
    //fmt.Printf("ROM: %d, VROM: %d, Mapper %s (%d:%d)\n", prgPageCount, chrPageCount, cart.MapperName, mapperCode, submapperCode)

    cart.Prg = readChunk(f, PRG_PAGE_SIZE * prgPageCount)
    if chrPageCount != 0 {
        cart.Chr = readChunk(f, CHR_PAGE_SIZE * chrPageCount)
    }

    if cart.Platform == PlatformPlayChoice {
        // "PlayChoice-10 (8KB of Hint Screen data stored after CHR data)"
        cart.HintScreen = readChunk(f, HINTSCREEN_SIZE)
    }

    // Decode the mapper code to a descriptive name (TODO: use INesMapperCode2Name)
    switch mapperCode {
    case 0:
        switch len(cart.Prg) {
        case 0x4000: cart.MapperName = "NES-NROM-128"
        case 0x8000: cart.MapperName = "NES-NROM-256"
        default: panic(fmt.Sprintf("invalid NROM size: %.4x\n", len(cart.Prg)))
        }
    case 3:
        cart.MapperName = "CNROM"       // TODO: official name?
    default:
        cart.MapperName = fmt.Sprintf("iNes Mapper %d:%d", mapperCode, submapperCode)
        //fmt.Printf("WARNING: nesfile: unrecognized iNES mapper %d (sub %d)\n", mapperCode, submapperCode)
    }

    //fmt.Printf("ROM: %d, VROM: %d, Mapper %s (%d:%d)\n", len(cart.Prg), len(cart.Chr), cart.MapperName, mapperCode, submapperCode)

    f.Close()

    return cart
}

// Read fixed-size chunk of data 
func readChunk(buffer *os.File, size int) ([]byte) {
    data := make([]byte, size)

    readLength, err := buffer.Read(data)
    if err != nil {
        panic(fmt.Sprintf("nesfile read error (when reading %d): %s", size, err))
    }
    if readLength != size { 
        panic(fmt.Sprintf("nesfile read incomplete: %d != %d", readLength, size))
    }

    return data
}

