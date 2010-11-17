(* Created:20101114
 * By Jeff Connelly
 *
 * iNES (.nes) file reader
 *)


let read filename =
    let io = IO.input_channel (open_in_bin filename) in    (* not using extlib IO because it lacks seek *)

    (* http://nesdev.parodius.com/neshdr20.txt *)
    let signature = 0x1a53454e in     (* NES^Z, little endian *)

    assert (IO.read_i32 io == signature);

    let prg_page_count = IO.read_byte io in   (* ROM = PRG (program) data, 16384 bytes/page *)
    let chr_page_count = IO.read_byte io in   (* VROM = CHR (character) data, 8192 bytes/page *)
    let mapper_info1 = IO.read_byte io in
    let mapper_info2 = IO.read_byte io in

    Printf.printf "ROM: %d, VROM: %d -- mapper info: %x %x\n" prg_page_count chr_page_count mapper_info1 mapper_info2;

    let _ram_pages = IO.read_byte io in   (* 8192 bytes/page *)
    let _pal_flag = IO.read_byte io in

    let _reserved = IO.really_nread io 6 in

    (* Read data PRG and CHR pages, in reverse order *)
    let rec read_pages io count size pages = 
        if count > 0 then 
        (   
            let page = IO.really_nread io size in
            page :: (read_pages io (count - 1) size pages) 
        )
        else
            pages
    in

    let prg_data = (read_pages io prg_page_count 16384 []) in
    let chr_data = (read_pages io chr_page_count 8192 []) in

    Printf.printf "Read %d prg, and %d chr\n" (List.length prg_data) (List.length chr_data);

    let prg0_io = (IO.input_string (List.nth prg_data 0)) in

    try
        while true do
            print_endline (Cpu6502.read_and_print prg0_io)
        done
    with IO.No_more_input -> ();

    print_endline (Cpu6502.read_and_print (IO.input_string "\xa9\x40"));;

