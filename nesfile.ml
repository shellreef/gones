(* Created:20101114
 * By Jeff Connelly
 *
 * iNES (.nes) file reader
 *)

let read filename =
    let io = IO.input_channel (open_in_bin filename) in    (* not using extlib IO because it lacks seek *)

    (* http://nesdev.parodius.com/neshdr20.txt *)
    let signature = 0x1a53454e in     (* NES^Z, little endian *)
    let page_size = 16384 in 

    assert (IO.read_i32 io == signature);

    let prg_pages = IO.read_byte io in
    let chr_pages = IO.read_byte io in
    let mapper_info1 = IO.read_byte io in
    let mapper_info2 = IO.read_byte io in

    Printf.printf "ROM: %d, VROM: %d -- mapper info: %x %x\n" prg_pages chr_pages mapper_info1 mapper_info2;

    let ext_info = IO.really_nread io 8 in

    (* TODO: read banks *)

    print_endline (Cpu6502.read_and_print io);
    print_endline (Cpu6502.read_and_print io);
    print_endline (Cpu6502.read_and_print io);
    print_endline (Cpu6502.read_and_print io);
    print_endline (Cpu6502.read_and_print io);
    print_endline (Cpu6502.read_and_print io);
    print_endline (Cpu6502.read_and_print io);
    print_endline (Cpu6502.read_and_print io);
    print_endline (Cpu6502.read_and_print io);
    print_endline (Cpu6502.read_and_print io);
    print_endline (Cpu6502.read_and_print io);
    print_endline (Cpu6502.read_and_print io);
    print_endline (Cpu6502.read_and_print io);
    print_endline (Cpu6502.read_and_print io);


    print_endline (Cpu6502.read_and_print (IO.input_string "\xa9\x40"));;

