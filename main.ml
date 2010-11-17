(* Created:20101114
 * By Jeff Connelly
 *
 *)

let game = Nesfile.read (Array.get Sys.argv 1);;

let prg0_io = (IO.input_string (List.nth game.Nesfile.prg_data 0));;

try
    while true do
        let instr = Cpu6502.read_instruction prg0_io in

        print_endline (Cpu6502.string_of_instruction instr)
    done
with IO.No_more_input -> ();;

print_endline (Cpu6502.read_and_print (IO.input_string "\xa9\x40"));;
