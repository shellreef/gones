#!/usr/bin/perl
# Created:20110107
# By Jeff Connelly

# Find what update the Game Genie code came from

use strict;

my $SOURCES = "sources/";

open(DB,"<all-nev.csv")||die;
while(<DB>) {
    chomp;
    my ($name, $id, $type, $no, $code, $title) = split /\t/;

    my @results = look($name);
}

sub look
{
    my ($match) = @_;

    opendir(D, $SOURCES)||die "cannot opendir sources";
    my @files=grep{!m/^\./}readdir(D);
    closedir(D);

    my @results;
    
    for my $file (@files) {
        open(FH, "<$SOURCES/$file") || die "cannot open source: $SOURCES/$file: $!";
        while(<FH>) {
            chomp;

            if (lc($_) eq lc($match)) {
                print "FOUND $match in $file\n";
                push @results, $file;
            }
        }
    }

    die "unable to locate $match!" if !@results;
}
