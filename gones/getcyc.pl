# Open http://www.obelisk.demon.co.uk/6502/reference.html 
# Select all, copy, pbpaste then pipe to this script to get CSV of opcode table
while(<>) {
    chomp;
    if ($_=~m/^[\$]([0-9A-F][0-9A-F])/) {
        $opcode = hex($1);
        $h = lc $1;
        chomp($bytes = <>);
        chomp($s = <>);
        $cycles = substr($s, 0, 1);
        $note = substr($s, 1);
        print "{0x$h, $cycles}, ";
        $note =~ tr/()//d;
        print " // $note" if length($note) > 2;
        print "\n";
    }
}
