#!/usr/bin/perl -I/home/groups/g/gg/ggdb/cgi-bin/lib
use CGI qw/:standard/;
require NES::GG;

print header, "<html>it works\n<pre>";
$gg = new NES::GG;
print $gg->encode(0, 1);
eval "use NES::PAR";
print $@;

