// hello1.c
// An introduction to keyboard input and text output

#include <allegro5/allegro5.h>
#include <allegro5/allegro_font.h>
#include <allegro5/allegro_ttf.h>
#include <stdio.h>

struct Data {
	ALLEGRO_FONT *f1, *f2;
	ALLEGRO_DISPLAY *display;
	ALLEGRO_EVENT_QUEUE *queue;
	
	ALLEGRO_COLOR bright_green;
} data;

const char *font_file = "times.ttf";


int main(int argc, char *argv[]) {
	// Initialize Allegro 5 and the font routines
	al_init();
	al_init_font_addon();
	
	// Create a window to display things on: 640x480 pixels
	data.display = al_create_display(640, 480);
	if(!data.display) {
		printf("Error creating display.\n");
		return 1;
	}
	
	
	// Install the keyboard handler
	if(!al_install_keyboard()) {
		printf("Error installing keyboard.\n");
		return 1;
	}
	
	
	// Load a font
	data.f1 = al_load_ttf_font(font_file, 48, 0);
	data.f2 = al_load_ttf_font(font_file, -48, 0);
	if(!data.f1 || !data.f2) {
		printf("Error loading \"%s\".\n", font_file);
		return 1;
	}
		
	// Make and set a color to draw with
	data.bright_green = al_map_rgba_f(0.5, 1.0, 0.5, 1.0);
	//al_set_blender(ALLEGRO_ONE, ALLEGRO_ALPHA, ALLEGRO_INVERSE_ALPHA, data.bright_green);

	// Draw a message to the backbuffer in the fonts we just loaded using the color we just set
	//al_draw_text(data.f1, 10, 10, -1, "Allegro 5 Rocks!");
	//al_draw_text(data.f2, 10, 60, -1, "Allegro 5 Rocks!");
	
	// Make the backbuffer visible
	al_flip_display();
	
	// Start the event queue to handle keyboard input
	data.queue = al_create_event_queue();
	al_register_event_source(data.queue, (ALLEGRO_EVENT_SOURCE*)al_get_keyboard_event_source());
	
	
	// Wait until the user presses escape
	ALLEGRO_EVENT event;
	while(1) {
		// Block until an event enters the queue
		al_wait_for_event(data.queue, &event);
		
		// If the key was escape, then break out of the loop
		if(event.type == ALLEGRO_EVENT_KEY_DOWN && event.keyboard.keycode == ALLEGRO_KEY_ESCAPE) break;
	}
	
	
	return 0;
}
