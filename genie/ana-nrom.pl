#!/usr/bin/perl
# Created:20110109
# By Jeff Connelly

# Analyze codes for NROM games

# TODO

use strict;
use Data::Dumper;

our $ROOT = "../roms/best";

open(M, "<galoob/gamelist-galoob.csv")||die;
my %gg2gn;
while(<M>) {
    chomp;
    my ($galoob, $id, $goodnes) = split /\t/;
    $gg2gn{$galoob} = $goodnes;
}

open(DB, "<galoob/all-nev.csv") || die;
my %codes;
while(<DB>) {
    chomp;
    my ($game, $id, $type, @rest) = split /\t/;
    if ($type eq "code") {
        # TODO: read more info
        $codes{$game} = {} if !exists $codes{$game};
        my ($source, $id, $code, $title) = @rest;
        $codes{$game}{$title} = $code;
    }
}

for my $game (sort keys %codes) {
    my %game_codes = %{$codes{$game}};
    for my $title (keys %game_codes) {
        my $code = $game_codes{$title};

        # we're only considering simple, NROM games here
        # but, Contra is UNROM, yet has 6-letter codes?
        next if length($code) != 6;

        my $dir = "$ROOT/$gg2gn{$game}";
        opendir(D, $dir) || die "cannot open: $dir";
        my @files = readdir(D);
        closedir(D);

        print "(@files) $dir,$title,$code\n";
    }
}

