#!/usr/bin/perl
use strict;
use warnings;

open(FH, "</tmp/a.nes") || die "cannot open: $!";
$/ = \16;
my $header = <FH>;
my $offset = 0x35da;
my $i = 0;
my $size = 16384;
while(<FH>) {
    $/ = \$size;
    my $bank = <FH>;
    last if !defined($bank);

    my $key = ord(substr($bank, $offset, 1));
    printf "%.6x\t%.2x\n", ($i * $size), $key;

    ++$i;
}

