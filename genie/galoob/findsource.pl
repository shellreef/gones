#!/usr/bin/perl
# Created:20110107
# By Jeff Connelly

# Find what update the Game Genie code came from

use strict;

my %sourcefiles = %{read_sources("sources/")};

open(DB,"<all-nev.csv")||die;
while(<DB>) {
    chomp;
    my ($name, $id, $type, $no, $code, $title) = split /\t/;

    my @results = look($name);
}

sub look
{
    my ($match) = @_;

    my $found = 0;
    for my $file (keys %sourcefiles) {
        my @lines = @{$sourcefiles{$file}};
        for my $line (@lines) {
            if (lc($match) eq lc($line)) {
                $found = 1;
                print "Found $match in $file\n";
            }
        }
    }
    die "unable to locate $match" if !$found;
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
