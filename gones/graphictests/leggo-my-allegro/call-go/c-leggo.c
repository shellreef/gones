// Created:20101214
// By Jeff Connelly

// Leggo my Allegro: an Allegro 5 wrapper for Go

// Turn off magic re#defining of main(), so we can wrap it manually
#define ALLEGRO_NO_MAGIC_MAIN
#include <allegro5/allegro5.h>

#include <stdio.h>
#include <stdlib.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <unistd.h>

#include "leggo.h"
#include "_cgo_export.h"

// The "user-defined" callback function given to al_run_main()
// Allegro will call this as part of its UI wrapper, so when it 
// returns, the program will exit. All code is invoked from here.
int leggo_user_main(int argc, char **argv) {
    printf("in leggo_user_main, trying to call back into Go code...\n");
    GoFoo();

    return 0;
}

// Wrap calling Allegro's al_run_main(), with our own leggo_main. 
// This is a C function so Go can call it using cgo, in leggo.Main().
void al_run_main_wrapper() {
    printf("in al_run_main_wrapper, about to call GoFoo()\n");
    GoFoo();
    
    printf("al_run_main_wrapper(): about to call al_run_main()\n");
    al_run_main(0, 0, leggo_user_main);
}

