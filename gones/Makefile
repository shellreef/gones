EXE=gones

$(EXE):
	gd . -o gones

test: $(EXE)
	./$(EXE) -start=0xc000 ../nestest.nes > log.actual
	perl tracediff.pl ../nestest.log log.actual -q
	grep 'Nestest automation' log.actual
	@grep 'Nestest automation: Pass' log.actual >/dev/null

	gd . -test

compare:
	vim ../nestest.log log

clean:
	gd . clean
