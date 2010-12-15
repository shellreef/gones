include $(GOROOT)/src/Make.inc

TARG=leggo

CGOFILES=leggo.go
CGO_CFLAGS=-DALLEGRO_NO_MAGIC_MAIN
CGO_LDFLAGS=-lallegro

include $(GOROOT)/src/Make.pkg

try: install try.go
	$(GC) try.go
	$(LD) -o $@ try.$O
