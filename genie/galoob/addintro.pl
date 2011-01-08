#!/usr/bin/perl
# Created:20110108
# By Jeff Connelly

# Add intro text from wip

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

open(DB,"<wip")||die;
my %newlines;  # lines of code/intro/info keyed by game
my %new_game2id;
while(<DB>) {
    chomp;
    my ($name, $id, @rest) = split /\t/;

    $newlines{$name} = [] if !exists $newlines{$name};
    $new_game2id{$name} = $id;
    push @{$newlines{$name}}, [@rest]; 
}

my $ourcount = scalar keys %ourlines;
my $newcount = scalar keys %newlines;
die "missing games: $ourcount != $newcount" if $ourcount != $newcount;

open(GL,"<order.txt")||die; # preserve order of games, don't sort
my@gamelist;while(<GL>){chomp;my($name,$id,$gn)=split/\t/;push@gamelist,$name;}

for my $game (@gamelist)
{
    my @ourlines = @{$ourlines{$game}};
    my @newlines = @{$newlines{$game}};

    # TODO: add intro lines

    for (@ourlines)
    {
        my @ourfields = @$_;
        my $ourline = join("\t", @ourfields);
        # Pass through unchanged
        print "$game\t$game2id{$game}\t$ourline\n";
    }

#    my $delta = @newlines - @ourlines;
#    #printf "%-80s %2d -> %2d    (%+2d)\n", $game, scalar(@ourlines), scalar(@newlines), $delta;
#
#    for (@newlines) {
#        my @fields = @$_;
#        shift @fields; # remove "filename" field so it can be compared
#        my $newline = join("\t", @fields);
#        print "$game\t$new_game2id{$game}\t$newline\n";
#    }
}

