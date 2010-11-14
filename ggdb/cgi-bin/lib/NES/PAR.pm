#!/usr/bin/perl -w
# Created:04-05-01
# By Jeff Connelly

# NES PAR code module

package NES::PAR;

=head1 NAME

NES::PAR - Manipulate NES Pro Action Replay codes

=head1 SYNOPSIS

    use NES::PAR;

    my $par = new NES::PAR;
    my ($type, $address, $value) = $par->decode("0007EE01");

    my ($code) = $par->encode(0, 0x07ee, 0x01);

=head1 DESCRIPTION

NES::PAR provides an object-oriented interface for manipulating NES Pro
Action Replay codes. Because PAR codes are not encrypted in any way, this
library is very small.

=cut

use strict;
#use warnings;

sub new
{
    my ($self, $game) = @_;

    $self = {};
    bless($self);
    return $self;
}

sub decode
{
    my ($self, $code) = @_;
    if (@_ == 1) { $code = $_[0] }
    my ($type, $address, $value) =
        $code =~ m/^(..)(....)(..)$/;       # todo: use unpack
    $type = hex $type;
    $address = hex $address;
    $value = hex $value;

    return ($type, $address, $value);
}

sub encode
{
    my ($self, $type, $addr, $value) = @_;
    $type &= 0xff;
    $value &= 0xff;
    $addr &= 0xffff;
    return sprintf "%.2X%.4X%.2X", $type, $addr, $value;
}

1;

