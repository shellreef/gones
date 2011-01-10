#!/usr/bin/perl
# Created:20110109
# By Jeff Connelly

# Analyze codes for NROM games

# TODO

use strict;
use Data::Dumper;

our $ROOT = "../roms/best";
our $EMU = "../gones/gones";

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
    my $dir = "$ROOT/$gg2gn{$game}";
    opendir(D, $dir) || die "cannot open: $dir";
    my @files = grep{!m/^\./}readdir(D);
    closedir(D);

    die "nothing for $game" if @files == 0;
    next if @files > 1;  # skip games w/ alternate versions

    my $file = pop @files;

    # Skip games with bankswitched PRG
    # TODO: allow bankswitched, since most games are
    # TODO: make GoNES scriptable so can use it, instead of rewriting all this
    # iNES header parsing in Perl, or using Perl at all..
    open(FH, "<$dir/$file") || die "cannot open $dir/$file: $!";
    $/ = \0x10;
    my $header = <FH>;
    close(FH);
    my $prgSize = ord(substr($header, 4, 1)) * 16384;
    my $chrSize = ord(substr($header, 5, 1)) * 8192;
    my $mapper = (ord(substr($header, 6, 1)) & 0xf0) >> 4 | ((ord(substr($header, 7, 1)) & 0x0f) << 4);
    next if ($mapper != 0   # NROM has fixed PRG (it has fixed everything)
        && $mapper != 3);   # CNROM has fixed PRG, too (only CHR can switch)

    print "$game\n";

    my %game_codes = %{$codes{$game}};
    for my $title (keys %game_codes) {
        my $code = $game_codes{$title};

        # not all non-bankswitched games have 6-letter codes, for some reason
        # Mario Bros has 8-letter codes, even though it is NROM, and Contra has 6-letter
        # codes, even though it is UNROM. Maybe Codemasters didn't know what bankswitching
        # the games used, but just went with what worked.
        #next if length($code) != 6;

        # skip multiple code codes
        next if length($code) != 6 && length($code) != 8;

        analyze("$dir/$file", $code, $title); 
    }
}

sub analyze
{
    my ($file, $code, $title) = @_;

    # Call into GoNES to find ROM offset
    my $cmd = qq($EMU "l $file" "c $code");
    my $out = `$cmd`;
    my ($ines) = $out =~ m/iNES offset: ([0-9a-fA-F]+)/;
    my ($value) = $out =~ m/Value: ([0-9a-fA-F]+)/;
    $value = hex $value;

    if (!defined($ines)) {
        # sometimes it has trouble
        print "$cmd\n";
        print "$out\n";
    }

    # Read context
    # TODO; have GoNES do this
    my $offset = hex($ines);
    open(FH, "<$file")||die"cannot open $file: $!";
    seek(FH, $offset - 0x10, 0);        # get some bytes before
    $/ = \0x10;
    my $before = <FH>;
    my $after = <FH>;
    
    my $context = $before . $after;
    my @bytes = split //, $context;
    for my $byte (@bytes) {
        printf "%.2x ", ord($byte);
    }
    printf "%s => %.8X:%.2X %s\n", $code, $offset, $value, $title;
}

