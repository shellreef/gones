(* Created:2010114
 * By Jeff Connelly
 *
 * Disassemble 6502 instructions
 *)

type opcode = ADC | AND | ASL | BCC | BCS | BEQ | BIT | BMI | BNE | BPL | BRK | BVC | 
    BVS | CLC | CLD | CLI | CLV | CMP | CPX | CPY | DEC | DEX | DEY | EOR | INC | INX | 
    INY | JMP | JSR | LDA | LDX | LDY | LSR | NOP | ORA | PHA | PHP | PLA | PLP | ROL | 
    ROR | RTI | RTS | SBC | SEC | SED | SEI | STA | STX | STY | TAX | TAY | TSX | TXA | 
    TXS | TYA |
    (* undefined / invalid / undocumented TODO: http://nesdev.parodius.com/undocumented_opcodes.txt *)
    U__;;

type addr_mode = Imd | Zpg | Zpx | Zpy | Abs | Abx | Aby | Ndx | Ndy | Imp | Acc | Ind | Rel;;

(* Opcode byte to opcode and addressing mode
Note: http://nesdev.parodius.com/6502.txt has several errors. 
http://www.akk.org/~flo/6502%20OpCode%20Disass.pdf is more correct, notably:
0x7d is ADC, Aby
0x8d is STA, Abs
0x90 is BCC, Rel
*)

let opcode_map = [|
(* Indexed by opcode, value is (mneumonic, addressing mode code) *)
(* x0         x1         x2         x3         x4        x5          x6         x7   *)
(* x8         x9         xa         xb         xc        xd          xe         xf   *)
(BRK, Imp);(ORA, Ndx);(U__, Imp);(U__, Imp);(U__, Imp);(ORA, Zpg);(ASL, Zpg);(U__, Imp); (* 0x *)
(PHP, Imp);(ORA, Imd);(ASL, Acc);(U__, Imp);(U__, Imp);(ORA, Abs);(ASL, Abs);(U__, Imp); 
(BPL, Rel);(ORA, Ndy);(U__, Imp);(U__, Imp);(U__, Imp);(ORA, Zpx);(ASL, Zpx);(U__, Imp); (* 1x *)
(CLC, Imp);(ORA, Aby);(U__, Imp);(U__, Imp);(U__, Imp);(U__, Imp);(ASL, Abx);(U__, Imp); 
(JSR, Abs);(AND, Ndx);(U__, Imp);(U__, Imp);(BIT, Zpg);(AND, Zpg);(ROL, Zpg);(U__, Imp); (* 2x *)
(PLP, Imp);(AND, Imd);(ROL, Acc);(U__, Imp);(BIT, Abs);(AND, Abs);(ROL, Abs);(U__, Imp); 
(BMI, Rel);(AND, Ndy);(U__, Imp);(U__, Imp);(U__, Imp);(AND, Zpx);(ROL, Zpx);(U__, Imp); (* 3x *)
(SEC, Imp);(AND, Aby);(U__, Imp);(U__, Imp);(U__, Imp);(AND, Abx);(ROL, Abx);(U__, Imp); 
(RTI, Imp);(EOR, Ndx);(U__, Imp);(U__, Imp);(U__, Imp);(EOR, Zpg);(LSR, Zpg);(U__, Imp); (* 4x *)
(PHA, Imp);(EOR, Imd);(LSR, Acc);(U__, Imp);(JMP, Abs);(EOR, Abs);(LSR, Abs);(U__, Imp);
(BVC, Rel);(EOR, Ndy);(U__, Imp);(U__, Imp);(U__, Imp);(EOR, Zpx);(LSR, Zpx);(U__, Imp); (* 5x *)
(CLI, Imp);(EOR, Aby);(U__, Imp);(U__, Imp);(U__, Imp);(U__, Imp);(LSR, Abx);(U__, Imp);
(RTS, Imp);(ADC, Ndx);(U__, Imp);(U__, Imp);(U__, Imp);(ADC, Zpg);(ROR, Zpg);(U__, Imp); (* 6x *)
(PLA, Imp);(ADC, Imd);(ROR, Acc);(U__, Imp);(JMP, Ind);(U__, Imp);(ROR, Abs);(U__, Imp);
(BVS, Rel);(ADC, Ndy);(U__, Imp);(U__, Imp);(U__, Imp);(ADC, Zpx);(ROR, Zpx);(U__, Imp); (* 7x *)
(SEI, Imp);(ADC, Aby);(U__, Imp);(U__, Imp);(U__, Imp);(ADC, Aby);(ROR, Abx);(U__, Imp);
(U__, Imp);(STA, Ndx);(U__, Imp);(U__, Imp);(STY, Zpg);(STA, Zpg);(STX, Zpg);(U__, Imp); (* 8x *)
(DEY, Imp);(U__, Imp);(TXA, Imp);(U__, Imp);(STY, Abs);(STA, Abs);(STX, Abs);(U__, Imp);
(BCC, Rel);(STA, Ndy);(U__, Imp);(U__, Imp);(STY, Zpx);(STA, Zpx);(STX, Zpx);(U__, Imp); (* 9x *)
(TYA, Imp);(STA, Aby);(TXS, Imp);(U__, Imp);(U__, Imp);(U__, Imp);(U__, Imp);(U__, Imp);
(LDY, Imd);(LDA, Ndx);(LDX, Imd);(U__, Imp);(LDY, Zpg);(LDA, Zpg);(LDX, Zpg);(U__, Imp); (* ax *)
(TAY, Imp);(LDA, Imd);(TAX, Imp);(U__, Imp);(LDY, Abs);(LDA, Abs);(LDX, Abs);(U__, Imp);
(BCS, Rel);(LDA, Ndy);(U__, Imp);(U__, Imp);(LDY, Zpx);(LDA, Zpx);(LDX, Zpy);(U__, Imp); (* bx *)
(CLV, Imp);(LDA, Aby);(TSX, Imp);(U__, Imp);(LDY, Abx);(LDA, Abx);(LDX, Aby);(U__, Imp);
(CPY, Imd);(CMP, Ndx);(U__, Imp);(U__, Imp);(CPY, Zpg);(CMP, Zpg);(DEC, Zpg);(U__, Imp); (* cx *)
(INY, Imp);(CMP, Imd);(DEX, Imp);(U__, Imp);(CPY, Abs);(CMP, Abs);(DEC, Abs);(U__, Imp);
(BNE, Rel);(CMP, Ndy);(U__, Imp);(U__, Imp);(U__, Imp);(CMP, Zpx);(DEC, Zpx);(U__, Imp); (* dx *)
(CLD, Imp);(CMP, Aby);(U__, Imp);(U__, Imp);(U__, Imp);(CMP, Abx);(DEC, Abx);(U__, Imp);
(CPX, Imd);(SBC, Ndx);(U__, Imp);(U__, Imp);(CPX, Zpg);(SBC, Zpg);(INC, Zpg);(U__, Imp); (* ex *)
(INX, Imp);(SBC, Imd);(NOP, Imp);(U__, Imp);(CPX, Zpx);(SBC, Abs);(INC, Abs);(U__, Imp);
(BEQ, Rel);(SBC, Ndy);(U__, Imp);(U__, Imp);(U__, Imp);(SBC, Zpx);(INC, Zpx);(U__, Imp); (* fx *)
(SED, Imp);(SBC, Aby);(U__, Imp);(U__, Imp);(U__, Imp);(SBC, Abx);(INC, Abx);(U__, Imp);
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

(* Bytes after opcode which operand requires for each addressing mode (not actually used) *)
let operand_bytes_for_mode addr_mode =
    match addr_mode with
    | Imd -> 1 | Zpg -> 1 | Zpx -> 1 | Zpy -> 1 | Abs -> 2 | Abx -> 2 | Aby -> 2
    | Ndx -> 1 | Ndy -> 1 | Imp -> 0 | Acc -> 0 | Ind -> 2 | Rel -> 1;;

let read_operand addr_mode io = 
    match addr_mode with
    | Imd -> IO.read_byte io
    | Zpg -> IO.read_byte io
    | Zpx -> IO.read_byte io
    | Zpy -> IO.read_byte io
    | Abs -> IO.read_ui16 io
    | Abx -> IO.read_ui16 io
    | Aby -> IO.read_ui16 io
    | Ndx -> IO.read_byte io
    | Ndy -> IO.read_byte io
    | Imp -> 0
    | Acc -> 0 
    | Ind -> IO.read_ui16 io
    | Rel -> IO.read_signed_byte io;;

let name_of_mode addr_mode =
    match addr_mode with
    | Imd -> "Imdediate"
    | Zpg -> "Zpgo Page"
    | Zpx -> "Indexed X Zpgo Page"
    | Zpy -> "Indexed Y Zpgo Page"
    | Abs -> "Absolute"
    | Abx -> "Indexed X"
    | Aby -> "Indexed Y"
    | Ndx -> "Ndx-indexed Indirect"
    | Ndy -> "Post-indexed Indirect"
    | Imp -> "Implied"
    | Acc -> "Accumulator"
    | Ind -> "Indirect"
    | Rel -> "Relative";;


let string_of_operand addr_mode operand =
    match addr_mode with
    | Imd -> Printf.sprintf "#$%.2X" operand
    | Zpg -> Printf.sprintf "$%.2X" operand
    | Zpx -> Printf.sprintf "$%.2X,X" operand
    | Zpy -> Printf.sprintf "$%.2X,X" operand
    | Abs -> Printf.sprintf "#$%.4X" operand
    | Abx -> Printf.sprintf "$%.4X,X" operand
    | Aby -> Printf.sprintf "$%.4X,Y" operand
    | Ndx -> Printf.sprintf "($%.2X,X)" operand
    | Ndy -> Printf.sprintf "($%.2X),Y" operand
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


