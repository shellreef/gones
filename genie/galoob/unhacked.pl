#!/usr/bin/perl
# Created:20110501
# By Jeff Connelly

use strict;
use Data::Dumper;

sub known_games
{
    open(FH, "<gamelist-goodnes314-interesting.csv")||die;
    chomp(my @known = <FH>);
    close(FH);
    return @known;
}

# Get list of games that Galoob made codes for
# Returns hashref of name->id (useless), hashref of Galoob's name->GoodNES name
sub galoobs_info
{
    my (%ids, %gg2gn);
    open(FH, "<gamelist-galoob.csv") || die;
    while(<FH>) {
        chomp;
        my ($galoob_name, $galoob_id, $goodnes_name) = split /\t/;
        $ids{$galoob_name} = $galoob_id;
        $gg2gn{$galoob_name} = $goodnes_name;
    }
    close(FH);

    return (\%ids, \%gg2gn);
}

# Get list of games Gallob made codes for, by their GoodNES name 
sub galoobs_games
{
    my ($gg2gn, %gn2gg);
    (undef, $gg2gn) = galoobs_info();

    %gn2gg = reverse(%$gg2gn);

    return %gn2gg;
}

my @known = known_games();
my %coded = galoobs_games();
#print scalar(@known), " known games\n";
#print scalar(keys %coded), " have codes\n";

my @uncoded;
for my $known (@known)
{
    push @uncoded, $known if !exists $coded{$known};
}

print join("\n", @uncoded);
print "\nNo official Galoob codes for ",scalar(@uncoded)," games\n";
