#include <stdlib.h>
#include <SDL/SDL.h>

int main(int main, char *argv[]) 
{
    SDL_Surface *screen;

    if (SDL_Init(SDL_INIT_VIDEO) < 0) {
        fprintf(stderr, "Unable to init: %s\n", SDL_GetError());
        exit(EXIT_FAILURE);
    }

    screen = SDL_SetVideoMode(640, 480, 8, SDL_SWSURFACE);
    if (!screen) {
        fprintf(stderr, "Unable to set video: %s\n", SDL_GetError());
        exit(EXIT_FAILURE);
    }

    int x;

    for (x = 0; x < 100; x += 1) {
        int y = x;

        Uint8 *ptr = screen->pixels + y * screen->pitch + x;
        *ptr = SDL_MapRGB(screen->format, 255, 0, 0);
        
        SDL_UpdateRect(screen, x, y, 1, 1);
    }

    SDL_Event event;
    while(SDL_WaitEvent(&event)) {
        switch (event.type) {
        case SDL_KEYDOWN:
        case SDL_QUIT:
            exit(0);
        }
    }
}
