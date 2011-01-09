#!/usr/bin/perl
# Created:20110107
# By Jeff Connelly

# Filter

# Inspired by Ungoodmerge

use strict;
use Data::Dumper;

our $VERBOSE = 0;
our $ROOT = "roms/3.14/extracted/";
#our $OUT = undef;       # link good stuff here; undef for a dry run
our $OUT = "roms/best/";
our $OUT_FLAT = "roms/bestflat/";  # non-hierarchical

# Plain strings to filter on
our @OMISSIONS = (
    "(PD)",     # "public domain", not commercial games (although some might be interesting)"
    "Hack)",    # ROM-hacked by someone, likely not very interesting
    #"Demo",    # don't want to match Demon
    #"BIOS",    # could be interesting
    "(VS)",    # Vs Unisystem, maybe include?
    "(GBA e-Reader)", # http://en.wikipedia.org/wiki/Nintendo_e-Reader  re-releases, usually more limited
    # Not applicable or present in NES enough to be worth filtering
    #"(MB)", 
    #" Sample",
    #"(MB2GBA)",
    #"Demo)",
    #"Preview Version)",
    #"(Debug",
    "(Prototype)",  # interesting, but not the latest release
    "(Older Beta)", #


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
    "W",        # World
    "B",
    "GC",
    "E-GC",
    "GBA e-Reader",
    "NSS",
    "PC10",
    "unknown",      # No region
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

# File type extensions to allow
our @ALLOWED_FILE_TYPES = (
    "nes",              # iNES
    #"unif", "unf",     # UNIF, an alternate format, not bad, but all games are also available in iNES
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

    # Nothing left for this game
    # TODO: option to maybe leave in one of the matches for each game?
    next if @maybe == 0;

    # Group files together by identified region, or unknown
    my %region2files;
    for my $file (@maybe) {
        # cannot simply extract (..) because it is not always the region
        my $identified_region = 0;
        for my $region (@REGION_PRIORITY) {
            if (index($file, "($region)") != -1) {
                $region2files{$region} = [] if !exists $region2files{$region};
                push @{$region2files{$region}}, $file;
                $identified_region = 1;
                last;
            }
        }
        if (!$identified_region) {
            $region2files{"unknown"} = [] if !exists $region2files{"unknown"};
            push @{$region2files{"unknown"}}, $file;
        }
    }

    # Sort regions found by desired priority
    my @regions_found = sort{index_of($b, @REGION_PRIORITY) <=> index_of($a, @REGION_PRIORITY)} keys %region2files;
    die "what? @regions_found $game" if @regions_found == 0;    # should not happen, found some above & should all be categorized

    # Add all from most desired region
    my $best_region = pop @regions_found;
    print "$game best region: $best_region\n";
    my @good = @{$region2files{$best_region}};
    $count_good += @good;

    print scalar(@good), " $game";
    for my $good (@good) {
        print "\t$good\n";

        next if !defined($OUT);

        mkdir($OUT) if !-e $OUT;
        mkdir($OUT_FLAT) if !-e $OUT_FLAT && defined($OUT_FLAT);
        mkdir("$OUT$game") if !-e "$OUT$game";
        system("ln", "$ROOT$game/$good", "$OUT$game/$good");
        die "failed to link? $OUT$game/$good" if $?;
        system("ln", "$ROOT$game/$good", "$OUT_FLAT$good");
        die "failed to link2? $OUT_FLAT$good" if $?;
    }
}

# Return the index of an element within an array, like index() but for arrays not strings
sub index_of
{
    my ($element, @array) = @_;
    for my $i (0..$#array) {
        return $i if $array[$i] eq $element;
    }
    return 999;   # high, if unindentifiable (not low, not -1)
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

    my $allowed = 0;
    my $extension = (split(/\./, $file))[-1];
    for (@ALLOWED_FILE_TYPES) {
        if ($extension eq $_) {
            $allowed = 1;
            last;
        }
    }
    return "unacceptable file type: $extension" if !$allowed;
    

    # don't filter
    return undef;
}
