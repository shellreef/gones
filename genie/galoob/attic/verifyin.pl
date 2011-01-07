#!/usr/bin/perl
# Created:20110105
# By Jeff Connelly

# Verify all games in gamelist-gg.csv can be found in nev8.txt

open(FH, "<gamelist-galoob.csv") || die;
my %galoobs;
my %ids;
while(<FH>) {
    chomp;
    my ($name, $id) = split /\t/, $_, 2;

    open(IN, "<all-nev.csv")||die;
    my $found = 0;
    while(<IN>) {
        chomp;
        my ($that_name) = split /\t/;
        if ($that_name eq $name) {
            $found = 1;
            last;
        }
    }
    $ids{$name} = $id;

    die "incorrect game name: |$name|" if !$found;

    $galoobs{$name} = 1;
}
close(FH);
print "gamelist-galoob.csv Galoob game names OK\n";

open(FH, "<gamelist-goodnes314.csv")||die;
chomp(my @goodnes = <FH>);
close(FH);
#oh
my %goodnes;
$goodnes{$_} = 1 for @goodnes;
close(FH);

open(FH, "<gamelist-galoob.csv") || die;
my %map;
my $i;
while(<FH>) {
    chomp;
    my ($galoob, $id, $goodnes) = split /\t/;

    die "incorrect Galoob game name: $galoob" if !$galoobs{$galoob};
    die "incorrect GoodNES game name: $goodnes" if !$goodnes{$goodnes};
    $i++;

    $map{$galoob} = $goodnes;
}
my $j = scalar keys %galoobs;
die "missing some games ($i != $j)" if $i != $j;

print "gamelist-galoob.csv GoodNES names mapping OK\n";
print "$i games matched\n";

open(MASTER,">gggg.csv")||die;
for my $galoob (sort keys %ids) {
    my $id = $ids{$galoob};
    die "no id?!" if !defined($id);
    my $gn = $map{$galoob};
    die "no gn?!" if !defined($gn);
    print MASTER "$galoob\t$id\t$gn\n";
}
close(MASTER);

