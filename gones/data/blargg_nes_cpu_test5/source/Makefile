all: 01-implied.nes 02-immediate.nes 03-zero_page.nes 04-zp_xy.nes 05-absolute.nes 06-abs_xy.nes 07-ind_x.nes 08-ind_y.nes 09-branches.nes 10-stack.nes 11-special.nes

%.o: %.a 
	ca65 -I common -o $@ $<

%.nes: %.o
	ld65 $< -o ../rom_singles/$@ -C nes.cfg 
