// Created:20101126
// By Jeff Connelly

// Disassemble 6502 instructions
// Provides data used by cpu6502 NextInstruction

package dis6502

import (
    "fmt"
)

// Operation code mneumonic
// Important: this is NOT the opcode byte, since that also encodes
// the addressing mode. It is only the operation, and the numerical
// value has no relation to the opcode byte in the instruction. I had
// it as a string, but changed it to an integer for convenience. To
// get the mneumonic value, use String().
// Note 2: some unofficial ops have multiple mneumonics; only one is here (see String())
type Opcode int
const (U__=iota;
ADC; AND; ASL; BCS; BEQ; BIT; BMI; BPL; BVC; BCC; BNE; BRK; BVS; CLC; CLD; CLI; CLV; 
CMP; CPX; CPY; DEC; DEX; DEY; EOR; INC; INX; INY; JMP; JSR; LDA; LDX; LDY; LSR; NOP; 
ORA; PHA; PHP; PLA; PLP; ROL; ROR; RTI; RTS; SBC; SEC; SED; SEI; STA; STX; STY; TAX; 
TAY; TSX; TXA; TXS; TYA; AAC; SAX; ARR; ASR; ATX; AXA; AXS; DCP; TOP; DOP; ISB; KIL; 
LAR; LAX; RLA; RRA; SLO; SRE; SXA; SYA; XAA; XAS;
)

// Addressing mode
type AddrMode int
const (
    Imd=iota;      // Immediate
    Zpg;      // Zero Page
    Zpx;      // Zero Page,X
    Zpy;      // Zero Page,Y
    Abs;      // Absolute
    Abx;      // Absolute,X
    Aby;      // Absolute, Y
    Ndx;      // (Indirect,X)
    Ndy;      // (Indirect),Y
    Imp;      // Implied
    Acc;      // Accumulator
    Ind;      // (Indirect)
    Rel;      // Relative
);
/*type AddrMode string
const (
    Imd="Imd";      // Immediate
    Zpg="Zpg";      // Zero Page
    Zpx="Zpx";      // Zero Page,X
    Zpy="Zpy";      // Zero Page,Y
    Abs="Abs";      // Absolute
    Abx="Abx";      // Absolute,X
    Aby="Aby";      // Absolute, Y
    Ndx="Ndx";      // (Indirect,X)
    Ndy="Ndy";      // (Indirect),Y
    Imp="Imp";      // Implied
    Acc="Acc";      // Accumulator
    Ind="Ind";      // (Indirect)
    Rel="Rel";      // Relative
);*/


// Opcode and addressing mode for opcode definition table
type OpcodeAddrMode struct { Opcode Opcode; AddrMode AddrMode }

// Instruction with operand
type Instruction struct { 
    Opcode Opcode       // Mneumonic opcode constant
    OpcodeByte uint8    // Value of opcode
    AddrMode AddrMode   // Addressing mode constant
    Operand int         // Operand value, if applicable
    Official bool       // Whether is an official opcode; otherwise undocumented
}

