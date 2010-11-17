(* Created:20101114
 * By Jeff Connelly
 *
 *)

let game = Nesfile.read (Array.get Sys.argv 1);;

let prg0_io = (IO.input_string (List.nth game.Nesfile.prg_data 0));;

let i_flag = ref false;;
let d_flag = ref false;;
let r_acc = ref 0;;

try
    while true do
        let instr = Dis6502.read_instruction prg0_io in

        print_endline (Dis6502.string_of_instruction instr);

        match instr.Dis6502.opcode with
        | Dis6502.SEI -> i_flag := true
        | Dis6502.CLD -> d_flag := false
        | _ -> ()

    done
with IO.No_more_input -> ();;

print_endline (Dis6502.string_of_instruction (Dis6502.read_instruction (IO.input_string "\xa9\x40")));;
