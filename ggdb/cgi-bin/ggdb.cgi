#!/usr/bin/perl -wI/home/groups/g/gg/ggdb/cgi-bin/lib
# GGDB
# Note: perl-5.005, not 5.6.0

use Apache::DBI;
use CGI qw/:standard fatalsToBrowser/;
use strict;
use DBI;
use NES::GG;
use NES::PAR;

print header, "<html>
<head>
<title>GGDB</title>
</head>
<body>
";

$SIG{__DIE__} = sub { print "fatal error: $_[0]\n"; }; 

my ($dbh, $sth, $cmd, $gg);

$gg = new NES::GG;

$dbh = DBI->connect("DBI:mysql:database=g24211_ggdb;host=mysql4-g.sourceforge.net", 
    "g24211rw", "notsosecret589",
    { RaiseError => 1, AutoCommit => 1 });

# /X/Y => game X, code Y
my ($gid, $cid) = split m#/#, substr($ENV{PATH_INFO}, 1);

if (!$gid) {
    list_games();
} elsif ($gid && !$cid) {
    list_codes($gid);
} else {
    show_code($gid, $cid);
}

$dbh->disconnect();

# Show a specific code
sub show_code
{
    my ($gid, $cid) = @_;
    my ($sth, @codes, @row);

    # Lookup name of game
    $sth = $dbh->prepare("SELECT title FROM games WHERE id=?");
    $sth->execute($gid);
    my ($game) = $sth->fetchrow_array;

    $sth = $dbh->prepare("SELECT name,effect,author,submitted,technotes,disbef,disaft " .
        "FROM ggcodes WHERE game=? AND id=?");
    $sth->execute($gid, $cid);
    my ($name, $effect, $author, $submitted, $exp, $before, $after) = $sth->fetchrow_array;

    $before =~ s/\n/<br>\n/g;
    $after  =~ s/\n/<br>\n/g;

    $sth = $dbh->prepare(
"SELECT address,value,akey FROM ggcode WHERE game=? AND id=?");
    $sth->execute($gid, $cid);
    $name = uc $name;
    $game = uc $game;
    $name = qq#<b><a href="$ENV{SCRIPT_NAME}/$gid">$game</b></a>: $name#;;
    print "<table border=1 width=100%><td align=center colspan=5>$name</td></tr>";
    print "<tr><td width=25%>"; 
    # Print the encrypted/decoded table
    print "<table width=100% height=100%><tr><th>Encrypted</th>" .
          "<th>Decoded</th><th>Edit</th></tr>\n";

    # Read individual codes from database
    push @codes, [ @row ] while @row = $sth->fetchrow_array;
    if (defined param('v')) {
        my ($a, $v, $k) = (param('a'), param('v'), param('k'));
        push @codes, [ $a, $v, $k ];
    }

    foreach my $row (@codes) {
        my ($addr, $val, $key) = @$row;
        my ($encrypted) = $gg->encode($addr, $val, $key);
        $addr &= 0x7fff;
        $val &= 0xff;
        $key &= 0xff if defined($key);
        printf "<tr><td><b>%s</b></td><td>%.4X:%.2X%s</td><td>%s</td></tr>\n",
         $encrypted, $addr, $val, defined($key) ? (sprintf "?%.2X", $key) : "",
         qq#<form action="$ENV{SCRIPT_NAME}/$gid/$cid" method=post># .
         qq#<input size=3 name=v><input type=hidden name=a value=$addr># .
         ((defined $key) ? 
         qq#<input type=hidden name=k value=$key>(k)# : "") . "</form>\n";
    }
    print "</table></td><td>";
    # Print the before/after table 
    print "<table border=1 width=100% height=100%><tr><th>Before</th><th>After</th></tr><tr>" .
         "<td>$before</td><td>$after</td></tr></table>\n";
    # Print the information table
    print "<tr><td colspan=5><table border=0 width=100%><tr><th>Creator</th></tr>" .
        "<tr><td>$author</td></tr><tr><th>Effect</th></tr><tr><td><p>$effect</tr>" .
        "</tr><th>Technical Explanation</th>" .
        "</tr><tr><td>$exp</td></tr></table></td></tr></table>\n"; 
}

# Show codes for a certain game, and game info
sub list_codes
{
    my ($gid) = @_;
    my ($sth);

    my ($name, $author, $effect, $exp, $before, $after, $codes) = (
	param('name'), param('author'), param('effect'), 
        param('exp'), param('before'), param('after'), param('codes'));
    if (defined($name)) {
        # Insert meta-info
        $sth = $dbh->prepare("INSERT INTO ggcodes (game,name,author" .
            ",effect,technotes,disbef,disaft,submitted) VALUES(?,?,?,?,?,?,?,?)");

        my ($now);
        my ($year, $month, $day);
        (undef, undef, undef, $day, $month, $year) = localtime(time);
        $year += 1900; $month++; $now = sprintf "%.4d-%.2d-%.2d", $year, $month, $day;       
 
        $sth->execute($gid, $name, $author, $effect, $exp, $before, $after, $now);

        # Find the auto-assigned code id
        $sth = $dbh->prepare("SELECT id FROM ggcodes WHERE name=? AND game=?");
        $sth->execute($name, $gid);
        my ($cid) = $sth->fetchrow_array();

        # Insert code for the meta-codes
        my @codes = split /\s+/, $codes;
        foreach my $code (@codes) {
            my ($address, $value, $key);

            if ($gg->looks_encoded($code)) {         # GG format
                ($address, $value, $key) = $gg->decode($code); 
            } elsif ($gg->looks_decoded($code)) {    # hex format
                ($address, $value, $key) = $gg->parse_gg($code);
            } else {
                print "<p>couldn't make sense of '$code'\n";
            }

            $sth = $dbh->prepare("INSERT INTO ggcode (game,id,address,value,akey) ".
                " VALUES(?,?,?,?,?)");
            $sth->execute($gid, $cid, $address, $value, $key);
        }
    }

    $sth = $dbh->prepare("SELECT title,publisher,developer,released,info,prg,chr,mappers"
	. " FROM games WHERE id=?");
    $sth->execute($gid);
    print "<table border=1 width=50%>\n";
    my ($game, $publisher, $developer, $date, $info, $prg, $chr, $mappers) = 
	$sth->fetchrow_array;
    print qq#<tr><th><a href="$ENV{SCRIPT_NAME}/">Game</a></th><td>$game</td></tr>
<tr><th>Publisher</th><td>$publisher</td></tr>
<tr><th>Developer</th><td>$developer</td></tr>#;

    print "<tr><th>PRG Banks</th><td>$prg</td></tr>\n"   if $prg;
    print "<tr><th>CHR Banks</th><td>$chr</td></tr>\n"   if $chr;
    print "<tr><th>Mappers</th><td>$mappers</td></tr>\n" if $mappers;
 
print "<tr><th>Info</th><td>$info</td></tr>
</table><p>\n";

# TODO: Ordering, with multiple levels, ascending and descending
#    my $order = param('ord');
#    my @order = split /-/, $ord;
#    my $desc = param('d');
#    my @desc = split //, $desc;
#
#    # New code order - place first
#    my ($code_order, $code_desc);
#    my (@code_order, @code_desc);
#    @code_order = uniqify("code", @order);
#    $code_order = join "-", @code_order; # XXXX 

    print <<EOH;
<table border=1 width=100%>
<tr><th>Game Genie Code Name</th>
<th>Code</th>
<th>Creator</th> 
<th>Date Submitted</th></tr>

EOH

    $sth = $dbh->prepare(
"SELECT id,name,author,submitted FROM ggcodes WHERE game=?");

#($sort ? " ORDER BY ?" : "") . ($desc ? " DESC" : ""));
    $sth->execute($gid);
 
    while (my @row = $sth->fetchrow_array)
    {
        my ($id, $name, $author, $submitted) = @row;
        my ($code, $codesth);
        $codesth = $dbh->prepare("SELECT address,value,akey FROM ggcode WHERE id=?");
        $codesth->execute($id);
        while (my @code_row = $codesth->fetchrow_array) {
            my ($address, $value, $key) = @code_row;
            $code .= $gg->encode($address, $value, $key) . " + ";
        } 
        $code =~ s/ \+ $//;
        print qq#<tr><td><a href="$ENV{SCRIPT_NAME}/$gid/$id">$name</a></td># .
              qq#<td>$code</td><td>$author</td><td>$submitted</td></tr>\n#;
    }
    print "</table>\n";

    # Code submission form
    print qq#<p><b>Submit a code for this game</b>
<form action="$ENV{SCRIPT_NAME}/$gid" method=post>
<input type=hidden name=game value="$gid">
<table>
<tr><th>Code Name</th><td><input name=name></td><td><input type=submit value=Submit></tr>
<tr><th>Creator</th><td><input name=author></td></tr>
<tr><th>Effect</th><td><textarea name=effect width=100% rows=5></textarea></td>
    <th>Before</th><td><textarea name=after width=100% rows=6></textarea></td></tr>
<tr><th>Technical Explaination</th><td><textarea name=exp width=100% rows=6></textarea></td>
    <th>After</th><td><textarea name=after width=100% rows=6></textarea></td></tr>
<tr><th>Code List</th><td colspan=2><textarea name=codes width=100% rows=3></textarea></td>
    <td>Enter the codes here in hexadecimal or Game Genie format.</td></tr>
</table></form>#;
}

sub list_games
{
    my ($sth);
    my ($game, $publisher, $developer, $rlsdate, $info, $prg, $chr, $mappers) = 
        (param('game'), param('publisher'), param('developer'), param('rlsdate'),
         param('info'), param('prg'), param('chr'), param('mappers'));
    if (defined($game)) {    # wants to add a new game
        $sth = $dbh->prepare("INSERT INTO games (title,publisher,developer,released,info," .
            "prg,chr,mappers) VALUES(?,?,?,?,?,?,?,?)");
        $sth->execute($game, $publisher, $developer, $rlsdate, $info, $prg, $chr, $mappers);
    } 

    $sth = $dbh->prepare(
"SELECT id,title,publisher,developer,released,info FROM games ORDER BY title");
    $sth->execute();
    print "<table border=1 width=100%>\n";
    print "<tr><th>Game</th><th>Publisher</th><th>Developer</th><th>Release Date</th>"
        . "<th># GG</th><th># PAR</th></tr>\n";
    while (my @row = $sth->fetchrow_array) 
    {
        my ($gid, $game, $publisher, $developer, $date, $info) = @row;
        my ($ggnum, $ggsth, $parnum, $parsth);
        $ggsth = $dbh->prepare("SELECT COUNT(*) FROM ggcodes WHERE game=?");
        $ggsth->execute($gid);
        ($ggnum) = $ggsth->fetchrow_array;

        $parsth = $dbh->prepare("SELECT COUNT(*) FROM parcodes WHERE game=?");
        $parsth->execute($gid);
        ($parnum) = $parsth->fetchrow_array;
        print qq#<tr><td><a href="$ENV{SCRIPT_NAME}/$gid">$game</a></td># .
              qq#<td>$publisher</td><td>$developer</td><td>$date</td>#.
              qq#<td align=right>$ggnum</td><td align=right>$parnum</td></tr>#;
    }
    print qq#</table><p>Want to add a game not listed here? Use the form below:
<form action="$ENV{SCRIPT_NAME}" method=post>
<table><tr><th>Game</th><td><input name=game></td></tr>
<tr><th>Publisher</th><td><input name=publisher></td></tr>
<tr><th>Developer</th><td><input name=developer></td></tr>
<tr><th>Release Date</th><td><input name=rlsdate></td></tr>
<tr><th>PRG Banks</th><td><input name=prg></td></tr>
<tr><th>CHR Banks</th><td><input name=chr></td></tr>
<tr><th>Mappers</th><td><input name=mappers></td></tr>
<tr><th>Short Description</th></td><td><textarea name=info width=100% rows=5></textarea></td></tr>
<tr><td colspan=2 align=center><input type=submit value="Add Game"></tr>
</table></form>#;
}

# Delete non-unique elements from an array
sub uniqify
{
    my (@ary) = @_;
    my (%ary);
    my (@ret);

    foreach my $elt (@ary) {
        next if $ary{$elt};
        $ary{$elt}++;
        push @ret, $elt;
    }
    return @ret;
}

