;;; Based on  http://www.dwedit.org/files/hello.asm mentioned on http://nesdev.com/bbs/viewtopic.php?p=71841
;;; ca65 hello.asm ; ld65 -o hello.nes hello.o  -C nes.cfg

    .org $BFF0

	.byte "NES",$1A,$01,$00,$20,$00
	.byte 0,0,0,0,0,0,0,0

;;; Memory mapped registers

OAM       = $0200

PPUCTRL   = $2000
PPUMASK   = $2001
PPUSTATUS = $2002
PPUSTAT = $2002
SPRADDR   = $2003  ; always write 0 here and use DMA from OAM
PPUSCROLL = $2005
PPUADDR   = $2006
PPUDATA   = $2007

SPRDMA    = $4014
SNDCHN    = $4015
JOY1      = $4016
JOY2      = $4017

PPUCTRL_NMI      = $80
PPUCTRL_8X8      = $00
PPUCTRL_8X16     = $20
PPUCTRL_BGHIPAT  = $10
PPUCTRL_SPRHIPAT = $08
PPUCTRL_WRDOWN   = $04  ; when set, PPU address increments by 32

PPUMASK_RED      = $80  ; when set, slightly darkens other colors
PPUMASK_GREEN    = $40
PPUMASK_BLUE     = $20
PPUMASK_SPR      = $14  ; SPR: show sprites in x=0-255
PPUMASK_SPRCLIP  = $10  ; SPRCLIP: show sprites in x=8-255
PPUMASK_BG0      = $0A  ; BG0: similarly
PPUMASK_BG0CLIP  = $08
PPUMASK_MONO     = $01  ; when set, zeroes the low nibble of palette values

PPUSTATUS_VBL  = $80  ; the PPU has entered a vblank since last $2002 read
PPUSTATUS_SPR0 = $40  ; sprite 0 has overlapped BG since ???
PPUSTATUS_OVER = $20  ; More than 64 sprite pixels on a scanline since ???

temp0 = $00
temp1 = $01
temp2 = $02
temp3 = $03
temp4 = $04
temp5 = $05
temp6 = $06
temp7 = $07
temp8 = $08
temp9 = $09
tempA = $0A
tempB = $0B
tempC = $0C
tempD = $0D
tempE = $0E
tempF = $0F

addy = temp2

frame = $10
timer = $11

;text_x = $12
;text_y = $13

vblanked = $7F

