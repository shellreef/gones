# Open http://www.obelisk.demon.co.uk/6502/reference.html 
# Select all, copy, pbpaste then pipe to this script to get CSV of opcode table
while(<>) {
    chomp;
    if ($_=~m/^[\$]([0-9A-F][0-9A-F])/) {
        $opcode = hex($1);
        chomp($bytes = <>);
        chomp($cycles = <>);
        print "$1,$bytes,$cycles\n";
    }
}
