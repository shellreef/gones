#!/usr/bin/perl
# Created:03-27-01
# By Jeff Connelly

# NES Game Genie conversion module

package NES::GG;

=head1 NAME

NES::GG - Decode and encode NES Game Genie codes

=head1 SYNOPSIS

  use NES::GG;

  my $gg = new NES::GG;
  my ($reladdr, $value, $key) = $gg->decode("SLXPLOVS");
  my @absaddrs = $gg->addrkey_to_addr($reladdr, $key);
  # @absaddrs can now be looked up in a disassembly of a NES ROM

  print $gg->encode(0x1123, 0xBD, 0xDE);
  print $gg->encode(0x45123, 0xBD);  

=head1 DESCRIPTION

NES::GG provides an object-oriented interface for decoding and
encoding Game Genie codes for the NES console system. Convert::NES::GG is
able to find the ROM addresses affected by any NES Game Genie code, aiding
in creation of new codes and understanding of existing ones.

=cut

# Usage: Type a Game Genie code to convert it to a hex code (that is,
#   in aaaa?kk:vv or aaaa:vv format), or a hex code to convert it to a
#   Game Genie code.

# Note that with Game Genie codes, the address can be at most 7FFF. To get
# around this limitation, Game Genie allows codes to contain a "key" --
# a byte that the ROM at the specified address must contain for the code to
# be applied. gg.pl can automatically discover and eliminate keys.

# KEY DISCOVERY: (required loaded ROM, type =..\filename.nes at prompt)
# Type a raw ROM (not .NES) address followed by a colon then a byte value to
# generate a code with a key. This code will only affect the absolute address,
# unless of course an address with the same offset but different bank
# contains the same value. gg.pl considers ROM addresses to start at 0x8000.
# Note that .NES files have a 0x10-byte header which needs to be considered.

# Warning: If you enter an address >7FFF without a ROM loaded, gg.pl will
#   ignore the high bits. If the game cart has only one CHR bank (such as
#   with Super Mario Brothers), this is no problem. Beware of unwanted
#   address modification with six-letter codes.

# KEY ELIMINATION: (required loaded ROM, type =..\filename.nes at prompt)
# gg.pl can also tell you what addresses a code-with-a-key affects by
# searching through the loaded ROM at the same relative address, but in a
# different bank. This is useful for finding out what existing codes do.

# Warning: Don't be too concerned about unwanted modifications in CHR banks.
#   PRG banks is what needs concerning. There is nothing that can be done
#   if a code-with-a-key affects unwanted PRG addresses, save finding a new
#   address to modify.

# Unload a ROM with -.

# Example of using key discovery and elimination:

# C:\>perl gg.pl
# =..\skordie2.nes
# Loaded ..\skordie2.nes -- key discovery/elimination enabled
# PAUYLLLE          ("9 skateboards" code from GG manual)
# 7333?03:09        (hex code with key)
# 27333:09          (affected address)
# 47333:09          (affected address)

# Disassembling skordie2.nes (with an orgin at $8000) reveals $27332 contains
# this instruction:
#
#   LDA #$03
#
# So now you know what PAUYLLLE does..(changes LDA #$03 to LDA #$09).

# Similarity, while making your own codes from a disassembling, ROM addresses
# can be fed to gg.pl and it's key can be discovered automagically:
#
# C:\>perl gg.pl
# =..\skordie2.nes
# Loaded ..\skordie2.nes -- key discovery/elimination enabled
# 8000:00
# Found 0000?A5     (gg.pl looked up $8000 in skordie2.nes, and found $A5)
# AAEAAASZ          (this is code 0000?A5:00)
# AAEAAASZ          (it should be, let's check..)
# 0000?A5:00        (eight digit code)
# 08000:00          (affected address)
# 30000:00          (affected address)
# 38000:00          (affected address)

use strict;
#use warnings;
use Carp;

# Usage: NES::GG->new("ROM");

=head1 new($filename)

Instantiate a new Convert::NES::GG object. If a filename is specified, it
is loaded. Required before any other routines are used.

Example:

  $gg = new Convert::NES::GG;

=cut

sub new
{
    my ($self, $game) = @_;

    $self = {};
    bless($self);

    if (defined($game))
    {
        $self->{ROM} = $game;
        $self->load_ROM($game);
    }
    return $self;
}

#our (@banks);

