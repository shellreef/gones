// Created:20101126
// By Jeff Connelly

// Disassemble 6502 instructions
// Provides data used by cpu6502 NextInstruction

package dis6502

import (
    "fmt"
)

// Operation code, a string for easy printing
type Opcode string
const (U__="???"; // undefined / invalid / undocumented TODO: http://nesdev.parodius.com/undocumented_opcodes.txt 
ADC="ADC"; AND="AND"; ASL="ASL"; BCS="BCS"; BEQ="BEQ"; BIT="BIT"; BMI="BMI"; 
BPL="BPL"; BVC="BVC"; BCC="BCC";
BNE="BNE"; BRK="BRK"; BVS="BVS"; CLC="CLC"; CLD="CLD"; CLI="CLI"; CLV="CLV"; 
CMP="CMP"; CPX="CPX"; CPY="CPY"; DEC="DEC"; DEX="DEX"; DEY="DEY"; EOR="EOR"; 
INC="INC"; INX="INX"; INY="INY"; JMP="JMP"; JSR="JSR"; LDA="LDA"; LDX="LDX"; 
LDY="LDY"; LSR="LSR"; NOP="NOP"; ORA="ORA"; PHA="PHA"; PHP="PHP"; PLA="PLA"; 
PLP="PLP"; ROL="ROL"; ROR="ROR"; RTI="RTI"; RTS="RTS"; SBC="SBC"; SEC="SEC"; 
SED="SED"; SEI="SEI"; STA="STA"; STX="STX"; STY="STY"; TAX="TAX"; TAY="TAY"; 
TSX="TSX"; TXA="TXA"; TXS="TXS"; TYA="TYA";

// Undocumented http://nesdev.parodius.com/undocumented_opcodes.txt
AAC="AAC"; AAX="AAX"; ARR="ARR"; ATX="ATX"; AXA="AXA"; AXS="AXS"; DCP="DCP";
DOP="DOP"; ISC="ISC"; KIL="KIL"; LAR="LAR"; LAX="LAX"; RLA="RLA"; RRA="RRA";
SLO="SLO"; SRE="SRE"; SXA="SXA"; SYA="SYA"; TOP="TOP"; XXA="XXA"; XAS="XAS";
)

// Addressing mode
type AddrMode string
const (Imd="Imd"; Zpg="Zpg"; Zpx="Zpx"; Zpy="Zpy"; Abs="Abs"; Abx="Abx"; 
Aby="Aby"; Ndx="Ndx"; Ndy="Ndy"; Imp="Imp"; Acc="Acc"; Ind="Ind"; Rel="Rel");

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
http://www.akk.org/~flo/6502%20OpCode%20Disass.pdf is more correct, notably:
0x7d is ADC, Aby
0x8d is STA, Abs
0x90 is BCC, Rel

