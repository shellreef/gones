gones::
	gd -o gones

test: nestest unittest timingtest

fulltest: test instr_test_v3

# CPU functionality test
nestest: gones
	rm -f /tmp/log.actual diff
	-timeout 5 ./gones V 'l tests/nestest.nes' 'g 0xc000' > /tmp/log.actual
	perl tracediff.pl tests/nestest.log /tmp/log.actual -q > diff
	! grep "^-" diff
	grep 'Nestest automation' /tmp/log.actual
	@grep 'Nestest automation: Pass' /tmp/log.actual >/dev/null

compare-nestest:
	vim tests/nestest.log /tmp/log.actual

timingtest:
	@echo "CPU timing test"
	perl cmptiming.pl

unittest: gones
	gd -test

# More comprehensive CPU test
instr_test_v3:
	./gones tests/instr_test-v3/rom_singles/01-implied.nes
	@# Unreliable opcode: $AB ATX #n
	-./gones tests/instr_test-v3/rom_singles/02-immediate.nes
	./gones tests/instr_test-v3/rom_singles/03-zero_page.nes
	./gones tests/instr_test-v3/rom_singles/04-zp_xy.nes
	./gones tests/instr_test-v3/rom_singles/05-absolute.nes
	@# Unreliable opcode: 9C SYA abs,X
	@# Unreliable opcode: 9E SXA abs,Y
	-./gones tests/instr_test-v3/rom_singles/06-abs_xy.nes
	./gones tests/instr_test-v3/rom_singles/07-ind_x.nes
	./gones tests/instr_test-v3/rom_singles/08-ind_y.nes
	./gones tests/instr_test-v3/rom_singles/09-branches.nes
	./gones tests/instr_test-v3/rom_singles/10-stack.nes
	./gones tests/instr_test-v3/rom_singles/11-jmp_jsr.nes
	./gones tests/instr_test-v3/rom_singles/12-rts.nes
	./gones tests/instr_test-v3/rom_singles/13-rti.nes
	./gones tests/instr_test-v3/rom_singles/14-brk.nes
	./gones tests/instr_test-v3/rom_singles/15-special.nes
	
clean::
	gd clean
	rm gones

loop:: gones
	./gones b g
