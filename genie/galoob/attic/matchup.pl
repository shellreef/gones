#!/usr/bin/perl
# Created:20110105
# By Jeff Connelly

# Match Galoob's game names to GoodNES's

use Text::Soundex;

# Read lists
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

# Build a mapping, based on exact matches and user input, as best as we can
my %gg2gn;

for my $galoob (@galoob) {
    my ($guess);

    $guess = $galoob;
    $guess =~ s/\(tm\)//g;
    $guess =~ s/Game$//;
    if ($guess =~ m/^The /) {
        $guess =~ s/^The //;
        $guess .= ", The";
    }
    if ($guess =~ m/^A /) {
        $guess =~ s/^A //;
        $guess .= ", A";
    }

    # First try our chances at an "exact" match, all same letters/numbers
    (my $exact = lc($guess)) =~ tr/a-z0-9//cd;
    if (exists $goodnes_exact{$exact}) {
        my $goodnes = $goodnes_exact{$exact};
        print "FOUND: $guess -> $goodnes\n";
        $gg2gn{$galoob} = $goodnes;
        next;
    }


    my $soundex = soundex($guess);
    my @approx = @{$goodnes_soundex{$soundex}};
    if (@approx) {
        print "No exact match for $galoob ($exact), guesses:\n";
        for my $i (0..$#approx) {
            my $approx = $approx[$i];
            my $index = $i + 1;
            print "$index. $approx\n";
        }
        chomp(my $index = <>);
        my $goodnes;
        if ($index == 0) {
            print "Nothing matches? Oh well..\n";
            $goodnes = undef;
        } else {
            $goodnes = $approx[$index - 1];
            print "Matching to: $goodnes\n";
        }
        $gg2gn{$galoob} = $goodnes;
    } else {
        print "No idea what this is: $galoob\n";
        $gg2gn{$galoob} = $undef;
    }
}

# Save matches
open(OUT, ">gg2gn.csv") || die;
for my $galoob (sort keys %gg2gn) {
    my $goodnes = $gg2gn{$galoob};

    $goodnes = "UNKNOWN" if !defined($goodnes);
    print OUT "$galoob\t$goodnes\n";
}
close(OUT);