Shoud also include undocumented opcodes, resources:
http://nesdev.parodius.com/undocumented_opcodes.txt
http://nesdev.parodius.com/extra_instructions.txt
http://www.nvg.org/bbc/doc/6502.txt (supersedes http://nesdev.parodius.com/6502_cpu.txt)
*/
// Indexed by opcode number, maps to decoded opcode and addressing mode
var Opcodes = [...]OpcodeAddrMode{
// Indexed by opcode, value is (mneumonic, addressing mode code) 
// x0         x1         x2         x3         x4        x5          x6         x7   
// x8         x9         xa         xb         xc        xd          xe         xf   
{BRK, Imp},{ORA, Ndx},{U__, Imp},{U__, Imp},{U__, Imp},{ORA, Zpg},{ASL, Zpg},{U__, Imp}, // 0x 
{PHP, Imp},{ORA, Imd},{ASL, Acc},{U__, Imp},{U__, Imp},{ORA, Abs},{ASL, Abs},{U__, Imp}, 
{BPL, Rel},{ORA, Ndy},{U__, Imp},{U__, Imp},{U__, Imp},{ORA, Zpx},{ASL, Zpx},{U__, Imp}, // 1x 
{CLC, Imp},{ORA, Aby},{U__, Imp},{U__, Imp},{U__, Imp},{U__, Imp},{ASL, Abx},{U__, Imp}, 
{JSR, Abs},{AND, Ndx},{U__, Imp},{U__, Imp},{BIT, Zpg},{AND, Zpg},{ROL, Zpg},{U__, Imp}, // 2x 
{PLP, Imp},{AND, Imd},{ROL, Acc},{U__, Imp},{BIT, Abs},{AND, Abs},{ROL, Abs},{U__, Imp}, 
{BMI, Rel},{AND, Ndy},{U__, Imp},{U__, Imp},{U__, Imp},{AND, Zpx},{ROL, Zpx},{U__, Imp}, // 3x 
{SEC, Imp},{AND, Aby},{U__, Imp},{U__, Imp},{U__, Imp},{AND, Abx},{ROL, Abx},{U__, Imp}, 
{RTI, Imp},{EOR, Ndx},{U__, Imp},{U__, Imp},{U__, Imp},{EOR, Zpg},{LSR, Zpg},{U__, Imp}, // 4x 
{PHA, Imp},{EOR, Imd},{LSR, Acc},{U__, Imp},{JMP, Abs},{EOR, Abs},{LSR, Abs},{U__, Imp},
{BVC, Rel},{EOR, Ndy},{U__, Imp},{U__, Imp},{U__, Imp},{EOR, Zpx},{LSR, Zpx},{U__, Imp}, // 5x 
{CLI, Imp},{EOR, Aby},{U__, Imp},{U__, Imp},{U__, Imp},{U__, Imp},{LSR, Abx},{U__, Imp},
{RTS, Imp},{ADC, Ndx},{U__, Imp},{RRA, Ndx},{U__, Imp},{ADC, Zpg},{ROR, Zpg},{RRA, Zpg}, // 6x 
{PLA, Imp},{ADC, Imd},{ROR, Acc},{U__, Imp},{JMP, Ind},{ADC, Abs},{ROR, Abs},{RRA, Abs},
{BVS, Rel},{ADC, Ndy},{U__, Imp},{RRA, Ndy},{U__, Imp},{ADC, Zpx},{ROR, Zpx},{RRA, Zpx}, // 7x 
{SEI, Imp},{ADC, Aby},{U__, Imp},{RRA, Aby},{U__, Imp},{ADC, Aby},{ROR, Abx},{RRA, Abx},
{U__, Imp},{STA, Ndx},{U__, Imp},{U__, Imp},{STY, Zpg},{STA, Zpg},{STX, Zpg},{U__, Imp}, // 8x 
{DEY, Imp},{U__, Imp},{TXA, Imp},{U__, Imp},{STY, Abs},{STA, Abs},{STX, Abs},{U__, Imp},
{BCC, Rel},{STA, Ndy},{U__, Imp},{U__, Imp},{STY, Zpx},{STA, Zpx},{STX, Zpx},{U__, Imp}, // 9x 
{TYA, Imp},{STA, Aby},{TXS, Imp},{U__, Imp},{U__, Imp},{U__, Imp},{U__, Imp},{U__, Imp},
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
{CLC, Imp},{ORA, Aby},{U__, Imp},{U__, Imp},{U__, Imp},{U__, Imp},{ASL, Abx},{U__, Imp}, 
{JSR, Abs},{AND, Ndx},{U__, Imp},{U__, Imp},{BIT, Zpg},{AND, Zpg},{ROL, Zpg},{U__, Imp}, // 2x 
{PLP, Imp},{AND, Imd},{ROL, Acc},{U__, Imp},{BIT, Abs},{AND, Abs},{ROL, Abs},{U__, Imp}, 
{BMI, Rel},{AND, Ndy},{U__, Imp},{U__, Imp},{U__, Imp},{AND, Zpx},{ROL, Zpx},{U__, Imp}, // 3x 
{SEC, Imp},{AND, Aby},{U__, Imp},{U__, Imp},{U__, Imp},{AND, Abx},{ROL, Abx},{U__, Imp}, 
{RTI, Imp},{EOR, Ndx},{U__, Imp},{U__, Imp},{U__, Imp},{EOR, Zpg},{LSR, Zpg},{U__, Imp}, // 4x 
{PHA, Imp},{EOR, Imd},{LSR, Acc},{U__, Imp},{JMP, Abs},{EOR, Abs},{LSR, Abs},{U__, Imp},
{BVC, Rel},{EOR, Ndy},{U__, Imp},{U__, Imp},{U__, Imp},{EOR, Zpx},{LSR, Zpx},{U__, Imp}, // 5x 
{CLI, Imp},{EOR, Aby},{U__, Imp},{U__, Imp},{U__, Imp},{U__, Imp},{LSR, Abx},{U__, Imp},
{RTS, Imp},{ADC, Ndx},{U__, Imp},{U__, Imp},{U__, Imp},{ADC, Zpg},{ROR, Zpg},{U__, Imp}, // 6x 
{PLA, Imp},{ADC, Imd},{ROR, Acc},{U__, Imp},{JMP, Ind},{ADC, Abs},{ROR, Abs},{U__, Imp},
{BVS, Rel},{ADC, Ndy},{U__, Imp},{U__, Imp},{U__, Imp},{ADC, Zpx},{ROR, Zpx},{U__, Imp}, // 7x 
{SEI, Imp},{ADC, Aby},{U__, Imp},{U__, Imp},{U__, Imp},{ADC, Aby},{ROR, Abx},{U__, Imp},
{U__, Imp},{STA, Ndx},{U__, Imp},{U__, Imp},{STY, Zpg},{STA, Zpg},{STX, Zpg},{U__, Imp}, // 8x 
{DEY, Imp},{U__, Imp},{TXA, Imp},{U__, Imp},{STY, Abs},{STA, Abs},{STX, Abs},{U__, Imp},
{BCC, Rel},{STA, Ndy},{U__, Imp},{U__, Imp},{STY, Zpx},{STA, Zpx},{STX, Zpx},{U__, Imp}, // 9x 
{TYA, Imp},{STA, Aby},{TXS, Imp},{U__, Imp},{U__, Imp},{U__, Imp},{U__, Imp},{U__, Imp},
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

    return fmt.Sprintf("%s %s", instr.Opcode, instr.AddrMode.formatOperand(instr.Operand))
}
