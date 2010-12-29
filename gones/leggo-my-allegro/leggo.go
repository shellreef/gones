// Created:20101212
// By Jeff Connelly

// Leggo my Allegro: an Allegro 5 wrapper for Go

package leggo

import ("fmt"
    "net"
    "os"
    "unsafe"
    "runtime"
    )

// #define ALLEGRO_NO_MAGIC_MAIN
// #include <allegro5/allegro.h>
// #include <sys/mman.h>
// #include "leggo.h"
import "C"

type Screen struct {
    data unsafe.Pointer
}

const SOCKET_FILE = "/tmp/leggo.sock"   // TODO: stop using insecure temporary directory

// TODO: all events in include/allegro5/events.h
const EVENT_KEY_DOWN = C.ALLEGRO_EVENT_KEY_DOWN
const EVENT_KEY_CHAR = C.ALLEGRO_EVENT_KEY_CHAR
const EVENT_KEY_UP = C.ALLEGRO_EVENT_KEY_UP

// include/allegro5/keycodes.h
// TODO: can these be imported from C.* directly?
const (KEY_A=C.ALLEGRO_KEY_A + iota;
KEY_B; KEY_C; KEY_D; KEY_E; KEY_F; KEY_G;
KEY_H; KEY_I; KEY_J; KEY_K; KEY_L; KEY_M;
KEY_N; KEY_O; KEY_P; KEY_Q; KEY_R; KEY_S;
KEY_T; KEY_U; KEY_V; KEY_W; KEY_X; KEY_Y;
KEY_Z; KEY_0; KEY_1; KEY_2; KEY_3; KEY_4;
KEY_5; KEY_6; KEY_7; KEY_8; KEY_9; KEY_PAD_0;
KEY_PAD_1; KEY_PAD_2; KEY_PAD_3; KEY_PAD_4; KEY_PAD_5;
KEY_PAD_6; KEY_PAD_7; KEY_PAD_8; KEY_PAD_9; KEY_F1;
KEY_F2; KEY_F3; KEY_F4; KEY_F5; KEY_F6; KEY_F7;
KEY_F8; KEY_F9; KEY_F10; KEY_F11; KEY_F12; KEY_ESCAPE;
KEY_TILDE; KEY_MINUS; KEY_EQUALS; KEY_BACKSPACE; KEY_TAB;
KEY_OPENBRACE; KEY_CLOSEBRACE; KEY_ENTER; KEY_SEMICOLON;
KEY_QUOTE; KEY_BACKSLASH; KEY_BACKSLASH2; KEY_COMMA; KEY_FULLSTOP;
KEY_SLASH; KEY_SPACE; KEY_INSERT; KEY_DELETE; KEY_HOME;
KEY_END; KEY_PGUP; KEY_PGDN; KEY_LEFT; KEY_RIGHT;
KEY_UP; KEY_DOWN; KEY_PAD_SLASH; KEY_PAD_ASTERISK; KEY_PAD_MINUS;
KEY_PAD_PLUS; KEY_PAD_DELETE; KEY_PAD_ENTER; KEY_PRINTSCREEN; KEY_PAUSE;
KEY_ABNT_C1; KEY_YEN; KEY_KANA; KEY_CONVERT; KEY_NOCONVERT;
KEY_AT; KEY_CIRCUMFLEX; KEY_COLON2; KEY_KANJI; KEY_EQUALS_PAD;
KEY_BACKQUOTE; KEY_SEMICOLON2; KEY_COMMAND; KEY_UNKNOWN;

KEY_MODIFIERS=iota + C.ALLEGRO_KEY_MODIFIERS;
KEY_LSHIFT; KEY_RSHIFT; KEY_LCTRL; KEY_RCTRL;
KEY_ALT; KEY_ALTGR; KEY_LWIN; KEY_RWIN;
KEY_MENU; KEY_SCROLLLOCK; KEY_NUMLOCK; KEY_CAPSLOCK;
KEY_MAX);

// TODO: more fields from Allegro's events, wrap directly if can
// (but would still need to serialize over socket)
type Event struct {
    Type uint8
    Keycode uint8
}
    

//export GoLeggoExit
func GoLeggoExit() {
    // Note: don't call this GoExit! It will cause dyld to fail to find anything here.
    fmt.Printf("in LeggoExit!\n");
}

/*
func CreateDisplay(width int, height int) (*C.ALLEGRO_DISPLAY) {
    return C.al_create_display(C.int(width), C.int(height))
}*/

