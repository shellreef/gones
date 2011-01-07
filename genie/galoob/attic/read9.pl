#!/usr/bin/perl
# Created:20110106
# By Jeff Connelly

# Read nev9.txt games
open(FH, "<nev9.txt")||die;
<FH>; <FH>;
while(<FH>)
{
    chomp;
    last if $_ eq "";
}

# Read index
my @games;#I wish
while(<FH>)
{
    chomp;
    my ($number, $name) = split /\t/;

    push @games, $name;
    last if $number == 481;
}
<FH>; <FH>;

#print join("\n", @games), "\n";

# Now read codes
while(<FH>) 
{
    chomp;
    my ($name) = $_;
    my @codelines;

    # Look for beginning of fields; read intro text if any
    my @intro;
    while(<FH>)
    {
        chomp;
        last if m/^CODE/i && m/KEY IN/i && m/EFFECT/i;
        if (m/^\d+\t/) {
            push @codelines, $_;        # some don't have fields and jump right into codes
            last
        }
        push @intro, $_;
    }
    my $id;
    # Abbreviation for game, if any (useless because it is not unique)
    $id = pop(@intro);

    #print "$name\n";
    #print join("\n|", @intro), "\n", "\n";

    # Actual lines of code
    while(<FH>)
    {
        chomp;
        last if $_ eq "";
        push @codelines, $_; 
    }
    if ($name !~ m/The Jungle Book/) {   # the one game missing a double blank
        chomp(my $blank = <FH>);
    }
    die "missing blank: $blank" if $blank ne "";

    # Show data in CSV format
    for my $intro (@intro) 
    {
        print join("\t", $name, $id, "intro", $intro), "\n";
    }
    for my $line (@codelines)
    {
        if ($line =~ m/^(\d+)\s+([APZLGITYEOXUKSVN]{6,8}[^\t]*)\t(.*)/) {
            my ($no, $code, $title) = ($1, $2, $3);

            print join("\t", $name, $id, "code", $no, $code, $title), "\n";
        } else {
            print join("\t", $name, $id, "info", $line, "\n");
        }
    }
}