font:
 .byte $00,$00,$00,$00,$00,$00,$00,$00
 .byte $18,$3C,$3C,$3C,$18,$18,$00,$18
 .byte $6C,$6C,$6C,$00,$00,$00,$00,$00
 .byte $6C,$6C,$FE,$6C,$FE,$6C,$6C,$00
 .byte $30,$7C,$C0,$78,$0C,$F8,$30,$00
 .byte $00,$C6,$CC,$18,$30,$66,$C6,$00
 .byte $38,$6C,$38,$76,$DC,$CC,$76,$00
 .byte $60,$60,$C0,$00,$00,$00,$00,$00
 .byte $18,$30,$60,$60,$60,$30,$18,$00
 .byte $60,$30,$18,$18,$18,$30,$60,$00
 .byte $00,$66,$3C,$FF,$3C,$66,$00,$00
 .byte $00,$30,$30,$FC,$30,$30,$00,$00
 .byte $00,$00,$00,$00,$00,$30,$30,$60
 .byte $00,$00,$00,$FC,$00,$00,$00,$00
 .byte $00,$00,$00,$00,$00,$30,$30,$00
 .byte $06,$0C,$18,$30,$60,$C0,$80,$00
 .byte $38,$4C,$C6,$C6,$C6,$64,$38,$00
 .byte $18,$38,$18,$18,$18,$18,$7E,$00
 .byte $7C,$C6,$0E,$3C,$78,$E0,$FE,$00
 .byte $7E,$0C,$18,$3C,$06,$C6,$7C,$00
 .byte $1C,$3C,$6C,$CC,$FE,$0C,$0C,$00
 .byte $FC,$C0,$FC,$06,$06,$C6,$7C,$00
 .byte $3C,$60,$C0,$FC,$C6,$C6,$7C,$00
 .byte $FE,$C6,$0C,$18,$30,$30,$30,$00
 .byte $7C,$C6,$C6,$7C,$C6,$C6,$7C,$00
 .byte $7C,$C6,$C6,$7E,$06,$0C,$78,$00
 .byte $00,$30,$30,$00,$00,$30,$30,$00
 .byte $00,$30,$30,$00,$00,$30,$30,$60
 .byte $18,$30,$60,$C0,$60,$30,$18,$00
 .byte $00,$00,$FC,$00,$00,$FC,$00,$00
 .byte $60,$30,$18,$0C,$18,$30,$60,$00
 .byte $78,$CC,$0C,$18,$30,$00,$30,$00
 .byte $7C,$C6,$DE,$DE,$DE,$C0,$78,$00
 .byte $38,$6C,$C6,$C6,$FE,$C6,$C6,$00
 .byte $FC,$C6,$C6,$FC,$C6,$C6,$FC,$00
 .byte $3C,$66,$C0,$C0,$C0,$66,$3C,$00
 .byte $F8,$CC,$C6,$C6,$C6,$CC,$F8,$00
 .byte $FE,$C0,$C0,$FC,$C0,$C0,$FE,$00
 .byte $FE,$C0,$C0,$FC,$C0,$C0,$C0,$00
 .byte $3E,$60,$C0,$CE,$C6,$66,$3E,$00
 .byte $C6,$C6,$C6,$FE,$C6,$C6,$C6,$00
 .byte $7E,$18,$18,$18,$18,$18,$7E,$00
 .byte $1E,$06,$06,$06,$C6,$C6,$7C,$00
 .byte $C6,$CC,$D8,$F0,$F8,$DC,$CE,$00
 .byte $60,$60,$60,$60,$60,$60,$7E,$00
 .byte $C6,$EE,$FE,$FE,$D6,$C6,$C6,$00
 .byte $C6,$E6,$F6,$FE,$DE,$CE,$C6,$00
 .byte $7C,$C6,$C6,$C6,$C6,$C6,$7C,$00
 .byte $FC,$C6,$C6,$C6,$FC,$C0,$C0,$00
 .byte $7C,$C6,$C6,$C6,$DE,$CC,$7A,$00
 .byte $FC,$C6,$C6,$CE,$F8,$DC,$CE,$00
 .byte $78,$CC,$C0,$7C,$06,$C6,$7C,$00
 .byte $7E,$18,$18,$18,$18,$18,$18,$00
 .byte $C6,$C6,$C6,$C6,$C6,$C6,$7C,$00
 .byte $C6,$C6,$C6,$EE,$7C,$38,$10,$00
 .byte $C6,$C6,$D6,$FE,$FE,$EE,$C6,$00
 .byte $C6,$EE,$7C,$38,$7C,$EE,$C6,$00
 .byte $66,$66,$66,$3C,$18,$18,$18,$00
 .byte $FE,$0E,$1C,$38,$70,$E0,$FE,$00
 .byte $78,$60,$60,$60,$60,$60,$78,$00
 .byte $C0,$60,$30,$18,$0C,$06,$02,$00
 .byte $78,$18,$18,$18,$18,$18,$78,$00
 .byte $10,$38,$6C,$C6,$00,$00,$00,$00
 .byte $00,$00,$00,$00,$00,$00,$00,$FF
 .byte $30,$30,$18,$00,$00,$00,$00,$00
 .byte $00,$00,$3C,$66,$66,$66,$3B,$00
 .byte $60,$60,$7C,$66,$66,$66,$7C,$00
 .byte $00,$00,$3E,$60,$60,$60,$3E,$00
 .byte $06,$06,$3E,$66,$66,$66,$3E,$00
 .byte $00,$00,$3C,$66,$7E,$60,$3E,$00
 .byte $0E,$18,$18,$7E,$18,$18,$18,$00
 .byte $00,$00,$3E,$66,$66,$3E,$06,$3C
 .byte $60,$60,$60,$7C,$66,$66,$66,$00
 .byte $00,$18,$00,$18,$18,$18,$18,$00
 .byte $00,$06,$00,$06,$06,$06,$66,$3C
 .byte $60,$60,$62,$64,$68,$7C,$66,$00
 .byte $18,$18,$18,$18,$18,$18,$18,$00
 .byte $00,$00,$76,$6B,$6B,$6B,$6B,$00
 .byte $00,$00,$7C,$66,$66,$66,$66,$00
 .byte $00,$00,$3C,$66,$66,$66,$3C,$00
 .byte $00,$00,$7C,$66,$66,$7C,$60,$60
 .byte $00,$00,$3E,$66,$66,$3E,$06,$06
 .byte $00,$00,$6E,$70,$60,$60,$60,$00
 .byte $00,$00,$3C,$40,$3C,$06,$7C,$00
 .byte $30,$30,$FC,$30,$30,$30,$1C,$00
 .byte $00,$00,$66,$66,$66,$66,$3C,$00
 .byte $00,$00,$66,$66,$66,$24,$18,$00
 .byte $00,$00,$63,$6B,$6B,$6B,$36,$00
 .byte $00,$00,$63,$36,$1C,$36,$63,$00
 .byte $00,$00,$66,$66,$2C,$18,$30,$60
 .byte $00,$00,$7E,$0C,$18,$30,$7E,$00
 .byte $1C,$30,$30,$E0,$30,$30,$1C,$00
 .byte $18,$18,$18,$00,$18,$18,$18,$00
 .byte $E0,$30,$30,$1C,$30,$30,$E0,$00
 .byte $76,$DC,$00,$00,$00,$00,$00,$00
 .byte $00,$10,$38,$6C,$C6,$C6,$FE,$00

