#!/usr/bin/perl
# Created:20110107
# By Jeff Connelly

# Find what update the Game Genie code came from

use strict;

my %sourcefiles = %{read_sources("sources/")};

# Read our lines of codes for each game
open(DB,"<all-nev.csv")||die;
my %ourlines;  # lines of code/intro/info keyed by game
my %game2id;
while(<DB>) {
    chomp;
    my ($name, $id, @rest) = split /\t/;

    $ourlines{$name} = [] if !exists $ourlines{$name};
    $game2id{$name} = $id;
    push @{$ourlines{$name}}, [@rest]; 
}

for my $game (sort keys %ourlines) {
    my @ourlines = @{$ourlines{$game}};
    print "\n$game\n";

    my $found_count = 0;
    for my $file (sort keys %sourcefiles) {
        my @lines = @{$sourcefiles{$file}};
        for my $i (0..$#lines) {
            my $line = $lines[$i];
            if (lc($game) eq lc($line)) {   # game header line
                ++$found_count;
                print "\t$file\n";
                my @theirlines;

                # Read more of lines from original source data
                my $j = 0;
                while($i < $#lines && $j < $#ourlines) {
                    ++$i;
                    my $theirline = $lines[$i];
                    next if $theirline eq $game2id{$game};      # skip game abbreviation
                    next if $theirline =~ m/CODE/i && $theirline =~ m/KEY IN/i && $theirline =~ m/EFFECT/i;   # skip field headers
                    
                    my @ourfields = @{$ourlines[$j]};
                    my $ourline = join("\t", @ourfields);

                    printf "%-80s %-80s\n", $ourline, $theirline;

                    ++$j;
                }
            }
        }
    }
    die "unable to locate $game" if !$found_count;


}


# Read all sources into a hash of arrays
sub read_sources
{
    my $SOURCES = $_[0];

    my %sourcefiles;
    opendir(D, $SOURCES)||die "cannot opendir sources";
    my @files=grep{!m/^\./}readdir(D);
    closedir(D);

    my @results;
    
    for my $file (@files) {
        open(FH, "<$SOURCES/$file") || die "cannot open source: $SOURCES/$file: $!";
        $sourcefiles{$file} = [];
        while(<FH>) {
            chomp;
            push @{$sourcefiles{$file}}, $_;
        }
    }

    return \%sourcefiles;
}
