#!/usr/bin/perl
# Created:20110107
# By Jeff Connelly

# Differentiate wip and all-nev.csv to find missing lines

open(WIP, "<wip")||die;
open(ALL,"<all-nev.csv")||die;
while()
{
    chomp(my $wip = <WIP>);
    chomp(my $all = <ALL>);

    last if !length $wip;

    my @wip = split /\t/, $wip;
    my @all = split /\t/, $all;


    my @comparable = @wip;
    splice @comparable, 2, 1; # remove "filename" field so it can be compared

    my $comparable = join("\t", @comparable);

    if ($all ne $comparable) {
        print "-$all\n";
        print "+$comparable\n";
    } else {
        print " $all\n";
    }
}
