#!/usr/bin/perl
# Created:20110107
# By Jeff Connelly

# Filter

# Inspired by Ungoodmerge

use strict;

our $VERBOSE = 0;
our $ROOT = "roms/3.14/extracted/";
# Plain strings to filter on
our @OMISSIONS = (
    "(PD)",     # "public domain", not commercial games (although some might be interesting)"
    "Hack)",    # ROM-hacked by someone, likely not very interesting
    #"Demo",    # don't want to match Demon
    #"BIOS",    # could be interesting
    "(VS)",    # Vs Unisystem, maybe include?
    # Not applicable or present in NES enough to be worth filtering
    #"(MB)", 
    #" Sample",
    #"(MB2GBA)",
    #"Demo)",
    #"Preview Version)",
    #"(Debug",

    #"(Sample)",    # actually would be cool to hack samples..
    #"(Unl)",       # very important to allow unlicensed games!
);

# Tags in [brackets] to omit
our @UNWANTED_TAGS = (
    qr/^b/,    # bad dump, who wants that??
    qr/^o/,    # overdumps, also not useful
    qr/^h/,    # hacked ROM dumps, usually for FFE or M03
    #qr/^x/,    # bad checksums - haven't found
    #qr/^BF/,   # Bung's flashcart
    qr/^t/,    # trainer
    qr/^p/,    # pirated (usually to remove copyrights)
    qr/^f/,    # fixed to run better on emulator/copier
    #qr/a./,   # alternate versions - keep these, codes might be different! (not PRG0/PRG1?)

    qr/^T[+-]/,     # translates, old(-) or new (+)
    # TODO: filter games with intros
    );

# Wanted region codes, highest priority first
our @REGION_PRIORITY = (
    "U",        # USA
    "UE",       # USA + Europe
    "JU",       # Japan + Europe
    "4",
    "5",
    "8",
    "E",        # Europe
    "JE",       # Japan + Europe
    "UK",       # United Kingdom
    "A",        # Australia
    "PAL",
    "W",
    "B",
    "GC",
    "E-GC",
    "GBA e-Reader",
    "NSS",
    "PC10",
    );

# Regions you don't want
our @REGION_EXCLUDE = (
    "J",        # Japan (only)
    "F",        # France
    "Ch",       # Chinese
    "G",        # German
    "S",        # Spanish
    "I",
    "ST",
    "R",
    "GR",
    "K",
    "J-AC",
    "KC",
    "SW",
    "iQue",
    "NL",
    "HK",
    "D",
    "FC",
    "C"
    );
    

opendir(D, $ROOT) || die;
my @games = grep{!m/^\./}readdir(D);
closedir(D);

our ($count_good, $count_bad, $count_total) = (0, 0);
for my $game (@games) {
    filter_game($game);
}

print "Total: $count_total (accepted $count_good, rejected $count_bad)\n";

sub filter_game
{
    my ($game) = @_;

    opendir(D,"$ROOT$game")||die "cannot open $ROOT$game: $!";
    my (@files) = grep{!m/^\./}readdir(D);
    closedir(D);

    my (@maybe, @bad);

    # Filename-based, independent filters
    for my $file (@files) {
        my $reason = filter_file($file);
        if (defined($reason)) {
            print "-$file  reason: $reason\n" if $VERBOSE;
            push @bad, [$file, $reason];
            $count_bad++;
        } else {
            print "~$file\n" if $VERBOSE;
            push @maybe, $file;
        }
        ++$count_total;
    }

    my (@good);

    if (@maybe == 0) {
        # Nothing
    } elsif (@maybe == 1) {
        # Only one choice
        push @good, pop @maybe;
        $count_good++;
    } else {
        # Take first from region, in priority given
        OUTER: for my $region (@REGION_PRIORITY) {
            for my $file (@maybe) {
                if (index($file, "($region)") != -1) {
                    push @good, $file;
                    $count_good++;
                    last OUTER;
                } 
            }
        }

        # did we filter everything?
        if (@good == 0) {
            die "can't figure out $game: @maybe\n@good\n";
        }
    }

    print "$game: @good\n";
}

# Return 
sub filter_file
{
    my ($file) = @_;

    for (@OMISSIONS) {
        return "string match: $_" if index($file, $_) != -1 
    }

    for my $re (@UNWANTED_TAGS) {
        while ($file =~ m/\[([^]]+)\]/g) {
            my $tag = $1;

            return "tag: $tag" if $tag =~ $re;
        }
    }
    
    for (@REGION_EXCLUDE) {
        return "region: $_" if index($file, "($_)") != -1;
    }

    # don't filter
    return undef;
}
