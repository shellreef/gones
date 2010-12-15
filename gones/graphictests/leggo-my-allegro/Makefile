
include $(GOROOT)/src/Make.inc

TARG=leggo

CGOFILES=leggo.go

LDPATH_freebsd=-Wl,-R,`pwd`
LDPATH_linux=-Wl,-R,`pwd`
LDPATH_darwin=

CGO_LDFLAGS=-lallegro _cgo_export.o c-leggo.so $(LDPATH_$(GOOS))
CGO_DEPS=_cgo_export.o c-leggo.so

CLEANFILES += leggo

include $(GOROOT)/src/Make.pkg

c-leggo.o: c-leggo.c _cgo_export.h
	gcc $(_CGO_CFLAGS_$(GOARCH)) -g -c -fPIC $(CFLAGS) c-leggo.c

c-leggo.so: c-leggo.o
	gcc $(_CGO_CFLAGS_$(GOARCH)) -o $@ c-leggo.o $(_CGO_LDFLAGS_$(GOOS))

leggo: install try.go
	$(GC) try.go
	$(LD) -o $@ try.$O

