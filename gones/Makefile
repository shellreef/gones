gones::
	gd -o gones

test: nestest unittest

# CPU functionality test
nestest: gones
	rm -f log.actual diff
	-./gones 'l ../nestest.nes' 'g 0xc000' > log.actual
	perl tracediff.pl ../nestest.log log.actual -q > diff
	! grep "^-" diff
	grep 'Nestest automation' log.actual
	@grep 'Nestest automation: Pass' log.actual >/dev/null

unittest: gones
	gd -test

	#@echo "CPU timing test"
	#perl cmptiming.pl


compare:
	vim ../nestest.log log

clean::
	gd clean
