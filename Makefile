COMPILER = ocamlopt
LINKER = ocamlopt

EXE = 6502
SOURCES = cpu6502.ml
NATIVE_OBJECTS = ${SOURCES:.ml=.cmx}
MACH_OBJECTS = ${SOURCES:.ml=.o}
INTERFACES = ${SOURCES:.ml=.cmi}

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
