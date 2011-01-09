#!/usr/bin/perl
# Created:20110108
# By Jeff Connelly

# Analyze codes

use strict;
our $ROOT = "../../roms/best";

open(M,"<gamelist-galoob.csv")||die;
my %gg2gn;
while(<M>) {
    chomp;
    my ($galoob, $id, $goodnes) = split /\t/;
    $gg2gn{$galoob} = $goodnes;

    my $dir = "$ROOT/$goodnes";
    die "can't find $dir" if !-e $dir;
    opendir(D, $dir) || die "can't opendir $dir: $!";
    my @files = grep{!m/^\./}readdir(D);
    die "nothing to be found in $dir" if @files == 0;
    closedir(D);

    #my $file;
    #if (@files == 1) {
    #    $file = pop @files;
    #} else {
        # Show games that have alternate versions
        # There 37 which Galoob made codes for
        # Most are PRG0/PRG1/REV differences, except 4 are different distributors
    printf "%d %-80s %s\n", scalar(@files), $galoob, join(" ** ", @files);
    #}
}
die;

open(FH, "<all-nev.csv")||die;
while(<FH>) {
    chomp;
    my ($game, $id, $type, @rest) = split /\t/;
    if ($type eq "code") {
        my ($source, $no, $code, $title) = @rest;
        print "$game,$code,$title\n";
    }
}