/* Opcode byte to opcode and addressing mode
Note: http://nesdev.parodius.com/6502.txt has several errors. 
http://www.zophar.net/fileuploads/2/10532krzvs/6502.txt is updated with fixes, notably
0x7d is ADC, Aby
0x8d is STA, Abs
0x90 is BCC, Rel
http://www.akk.org/~flo/6502%20OpCode%20Disass.pdf is also correct
Chart, but doesn't have undoc: http://e-tradition.net/bytes/6502/6502_instruction_set.html
Comprehensive chart with undoc: http://www.xmission.com/~trevin/atari/6502_opcode_table.html - uses different mneumonics

Shoud also include undocumented opcodes, resources:
http://nesdev.parodius.com/undocumented_opcodes.txt - I'm using the first opcode mneumonic from here when possible
http://nesdev.parodius.com/extra_instructions.txt - has better details on operation, but not as popular mneumonics
http://www.nvg.org/bbc/doc/6502.txt (supersedes http://nesdev.parodius.com/6502_cpu.txt)
Nice succinct tables: http://apple1.chez.com/Apple1project/Docs/m6502/6502-6510-8500-8502%20Opcodes.htm
*/
// Indexed by opcode number, maps to decoded opcode and addressing mode
var Opcodes = [...]OpcodeAddrMode{
// Indexed by opcode, value is (mneumonic, addressing mode code) 
// x0         x1         x2         x3         x4        x5          x6         x7   
// x8         x9         xa         xb         xc        xd          xe         xf   
{BRK, Imp},{ORA, Ndx},{KIL, Imp},{SLO, Ndx},{DOP, Zpg},{ORA, Zpg},{ASL, Zpg},{SLO, Zpg}, // 0x 
{PHP, Imp},{ORA, Imd},{ASL, Acc},{AAC, Imd},{TOP, Abs},{ORA, Abs},{ASL, Abs},{SLO, Abs}, 
{BPL, Rel},{ORA, Ndy},{KIL, Imp},{SLO, Ndy},{DOP, Zpx},{ORA, Zpx},{ASL, Zpx},{SLO, Zpx}, // 1x 
{CLC, Imp},{ORA, Aby},{NOP, Imp},{SLO, Aby},{TOP, Abx},{ORA, Abx},{ASL, Abx},{SLO, Abx}, 
{JSR, Abs},{AND, Ndx},{KIL, Imp},{RLA, Ndx},{BIT, Zpg},{AND, Zpg},{ROL, Zpg},{RLA, Zpg}, // 2x 
{PLP, Imp},{AND, Imd},{ROL, Acc},{AAC, Imd},{BIT, Abs},{AND, Abs},{ROL, Abs},{RLA, Abs}, 
{BMI, Rel},{AND, Ndy},{KIL, Imp},{RLA, Ndy},{DOP, Zpx},{AND, Zpx},{ROL, Zpx},{RLA, Zpx}, // 3x 
{SEC, Imp},{AND, Aby},{NOP, Imp},{RLA, Aby},{TOP, Abx},{AND, Abx},{ROL, Abx},{RLA, Abx}, 
{RTI, Imp},{EOR, Ndx},{KIL, Imp},{SRE, Ndx},{DOP, Zpg},{EOR, Zpg},{LSR, Zpg},{SRE, Zpg}, // 4x 
{PHA, Imp},{EOR, Imd},{LSR, Acc},{ASR, Imd},{JMP, Abs},{EOR, Abs},{LSR, Abs},{SRE, Abs},
{BVC, Rel},{EOR, Ndy},{KIL, Imp},{SRE, Ndy},{DOP, Zpx},{EOR, Zpx},{LSR, Zpx},{SRE, Zpx}, // 5x 
{CLI, Imp},{EOR, Aby},{NOP, Imp},{SRE, Aby},{TOP, Abx},{EOR, Abx},{LSR, Abx},{SRE, Abx},
{RTS, Imp},{ADC, Ndx},{KIL, Imp},{RRA, Ndx},{DOP, Zpg},{ADC, Zpg},{ROR, Zpg},{RRA, Zpg}, // 6x 
{PLA, Imp},{ADC, Imd},{ROR, Acc},{ARR, Imd},{JMP, Ind},{ADC, Abs},{ROR, Abs},{RRA, Abs},
{BVS, Rel},{ADC, Ndy},{KIL, Imp},{RRA, Ndy},{DOP, Zpx},{ADC, Zpx},{ROR, Zpx},{RRA, Zpx}, // 7x 
{SEI, Imp},{ADC, Aby},{NOP, Imp},{RRA, Aby},{TOP, Abx},{ADC, Abx},{ROR, Abx},{RRA, Abx},
{DOP, Imd},{STA, Ndx},{DOP, Imd},{SAX, Ndx},{STY, Zpg},{STA, Zpg},{STX, Zpg},{SAX, Zpg}, // 8x 
{DEY, Imp},{DOP, Imd},{TXA, Imp},{XAA, Imd},{STY, Abs},{STA, Abs},{STX, Abs},{SAX, Abs},
{BCC, Rel},{STA, Ndy},{KIL, Imp},{AXA, Ndy},{STY, Zpx},{STA, Zpx},{STX, Zpy},{SAX, Zpy}, // 9x 
{TYA, Imp},{STA, Aby},{TXS, Imp},{XAS, Aby},{SYA, Abx},{STA, Abx},{SXA, Aby},{AXA, Aby},
{LDY, Imd},{LDA, Ndx},{LDX, Imd},{LAX, Ndx},{LDY, Zpg},{LDA, Zpg},{LDX, Zpg},{LAX, Zpg}, // ax 
{TAY, Imp},{LDA, Imd},{TAX, Imp},{ATX, Imd},{LDY, Abs},{LDA, Abs},{LDX, Abs},{LAX, Abs},
{BCS, Rel},{LDA, Ndy},{KIL, Imp},{LAX, Ndy},{LDY, Zpx},{LDA, Zpx},{LDX, Zpy},{LAX, Zpy}, // bx 
{CLV, Imp},{LDA, Aby},{TSX, Imp},{LAR, Aby},{LDY, Abx},{LDA, Abx},{LDX, Aby},{LAX, Aby},
{CPY, Imd},{CMP, Ndx},{DOP, Imd},{DCP, Ndx},{CPY, Zpg},{CMP, Zpg},{DEC, Zpg},{DCP, Zpg}, // cx 
{INY, Imp},{CMP, Imd},{DEX, Imp},{AXS, Imd},{CPY, Abs},{CMP, Abs},{DEC, Abs},{DCP, Abs},
{BNE, Rel},{CMP, Ndy},{KIL, Imp},{DCP, Ndy},{DOP, Zpx},{CMP, Zpx},{DEC, Zpx},{DCP, Zpx}, // dx 
{CLD, Imp},{CMP, Aby},{NOP, Imp},{DCP, Aby},{TOP, Abx},{CMP, Abx},{DEC, Abx},{DCP, Abx},
{CPX, Imd},{SBC, Ndx},{DOP, Imd},{ISB, Ndx},{CPX, Zpg},{SBC, Zpg},{INC, Zpg},{ISB, Zpg}, // ex 
{INX, Imp},{SBC, Imd},{NOP, Imp},{SBC, Imd},{CPX, Abs},{SBC, Abs},{INC, Abs},{ISB, Abs},
{BEQ, Rel},{SBC, Ndy},{KIL, Imp},{ISB, Ndy},{DOP, Zpx},{SBC, Zpx},{INC, Zpx},{ISB, Zpx}, // fx 
{SED, Imp},{SBC, Aby},{NOP, Imp},{ISB, Aby},{TOP, Abx},{SBC, Abx},{INC, Abx},{ISB, Abx},
}

