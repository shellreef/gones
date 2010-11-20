(* Created:20101116
 * By Jeff Connelly
 *
 *)

(* TODO *)

(* Unfortunately you have to constrain these manually using the ops in int8.ml *)
type int16 = int;;
type int8 = int;;

type cpu_registers = {
	pc:int16;   (* program counter *)
	s:int8;     (* stack pointer *)
	p:int8;     (* processor status flag *)
		(* N: negative flag *)
		(* V: overflow flag *)
		(* -: unused *)
		(* B: break flag *)
		(* D: decimal flag *)
		(* I: interrupt disable *)
		(* Z: zero flag *)
		(* C: carry flag *)
	x:int8;	    (* index register X, has (####,X) mode *)
	y:int8;	    (* index register Y, has (####),Y mode *)
};;
	
