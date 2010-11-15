let f = open_in_bin "SMARIO.NES";;

let iNes_signature = 0x4E45531A;;
let page_size = 16384;;

(* http://nesdev.parodius.com/neshdr20.txt *)
assert (input_binary_int f == iNes_signature);;

let prg_pages = input_byte f;;
let chr_pages = input_byte f;;
let mapper_info1 = input_byte f;;
let mapper_info2 = input_byte f;;
(* TODO: read extended info, bytes 8-15 *)
seek_in f 16;;

Printf.printf "%x\n" (input_byte f);;
Printf.printf "%x\n" (input_byte f);;
Printf.printf "%x\n" (input_byte f);;
Printf.printf "%x\n" (input_byte f);;

Printf.printf "ROM: %d, VROM: %d -- mapper info: %x %x\n" prg_pages chr_pages mapper_info1 mapper_info2;;

print_endline (Cpu6502.readInstruction (IO.input_string "\xa9\x40"));;