// Excludes http://nesdev.parodius.com/undocumented_opcodes.txt
// Useful to guess whether the code being executed is in fact code.. or an "illegal" opcode.
// However, undocumented/unofficial/"illegal" opcodes are functional so should be emulated,
// especially since some games (though rarely) or maybe Game Genie codes may rely on them.
var OfficialOpcodes = [...]OpcodeAddrMode{
// x0         x1         x2         x3         x4        x5          x6         x7   
// x8         x9         xa         xb         xc        xd          xe         xf   
{BRK, Imp},{ORA, Ndx},{U__, Imp},{U__, Imp},{U__, Imp},{ORA, Zpg},{ASL, Zpg},{U__, Imp}, // 0x 
{PHP, Imp},{ORA, Imd},{ASL, Acc},{U__, Imp},{U__, Imp},{ORA, Abs},{ASL, Abs},{U__, Imp}, 
{BPL, Rel},{ORA, Ndy},{U__, Imp},{U__, Imp},{U__, Imp},{ORA, Zpx},{ASL, Zpx},{U__, Imp}, // 1x 
{CLC, Imp},{ORA, Aby},{U__, Imp},{U__, Imp},{U__, Imp},{ORA, Abx},{ASL, Abx},{U__, Imp}, 
{JSR, Abs},{AND, Ndx},{U__, Imp},{U__, Imp},{BIT, Zpg},{AND, Zpg},{ROL, Zpg},{U__, Imp}, // 2x 
{PLP, Imp},{AND, Imd},{ROL, Acc},{U__, Imp},{BIT, Abs},{AND, Abs},{ROL, Abs},{U__, Imp}, 
{BMI, Rel},{AND, Ndy},{U__, Imp},{U__, Imp},{U__, Imp},{AND, Zpx},{ROL, Zpx},{U__, Imp}, // 3x 
{SEC, Imp},{AND, Aby},{U__, Imp},{U__, Imp},{U__, Imp},{AND, Abx},{ROL, Abx},{U__, Imp}, 
{RTI, Imp},{EOR, Ndx},{U__, Imp},{U__, Imp},{U__, Imp},{EOR, Zpg},{LSR, Zpg},{U__, Imp}, // 4x 
{PHA, Imp},{EOR, Imd},{LSR, Acc},{U__, Imp},{JMP, Abs},{EOR, Abs},{LSR, Abs},{U__, Imp},
{BVC, Rel},{EOR, Ndy},{U__, Imp},{U__, Imp},{U__, Imp},{EOR, Zpx},{LSR, Zpx},{U__, Imp}, // 5x 
{CLI, Imp},{EOR, Aby},{U__, Imp},{U__, Imp},{U__, Imp},{EOR, Abx},{LSR, Abx},{U__, Imp},
{RTS, Imp},{ADC, Ndx},{U__, Imp},{U__, Imp},{U__, Imp},{ADC, Zpg},{ROR, Zpg},{U__, Imp}, // 6x 
{PLA, Imp},{ADC, Imd},{ROR, Acc},{U__, Imp},{JMP, Ind},{ADC, Abs},{ROR, Abs},{U__, Imp},
{BVS, Rel},{ADC, Ndy},{U__, Imp},{U__, Imp},{U__, Imp},{ADC, Zpx},{ROR, Zpx},{U__, Imp}, // 7x 
{SEI, Imp},{ADC, Aby},{U__, Imp},{U__, Imp},{U__, Imp},{ADC, Abx},{ROR, Abx},{U__, Imp},
{U__, Imp},{STA, Ndx},{U__, Imp},{U__, Imp},{STY, Zpg},{STA, Zpg},{STX, Zpg},{U__, Imp}, // 8x 
{DEY, Imp},{U__, Imp},{TXA, Imp},{U__, Imp},{STY, Abs},{STA, Abs},{STX, Abs},{U__, Imp},
{BCC, Rel},{STA, Ndy},{U__, Imp},{U__, Imp},{STY, Zpx},{STA, Zpx},{STX, Zpy},{U__, Imp}, // 9x 
{TYA, Imp},{STA, Aby},{TXS, Imp},{U__, Imp},{U__, Imp},{STA, Abx},{U__, Imp},{U__, Imp},
{LDY, Imd},{LDA, Ndx},{LDX, Imd},{U__, Imp},{LDY, Zpg},{LDA, Zpg},{LDX, Zpg},{U__, Imp}, // ax 
{TAY, Imp},{LDA, Imd},{TAX, Imp},{U__, Imp},{LDY, Abs},{LDA, Abs},{LDX, Abs},{U__, Imp},
{BCS, Rel},{LDA, Ndy},{U__, Imp},{U__, Imp},{LDY, Zpx},{LDA, Zpx},{LDX, Zpy},{U__, Imp}, // bx 
{CLV, Imp},{LDA, Aby},{TSX, Imp},{U__, Imp},{LDY, Abx},{LDA, Abx},{LDX, Aby},{U__, Imp},
{CPY, Imd},{CMP, Ndx},{U__, Imp},{U__, Imp},{CPY, Zpg},{CMP, Zpg},{DEC, Zpg},{U__, Imp}, // cx 
{INY, Imp},{CMP, Imd},{DEX, Imp},{U__, Imp},{CPY, Abs},{CMP, Abs},{DEC, Abs},{U__, Imp},
{BNE, Rel},{CMP, Ndy},{U__, Imp},{U__, Imp},{U__, Imp},{CMP, Zpx},{DEC, Zpx},{U__, Imp}, // dx 
{CLD, Imp},{CMP, Aby},{U__, Imp},{U__, Imp},{U__, Imp},{CMP, Abx},{DEC, Abx},{U__, Imp},
{CPX, Imd},{SBC, Ndx},{U__, Imp},{U__, Imp},{CPX, Zpg},{SBC, Zpg},{INC, Zpg},{U__, Imp}, // ex 
{INX, Imp},{SBC, Imd},{NOP, Imp},{U__, Imp},{CPX, Abs},{SBC, Abs},{INC, Abs},{U__, Imp},
{BEQ, Rel},{SBC, Ndy},{U__, Imp},{U__, Imp},{U__, Imp},{SBC, Zpx},{INC, Zpx},{U__, Imp}, // fx 
{SED, Imp},{SBC, Aby},{U__, Imp},{U__, Imp},{U__, Imp},{SBC, Abx},{INC, Abx},{U__, Imp},
}

