COMPILER = 6g
LINKER = 6l

EXE = gones
SRCS = nesfile.go dis6502.go cpu6502.go main.go
OBJS = ${SRCS:.go=.6}

.SUFFIXES: .go .6

.go.6:
	$(COMPILER) -c $<

all: $(EXE)

test: $(EXE)
	./$(EXE) ../nestest.nes

$(EXE): $(OBJS)
	$(LINKER) -o $@ main.6

$(OBJS): 

clean:
	-$(RM) -f $(EXE) $(OBJS)
