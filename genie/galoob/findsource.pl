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

    my %foundfiles;
    for my $file (sort keys %sourcefiles) {
        my @lines = @{$sourcefiles{$file}};
        for my $i (0..$#lines) {
            my $line = $lines[$i];
            if (lc($game) eq lc($line)) {   # game header line
                next if $file eq "nev9.txt" && $foundfiles{"nev8.txt"};    # nev8=nev9+additions; only match nev9 for those additions

                $foundfiles{$file} = 1;
                print "\t$file\n";
                my @theirlines;

                # Read more of lines from original source data
                my $j = 0;
                my $startedcodes = 0;
                while($i < $#lines && $j < $#ourlines) {
                    ++$i;
                    my $theirline = $lines[$i];
                    next if $theirline eq $game2id{$game};      # skip game abbreviation
                    next if $theirline =~ m/CODE/i && $theirline =~ m/KEY IN/i && $theirline =~ m/EFFECT/i;   # skip field headers
                    my ($theirtype, @theirfields);

                    if ($theirline =~ m/^\d+/) {
                        $theirtype = "code";
                        $startedcodes = 1;
                        @theirfields = split /\t/, $theirline;
                    } else {
                        $theirtype = $startedcodes ? "info" : "intro";
                        @theirfields = $theirline;
                    }

                    my @ourfields = @{$ourlines[$j]};
                    my $ourtype = shift @ourfields;     # what we think this line is (code, intro, info)
                    my $ourline = join("\t", @ourfields);

                    my $match;
                    if (basicallyequal($ourline, $theirline)) {
                        $match = "+";
                    } else {
                        $match = "-";

                        # Perhaps it wrapped?
                        my $next = $lines[$i + 1];
                        if (basicallyequal($ourline, "$theirline $next")) {
                            ++$i;
                            $match = "*";
                            $theirline = "$theirline $next";
                        }
                    }

                    printf "%1s%-8s<%-100s|%-100s>%s\n", $match, $ourtype, $ourline, $theirline, $theirtype;
                    ++$j;
                }
            }
        }
    }
    die "unable to locate $game" if (scalar keys %foundfiles) == 0;


}

# Return whether two strings are equal except punctuation
sub basicallyequal
{
    my ($a, $b) = @_;
    my ($a2, $b2);
    ($a2 = lc($a)) =~ tr/A-Za-z0-9//cd;
    ($b2 = lc($b)) =~ tr/A-Za-z0-9//cd;

    return $a2 eq $b2;
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
