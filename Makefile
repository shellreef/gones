EXE = 6502
SOURCES = cpu6502.ml nesfile.ml main.ml
NATIVE_OBJECTS = ${SOURCES:.ml=.cmx}
MACH_OBJECTS = ${SOURCES:.ml=.o}
INTERFACES = ${SOURCES:.ml=.cmi}

# For IO, from http://code.google.com/p/ocaml-extlib/
PACKAGES = extlib

COMPILER = ocamlfind ocamlopt -package $(PACKAGES)
LINKER = ocamlfind ocamlopt -package $(PACKAGES) -linkpkg


.SUFFIXES: .ml .cmx

.ml.cmx:
	$(COMPILER) -c $<

all: $(EXE)

run: $(EXE)
	./$(EXE)

$(EXE): $(NATIVE_OBJECTS)
	$(LINKER) -o $@ $(NATIVE_OBJECTS)

$(NATIVE_OBJECTS): 

clean:
	-$(RM) -f $(EXE) $(NATIVE_OBJECTS) $(MACH_OBJECTS) $(INTERFACES) 
