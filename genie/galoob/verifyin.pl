#!/usr/bin/perl
# Created:20110105
# By Jeff Connelly

# Verify all games in gamelist-gg.csv can be found in nev8.txt

open(FH, "<gamelist-gg.csv") || die;
while(<FH>) {
    chomp;
    my ($id, $name) = split /,/, $_, 2;

    open(IN, "<nev8.txt")||die;
    my $found = 0;
    while(<IN>) {
        chomp;
        if ($_ eq $name) {
            $found = 1;
            last;
        }
    }

    print "$found,$name\n";
    die "incorrect game name: $name" if !$found;
}

