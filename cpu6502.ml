(* Created:2010114
 * By Jeff Connelly
 *
 * 6502 microprocessor
 *)

type opcode = ADC | AND | ASL | BCC | BCS | BEQ | BIT | BMI | BNE | BPL | BRK | BVC | 
    BVS | CLC | CLD | CLI | CLV | CMP | CPX | CPY | DEC | DEX | DEY | EOR | INC | INX | 
    INY | JMP | JSR | LDA | LDX | LDY | LSR | NOP | ORA | PHA | PHP | PLA | PLP | ROL | 
    ROR | RTI | RTS | SBC | SEC | SED | SEI | STA | STX | STY | TAX | TAY | TSX | TXA | 
    TXS | TYA |
    (* undefined / invalid / undocumented TODO: http://nesdev.parodius.com/undocumented_opcodes.txt *)
    U__;;

type addr_mode = Imm | Zer | Ixz | Iyz | Abs | Inx | Iny | Pre | Pst | Imp | Acc | Ind | Rel;;

(* Opcode byte to opcode and addressing mode
Note: http://nesdev.parodius.com/6502.txt has several errors. 
http://www.akk.org/~flo/6502%20OpCode%20Disass.pdf is more correct, notably:
0x7d is ADC, Iny
0x8d is STA, Abs
0x90 is BCC, Rel
*)

let opcode_map = [|
(* Indexed by opcode, value is (mneumonic, addressing mode code) *)
(* x0         x1         x2         x3         x4        x5          x6         x7   *)
(* x8         x9         xa         xb         xc        xd          xe         xf   *)
(BRK, Imp);(ORA, Pre);(U__, Imp);(U__, Imp);(U__, Imp);(ORA, Zer);(ASL, Zer);(U__, Imp); (* 0x *)
(PHP, Imp);(ORA, Imm);(ASL, Acc);(U__, Imp);(U__, Imp);(ORA, Abs);(ASL, Abs);(U__, Imp); 
(BPL, Rel);(ORA, Pst);(U__, Imp);(U__, Imp);(U__, Imp);(ORA, Ixz);(ASL, Ixz);(U__, Imp); (* 1x *)
(CLC, Imp);(ORA, Iny);(U__, Imp);(U__, Imp);(U__, Imp);(U__, Imp);(ASL, Inx);(U__, Imp); 
(JSR, Abs);(AND, Pre);(U__, Imp);(U__, Imp);(BIT, Zer);(AND, Zer);(ROL, Zer);(U__, Imp); (* 2x *)
(PLP, Imp);(AND, Imm);(ROL, Acc);(U__, Imp);(BIT, Abs);(AND, Abs);(ROL, Abs);(U__, Imp); 
(BMI, Rel);(AND, Pst);(U__, Imp);(U__, Imp);(U__, Imp);(AND, Ixz);(ROL, Ixz);(U__, Imp); (* 3x *)
(SEC, Imp);(AND, Iny);(U__, Imp);(U__, Imp);(U__, Imp);(AND, Inx);(ROL, Inx);(U__, Imp); 
(RTI, Imp);(EOR, Pre);(U__, Imp);(U__, Imp);(U__, Imp);(EOR, Zer);(LSR, Zer);(U__, Imp); (* 4x *)
(PHA, Imp);(EOR, Imm);(LSR, Acc);(U__, Imp);(JMP, Abs);(EOR, Abs);(LSR, Abs);(U__, Imp);
(BVC, Rel);(EOR, Pst);(U__, Imp);(U__, Imp);(U__, Imp);(EOR, Ixz);(LSR, Ixz);(U__, Imp); (* 5x *)
(CLI, Imp);(EOR, Iny);(U__, Imp);(U__, Imp);(U__, Imp);(U__, Imp);(LSR, Inx);(U__, Imp);
(RTS, Imp);(ADC, Pre);(U__, Imp);(U__, Imp);(U__, Imp);(ADC, Zer);(ROR, Zer);(U__, Imp); (* 6x *)
(PLA, Imp);(ADC, Imm);(ROR, Acc);(U__, Imp);(JMP, Ind);(U__, Imp);(ROR, Abs);(U__, Imp);
(BVS, Rel);(ADC, Pst);(U__, Imp);(U__, Imp);(U__, Imp);(ADC, Ixz);(ROR, Ixz);(U__, Imp); (* 7x *)
(SEI, Imp);(ADC, Iny);(U__, Imp);(U__, Imp);(U__, Imp);(ADC, Iny);(ROR, Inx);(U__, Imp);
(U__, Imp);(STA, Pre);(U__, Imp);(U__, Imp);(STY, Zer);(STA, Zer);(STX, Zer);(U__, Imp); (* 8x *)
(DEY, Imp);(U__, Imp);(TXA, Imp);(U__, Imp);(STY, Abs);(STA, Abs);(STX, Abs);(U__, Imp);
(BCC, Rel);(STA, Pst);(U__, Imp);(U__, Imp);(STY, Ixz);(STA, Ixz);(STX, Ixz);(U__, Imp); (* 9x *)
(TYA, Imp);(STA, Iny);(TXS, Imp);(U__, Imp);(U__, Imp);(U__, Imp);(U__, Imp);(U__, Imp);
(LDY, Imm);(LDA, Pre);(LDX, Imm);(U__, Imp);(LDY, Zer);(LDA, Zer);(LDX, Zer);(U__, Imp); (* ax *)
(TAY, Imp);(LDA, Imm);(TAX, Imp);(U__, Imp);(LDY, Abs);(LDA, Abs);(LDX, Abs);(U__, Imp);
(BCS, Rel);(LDA, Pst);(U__, Imp);(U__, Imp);(LDY, Ixz);(LDA, Ixz);(LDX, Iyz);(U__, Imp); (* bx *)
(CLV, Imp);(LDA, Iny);(TSX, Imp);(U__, Imp);(LDY, Inx);(LDA, Inx);(LDX, Iny);(U__, Imp);
(CPY, Imm);(CMP, Pre);(U__, Imp);(U__, Imp);(CPY, Zer);(CMP, Zer);(DEC, Zer);(U__, Imp); (* cx *)
(INY, Imp);(CMP, Imm);(DEX, Imp);(U__, Imp);(CPY, Abs);(CMP, Abs);(DEC, Abs);(U__, Imp);
(BNE, Rel);(CMP, Pst);(U__, Imp);(U__, Imp);(U__, Imp);(CMP, Ixz);(DEC, Ixz);(U__, Imp); (* dx *)
(CLD, Imp);(CMP, Iny);(U__, Imp);(U__, Imp);(U__, Imp);(CMP, Inx);(DEC, Inx);(U__, Imp);
(CPX, Imm);(SBC, Pre);(U__, Imp);(U__, Imp);(CPX, Zer);(SBC, Zer);(INC, Zer);(U__, Imp); (* ex *)
(INX, Imp);(SBC, Imm);(NOP, Imp);(U__, Imp);(CPX, Ixz);(SBC, Abs);(INC, Abs);(U__, Imp);
(BEQ, Rel);(SBC, Pst);(U__, Imp);(U__, Imp);(U__, Imp);(SBC, Ixz);(INC, Ixz);(U__, Imp); (* fx *)
(SED, Imp);(SBC, Iny);(U__, Imp);(U__, Imp);(U__, Imp);(SBC, Inx);(INC, Inx);(U__, Imp);
|];;

(* This is lame, but OCaml doesn't have reflection like Python f.func_name *)
let string_of_opcode opcode = 
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
let operandBytesForMode addr_mode =
    match addr_mode with
    | Imm -> 1 | Zer -> 1 | Ixz -> 1 | Iyz -> 1 | Abs -> 2 | Inx -> 2 | Iny -> 2
    | Pre -> 1 | Pst -> 1 | Imp -> 0 | Acc -> 0 | Ind -> 2 | Rel -> 1;;

let read_operand addr_mode io = 
    match addr_mode with
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

let nameOfMode addr_mode =
    match addr_mode with
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


let string_of_operand addr_mode operand =
    match addr_mode with
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

type instruction = {opcode: opcode; addr_mode: addr_mode; opcode_byte: int; operand: int};;

(* Read and decode one instruction *)
let read_instruction io = 
    let opcode_byte = IO.read_byte io in                         (* Integer of opcode *)
    let opcode, addr_mode = Array.get opcode_map opcode_byte in  (* opcode and addr_mode variant types *)
    let operand = read_operand addr_mode io in                   (* Integer of operand *)

    {opcode=opcode; addr_mode=addr_mode; opcode_byte=opcode_byte; operand=operand};;

let string_of_instruction instr =
    if instr.opcode != U__ then
        (string_of_opcode instr.opcode) ^ " " ^ (string_of_operand instr.addr_mode instr.operand)
    else
        Printf.sprintf ".DB #$%.2X" instr.opcode_byte
    ;;


