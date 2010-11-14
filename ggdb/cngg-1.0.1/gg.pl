use Convert::NES::GG;

$gg = new Convert::NES::GG;
while(<>)
{
    chomp;
    my ($reladdr, $value, $key);
    if (m/=(.*)/)                       # Load ROM
    { 
        if (my $prg_size = $gg->load_ROM($1))       
        {
            printf "Loaded $1 (%d banks, max %.5X) -- key discovery/elimination enabled\n",
                    $prg_size, $prg_size * 0x4000;    
        } else {
            print "Couldn't load $1\n";
        }
    } elsif ($_ eq '-') {               # Unload ROM
        $gg->unload_ROM();
    } elsif ($gg->looks_decoded($_)) {  # 0000:00
        ($reladdr, $value, $key) = $gg->parse_gg($_);
        print scalar $gg->encode($reladdr, $value, $key), "\n";
    } elsif ($gg->looks_encoded($_)) {      # AAAAAA
        ($reladdr, $value, $key) = $gg->decode($_);
        if (!defined($key))
        {
            printf "%.4X:%.2X\n", $reladdr, $value;
        } else {
            printf "%.4X?%.2X:%.2X\n", $reladdr, $key, $value;
        }
    } 
        

    if (defined($reladdr) && $gg->ROM_loaded())
    {
        my @absaddr = $gg->addrkey_to_addr($reladdr, $key);
        foreach my $absaddr (@absaddr)
        {
            printf "%.5X:%.2X\n", $absaddr, $value;
        }
    }
}
