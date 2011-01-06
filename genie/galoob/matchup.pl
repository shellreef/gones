#!/usr/bin/perl
# Created:20110105
# By Jeff Connelly

# Match Galoob's game names to GoodNES's

use Text::Soundex;

open(FH, "<gamelist-gg.csv")||die;
my @galoob;
while(<FH>) {
    chomp;
    my ($id, $name) = split /,/, $_, 2;
    push @galoob, $name;
}
close(FH);

open(FH, "<gamelist-goodnes314.csv")||die;
chomp(@goodnes = <FH>);
my %goodnes_soundex;
my %goodnes_exact;
for my $goodnes (@goodnes) {
    # Soundex allows matching on what sounds approximately the same
    my $soundex = soundex($goodnes);
    $goodnes_soundex{$soundex} = [] if !exists $goodnes_soundex{$soundex};
    push @{$goodnes_soundex{$soundex}}, $goodnes;

    # Remove all but letters/digits, and lowercase them, for an exact (well, near to it) one-to-one match
    (my $exact = lc($goodnes)) =~ tr/a-z0-9//cd;
    $goodnes_exact{$exact} = $goodnes;
}
close(FH);

my %gg2gn;

for my $galoob (@galoob) {
    my ($guess);

    $guess = $galoob;
    $guess =~ s/\(tm\)//;
    $guess =~ s/Game$//;

    # First try our chances at an "exact" match, all same letters/numbers
    (my $exact = lc($guess)) =~ tr/a-z0-9//cd;
    if (exists $goodnes_exact{$exact}) {
        my $goodnes = $goodnes_exact{$exact};
        print "FOUND: $guess -> $goodnes\n";
        $gg2n{$galoob} = $goodnes;
        next;
    }

    print "No exact match for $guess ($exact)\n";

    my $soundex = soundex($guess);
    my @approx = @{$goodnes_soundex{$soundex}};
    for my $i (0..$#approx) {
        my $approx = $approx[$i];
        print "$i. $approx\n";
    }
}
