#!/usr/bin/perl
use strict;
use warnings;

open(FH, "</Users/jeff/games/nese/roms/own/Super\ Mario\ Bros.\ 3\ \(U\)\ \(PRG0\)\ \[\!\].nes") || die "cannot open: $!";
my $offset = 0x35da;
my $size = 0x2000;  # 8kb
my $banks = (16 * 0x4000 / 0x2000);   # 16 x $4000 PRG in header, but mapper switches out $2000
my %keys;
my %tried = (
0xC8 => 'AOSUZSEG',
0x0E => 'AOSUZSTA',
0x29 => 'AOSUZSPZ',
0xBB => 'AOSUZSUL',
0x64 => 'AOSUZIGT',
0xFF => 'AOSUZSNY',
0xA2 => 'AOSUZIXZ',
0x8E => 'AOSUZSVA',
0xFD => 'AOSUZSSY',
0x02 => 'AOSUZIZA',
0x36 => 'AOSUZITL',
0x02 => 'AOSUZIZA',
0x09 => 'AOSUZSPA',
0x30 => 'AOSUZIAL',
0x20 => 'AOSUZIAZ',
0x7C => 'AOSUZSGY',
0xA0 => 'AOSUZIEZ',
0x04 => 'AOSUZIGA',
0x01 => 'AOSUZIPA',
0xBC => 'AOSUZSKL',
0x05 => 'AOSUZIIA',
0x12 => 'AOSUZIZP',
0xFC => 'AOSUZSKY',
0xF7 => 'AOSUZINY',
0x1E => 'AOSUZSTP',
#0x6C => 'AOSUZSGT',   # this one has a subset of the effects
);

for my $i (0..$banks) {
    my $bank = <FH>;
    last if !defined($bank);

    my $offset = 0x10 + ($size * $i) + $offset;
    seek(FH, $offset, 0);
    $/ = \1;
    my $key = ord(<FH>);

    #next if exists($tried{$key}); # already
    #print exists($keys{$key}) ? "---" : "   "; # already saw this key (but affects this bank too!)

    printf "%.6x\t%.2x\n", ($i * $size), $key;

    $keys{$key}++;

    ++$i;
}

