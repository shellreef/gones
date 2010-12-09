#!/usr/bin/perl
# Created:20101203
# By Jeff Connelly

# Diff two instruction traces, only on the processor state (not inconsequential data representation)

open(A, "<$ARGV[0]") || die "cannot open $ARGV[0]: $!";   # expected
open(B, "<$ARGV[1]") || die "cannot open $ARGV[1]: $!";   # actual

$quiet = $ARGV[2] eq "-q";   # only show failures

my $differences = 0;
while()
{
    my ($a, $b);

    chomp($b = <B>);
    last if !length($b);
    if (length($b) < 80) {
        # if lines are not traces, then they're probably informative debugging..print and resynchronize
        print "?$b\n";
        next
    }

    chomp($a = <A>);
    last if !length($a);


    # Full trace is:
    #                   only extract beginning here:  v-----------------------|-------|
    # v------------< well, also here too
    # C000  4C F5 C5  JMP $C5F5                       A:00 X:00 Y:00 P:24 SP:FD CYC:  0 SL:241
    my $state_a = substr($a, 48, 25+8) . substr($a, 0, 14);
    my $state_b = substr($b, 48, 25+8) . substr($b, 0, 14);

    if ($state_a eq $state_b) {
        print " $a\n" if !$quiet;   # always show line, for reference (or $b??)
    } else {
        # Not as expected! Log what we got instead
        print "-$a\n+$b\n";
        $differences += 1;
    }
}

exit($differences ? 1 : 0);
