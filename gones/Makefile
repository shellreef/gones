EXE=gones

$(EXE)::
	gd . -o gones

test: clean
	gd . -test

	./$(EXE) 'l ../nestest.nes' 'g 0xc000' > log.actual
	perl tracediff.pl ../nestest.log log.actual -q
	grep 'Nestest automation' log.actual
	@grep 'Nestest automation: Pass' log.actual >/dev/null

	#@echo "CPU timing test"
	#perl cmptiming.pl


compare:
	vim ../nestest.log log

clean::
	gd . clean
	rm -rf tmp*