// Setup a Unix domain socket to communicate with the Allegro side
// Calling an //export'd Go function from C in Allegro fails with a
// null pointer exception, which is why I'm doing this. Another alternative
// may be using pthread condition variables as described on 
// http://blog.labix.org/2010/12/10/integrating-go-with-c-the-zookeeper-binding-experience
// and https://github.com/0xe2-0x9a-0x9b/Go-SDL/tree/master/sdl/audio/
// and http://bazaar.launchpad.net/%7Eensemble/gozk/trunk/annotate/head%3A/helpers.c
func LeggoServer(start func(), process func(chan Event)) {
    os.Remove(SOCKET_FILE)

    listener, err := net.Listen("unix", SOCKET_FILE)
    if err != nil {
        fmt.Printf("LeggoServer: failed to listen on socket %s: %s\n", SOCKET_FILE, err)
        return
    }
    //fmt.Printf("listener = %s\n", listener)
    go start()

    ch := make(chan Event)
    go process(ch)

    for {
        //fmt.Printf("LeggoServer: waiting for connection\n")
        conn, err := listener.Accept()
        if err != nil {
            panic(fmt.Sprintf("LeggoServer: Accept() failed: %s", err))
        }

        //fmt.Printf("Established connection: %s\n", conn)
   
        for {
            var buffer [2]byte
            bytesRead, err := conn.Read(buffer[:])
            if err != nil {
                panic(fmt.Sprintf("Error reading from client: %s\n", err))
            }

            if bytesRead != 2 {
                panic(fmt.Sprintf("Read unexpected number of bytes from client: %d\n", bytesRead));
            } 

            kind := byte(buffer[0])
            keycode := byte(buffer[1])
            e := Event{kind, keycode}

            ch <- e

            //fmt.Printf(">>> Read event from client: %d,%d\n", kind, keycode)

            // TODO: dispatch events to Go channel instead of callback?
            //event(int(kind), int(keycode))

        }
    }
}

func LeggoSetup() (unsafe.Pointer) {
    // We use mmap'd memory to communicate what to display;
    // c-leggo.c refresh() will copy this to the screen each frame
    
    size := (C.RESOLUTION_H + 1) * (C.RESOLUTION_W + 1) * 4 // RGBA
    screenMap, err := C.mmap(nil, C.size_t(size), C.PROT_READ | C.PROT_WRITE, 
                       C.MAP_ANON | C.MAP_SHARED, 
        // TODO: on Mach, use VM_MAKE_TAG() in fd so vmmap can distinguish it
                       -1, 0)
    if err != nil {
        panic(fmt.Sprintf("LeggoSetup(): failed to mmap: %s\n", err))
    }

    C.set_screen_map(screenMap, C.size_t(size))

    //fmt.Printf("sm = %s\n", screenMap)

    return screenMap
}

// Get things going. start() will be called when setup.
//export LeggoMain
func LeggoMain(start func(), process func(chan Event)) {
    _ = LeggoSetup()

    go LeggoServer(start, process)

    //fmt.Printf("LeggoMain: about to call al_run_main_wrapper\n")
    C.al_run_main_wrapper()
}

// Write to an arbitrary byte in memory
// This has been moved to C.write_byte because it is too complex in Go
/*
func WriteByte(screen unsafe.Pointer, offset int, value byte) {
    ptr := unsafe.Pointer(uintptr(screen) + uintptr(offset))
    pixel := (*uintptr)(ptr) 
    *pixel = uintptr(value)
}*/

// Set a pixel through direct memory access
// TODO: higher-level functions
func WritePixel(x int, y int, r byte, g byte, b byte, a byte) {
    if x > C.RESOLUTION_W || y > C.RESOLUTION_H {
        // This is important, because write_byte has no bounds checking!
        panic(fmt.Sprintf("WritePixel out of range: (%d,%d)\n", x, y))
    }
    offset := 4*(x + y*C.RESOLUTION_W)
    C.write_byte(C.off_t(offset + 0), C.uint8_t(r))
    C.write_byte(C.off_t(offset + 1), C.uint8_t(g))
    C.write_byte(C.off_t(offset + 2), C.uint8_t(b))
    C.write_byte(C.off_t(offset + 3), C.uint8_t(a))
}

// Get frames per second
func FPS() (float) {
    return 1.0 / float(C.get_seconds_per_frame())
}

// Get width and height (TODO: user defined)
func Dimensions() (int, int) {
    return C.RESOLUTION_W, C.RESOLUTION_H
}

func init() {
    runtime.GOMAXPROCS(2)
}
