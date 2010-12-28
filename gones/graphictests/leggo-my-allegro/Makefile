
include $(GOROOT)/src/Make.inc

TARG=leggo

CGOFILES=\
	leggo.go

CGO_OFILES=\
	c-leggo.o

CGO_LDFLAGS=-lallegro

CLEANFILES += leggo

include $(GOROOT)/src/Make.pkg

demo: install demo.go
	$(GC) demo.go
	$(LD) -o $@ demo.$O

# Pure C test program, for comparison
hello1::
	gcc hello1.c -o hello1 -lallegro -lallegro_main 
