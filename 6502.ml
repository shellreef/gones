let opADC = None;;
let opAND = None;;
let opASL = None;;
let opBCC = None;;
let opBCS = None;;
let opBEQ = None;;
let opBIT = None;;
let opBMI = None;;
let opBNE = None;;
let opBPL = None;;
let opBRK = None;;
let opBVC = None;;
let opBVS = None;;
let opCLC = None;;
let opCLD = None;;
let opCLI = None;;
let opCLV = None;;
let opCMP = None;;
let opCPX = None;;
let opCPY = None;;
let opDEC = None;;
let opDEX = None;;
let opDEY = None;;
let opEOR = None;;
let opINC = None;;
let opINX = None;;
let opINY = None;;
let opJMP = None;;
let opJSR = None;;
let opLDA = None;;
let opLDX = None;;
let opLDY = None;;
let opLSR = None;;
let opNOP = None;;
let opORA = None;;
let opPHA = None;;
let opPHP = None;;
let opPLA = None;;
let opPLP = None;;
let opROL = None;;
let opROR = None;;
let opRTI = None;;
let opRTS = None;;
let opSBC = None;;
let opSEC = None;;
let opSED = None;;
let opSEI = None;;
let opSTA = None;;
let opSTX = None;;
let opSTY = None;;
let opTAX = None;;
let opTAY = None;;
let opTSX = None;;
let opTXA = None;;
let opTXS = None;;
let opTYA = None;;
let op___ = None;;    (* undefined / future expansion / undocumented TODO: http://nesdev.parodius.com/undocumented_opcodes.txt *)

(* http://nesdev.parodius.com/6502.txt *)
let opcodes = [|
(* Indexed by opcode, value is (mneumonic, addressing mode code) *)
(*  x0        x1         x2         x3        x4          x5         x6        x7   *)
(*  x8        x9         xa         xb        xc          xd         xe        xf   *)
(opBRK, 9);(opORA, 7);(op___, 0);(op___, 0);(op___, 0);(opORA, 1);(opASL, 1);(op___, 0); (* 0x *)
(opPHP, 9);(opORA, 0);(opASL,10);(op___, 0);(op___, 0);(opORA, 4);(opASL, 4);(op___, 0); 
(opBPL, 9);(opORA, 8);(op___, 0);(op___, 0);(op___, 0);(opORA, 2);(opASL, 2);(op___, 0); (* 1x *)
(opCLC, 9);(opORA, 6);(op___, 0);(op___, 0);(op___, 0);(op___, 0);(opASL, 5);(op___, 0); 
(opJSR, 4);(opAND, 7);(op___, 0);(op___, 0);(opBIT, 1);(opAND, 1);(opROL, 1);(op___, 0); (* 2x *)
(opPLP, 9);(opAND, 0);(opROL,10);(op___, 0);(opBIT, 4);(opAND, 4);(opROL, 4);(op___, 0); 
(opBMI,12);(opAND, 8);(op___, 0);(op___, 0);(op___, 0);(opAND, 2);(opROL, 2);(op___, 0); (* 3x *)
(opSEC, 9);(opAND, 6);(op___, 0);(op___, 0);(op___, 0);(opAND, 5);(opROL, 5);(op___, 0); 
(opEOR, 4);(opEOR, 7);(op___, 0);(op___, 0);(op___, 0);(opEOR, 1);(opLSR, 1);(op___, 0); (* 4x *)
(opPHA, 9);(opEOR, 0);(opLSR,10);(op___, 0);(opJMP, 4);(opRTI, 9);(opLSR, 4);(op___, 0);
(opEOR, 5);(opEOR, 8);(op___, 0);(op___, 0);(op___, 0);(opEOR, 2);(opLSR, 2);(op___, 0); (* 5x *)
(opCLI, 9);(opEOR, 6);(op___, 0);(op___, 0);(op___, 0);(op___, 0);(opLSR, 5);(op___, 0);
(opRTS, 9);(opADC, 7);(op___, 0);(op___, 0);(op___, 0);(opADC, 1);(opROR, 1);(op___, 0); (* 6x *)
(opPLA, 9);(opADC, 0);(opROR,10);(op___, 0);(opJMP,12);(op___, 0);(opROR, 4);(op___, 0);
(opBVS,12);(opADC, 8);(op___, 0);(op___, 0);(op___, 0);(opADC, 2);(opROR, 2);(op___, 0); (* 7x *)
(opSEI, 9);(opADC, 6);(op___, 0);(op___, 0);(op___, 0);(op___, 0);(opROR, 5);(op___, 0);
(opSTA, 4);(opSTA, 7);(op___, 0);(op___, 0);(opSTY, 1);(opSTA, 1);(opSTX, 1);(op___, 0); (* 8x *)
(opDEY, 9);(op___, 0);(opTXA, 9);(op___, 0);(opSTY, 4);(op___, 0);(opSTX, 4);(op___, 0);
(opSTA, 5);(opSTA, 8);(op___, 0);(op___, 0);(opSTY, 2);(opSTA, 2);(opSTX, 2);(op___, 0); (* 9x *)
(opTYA, 9);(opSTA, 6);(opTXS, 9);(op___, 0);(op___, 0);(op___, 0);(op___, 0);(op___, 0);
(opLDY, 0);(opLDA, 7);(opLDX, 0);(op___, 0);(opLDY, 1);(opLDA, 1);(opLDX, 1);(op___, 0); (* ax *)
(opTAY, 9);(opLDA, 0);(opTAX, 9);(op___, 0);(opLDY, 4);(opLDA, 4);(opLDX, 4);(op___, 0);
(opBCS,12);(opLDA, 8);(op___, 0);(op___, 0);(opLDY, 2);(opLDA, 2);(opLDX, 3);(op___, 0); (* bx *)
(opCLV, 9);(opLDA, 6);(opTSX, 9);(op___, 0);(opLDY, 5);(opLDA, 5);(opLDX, 6);(op___, 0);
(opCPY, 0);(opCMP, 7);(op___, 0);(op___, 0);(opCPY, 1);(opCMP, 1);(opDEC, 1);(op___, 0); (* cx *)
(opINY, 9);(opCMP, 0);(opDEX, 9);(op___, 0);(opCPY, 2);(opCMP, 4);(opDEC, 4);(op___, 0);
(opBNE,12);(opCMP, 8);(op___, 0);(op___, 0);(op___, 0);(opCMP, 2);(opDEC, 2);(op___, 0); (* dx *)
(opCLD, 9);(opCMP, 6);(op___, 0);(op___, 0);(op___, 0);(opCMP, 5);(opDEC, 5);(op___, 0);
(opCPX, 0);(opSBC, 7);(op___, 0);(op___, 0);(opCPX, 1);(opSBC, 1);(opINC, 1);(op___, 0); (* ex *)
(opINX, 9);(opSBC, 0);(opNOP, 9);(op___, 0);(opCPX, 2);(opSBC, 4);(opINC, 4);(op___, 0);
(opBEQ,12);(opSBC, 8);(op___, 0);(op___, 0);(op___, 0);(opSBC, 2);(opINC, 2);(op___, 0); (* fx *)
(opSED, 9);(opSBC, 6);(op___, 0);(op___, 0);(op___, 0);(opSBC, 5);(opINC, 5);(op___, 0);
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

