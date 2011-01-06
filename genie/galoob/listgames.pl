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
        my $gameid = "";
        while (!length($gameid)) {
            chomp($gameid = <FH>);
        }

        print "$gameid,$name\n";
        ++$i;
    }
}
