Leggo my Allegro
================

Leggo my Allegro is a simple (and incomplete) wrapper of the Allegro 5
game programming library for the Go programming language. Most of the
library is not yet wrapped, but pixel-level graphical output and
keyboard input using Go channels is supported. 

There is another Allegro wrapper for Go from 2009/11 available at
<https://github.com/lasarux/go-spanish/tree/master/allegro>. Leggo is
different in that it properly calls al_run_main(), so it works on OS X.

Installation
------------
1. Install Go release.2010-12-22 or later from http://golang.org/
2. Install Allegro 5.1+ from http://alleg.sourceforge.net/
If is it not yet released, check it out from Subversion:
svn co https://alleg.svn.sourceforge.net/svnroot/alleg/allegro/branches/5.1/
or use Allegro 5.0.0rc3 with patch: <https://sourceforge.net/tracker/?func=detail&aid=3143062&group_id=5665&atid=105665>
Install per directions in README

3. Compile and run demo
make
make demo
./demo

It should display white noise. Use up/down to view a couple other modes,
space to print frame rate, and escape to exit.

Platforms Supported
-------------------
Tested on Mac OS X 10.6.5 with amd64 Go compiler.

In principle this code should work on other platforms supported by
Go and Allegro, but it is largely untested. 

On Ubuntu Linux 10.10 amd64, it compiles but fails to link:

6l -o demo demo.6
/home/jeff/go/pkg/linux_amd64/leggo.a(c-leggo.o)(.text): atexit: not defined
atexit: not defined
make: *** [demo] Error 1

I'm not sure why. Patches welcome. 

Design Notes
------------
The challenge in wrapping Allegro -- and other graphical libraries,
including SDL -- on OS X is that they require the main application
loop to run from a Cocoa/Objective C application delegate, which itself
calls your application-specific code. 

See call-go/ for an attempt at calling Go from Allegro's main C 
callback function. Go code can be called from C using cgo's //export
in some cases, but not from the Cocoa wrapper. The program hangs with
a null dereference, at least on Go release.2010-12-22. (If this can be 
made to work, the design of this library could be a lot simpler.)

Until then, Leggo my Allegro uses an anonymous shared memory map to
(see mmap(2)) communicate the graphical data from Go to C, which refreshes
at 60 Hz, copying the data to the ALLEGRO_LOCKED_REGION buffer. The Allegro
event loop is written in C, and besides handling the screen refresh, also
dispatches keystroke events over a Unix domain socket to Go, which then
sends them over a Go channel. This complexity is hidden by the library.

Contact
-------
Jeff Connelly <leggo-my-allegro@xyzzy.cjb.net>

