#!/usr/bin/perl
# Created:20110105
# By Jeff Connelly
#
# Parse Game Genie codebook for game listing
open(FH, "<nev8.txt")||die;
my $i = 1;
while(<FH>) {
    chomp;
    if (m/Game$/ || $_ eq 'Adventures of Lolo 3(tm)') {
        my $name = $_;
        chomp(my $gameid = <FH>);

        print "$i,$gameid,$name\n";
        ++$i;
    }
}
