
include $(GOROOT)/src/Make.inc

TARG=leggo

CGOFILES=\
	leggo.go

CGO_OFILES=\
	c-leggo.o

CGO_LDFLAGS=-lallegro

CLEANFILES += leggo

include $(GOROOT)/src/Make.pkg

leggo: install main.go
	$(GC) main.go
	$(LD) -o $@ main.$O

try: leggo
	./leggo
