(* Created:20101114
 * By Jeff Connelly
 *
 * iNES (.nes) file reader
 *)

let read filename =
    let f = open_in_bin filename in    (* not using extlib IO because it lacks seek *)

    (* http://nesdev.parodius.com/neshdr20.txt *)
    let signature = 0x4E45531A in
    (*let page_size = 16384 in *)

    assert (input_binary_int f == signature);

    let prg_pages = input_byte f in
    let chr_pages = input_byte f in
    let mapper_info1 = input_byte f in
    let mapper_info2 = input_byte f in

    Printf.printf "ROM: %d, VROM: %d -- mapper info: %x %x\n" prg_pages chr_pages mapper_info1 mapper_info2;

    (* TODO: read extended info, bytes 8-15 *)
    seek_in f 16;

    (* TODO: read banks *)

    print_endline (Cpu6502.readAndPrint (IO.input_channel f));
    print_endline (Cpu6502.readAndPrint (IO.input_channel f));
    print_endline (Cpu6502.readAndPrint (IO.input_channel f));
    print_endline (Cpu6502.readAndPrint (IO.input_channel f));
    print_endline (Cpu6502.readAndPrint (IO.input_channel f));
    print_endline (Cpu6502.readAndPrint (IO.input_channel f));
    print_endline (Cpu6502.readAndPrint (IO.input_channel f));
    print_endline (Cpu6502.readAndPrint (IO.input_channel f));
    print_endline (Cpu6502.readAndPrint (IO.input_channel f));
    print_endline (Cpu6502.readAndPrint (IO.input_channel f));
    print_endline (Cpu6502.readAndPrint (IO.input_channel f));
    print_endline (Cpu6502.readAndPrint (IO.input_channel f));
    print_endline (Cpu6502.readAndPrint (IO.input_channel f));
    print_endline (Cpu6502.readAndPrint (IO.input_channel f));
    print_endline (Cpu6502.readAndPrint (IO.input_channel f));
    print_endline (Cpu6502.readAndPrint (IO.input_channel f));
    print_endline (Cpu6502.readAndPrint (IO.input_channel f));


    print_endline (Cpu6502.readAndPrint (IO.input_string "\xa9\x40"));;