// Get opcode mneumonic as a string
// This is admittedly ugly.. any better way?
// (besides making Opcode a string type)
func (op Opcode) String() (string) {
    switch op {
    case ADC: return "ADC" 
    case AND: return "AND" 
    case ASL: return "ASL" 
    case BCS: return "BCS" 
    case BEQ: return "BEQ" 
    case BIT: return "BIT" 
    case BMI: return "BMI" 
    case BPL: return "BPL" 
    case BVC: return "BVC" 
    case BCC: return "BCC" 
    case BNE: return "BNE" 
    case BRK: return "BRK" 
    case BVS: return "BVS" 
    case CLC: return "CLC" 
    case CLD: return "CLD" 
    case CLI: return "CLI" 
    case CLV: return "CLV" 
    case CMP: return "CMP" 
    case CPX: return "CPX" 
    case CPY: return "CPY" 
    case DEC: return "DEC" 
    case DEX: return "DEX" 
    case DEY: return "DEY" 
    case EOR: return "EOR" 
    case INC: return "INC" 
    case INX: return "INX" 
    case INY: return "INY" 
    case JMP: return "JMP" 
    case JSR: return "JSR" 
    case LDA: return "LDA" 
    case LDX: return "LDX" 
    case LDY: return "LDY" 
    case LSR: return "LSR" 
    case NOP: return "NOP" 
    case ORA: return "ORA" 
    case PHA: return "PHA" 
    case PHP: return "PHP" 
    case PLA: return "PLA" 
    case PLP: return "PLP" 
    case ROL: return "ROL" 
    case ROR: return "ROR" 
    case RTI: return "RTI" 
    case RTS: return "RTS" 
    case SBC: return "SBC" 
    case SEC: return "SEC" 
    case SED: return "SED" 
    case SEI: return "SEI" 
    case STA: return "STA" 
    case STX: return "STX" 
    case STY: return "STY" 
    case TAX: return "TAX" 
    case TAY: return "TAY" 
    case TSX: return "TSX" 
    case TXA: return "TXA" 
    case TXS: return "TXS" 
    case TYA: return "TYA" 
    // Unofficial
    case AAC: return "AAC" // aka ANC
    case SAX: return "SAX" // aka AAX
    case ARR: return "ARR" 
    case ASR: return "ASR" // aka ALR
    case ATX: return "ATX" // aka LXA, OAL
    case AXA: return "AXA" // aka SHA
    case AXS: return "AXS" // aka SBX
    case DCP: return "DCP" 
    case TOP: return "TOP" // aka NOP
    case DOP: return "DOP" // aka NOP
    case ISB: return "ISB" // aka ISC
    case KIL: return "KIL" // aka HLT
    case LAR: return "LAR" 
    case LAX: return "LAX" 
    case RLA: return "RLA" 
    case RRA: return "RRA" 
    case SLO: return "SLO" 
    case SRE: return "SRE" 
    case SXA: return "SXA" // aka SHX, XAS
    case SYA: return "SYA" // aka SHY
    case XAA: return "XAA" 
    case XAS: return "XAS" // aka SHS
    }
    return "???"
}

