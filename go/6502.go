// Created:20101126
// By Jeff Connelly

//

package main

import "fmt"

// TODO: should these map to the opcodes themselves instead of opaque integers?? Probably.
type Opcode int 
const (___ = 0;
ADC; AND; ASL; BCS; BEQ; BIT; BMI; BNE; BRK; BVS; CLC; CLD; CLI; CLV; CMP; CPX; 
CPY; DEC; DEX; DEY; EOR; INC; INX; INY; JMP; JSR; LDA; LDX; LDY; LSR; NOP; ORA; 
PHA; PHP; PLA; PLP; ROL; ROR; RTI; RTS; SBC; SEC; SED; SEI; STA; STX; STY; TAX; 
TAY; TSX; TXA; TXS; TYA)

// TODO: AddrMode

// Opcode and addressing mode
type OpcodeAddrMode struct { opcode, addrmode int }

// Indexed by opcode number, maps to decoded opcode and addressing mode
var opcodes = [...]OpcodeAddrMode{
// ==0==     ==1==     ==2==    ==3==     ==4==     ==5==     ==6==     ==7==
{BRK,  9},{ORA,  7},{___,  0},{___,  0},{___,  0},{ORA,  1},{ASL,  1},{___,  0},
{PHP,  9},{ORA,  0},{ASL, 10},{___,  0},{___,  0},{ORA,  4},{ASL,  4},{___,  0},
{ORA,  5},{ORA,  8},{___,  0},{___,  0},{___,  0},{ORA,  2},{ASL,  2},{___,  0},
{CLC,  9},{ORA,  6},{___,  0},{___,  0},{___,  0},{___,  0},{ASL,  5},{___,  0},
{JSR,  4},{AND,  7},{___,  0},{___,  0},{BIT,  1},{AND,  1},{ROL,  1},{___,  0},
{PLP,  9},{AND,  0},{ROL, 10},{___,  0},{BIT,  4},{AND,  4},{ROL,  4},{___,  0},
{BMI, 12},{AND,  8},{___,  0},{___,  0},{___,  0},{AND,  2},{ROL,  2},{___,  0},
{SEC,  9},{AND,  6},{___,  0},{___,  0},{___,  0},{AND,  5},{ROL,  5},{___,  0},
{EOR,  4},{EOR,  7},{___,  0},{___,  0},{___,  0},{EOR,  1},{LSR,  1},{___,  0},
{PHA,  9},{EOR,  0},{LSR, 10},{___,  0},{JMP,  4},{RTI,  9},{LSR,  4},{___,  0},
{EOR,  5},{EOR,  8},{___,  0},{___,  0},{___,  0},{EOR,  2},{LSR,  2},{___,  0},
{CLI,  9},{EOR,  6},{___,  0},{___,  0},{___,  0},{___,  0},{LSR,  5},{___,  0},
{RTS,  9},{ADC,  7},{___,  0},{___,  0},{___,  0},{ADC,  1},{ROR,  1},{___,  0},
{PLA,  9},{ADC,  0},{ROR, 10},{___,  0},{JMP, 12},{___,  0},{ROR,  4},{___,  0},
{BVS, 12},{ADC,  8},{___,  0},{___,  0},{___,  0},{ADC,  2},{ROR,  2},{___,  0},
{SEI,  9},{ADC,  6},{___,  0},{___,  0},{___,  0},{___,  0},{ROR,  5},{___,  0},
{STA,  4},{STA,  7},{___,  0},{___,  0},{STY,  1},{STA,  1},{STX,  1},{___,  0},
{DEY,  9},{___,  0},{TXA,  9},{___,  0},{STY,  4},{___,  0},{STX,  4},{___,  0},
{STA,  5},{STA,  8},{___,  0},{___,  0},{STY,  2},{STA,  2},{STX,  2},{___,  0},
{TYA,  9},{STA,  6},{TXS,  9},{___,  0},{___,  0},{___,  0},{___,  0},{___,  0},
{LDY,  0},{LDA,  7},{LDX,  0},{___,  0},{LDY,  1},{LDA,  1},{LDX,  1},{___,  0},
{TAY,  9},{LDA,  0},{TAX,  9},{___,  0},{LDY,  4},{LDA,  4},{LDX,  4},{___,  0},
{BCS, 12},{LDA,  8},{___,  0},{___,  0},{LDY,  2},{LDA,  2},{LDX,  3},{___,  0},
{CLV,  9},{LDA,  6},{TSX,  9},{___,  0},{LDY,  5},{LDA,  5},{LDX,  6},{___,  0},
{CPY,  0},{CMP,  7},{___,  0},{___,  0},{CPY,  1},{CMP,  1},{DEC,  1},{___,  0},
{INY,  9},{CMP,  0},{DEX,  9},{___,  0},{CPY,  2},{CMP,  4},{DEC,  4},{___,  0},
{BNE, 12},{CMP,  8},{___,  0},{___,  0},{___,  0},{CMP,  2},{DEC,  2},{___,  0},
{CLD,  9},{CMP,  6},{___,  0},{___,  0},{___,  0},{CMP,  5},{DEC,  5},{___,  0},
{CPX,  0},{SBC,  7},{___,  0},{___,  0},{CPX,  1},{SBC,  1},{INC,  1},{___,  0},
{INX,  9},{SBC,  0},{NOP,  9},{___,  0},{CPX,  2},{SBC,  4},{INC,  4},{___,  0},
{BEQ, 12},{SBC,  8},{___,  0},{___,  0},{___,  0},{SBC,  2},{INC,  2},{___,  0},
{SED,  9},{SBC,  6},{___,  0},{___,  0},{___,  0},{SBC,  5},{INC,  5},{___,  0},
}

func main() {
    fmt.Printf("Hello, world!\n")
}