load_font:
	lda #font&255
	sta addy
	lda #font/256
	sta addy+1
	ldx #3
	ldy #0
fontchar_loop:
	lda (addy),y
	sta $2007
	iny
	lda (addy),y
	sta $2007
	iny
	lda (addy),y
	sta $2007
	iny
	lda (addy),y
	sta $2007
	iny
	lda (addy),y
	sta $2007
	iny
	lda (addy),y
	sta $2007
	iny
	lda (addy),y
	sta $2007
	iny
	lda (addy),y
	sta $2007
	lda #0
	sta $2007
	sta $2007
	sta $2007
	sta $2007
	sta $2007
	sta $2007
	sta $2007
	sta $2007
	iny
	bne fontchar_loop
	inc addy+1
	dex
	bne fontchar_loop
	rts

print:
	ldy #0
print_loop:
	lda (addy),y
	beq print_quit
	sta $2007
	iny
	bne print_loop
print_quit:
	rts

hello:
	.byte "Hello World!",0

write_vram_x:
	ldy #0
wvxloop:
	lda (addy),y
	sta $2007
	iny
	dex
	bne wvxloop
	rts

clear_vram:
	lda #$00
	sta PPUADDR
	sta PPUADDR
	tay
	ldx #6
clear_loop:
	sta $2007
	sta $2007
	sta $2007
	sta $2007
	sta $2007
	sta $2007
	sta $2007
	sta $2007
	dey
	bne clear_loop
	dex
	bne clear_loop
	rts

clear_nt:
	lda #$20
	sta PPUADDR
	lda #$00
	sta PPUADDR
	jsr zero_nt
	lda #$2C
	sta PPUADDR
	lda #$00
	sta PPUADDR
zero_nt:
	lda #0
	ldy #0
zero_nt_loop:
	sta $2007
	sta $2007
	sta $2007
	sta $2007
	dey
	bne zero_nt_loop
	rts
	
nmihandler:
	inc vblanked
	rti

irqhandler:
	rti

main:

;;; Init CPU
	sei
	cld
	ldx #$c0
	stx JOY2
	ldx #$00
	stx SNDCHN
	
	ldx #$ff
	txs
	inx
	stx PPUCTRL
	stx PPUMASK

	lda #$08
	sta $4011

;;; Init machine
wait0001:
	bit PPUSTATUS
	bpl wait0001

	lda #$10
	sta $4011

	ldx #0
	ldy #$ef
	txa
ramclrloop:
	sta 0,x
	sta $100,x
	sta $200,x
	sta $300,x
	sta $400,x
	sta $500,x
	sta $600,x
	sta $700,x
	inx
	bne ramclrloop
	
	;initialize sprites
	lda #$EF
ramclrloop2:
	sta $200,x
	inx
	bne ramclrloop2

	lda #$14
	sta $4011

wait0002:
	bit PPUSTATUS
	bpl wait0002

	lda #%10000000
	sta PPUCTRL
	
	jsr clear_vram
	
	;sprite DMA
	lda #$02
	sta SPRDMA

	
	;load font
	lda #$02
	sta PPUADDR
	lda #$00
	sta PPUADDR
	jsr load_font
	
	lda #$20
	sta PPUADDR
	lda #$84
	sta PPUADDR
	
	lda #hello&255
	sta addy
	lda #hello/256
	sta addy+1
	jsr print
	
	lda #$3F
	sta PPUADDR
	lda #$00
	sta PPUADDR
	
	lda #pal&255
	sta addy
	lda #pal/256
	sta addy+1
	ldx #4
	jsr write_vram_x
	
	
main_loop:
	jsr waitframe
	
	lda #0
	sta $2005
	sta $2005
	lda #$18
	sta $2001
	
	lda #$80
	sta $2000

	jmp main_loop
	
waitframe:
	lda #$00
	sta vblanked
waitloop:
	lda vblanked
	beq waitloop
	rts

pal:
	.byte $0F, $30, $0F, $0F

    .org $FFFA
    .word nmihandler,main,irqhandler

.end


