#!/usr/bin/perl
# Created:20110107
# By Jeff Connelly

# Find duplicate codes

use strict;
use Data::Dumper;

open(DB,"<all-nev.csv")||die;
my %game2id;
my %code2game;
while(<DB>) {
    chomp;
    my ($game, $id, @rest) = split /\t/;

    $game2id{$game} = $id;
    my $type = shift @rest;
    
    if ($type eq "code") {
        my ($no, $code, $title) = @rest;
        die "$code is for both $game and $code2game{$code}\n" if exists $code2game{$game};
        $code2game{$code} = $game;
    }
}

for my $code (sort keys %code2game) {
    my $game = $code2game{$code};
    printf "%-40s %s\n", $code, $game;
}
