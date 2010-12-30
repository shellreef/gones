#!/usr/bin/perl
use strict;
use warnings;

open(FH, "</Users/jeff/games/nese/roms/own/Super\ Mario\ Bros.\ 2\ \(U\)\ \(PRG0\)\ \[\!\].nes") || die "cannot open: $!";
my $offset = 0x03db;
my $size = 0x4000;
for my $i (0..16) {
    my $bank = <FH>;
    last if !defined($bank);

    my $offset = 0x10 + ($size * $i) + $offset;
    seek(FH, $offset, 0);
    $/ = \1;
    my $key = ord(<FH>);
    printf "%.6x\t%.2x\n", ($i * $size), $key;

    ++$i;
}

