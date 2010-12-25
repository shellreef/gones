include $(GOROOT)/src/Make.inc

TARG=yasdl

CGOFILES=yasdl.go
CGO_OFILES=SDLMain.o
CGO_CFLAGS=`sdl-config --cflags` 
CGO_LDFLAGS=`sdl-config --libs`

CLEANFILES+=tryit

include $(GOROOT)/src/Make.pkg

c-hello::
	gcc c-hello.c `sdl-config --static-libs` -o c-hello

gohello: install
	6g gohello.go
	6l -o gohello gohello.6

try: gohello
	./gohello
