// Created:20110203
// By Jeff Connelly

// Game cartridge abstraction

package cartridge

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
    Platform Platform               // Platform to run on
    DisplayType DisplayType
    MapperName string               // Name of mapper, authoritive over MapperCode
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