{
    my %map = (
        A => 0, P => 1, Z => 2, L => 3, G => 4, I => 5, T => 6, Y => 7,
        E => 8, O => 9, X =>10, U =>11, K =>12, S => 13,V => 14,N => 15);
    sub tonum
    {
        my @numbers;
        for my $letter (@_)
        {
            die "invalid letter $letter" if !exists $map{uc $letter};
            push @numbers, $map{uc $letter};
        }
        return @numbers;
    }
}

{
    my @map = split //, 'APZLGITYEOXUKSVN';
    sub toletter
    {
        my @letters;
        for my $digit (@_)
        {
            die "invalid digit $digit"
                if $digit < 0 || $digit > 15;
            push @letters, $map[$digit];
        }
        return @letters;
    }
}

# Decode a code like AAAAAA into address, value, and key (list context),
# or addr?kk:vv (scalar context)
# Example: $gg->decode("AAAAAA");

=head2 ($address, $value, $key, $warning) = decode($ggcode)

=head2 $hex_code = decode($ggcode);

Decode an encoded Game Genie code into it's relative address, value, and
possibly key. In scalar context, returns $warning$address?$key:$value.
$warning is set to "~" if the code really should be eight digits, a space
otherwise.

Example:

    ($reladdr, $value, $key) = $gg->decode("SLXPLOVS");
    print "SLXPLOVS is:", scalar $gg->decode("SLXPLOVS");

=cut
sub decode
{
    my ($self, $ggcode) = @_;
    my (@code, @num, $value, $address, $key, $hex_code, $warning);
    $warning = " ";

    @code = split //, $ggcode;
    @num = tonum(@code);

    $value = (($num[0]&8)<<4)+(($num[1]&7)<<4)+($num[0]&7);
    $address = (($num[3]&7)<<12)+(($num[4]&8)<<8)+(($num[5]&7)<<8)+
        (($num[1]&8)<<4)+(($num[2]&7)<<4)+($num[3]&8)+($num[4]&7);
    if (@num == 8)
    {
        $value+=($num[7]&8);
        $key = (($num[6]&8)<<4)+(($num[7]&7)<<4)+($num[5]&8)+($num[6]&7);
        #printf "Key:\t\t %x\n", $key;
    } else {
        $value+=($num[5]&8);
    }

    if ($num[2] >> 3)       # Check for eightletterness
    {
        # Codes like this don't automagically go to the next line when
        # the sixth letter is typed, although they still work. Make the
        # correct, six-letter code that does the same thing.
        my $right_code;
        $right_code = $self->encode($address, $value);

        $warning = "~" if @num != 8;        # Expecting eight letters
        #print "is eight letter code, okay\n";
    }


    return ($address, $value, $key, $warning)
        if wantarray();

    # Doesn't want array. Give hex code.
    $hex_code = sprintf "%.4X", $address;
    if (@num == 8)
    {
        $hex_code .= sprintf "?%.2X", $key;
    }
    $hex_code .= sprintf ":%.2X", $value;
    #warn "SCALAR CONTEXT";
    return "$warning$hex_code";
}

# Parse a GG patch such as 7333?03:00 into address, value, and key

=head2 ($address, $value, $key) = $gg->parse_gg($ggcode)

Parses a GG code in canonical format (that is, aaaa?kk:vv or aaaa:vv)
into it's relative address, value, and possibily key.

=cut

sub parse_gg
{
    my ($self, $code) = @_;
    my ($key);
    my ($address, $value) = $code =~ m/
        ^([A-Fa-f0-9?]*)        # Address (+optional key) portion
        :([A-Fa-f0-9]*)         # Byte value
    /x;

    if (!defined($address) || !defined($value))
    {
        carp "$code is in an invalid code\n";
        return undef;
    }

    if ($address =~ m/\?/)
    {
        ($address, $key) = split /\?/, $address;
    }

    
    $address = hex $address;
    $value   = hex $value;
    $key     = hex $key         unless !defined($key);

    return ($address, $value, $key);
}

=head2 $ggcode = $gg->encode($address, $value, $key)

Encode a relative address, value, and a possibily a key into a Game Genie
code. If $key is defined, the resulting code will be eight letters. Else it
will be six.

If $address is between 0x0000 and 0x7FFF, it is assumed to be a relative
address. Addresses over 0x7FFF are forced to be be under 0x7FFF if no ROM
is loaded. If a ROM, however, is loaded, $address is treated as an absolute
address starting at 0x8000. A key will automagically be discovered if
possible.

