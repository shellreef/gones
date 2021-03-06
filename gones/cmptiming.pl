#!/usr/bin/perl
# Created:20101208
# By Jeff Connelly

# Ugly hack to compare CPU timing differences of nestest
# This is different than tracediff.pl, which performs a 
# line-by-line comparison, so an error in cycle count will
# propagate indefinitely. This one compares incremental cycle
# counts, so you can isolate the errors more easily.

open(A, "<data/nestest.log")||die;
open(B, "<data/nestest.log-actual")||die;

$failed = 0;
while()
{
    chomp($a = <A>);
    last if length($a) == 0;
    $comments = "";
    while()
    {
        chomp($b = <B>);
        if (length($b) >= 80) {
            last; 
        } elsif (length($b) == 0) {
            die "unexpected end-of-file in log.actual: missing/truncated instruction trace?";
        } else {
            $comments .= " $b\n";
        }
    }

    # extract running cycle count
    ($ca) = $a =~ m/CYC:([\d ]+)/;
    ($cb) = $b =~ m/CYC:([\d ]+)/;

    # compare from previous instruction to find how much it took
    $delta_ca = $ca - $prev_ca;
    $delta_cb = $cb - $prev_cb;

    # wrapped around
    $delta_ca += 341 if $delta_ca < 0;
    $delta_cb += 341 if $delta_cb < 0;

    # assuming NTSC, translate PPU pixel clock to CPU cycle count
    $cpu_ca = $delta_ca / 3;
    $cpu_cb = $delta_cb / 3;

    if ($cpu_ca == $cpu_cb) {
        #print " $a\n";
    } else {
        printf "-%-88s CPU cycles: $cpu_ca\n", $prev_a;
        $diff = $cpu_cb - $cpu_ca;
        printf "+%-88s CPU cycles: $cpu_cb ($diff)\n", $prev_b;
        print "$comments\n";
        $failed += 1;
    }

    # to compare to next
    $prev_ca = $ca;
    $prev_cb = $cb;

    $prev_a = $a;
    $prev_b = $b;
}

print "Timing: Failed $failed\n" if $failed != 0;
print "Timing: Passed all tests\n" if $failed == 0;

exit($failed != 0);
