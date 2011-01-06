#!/usr/bin/perl
# Created:20110105
# By Jeff Connelly

# Verify all games in gamelist-gg.csv can be found in nev8.txt

open(FH, "<gamelist-gg.csv") || die;
my %galoobs;
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

    die "incorrect game name: $name" if !$found;

    $galoobs{$name} = 1;
}
close(FH);
print "gamelist-gg.csv OK\n";

open(FH, "<gamelist-goodnes314.csv")||die;
chomp(my @goodnes = <FH>);
close(FH);
#oh
my %goodnes;
$goodnes{$_} = 1 for @goodnes;
close(FH);

open(FH, "<gg2gn.csv") || die;
my $i;
while(<FH>) {
    chomp;
    my ($galoob, $goodnes) = split /\t/;

    die "incorrect Galoob game name: $galoob" if !$galoobs{$galoob};
    die "incorrect GoodNES game name: $goodnes" if !$goodnes{$goodnes};
    $i++;
}
my $j = scalar keys %galoobs;
die "missing some games from gg2gn.csv ($i != $j)" if $i != $j;

print "gg2gn.csv OK\n";
print "$i games matched\n";

