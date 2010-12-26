// Created:20101214
// By Jeff Connelly

// Leggo my Allegro: an Allegro 5 wrapper for Go

// Turn off magic re#defining of main(), so we can wrap it manually
#define ALLEGRO_NO_MAGIC_MAIN
#include <allegro5/allegro5.h>

#include <stdio.h>

#include "_cgo_export.h"

int user_main(int argc, char **argv) {
    printf("user_main(): about to call GoFoo()\n");
    GoFoo();

    return 0;
}

void wrapper() {
    printf("wrapper():, about to call GoFoo()\n");
    GoFoo();
    
    printf("wrapper(): about to call al_run_main()\n");
    al_run_main(0, 0, user_main);
}

