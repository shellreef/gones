#!/usr/bin/perl
# Created:20110107
# By Jeff Connelly

# Find where codes came form

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

        my @sources = where_from($code);
        print scalar(@sources), "\t", $code, "\t", join("\t", @sources), "\n";
    }
}

for my $code (sort keys %code2game) {
    my $game = $code2game{$code};
    printf "%-40s %s\n", $code, $game;
}

# Find where a code came from
sub where_from  # not to be confused with COME FROM
{
    my $code = $_[0];
    my $cmd = qq(fgrep "$code" sources/* | awk -F: '{print \$1}');
    chomp(my @files = `$cmd`);
    @files = map{s/sources\///;$_}@files;
    my %files;
    $files{$_} = 1 for @files;

    delete $files{"nev9.txt"} if @files > 1 && $files{"nev8.txt"};  # nev9 superset of nev8
    @files = keys %files;
    delete $files{"nev8.txt"} if @files > 1;  # nev8 superset of other updates+original codebook
    @files = keys %files;
    delete $files{"nev9.txt"} if @files > 1;
    @files = keys %files;

    if (@files == 0) {
        if ($code eq "SZEAYXVK") {
            @files = ("nev8");      # really from original codebook, but misprinted
        } else {
            die "unable to locate $code" if @files == 0;
        }
    }

    return @files;
}
