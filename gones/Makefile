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
	./$(EXE) -start=0xc000 ../nestest.nes > log.actual
	perl tracediff.pl ../nestest.log log.actual -q
	grep 'Nestest automation' log.actual
	@grep 'Nestest automation: Pass' log.actual >/dev/null

compare:
	vim ../nestest.log log

$(EXE): $(OBJS)
	$(LINKER) -o $@ main.6

$(OBJS): 

clean:
	-$(RM) -f $(EXE) $(OBJS)
