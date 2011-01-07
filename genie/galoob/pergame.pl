#!/usr/bin/perl
# Created:20110106
# By Jeff Connelly

# Count codes created per game

open(FH, "<all-nev.csv")||die;
my %pergame;
while(<FH>)
{
    chomp;
    my @fields = split /\t/;

    my ($game, $id, $type, $no, $code, $title) = @fields;
    next if $type ne "code";

    # Approximate http://www.gshi.org/?s=v2&sys=21&pp=all
    if ($game =~ m/^The /) {
        $game =~ s/^The //;
        $game .= ", The";
    }
    if ($game =~ m/^A /) {
        $game =~ s/^A //;
        $game .= ", A";
    }
    if ($game =~ m/Hudson's /) {
        $game =~ s/Hudson's //;
        $game .= ", Hudson's";
    }
    $game =~ s/\(tm\)//g; 


    $pergame{$game} = [] if !exists $pergame{$game};
    push @{$pergame{$game}}, "$no $code $title";
}

for my $game (sort keys %pergame) 
{
    my $count = @{$pergame{$game}};


    printf "%-80s %d\n", $game, $count;
}
