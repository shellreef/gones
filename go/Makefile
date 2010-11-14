COMPILER = 6g
LINKER = 6l

EXE = 6502
SRCS = 6502.go
OBJS = ${SRCS:.go=.6}

.SUFFIXES: .go .6

.go.6:
	$(COMPILER) -c $<

all: $(EXE)

run: $(EXE)
	./$(EXE)

$(EXE): $(OBJS)
	$(LINKER) -o $@ $(OBJS)

$(OBJS): 

clean:
	-$(RM) -f $(EXE) $(OBJS)
