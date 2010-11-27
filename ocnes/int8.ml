(* Created:20101116
 * By Jeff Connelly
 *
 * 8-bit integer operations
 *)

(*

OCaml has int32, int64 (including literals): http://caml.inria.fr/pub/docs/manual-ocaml/manual021.html

extlib has BitSet, UChar: http://ocaml-lib.sourceforge.net/doc/index.html

Jane Street Core has int,int32,int32,int64/native http://www.janestreet.com/ocaml/janestreet-ocamldocs/core/index.html

Bigarray has int8_unsigned_elt: http://caml.inria.fr/pub/docs/manual-ocaml/libref/Bigarray.html#TYPEint8_signed_elt - but it is accessed as an int

*)


(* 8-bit operations *)
let ( +$ ) a b = (a + b) mod 256;;
let ( -$ ) a b = (a - b) mod 256;;
let ( *$ ) a b = (a * b) mod 256;;
let ( /$ ) a b = (a / b) mod 256;;

(* 16-bit operations *)
let ( +$$ ) a b = (a + b) mod 65536;;
let ( -$$ ) a b = (a - b) mod 65536;;
let ( *$$ ) a b = (a * b) mod 65536;;
let ( /$$ ) a b = (a / b) mod 65536;;

(* TODO: bitshift << >> and rotate *)

