

type opcode = ADC | AND | ASL | BCC | BCS | BEQ | BIT | BMI | BNE | BPL | BRK | BVC | 
    BVS | CLC | CLD | CLI | CLV | CMP | CPX | CPY | DEC | DEX | DEY | EOR | INC | INX | 
    INY | JMP | JSR | LDA | LDX | LDY | LSR | NOP | ORA | PHA | PHP | PLA | PLP | ROL | 
    ROR | RTI | RTS | SBC | SEC | SED | SEI | STA | STX | STY | TAX | TAY | TSX | TXA | 
    TXS | TYA |
    (* undefined / invalid / undocumented TODO: http://nesdev.parodius.com/undocumented_opcodes.txt *)
    U__;;


(* http://nesdev.parodius.com/6502.txt *)
let opcodeMap = [|
(* Indexed by opcode, value is (mneumonic, addressing mode code) *)
(* x0      x1       x2       x3       x4      x5        x6       x7   *)
(* x8      x9       xa       xb       xc      xd        xe       xf   *)
(BRK, 9);(ORA, 7);(U__, 0);(U__, 0);(U__, 0);(ORA, 1);(ASL, 1);(U__, 0); (* 0x *)
(PHP, 9);(ORA, 0);(ASL,10);(U__, 0);(U__, 0);(ORA, 4);(ASL, 4);(U__, 0); 
(BPL, 9);(ORA, 8);(U__, 0);(U__, 0);(U__, 0);(ORA, 2);(ASL, 2);(U__, 0); (* 1x *)
(CLC, 9);(ORA, 6);(U__, 0);(U__, 0);(U__, 0);(U__, 0);(ASL, 5);(U__, 0); 
(JSR, 4);(AND, 7);(U__, 0);(U__, 0);(BIT, 1);(AND, 1);(ROL, 1);(U__, 0); (* 2x *)
(PLP, 9);(AND, 0);(ROL,10);(U__, 0);(BIT, 4);(AND, 4);(ROL, 4);(U__, 0); 
(BMI,12);(AND, 8);(U__, 0);(U__, 0);(U__, 0);(AND, 2);(ROL, 2);(U__, 0); (* 3x *)
(SEC, 9);(AND, 6);(U__, 0);(U__, 0);(U__, 0);(AND, 5);(ROL, 5);(U__, 0); 
(EOR, 4);(EOR, 7);(U__, 0);(U__, 0);(U__, 0);(EOR, 1);(LSR, 1);(U__, 0); (* 4x *)
(PHA, 9);(EOR, 0);(LSR,10);(U__, 0);(JMP, 4);(RTI, 9);(LSR, 4);(U__, 0);
(EOR, 5);(EOR, 8);(U__, 0);(U__, 0);(U__, 0);(EOR, 2);(LSR, 2);(U__, 0); (* 5x *)
(CLI, 9);(EOR, 6);(U__, 0);(U__, 0);(U__, 0);(U__, 0);(LSR, 5);(U__, 0);
(RTS, 9);(ADC, 7);(U__, 0);(U__, 0);(U__, 0);(ADC, 1);(ROR, 1);(U__, 0); (* 6x *)
(PLA, 9);(ADC, 0);(ROR,10);(U__, 0);(JMP,12);(U__, 0);(ROR, 4);(U__, 0);
(BVS,12);(ADC, 8);(U__, 0);(U__, 0);(U__, 0);(ADC, 2);(ROR, 2);(U__, 0); (* 7x *)
(SEI, 9);(ADC, 6);(U__, 0);(U__, 0);(U__, 0);(U__, 0);(ROR, 5);(U__, 0);
(STA, 4);(STA, 7);(U__, 0);(U__, 0);(STY, 1);(STA, 1);(STX, 1);(U__, 0); (* 8x *)
(DEY, 9);(U__, 0);(TXA, 9);(U__, 0);(STY, 4);(U__, 0);(STX, 4);(U__, 0);
(STA, 5);(STA, 8);(U__, 0);(U__, 0);(STY, 2);(STA, 2);(STX, 2);(U__, 0); (* 9x *)
(TYA, 9);(STA, 6);(TXS, 9);(U__, 0);(U__, 0);(U__, 0);(U__, 0);(U__, 0);
(LDY, 0);(LDA, 7);(LDX, 0);(U__, 0);(LDY, 1);(LDA, 1);(LDX, 1);(U__, 0); (* ax *)
(TAY, 9);(LDA, 0);(TAX, 9);(U__, 0);(LDY, 4);(LDA, 4);(LDX, 4);(U__, 0);
(BCS,12);(LDA, 8);(U__, 0);(U__, 0);(LDY, 2);(LDA, 2);(LDX, 3);(U__, 0); (* bx *)
(CLV, 9);(LDA, 6);(TSX, 9);(U__, 0);(LDY, 5);(LDA, 5);(LDX, 6);(U__, 0);
(CPY, 0);(CMP, 7);(U__, 0);(U__, 0);(CPY, 1);(CMP, 1);(DEC, 1);(U__, 0); (* cx *)
(INY, 9);(CMP, 0);(DEX, 9);(U__, 0);(CPY, 2);(CMP, 4);(DEC, 4);(U__, 0);
(BNE,12);(CMP, 8);(U__, 0);(U__, 0);(U__, 0);(CMP, 2);(DEC, 2);(U__, 0); (* dx *)
(CLD, 9);(CMP, 6);(U__, 0);(U__, 0);(U__, 0);(CMP, 5);(DEC, 5);(U__, 0);
(CPX, 0);(SBC, 7);(U__, 0);(U__, 0);(CPX, 1);(SBC, 1);(INC, 1);(U__, 0); (* ex *)
(INX, 9);(SBC, 0);(NOP, 9);(U__, 0);(CPX, 2);(SBC, 4);(INC, 4);(U__, 0);
(BEQ,12);(SBC, 8);(U__, 0);(U__, 0);(U__, 0);(SBC, 2);(INC, 2);(U__, 0); (* fx *)
(SED, 9);(SBC, 6);(U__, 0);(U__, 0);(U__, 0);(SBC, 5);(INC, 5);(U__, 0);
|];;

let addressingModes = [|
    (* (* index *) operand bytes, name, TODO: formatter *)
    (* 0 *) (1, "Immediate");               (* sprintf '#$%.2X',$_[0]} *)
    (* 1 *) (1, "Zero Page");               (* sprintf '$%.2X',$_[0]} *)
    (* 2 *) (1, "Indexed X Zero Page");     (* sprintf '$%.2X,X',$_[0]} *)
    (* 3 *) (1, "Indexed Y Zero Page");     (* sprintf '$%.2X,Y',$_[0]} *)
    (* 4 *) (2, "Absolute");                (* sprintf '$%.4X',($_[0])+(($_[1])*0x100)} *)
    (* 5 *) (2, "Indexed X");               (* sprintf '$%.4X,X',($_[0])+(($_[1])*0x100)} *)
    (* 6 *) (2, "Indexed Y");               (* sprintf '$%.4X,Y',($_[0])+(($_[1])*0x100)} *)
    (* 7 *) (1, "Pre-indexed Indirect");    (* sprintf '($%.2X,X)', $_[0] } *)
    (* 8 *) (1, "Post-indexed indirect");   (* sprintf '($%.2X),Y', $_[0] } *)
    (* 9 *) (0, "Implied");                 (* '' *)
    (*10 *) (0, "Accumulator");             (* 'A' *)
    (*11 *) (2, "Indirect");                (* sprintf '($%.4X)', ($_[0])+(($_[1])*0x100)} *)  (* JMP only *)
    (*12 *) (1, "Relative");                (* sprintf '$%.4X', sign_num($_[0])+$_[1] } *) 
|];;

