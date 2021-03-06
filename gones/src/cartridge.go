// Created:20110203
// By Jeff Connelly

// Game cartridge abstraction

package cartridge

import (
        "fmt"
        "os"
        "encoding/binary"

        "ppu2c02"
        "cpu6502"
        )

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
    Prg []byte                      // Program data 
    PrgRam [0x7ff]byte              // Cartridge RAM (could be save RAM) TODO: abstract out of nesfile
    Chr []byte                      // Character data
    SHA1 string                     // SHA-1 of PRG+CHR for identification purposes

    Platform Platform               // Platform to run on
    DisplayType DisplayType
    MapperName string               // Name of mapper
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

// Load program code into CPU
func (cart *Cartridge) LoadPRG(cpu *cpu6502.CPU) {
    if len(cart.Prg) == 0 {
        panic("No PRG found")
    }
   
    // TODO: use cartdb (see main.Load()) if possible; fall back to mapper from file otherwise
   
    // Wire up mappers
    switch cart.MapperName {
    case "NES-NROM-128":
        cpu.MapROM(0x8000, 0xbfff, cart.Prg, "NROM-128",          0x3fff, nil, nil)
        cpu.MapROM(0xc000, 0xffff, cart.Prg, "NROM-128 (mirror)", 0x3fff, nil, nil)

    case "NES-NROM-256":
        cpu.MapROM(0x8000, 0xffff, cart.Prg, "NROM-256", 0x7fff, nil, nil)

    case "CNROM":
        // Fixed PRG 
        cpu.MapROM(0x8000, 0xffff, cart.Prg, "CNROM", 0x7fff, nil, nil)
        // TODO: bankable CHR

    default:
        panic(fmt.Sprintf("LoadPRG: no support for mapper: %s!\n", cart.MapperName))

        // Take a wild guess, which is probably wrong
        cpu.MapROM(0x8000, 0xffff, cart.Prg, fmt.Sprintf("Unknown mapper %s", cart.MapperName), 0x7fff, nil, nil)
    }

    // SRAM (assumed)  TODO: assume nothing
    cpu.Map(0x6000, 0x7fff,
            func(address uint16)(value uint8) { return cart.PrgRam[address & 0x7ff] },
            func(address uint16, value uint8) { cart.PrgRam[address & 0x7ff] = value },
            "Cartridge RAM")

    // Initialize to reset vector.. note, don't use ReadUInt16 since it adds CPU cycles!
    // TODO: instead, really should call cpu.Reset()
    pcl := cpu.ReadFrom(cpu6502.RESET_VECTOR)
    pch := cpu.ReadFrom(cpu6502.RESET_VECTOR + 1)
    cpu.PC = uint16(pch) << 8 + uint16(pcl)

}

// Initialize the PPU with a game cartridge
func (cart *Cartridge) LoadCHR(ppu *ppu2c02.PPU) {
    if len(cart.Chr) > 0 { 
        // Load CHR ROM pattern tables #0 and #1 from cartridge, if present, into $0000-1FFF
        // Note that not all games have this, some use CHR RAM instead, but 
        // in either case it is wired in the cartridge, not the NES itself
        // TODO: pointers instead of copying, so can switch CHR banks easily
        copy(ppu.Memory[0:0x2000], cart.Chr[0:0x2000])
        // TODO: what if there are more CHR banks, which one to load first?
    }
}

// Load a game cartridge from a file on disk
func LoadFile(filename string) (cart *Cartridge) {
    f, err := os.Open(filename, os.O_RDONLY, 0)
    if f == nil {
        panic(fmt.Sprintf("cannot open %s: %s", filename, err))
    }

    // Read beginning of file to identify what it is
    signature := make([]byte, 4) 
    //var signature [4]byte     // for some reason: binary.Read: invalid type [4]uint8
    err = binary.Read(f, binary.BigEndian, signature)
    if err != nil {
        panic(fmt.Sprintf("cannot read signature from %s: %s", filename, err))
    }

    // Go to beginning of file and read it completely
    f.Seek(0, 0)
    switch string(signature[:]) {
    case "NES\x1a": cart = OpenINES(f)
    case "UNIF": cart = OpenUNIF(f)
    // TODO: support zip, 7z, fds, split prg/chr
    default: panic(fmt.Sprintf("unrecognizable file format: %s, signature=%.8x", filename, signature))
    }

    if cart == nil {
        panic(fmt.Sprintf("failed to load: %s", filename))
    }


    return cart
}
