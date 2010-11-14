#!/usr/bin/perl -I/home/groups/g/gg/ggdb/cgi-bin/lib
# Created:2003-04-03
# By Jeff Connelly

use NES::Dis6502;
use CGI qw/:standard/;

$cpu = new NES::Dis6502;

print header,qq#<html><head><title>Disassemble 6502</title></head><body>\n<pre>#;

$in_code = param('code');
$in_code = "EA 20 22 44 EA EA EA" if !defined($in_code);
($code = $in_code) =~ s/\s*([0-9A-Za-z][0-9A-Za-z])\s*/chr(oct"0x$1")/ge;
$code = [ map { ord } split //, $code ];
$offset = 0;
while($offset <= @$code) {
    ($asm, $offset) = $cpu->op2asm($code, $offset);
    print "$offset\t$asm\n";
    last if $offset >= 10;
}

print qq#</pre><form action="$ENV{SCRIPT_NAME}" method=get>
<input type=submit value="Disassemble"><br>
<textarea name=code rows=20 width=100%>$in_code</textarea>
</form>#;