=cut
sub encode
{
    my ($self, $address, $value, $key) = @_;
    my (@code, @num);

    # If there is only one argument, it's in aaaa?kk:vv format
    if (@_ == 2)
    {
        ($address, $value, $key) = $self->parse_gg($_[1]);
    }

    if ($address > 0x7FFF)
    {
        # User entered something like 8000:00, but the maximum
        # address is 7FFF. Keys are required to access addresses in
        # bank 2 or higher.
        if (@{$self->{banks}})         # ROM loaded
        {
            # We can generate a key and reladdress by reading the ROM.
            # Note that the ROM starts at $8000 with these abs addresses,not 0
            $address -= 0x8000;
            # WARNING: THIS MIGHT NOT WORK -- DEPENDS ON MAPPERS
            my $bank_no = int($address / 0x4000);
               $address =     $address % 0x4000;       # Now relative
            $key = ord(substr($self->{banks}[$bank_no], $address, 1));
            printf "Found %.4X?%.2X\n", $address, $key;
        } else {
            # ROM not loaded, so just lop off high bits
            $address &= 0x7FFF;
            printf "warning: %.4X > 0x7FFF, forcing < 0x8000\n", $address;
            undef $key;
        }
    }

    $num[0]=($value&7)+(($value>>4)&8);
    $num[1]=(($value>>4)&7)+(($address>>4)&8);
    $num[2]=(($address>>4)&7);
    $num[3]=($address>>12)+($address&8);
    $num[4]=($address&7)+(($address>>8)&8);
    $num[5]=(($address>>8)&7);
  
    if (!defined($key))         # Length is six
    {
        $num[5]+=$value&8;
    } else {                    # Includes a key
        $num[2]+=8;
        $num[5]+=$key&8;
        $num[6]=($key&7)+(($key>>4)&8);
        $num[7]=(($key>>4)&7)+($value&8);
    }  
 
    @code = toletter(@num);
    return join("", @code);
}

=head2 $gg->ROM_loaded()

Return true if a ROM is loaded, meaning that special absolute address
functions are available. False if not.

=cut

# True if a ROM is loaded
sub ROM_loaded
{
    my ($self) = @_;
    return 1 if $self->{banks} && @{$self->{banks}};
}

=head2 ($prg_banks) = $gg->load_ROM($filename)

Load an NES ROM in iNES format, enabling special functions to be used.
Returns the number of 0x4000-byte PRG banks.

=cut

# Load a ROM for discovering and destroying keys
sub load_ROM
{
    my ($self, $filename) = @_;
    my ($ROM, $header, $prg_size);
    open($ROM, "<$filename") || die "cannot open $filename: $!";
    binmode($ROM);

    # Read 0x10-byte ROM header
    local $/;
    $/ = \0x10;
    $header = <$ROM>;
    $prg_size = ord(substr($header, 4, 1));
     
    # Read all 0x4000-byte PRG banks
    $/ = \0x4000;
    while(my $bank = <$ROM>)
    {
        push @{$self->{banks}}, $bank;
        printf "Read bank %.5X\n", length($bank);
        if (length($bank) != 0x4000)
        {
            warn "Partially read bank\n";
            warn sprintf "Next bank: %.5X\n", length(<$ROM>);
        }
    }
    close($ROM);
    return $prg_size;
}

# Discover which addresses are affected when applying a "key" to a rel addr.
#memoize('addrkey_to_addr');        # Won't work between loading new ROMs

=head2 (@affected) = $gg->addrkey_to_addr($reladdr, $key)

Only available is a ROM is loaded. Searches through all PRG banks of the
loaded ROM for the $key at $reladdr. Returns a list of affected addresses.
In scalar context, merely returns the first affected address.

Note that this code is not needed for six-letter keyless codes, as they
affect all the PRG banks.

=cut
sub addrkey_to_addr
{
    my ($self, $address, $key) = @_;
    my ($ROM, @absaddrs);

    return if !$self->ROM_loaded();
    return if !defined($key);       # If no key, all passes through

    $address &= 0x3FFF;         # Make address relative
    
    my $bank_addr = 0x8000;     # ROM loaded into memory at $8000   
                                # ^Subtract $8000 - $10 to get .NES address
    foreach my $bank (@{$self->{banks}})
    {
        if ($address > length($bank))
        {
            warn sprintf "%.4X > %.4X (key=$key, bank=%.5X)",
                $address, length($bank), $bank_addr;
            foreach my $bank (@{$self->{banks}})
            {
                printf "Bank: %.5X\t", length($bank);
            }
            die;
        }
        my $value = substr($bank, $address, 1);
        if (!defined($key))   # undef=no key, let all pass through
        {
            #printf "PUSHING UNDEF %.5X\n", $address+$bank_addr;
            push @absaddrs, $address + $bank_addr;
        } else {
            if ($value eq chr($key))
            {
                # Found affected address. Others may still be affected.
                my $abs_address = $address + $bank_addr;
                push @absaddrs, $abs_address;
            }
        }
        $bank_addr += 0x4000;
    }
    return wantarray() ? @absaddrs : $absaddrs[0];
}

