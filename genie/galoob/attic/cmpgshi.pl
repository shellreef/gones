#!/usr/bin/perl
# Created:20110106
# By Jeff Connelly

# Compare NES GG official code counts from GSHI with us

open(GSHI, "<count-gshi.csv")||die;
open(OURS, "<count-ours.csv")||die;

my ($total_gshi, $total_ours, $their_missing, $their_extra);

while()
{
    chomp(my $gshi = <GSHI>);
    chomp(my $ours = <OURS>);

    last if !length($gshi) && !length($ours);

    my ($gshi_name, $gshi_count) = $gshi =~ m/^(.{80}) (.*)/;
    my ($ours_name, $ours_count) = $ours =~ m/^(.{80}) (.*)/;

    $total_gshi += ($gshi_count+0);
    $total_ours += $ours_count;

    if ($gshi_count ne $ours_count) {
        print "-$gshi\n";
        print "+$ours\n";
        print "\n";

        my $delta = $ours_count - $gshi_count;
        $their_missing += $delta if $delta > 0;
        $their_extra += $delta if $delta < 0;

        print "\t\t$their_missing/$their_extra\n";

    }

}

print "GSHI has $total_gshi, you have $total_ours\n";
print "GSHI is missing $their_missing, and has $their_extra extra\n";
