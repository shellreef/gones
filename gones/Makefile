COMPILER = 6g
LINKER = 6l

EXE = 6502
SRCS = nesfile.go 6502.go
OBJS = ${SRCS:.go=.6}

.SUFFIXES: .go .6

.go.6:
	$(COMPILER) -c $<

all: $(EXE)

test: $(EXE)
	./$(EXE) ../*.NES

$(EXE): $(OBJS)
	$(LINKER) -o $@ $(EXE).6

$(OBJS): 

clean:
	-$(RM) -f $(EXE) $(OBJS)
