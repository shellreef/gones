// Created:20101212
// By Jeff Connelly

// Leggo my Allegro: an Allegro 5 run_main for Go

package leggo

import "fmt"

// extern int user_main(int argc, char **argv);
// extern void run_main();
// #include <allegro5/allegro.h>
import "C"

//export GoFoo
func GoFoo() {
    fmt.Printf("*** It worked! In GoFoo()\n")
}

func GoRunMain() {
    fmt.Printf("GoRunMain: about to call run_main\n")
    C.run_main()
}