func (addrMode AddrMode) formatOperand(operand int) (string) {
    switch addrMode {
    case Imd: return fmt.Sprintf("#$%.2X", operand)
    case Zpg: return fmt.Sprintf("$%.2X", operand)
    case Zpx: return fmt.Sprintf("$%.2X,X", operand)
    case Zpy: return fmt.Sprintf("$%.2X,Y", operand)
    case Abs: return fmt.Sprintf("$%.4X", operand)
    case Abx: return fmt.Sprintf("$%.4X,X", operand)
    case Aby: return fmt.Sprintf("$%.4X,Y", operand)
    case Ndx: return fmt.Sprintf("($%.2X,X)", operand)
    case Ndy: return fmt.Sprintf("($%.2X),Y", operand)
    case Imp: return ""
    case Acc: return "A"
    case Ind: return fmt.Sprintf("($%.4X)", operand)
    case Rel: return fmt.Sprintf("$%.4X", operand)   // TODO: return sign_num(operand)+offset, it really needs to be relative current offset 
    }
    panic(fmt.Sprintf("fotmatOperand unknown addressing mode: %s", addrMode))
}

// Get number of bytes an operand addressing mode requires, useful for disassembly
// These are read by cpu.NextOperand()
func (addrMode AddrMode) OperandSize() (int) {
    switch addrMode {
    case Imd, Zpg, Zpx, Zpy, Ndx, Ndy, Rel: 
        return 1
    case Abs, Abx, Aby, Ind:
        return 2
    case Imp, Acc: 
        return 0
    }
    panic(fmt.Sprintf("readOperand unknown addressing mode: %s", addrMode))
} 


func (instr Instruction) String() (string) {
     if instr.Opcode == U__ {
        return fmt.Sprintf(".DB #$%.2X", instr.OpcodeByte)
    } 

    undoc := ""
    if !instr.Official {
        // Denote undocumented opcodes with asterisk since they may be unintentional
        undoc = "*"
    }

    return fmt.Sprintf("%s%s %s", undoc, instr.Opcode, instr.AddrMode.formatOperand(instr.Operand))
}

// Find the opcode byte for an opcode mneumonic/addressing mode pair (basically, assemble)
// NOTE: since ops are not ambiguous - instead use hex in code to be explicit w/ a comment
func OpcodeByteFor(opcode Opcode, addrMode AddrMode) (uint8) {
    // TODO: find multiple opcodes, for unofficial?
    for op, entry := range Opcodes {
        if entry.Opcode == opcode && entry.AddrMode == addrMode {
            return uint8(op)
        }
    }
    return 0
}