=head2 $gg->looks_encoded($ggcode)

Return true if $ggcode is a decodable eight- or six-letter Game Genie code.

=cut

# True if $ggcode looks like a valid encoded GG code. Can be decoded.
sub looks_encoded
{
    my ($self, $ggcode) = @_;
    return 1 if $ggcode =~  #/ minimum of six letters, most of eight
        m/^[APZLGITYEOXUKSVN]{6}[APZLGITYEOXUKSVN]?[APZLGITYEOXUKSVN]?$/i;
}

# True if code is in the form a:b or a?c:b. Can be encoded.
=head2 $gg->looks_decoded($patch)

Return true if $patch is a encodable Game Genie code in canonical format.

=cut
sub looks_decoded
{
    my ($self, $patch) = @_;
    return 1 if $patch =~ m/:/;
}

=head2 $gg->unload_ROM()

Unload any loaded ROM.

=cut
sub unload_ROM
{
    my ($self) = @_;
    @{$self->{banks}} = ();
}

# Code or decode (OBSOLETE)
sub codec
{
    my $self = shift;
    ($_) = @_;
    my ($raw, $cooked, $ret);
    chomp($raw = $_);

    return if $raw eq '';
    $ret = "";

    eval
    {
        if ($raw =~ m<=(.*)>)                       # Set ROM filename
        {
            load_ROM($1);
            $ret .= "Loaded $1 -- key discovery/elimination enabled\n";
        } elsif ($raw eq '-') {                     # Unload ROM
            @{$self->{banks}} = ();
            $ret .= "ROM unloaded\n";
        }
        elsif ($raw =~ m/^[APZLGITYEOXUKSVN]*$/i)   # Game Genie-coded code
        {
            my ($reladdr, $value, $key) = decode($raw);
            if (!defined($key))
            {
                $ret .= sprintf "%.4X:%.2X\n", $reladdr, $value;
            } else {
                $ret .= sprintf "%.4X?%.2X:%.2X\n", $reladdr, $key, $value;
            }
            return $ret if !@{$self->{banks}};

            # If a ROM filename is given, we can eliminate the key
            my @absaddrs = addrkey_to_addr($reladdr, $key);
            foreach my $absaddr (@absaddrs)
            {
                $ret .= sprintf "%.5X:%.2X\n", $absaddr, $value;
            }
        } elsif ($raw =~ m/[:]/) {                  # Hex code
            my ($reladdr, $value, $key) = parse_gg($raw);
            $ret .= encode($reladdr, $value, $key) . "\n";
            return $ret if !@{$self->{banks}};
            
            # If a ROM filename is given, we can eliminate the key
            my @absaddrs = addrkey_to_addr($reladdr, $key);
            foreach my $absaddr (@absaddrs)
            {
                $ret = sprintf "%.5X:%.2X\n", $absaddr, $value;
            }
        } else {
            $ret = "Unrecognized code format '$raw'\n";
        }
    };
    $ret .= $@ if $@;
    return $ret;
}

=head1 EXAMPLE

gg.pl is provided as an example of how to use Convert::NES::GG. Usage is
simple:

=over

=item SLXPLOVS

Type a Game Genie code to decode it into address, key, value format.
If a ROM is loaded, affected addresses will be listed.

=item FFFF:BD

=item A23?DE:BD

Type a address/value pair to encode it as a Game Genie code. If the address
is <0x8000, treated as relative. If >0x8000 and a ROM is loaded, the ROM is
used to determine the key, else, high bits are lopped off.

If a ROM is loaded, affected addresses will be listed.

=item =filename.nes

Loads an iNES ROM for advanced features.

=item -

Unloads any currently loaded ROM.
              

=back


=head1 CREDITS

Encode/decode routines translated from Chris Covell's AmiGenie.cpp.

=head1 AUTHOR INFORMATION

Copyright (C) 2001, Jeff Connelly.

This library is free software; you can redistribute it and/or modify it
under the same terms as Perl itself.

Official web site: http://ggdb.sourceforge.net/

=cut

1;
