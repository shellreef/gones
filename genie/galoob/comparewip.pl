#!/usr/bin/perl
# Created:20110107
# By Jeff Connelly

# Differentiate wip and all-nev.csv to find missing lines

open(DB,"<all-nev.csv")||die;
my %ourlines;  # lines of code/intro/info keyed by game
my %game2id;
while(<DB>) {
    chomp;
    my ($name, $id, @rest) = split /\t/;

    $ourlines{$name} = [] if !exists $ourlines{$name};
    $game2id{$name} = $id;
    push @{$ourlines{$name}}, [@rest]; 
}

open(DB,"<wip")||die;
my %newlines;  # lines of code/intro/info keyed by game
my %new_game2id;
while(<DB>) {
    chomp;
    my ($name, $id, @rest) = split /\t/;

    $newlines{$name} = [] if !exists $newlines{$name};
    $new_game2id{$name} = $id;
    push @{$newlines{$name}}, [@rest]; 
}

my $ourcount = scalar keys %ourlines;
my $newcount = scalar keys %newlines;
die "missing games: $ourcount != $newcount" if $ourcount != $newcount;

for my $game (sort keys %ourlines) 
{
    my @ourlines = @{$ourlines{$game}};
    my @newlines = @{$newlines{$game}};

    my $delta = @newlines - @ourlines;
    printf "%-80s %2d -> %2d    (%+2d)\n", $game, scalar(@ourlines), scalar(@newlines), $delta;

    #splice @comparable, 2, 1; # remove "filename" field so it can be compared
}

