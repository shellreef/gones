gones::
	gd -o gones

test: nestest unittest timingtest

# CPU functionality test
nestest: gones
	rm -f log.actual diff
	-timeout 2 ./gones 'l tests/nestest.nes' 'g 0xc000' > log.actual
	perl tracediff.pl tests/nestest.log log.actual -q > diff
	! grep "^-" diff
	grep 'Nestest automation' log.actual
	@grep 'Nestest automation: Pass' log.actual >/dev/null

timingtest:
	@echo "CPU timing test"
	perl cmptiming.pl


unittest: gones
	gd -test


compare:
	vim tests/nestest.log log.actual

clean::
	gd clean
