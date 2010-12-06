all:
	gd . -o gones

test: all
	./$(EXE) -start=0xc000 ../nestest.nes > log.actual
	perl tracediff.pl ../nestest.log log.actual -q
	grep 'Nestest automation' log.actual
	@grep 'Nestest automation: Pass' log.actual >/dev/null

compare:
	vim ../nestest.log log

clean:
	gd . clean
