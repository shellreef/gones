// An introduction to graphics output
// Based on http://wiki.allegro.cc/pub/3/3a/Tut2Hello2.c

/* 
 gcc hello1.c -o hello1 -lallegro  -lallegro_main
 ./hello1
 */

#include <allegro5/allegro5.h>
#include <stdio.h>

int main(int argc, char *argv[]) {
	ALLEGRO_DISPLAY *display;
	ALLEGRO_EVENT_QUEUE *queue;
	
	ALLEGRO_COLOR bright_green;

	// Initialize Allegro 5 and the font routines
	al_init();
	
	// Install the keyboard handler
	if(!al_install_keyboard()) {
		printf("Error installing keyboard.\n");
		return 1;
	}
	
	
	// Create a window to display things on: 640x480 pixels
	display = al_create_display(640, 480);
	if(!display) {
		printf("Error creating display.\n");
		return 1;
	}


	// Draw a green background as a test
    al_set_target_bitmap(al_get_backbuffer(display));
    al_clear_to_color(al_map_rgb(128, 255, 128));
    al_flip_display();

	// Start the event queue to handle keyboard input
	queue = al_create_event_queue();
	al_register_event_source(queue, (ALLEGRO_EVENT_SOURCE*)al_get_keyboard_event_source());

    // Update timer
    int rate = al_get_display_refresh_rate(display);
    printf("display refresh rate = %d\n", rate);
    if (rate == 0) rate = 60;
    ALLEGRO_TIMER *timer = al_create_timer(1.0 / rate);
    printf("timer = %p\n", (void *)timer);
    al_start_timer(timer);
    al_register_event_source(queue, (ALLEGRO_EVENT_SOURCE*)timer);
	
	// Wait until the user presses escape
	ALLEGRO_EVENT event;
	while(1) {
		// Block until an event enters the queue
		al_wait_for_event(queue, &event);
		
        printf("event type = %d\n", event.type);

		// If the key was escape, then break out of the loop
		if(event.type == ALLEGRO_EVENT_KEY_DOWN && event.keyboard.keycode == ALLEGRO_KEY_ESCAPE) break;
	}
	
	
	return 0;
}
