type opcode = ADC | AND | ASL | BCC | BCS | BEQ | BIT | BMI | BNE | BPL | BRK | BVC | 
    BVS | CLC | CLD | CLI | CLV | CMP | CPX | CPY | DEC | DEX | DEY | EOR | INC | INX | 
    INY | JMP | JSR | LDA | LDX | LDY | LSR | NOP | ORA | PHA | PHP | PLA | PLP | ROL | 
    ROR | RTI | RTS | SBC | SEC | SED | SEI | STA | STX | STY | TAX | TAY | TSX | TXA | 
    TXS | TYA |
    (* undefined / invalid / undocumented TODO: http://nesdev.parodius.com/undocumented_opcodes.txt *)
    U__;;

type addrMode = Imm | Zer | Ixz | Iyz | Abs | Inx | Iny | Pre | Pst | Imp | Acc | Ind | Rel;;

(* http://nesdev.parodius.com/6502.txt *)
let opcodeMap = [|
(* Indexed by opcode, value is (mneumonic, addressing mode code) *)
(* x0         x1         x2         x3         x4        x5          x6         x7   *)
(* x8         x9         xa         xb         xc        xd          xe         xf   *)
(BRK, Imp);(ORA, Pre);(U__, Imm);(U__, Imm);(U__, Imm);(ORA, Zer);(ASL, Zer);(U__, Imm); (* 0x *)
(PHP, Imp);(ORA, Imm);(ASL, Acc);(U__, Imm);(U__, Imm);(ORA, Abs);(ASL, Abs);(U__, Imm); 
(BPL, Imp);(ORA, Pst);(U__, Imm);(U__, Imm);(U__, Imm);(ORA, Ixz);(ASL, Ixz);(U__, Imm); (* 1x *)
(CLC, Imp);(ORA, Iny);(U__, Imm);(U__, Imm);(U__, Imm);(U__, Imm);(ASL, Inx);(U__, Imm); 
(JSR, Abs);(AND, Pre);(U__, Imm);(U__, Imm);(BIT, Zer);(AND, Zer);(ROL, Zer);(U__, Imm); (* 2x *)
(PLP, Imp);(AND, Imm);(ROL, Acc);(U__, Imm);(BIT, Abs);(AND, Abs);(ROL, Abs);(U__, Imm); 
(BMI, Rel);(AND, Pst);(U__, Imm);(U__, Imm);(U__, Imm);(AND, Ixz);(ROL, Ixz);(U__, Imm); (* 3x *)
(SEC, Imp);(AND, Iny);(U__, Imm);(U__, Imm);(U__, Imm);(AND, Inx);(ROL, Inx);(U__, Imm); 
(EOR, Abs);(EOR, Pre);(U__, Imm);(U__, Imm);(U__, Imm);(EOR, Zer);(LSR, Zer);(U__, Imm); (* 4x *)
(PHA, Imp);(EOR, Imm);(LSR, Acc);(U__, Imm);(JMP, Abs);(RTI, Imp);(LSR, Abs);(U__, Imm);
(EOR, Inx);(EOR, Pst);(U__, Imm);(U__, Imm);(U__, Imm);(EOR, Ixz);(LSR, Ixz);(U__, Imm); (* 5x *)
(CLI, Imp);(EOR, Iny);(U__, Imm);(U__, Imm);(U__, Imm);(U__, Imm);(LSR, Inx);(U__, Imm);
(RTS, Imp);(ADC, Pre);(U__, Imm);(U__, Imm);(U__, Imm);(ADC, Zer);(ROR, Zer);(U__, Imm); (* 6x *)
(PLA, Imp);(ADC, Imm);(ROR, Acc);(U__, Imm);(JMP, Ind);(U__, Imm);(ROR, Abs);(U__, Imm);
(BVS, Rel);(ADC, Pst);(U__, Imm);(U__, Imm);(U__, Imm);(ADC, Ixz);(ROR, Ixz);(U__, Imm); (* 7x *)
(SEI, Imp);(ADC, Iny);(U__, Imm);(U__, Imm);(U__, Imm);(U__, Imm);(ROR, Inx);(U__, Imm);
(STA, Abs);(STA, Pre);(U__, Imm);(U__, Imm);(STY, Zer);(STA, Zer);(STX, Zer);(U__, Imm); (* 8x *)
(DEY, Imp);(U__, Imm);(TXA, Imp);(U__, Imm);(STY, Abs);(U__, Imm);(STX, Abs);(U__, Imm);
(STA, Inx);(STA, Pst);(U__, Imm);(U__, Imm);(STY, Ixz);(STA, Ixz);(STX, Ixz);(U__, Imm); (* 9x *)
(TYA, Imp);(STA, Iny);(TXS, Imp);(U__, Imm);(U__, Imm);(U__, Imm);(U__, Imm);(U__, Imm);
(LDY, Imm);(LDA, Pre);(LDX, Imm);(U__, Imm);(LDY, Zer);(LDA, Zer);(LDX, Zer);(U__, Imm); (* ax *)
(TAY, Imp);(LDA, Imm);(TAX, Imp);(U__, Imm);(LDY, Abs);(LDA, Abs);(LDX, Abs);(U__, Imm);
(BCS, Rel);(LDA, Pst);(U__, Imm);(U__, Imm);(LDY, Ixz);(LDA, Ixz);(LDX, Iyz);(U__, Imm); (* bx *)
(CLV, Imp);(LDA, Iny);(TSX, Imp);(U__, Imm);(LDY, Inx);(LDA, Inx);(LDX, Iny);(U__, Imm);
(CPY, Imm);(CMP, Pre);(U__, Imm);(U__, Imm);(CPY, Zer);(CMP, Zer);(DEC, Zer);(U__, Imm); (* cx *)
(INY, Imp);(CMP, Imm);(DEX, Imp);(U__, Imm);(CPY, Ixz);(CMP, Abs);(DEC, Abs);(U__, Imm);
(BNE, Rel);(CMP, Pst);(U__, Imm);(U__, Imm);(U__, Imm);(CMP, Ixz);(DEC, Ixz);(U__, Imm); (* dx *)
(CLD, Imp);(CMP, Iny);(U__, Imm);(U__, Imm);(U__, Imm);(CMP, Inx);(DEC, Inx);(U__, Imm);
(CPX, Imm);(SBC, Pre);(U__, Imm);(U__, Imm);(CPX, Zer);(SBC, Zer);(INC, Zer);(U__, Imm); (* ex *)
(INX, Imp);(SBC, Imm);(NOP, Imp);(U__, Imm);(CPX, Ixz);(SBC, Abs);(INC, Abs);(U__, Imm);
(BEQ, Rel);(SBC, Pst);(U__, Imm);(U__, Imm);(U__, Imm);(SBC, Ixz);(INC, Ixz);(U__, Imm); (* fx *)
(SED, Imp);(SBC, Iny);(U__, Imm);(U__, Imm);(U__, Imm);(SBC, Inx);(INC, Inx);(U__, Imm);
|];;

(* This is lame, but OCaml doesn't have reflection like Python f.func_name *)
let stringOfOpcode opcode = 
    match opcode with
    | ADC -> "ADC" | AND -> "AND" | ASL -> "ASL" | BCC -> "BCC" | BCS -> "BCS" | BEQ -> "BEQ" | BIT -> "BIT" | BMI -> "BMI"
    | BNE -> "BNE" | BPL -> "BPL" | BRK -> "BRK" | BVC -> "BVC" | BVS -> "BVS" | CLC -> "CLC" | CLD -> "CLD" | CLI -> "CLI"
    | CLV -> "CLV" | CMP -> "CMP" | CPX -> "CPX" | CPY -> "CPY" | DEC -> "DEC" | DEX -> "DEX" | DEY -> "DEY" | EOR -> "EOR"
    | INC -> "INC" | INX -> "INX" | INY -> "INY" | JMP -> "JMP" | JSR -> "JSR" | LDA -> "LDA" | LDX -> "LDX" | LDY -> "LDY"
    | LSR -> "LSR" | NOP -> "NOP" | ORA -> "ORA" | PHA -> "PHA" | PHP -> "PHP" | PLA -> "PLA" | PLP -> "PLP" | ROL -> "ROL"
    | ROR -> "ROR" | RTI -> "RTI" | RTS -> "RTS" | SBC -> "SBC" | SEC -> "SEC" | SED -> "SED" | SEI -> "SEI" | STA -> "STA"
    | STX -> "STX" | STY -> "STY" | TAX -> "TAX" | TAY -> "TAY" | TSX -> "TSX" | TXA -> "TXA" | TXS -> "TXS" | TYA -> "TYA"
    | U__ -> "???";;

(* Bytes after opcode which operand requires for each addressing mode *)
let operandBytesForMode mode =
    match mode with
    | Imm -> 1 | Zer -> 1 | Ixz -> 1 | Iyz -> 1 | Abs -> 2 | Inx -> 2 | Iny -> 2
    | Pre -> 1 | Pst -> 1 | Imp -> 0 | Acc -> 0 | Ind -> 2 | Rel -> 1;;

let readOperandForMode mode io = 
    match mode with
    | Imm -> IO.read_byte io
    | Zer -> IO.read_byte io
    | Ixz -> IO.read_byte io
    | Iyz -> IO.read_byte io
    | Abs -> IO.read_ui16 io
    | Inx -> IO.read_ui16 io
    | Iny -> IO.read_ui16 io
    | Pre -> IO.read_byte io
    | Pst -> IO.read_byte io
    | Imp -> 0
    | Acc -> 0 
    | Ind -> IO.read_ui16 io
    | Rel -> IO.read_signed_byte io;;

let nameOfMode mode =
    match mode with
    | Imm -> "Immediate"
    | Zer -> "Zero Page"
    | Ixz -> "Indexed X Zero Page"
    | Iyz -> "Indexed Y Zero Page"
    | Abs -> "Absolute"
    | Inx -> "Indexed X"
    | Iny -> "Indexed Y"
    | Pre -> "Pre-indexed Indirect"
    | Pst -> "Post-indexed Indirect"
    | Imp -> "Implied"
    | Acc -> "Accumulator"
    | Ind -> "Indirect"
    | Rel -> "Relative";;


let formatOperand mode operand =
    match mode with
    | Imm -> Printf.sprintf "#$%.2X" operand
    | Zer -> Printf.sprintf "$%.2X" operand
    | Ixz -> Printf.sprintf "$%.2X,X" operand
    | Iyz -> Printf.sprintf "$%.2X,X" operand
    | Abs -> Printf.sprintf "#$%.4X" operand
    | Inx -> Printf.sprintf "$%.4X,X" operand
    | Iny -> Printf.sprintf "$%.4X,Y" operand
    | Pre -> Printf.sprintf "($%.2X,X)" operand
    | Pst -> Printf.sprintf "($%.2X),Y" operand
    | Imp -> ""
    | Acc -> "A"
    | Ind -> Printf.sprintf "($%.4X)" operand
    | Rel -> Printf.sprintf "$%.4X" operand;;    (* TODO: sign_num(operand)+offset, it really needs to be relative current offset *)

(* This doesn't work because OCaml doesn't infer the match result is a format 
let formatOperand mode operand =
    Printf.sprintf (match mode with
    | Imm -> "#$%.2X" 
    | Zer -> "$%.2X"
    | Ixz -> "$%.2X,X"
    | Iyz -> "$%.2X,X" 
    | Abs -> "#$%.4X" 
    | Inx -> "$%.4X,X"
    | Iny -> "$%.4X,Y"
    | Pre -> "($%.2X,X)"
    | Pst -> "($%.2X),Y" 
    | Imp -> " [ignore: %x]"
    | Acc -> "A [ignore: %x]"
    | Ind -> "($%.4X)" 
    | Rel -> "$%.4X"
    ) operand;;
*)


(* Read and decode one instruction 
 * TODO: IO module file:///Users/jeff/Downloads/extlib-1.5.1/extlib-doc/IO.html
   so can easily read from string and keep track of position!
   *)
let readInstruction io = 
    let opcode, mode = Array.get opcodeMap (Char.code (IO.read io)) in
    let operand = readOperandForMode mode io in

    (stringOfOpcode opcode) ^ " " ^ (formatOperand mode operand);;

(* print_endline ((Cpu6502.stringOfOpcode (fst (Array.get Cpu6502.opcodeMap 0xa9))) ^ " " ^ (Cpu6502.formatOperand (snd (Array.get Cpu6502.opcodeMap 0xa9 )) 0x40));; *)
