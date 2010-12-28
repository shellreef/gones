gones::
	gd -o gones

test: nestest unittest timingtest

fulltest: test instr_test_v3

# Run without GUI
EMU=./gones h

# CPU functionality test
nestest: gones
	rm -f data/nestest.log-actual data/nestest/nestest.log-diff
	-timeout 5 $(EMU) V 'l data/nestest.nes' 'g 0xc000' > data/nestest.log-actual
	perl tracediff.pl data/nestest.log data/nestest.log-actual -q > data/nestest.log-diff
	! grep "^-" data/nestest.log-diff
	grep 'Nestest automation' data/nestest.log-actual
	@grep 'Nestest automation: Pass' data/nestest.log-actual >/dev/null

compare-nestest:
	vim data/nestest.log data/nestest.log-actual

timingtest:
	@echo "CPU timing test"
	perl cmptiming.pl

unittest: gones
	gd -test

# More comprehensive CPU test
instr_test_v3:
	$(EMU) data/instr_test-v3/rom_singles/01-implied.nes
	@# Unreliable opcode: $AB ATX #n
	-$(EMU) data/instr_test-v3/rom_singles/02-immediate.nes
	$(EMU) data/instr_test-v3/rom_singles/03-zero_page.nes
	$(EMU) data/instr_test-v3/rom_singles/04-zp_xy.nes
	$(EMU) data/instr_test-v3/rom_singles/05-absolute.nes
	@# Unreliable opcode: 9C SYA abs,X
	@# Unreliable opcode: 9E SXA abs,Y
	-$(EMU) data/instr_test-v3/rom_singles/06-abs_xy.nes
	$(EMU) data/instr_test-v3/rom_singles/07-ind_x.nes
	$(EMU) data/instr_test-v3/rom_singles/08-ind_y.nes
	$(EMU) data/instr_test-v3/rom_singles/09-branches.nes
	$(EMU) data/instr_test-v3/rom_singles/10-stack.nes
	$(EMU) data/instr_test-v3/rom_singles/11-jmp_jsr.nes
	$(EMU) data/instr_test-v3/rom_singles/12-rts.nes
	$(EMU) data/instr_test-v3/rom_singles/13-rti.nes
	$(EMU) data/instr_test-v3/rom_singles/14-brk.nes
	$(EMU) data/instr_test-v3/rom_singles/15-special.nes
	
clean::
	gd clean
	rm gones

loop:: gones
	$(EMU) b g
