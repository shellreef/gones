#!/usr/bin/python
# Created:20110113
# By Jeff Connelly

import csv
import xml.dom.minidom

import os
import hashlib
import sys

ROM_ROOT = "../../roms/best"

# Map alternate codes to their cartridge versions
# (Surprisingly, these are the only 4 games with alternate codes)
ALT2VARIANT = {
    # Name of game as Galoob calls it: [first version, second version] # code tested with to determine
    "Micro Machines(tm) The Official Video Game": ["NES Cart", "Aladdin Cart"], # AOKNIYAE / APXNIYAE Start on race 25 (Final Race!)
    "RC Pro Am(tm) Game": ["PRG0", "PRG1"], # AEXEPPZA / AAUAGZZA No continues (I think)
    "The Simpsons(tm): Bart(tm) vs. The Space Mutants Game": ["PRG0", "PRG1"], # IPKYXUGA / IPUYVUGA Super-jumping Bart
    "Super Mario Bros.(tm) 2 Game": ["PRG1", "PRG0"], # AAVENYZA / AAVEUYZA Weak Birdetta (level 1 boss kill with 1 egg hit instead of 3)
    "Ultimate Stuntman(tm) Game": ["Aladdin Cart (unreleased)", "NES Cart"], # SZEIPUVK / SXNSYXVK Infinite time
}

# ROM variants that are not PRG0 or PRG1
VARIANTS = {
"58011155EA6D62AF65C8A9D776DD8363FF40EDDD": "REV0",
"F18C7D82F624C6EFEB73E0A212997A28B90BDF85": "REVA",
"0EE24D0A864845449EF7434822C12FA0F063DE56": "REV1.1",
"12F58963DD70D32CCA2784BC1D83EC889FEC8E5D": "REV1.x",
"D1F279C5EBBB9069887CC6F2534A362286801CD1": "REV1.x [a1]",
"36EC0A750888DB2BAAA21651528807D70CA97C6B": "Taito",
"DD7B6084032EDCE204862153FD32039E32DA884C": "UBI Soft",
"C7FD43041FC139DC8440C95C28A0115DC79E2691": "Aladdin Cart",
"84908DC67C29BE8600184FC5525B9227C8AFF830": "NES Cart",
"92C3361B9E3B28A51FD30E7845C988A6D576EE65": "Namco",
"A34E68372082513209A795786C8EEA493CC2CD14": "Tengen",
"06990C8573128E5548C5DCD39479FABF67234926": "Aladdin Cart",
"6A6C235B96C5CC51A5BF6D6FBAF30E77AD789FC7": "NES Cart",
"0C4992FC08D2278697339D3B48066E7B5F943598": "REVA",
"DB295C6BAD1B58BC1170C4B300C1C8D2A6BC1A87": "REVB",
"C87B3E1F17670C028CE60AF3BBC7D688DC0F9DF3": "REV0",
"712983EAA00029C307688DE015C1B698CC4BF064": "REVA",
"102BD0C46C5718C979EB1AC387DADE6F6EB70EE4": "Aladdin Cart",
"90196DBFC5337B56106B33891C5FA4B2267F3732": "NES Cart",
"42F15207D202B43802E92AF1F89300CEB9C99F12": "REV0",
"FCE0C7B0A152DBC3B5992320211CC674E8A1622B": "REV1",
"ED281797EFF64CBA96897B59D85AE5E61F67353A": "REVB",
"48100033895E83877F554AB539CB028ACAAC44AC": "Family Edition",
"5BCF47901533372B7D9828380FCF32F11C6F9CE8": "Junior Edition",
}

# Mapping of what Galoob calls a game, to its abbreviation (NOT unique), and GoodNES name
game2id = {}
game2gn = {}
game2carts = {}
for row in csv.reader(file("gamelist-galoob.csv", "rb"), delimiter="\t"):
    galoob, id, goodnes = row
    game2id[galoob] = id
    game2gn[galoob] = goodnes

    # Cartridges (can have >1 per game, different variants)
    dir = os.path.join(ROM_ROOT, goodnes)
    filenames = os.listdir(dir)
    carts = []
    assert len(filenames) > 0, "nothing in %s" % (dir,)
    for filename in filenames:
        # Calculate identifying hash. This is technically defined as hash(PRG+CHR),
        # but to avoid parsing the whole header, we just skip it and assume PRG+CHR
        # directly follows. Though the format does allow a "trainer" header, and a
        # hintscreen trailer, almost no files have either, so this usually works.
        f = file(os.path.join(dir, filename), "rb")
        # UNIF could be supported, but would require parsing out PRGx+CHRx
        assert filename.endswith(".nes"), "can only analyze iNes, sorry %s" % (filename,)
        f.read(0x10)
        hash = hashlib.sha1(f.read()).hexdigest().upper()

        if len(filenames) == 1:
            variant = None  # no variant name needed
        else:
            if "PRG0" in filename:
                variant = "PRG0"
            elif "PRG1" in filename:
                variant = "PRG1"
            else:
                variant = VARIANTS[hash]

        carts.append((hash, filename, variant))

    game2carts[galoob] = carts

# Read comprehensive code file, parsed
rows = []
for row in csv.reader(file("all-nev.csv", "rb"), delimiter="\t"):
    rows.append(row)

# Build structure
i = 0
game_lines = {}
game_intro = {}
game_order = []
while i < len(rows):
    row = rows[i]
    game, id, type = row[0:3]
    rest = row[3:]

    if not game_lines.has_key(game):
        game_lines[game] = []
        game_order.append(game)

    if type == "intro":
        if not game_intro.has_key(game):
            game_intro[game] = []
        game_intro[game].append("\t".join(rest))
    elif type == "info":
        # Suck up multiple info lines into one
        info_lines = []
        while rows[i][2] == "info":
            rest = rows[i][3:]
            info_lines.append("\t".join(rest))
            i += 1
            if i >= len(rows):
                break
        if i < len(rows): i -= 1    # Spit out unintended line
        game_lines[game].append(("info", ["\n".join(info_lines)]))
    elif type == "code":
        game_lines[game].append((type, rest))
    else:
        print "unknown type: ", type
        raise SystemExit

    i += 1

# Read game variants


# Write
doc = xml.dom.minidom.Document()
root = doc.createElement("cheats")
for game in game_order:
    lines = game_lines[game]
    game_node = doc.createElement("game")
    game_node.setAttribute("galoob-name", game)
    game_node.setAttribute("galoob-id", game2id[game])
    game_node.setAttribute("name", game2gn[game])

    # Cartridge info
    for hash, filename, variant in game2carts[game]:
        cart_node = doc.createElement("cartridge")
        cart_node.setAttribute("sha1", hash)
        cart_node.setAttribute("filename", filename)
        if variant is not None:
            cart_node.setAttribute("name", variant)
        game_node.appendChild(cart_node)

    # Intro text
    if game_intro.has_key(game):
        intro_node = doc.createElement("intro")
        intro_node.appendChild(doc.createTextNode("\n".join(game_intro[game])))
        game_node.appendChild(intro_node)

    # Code and info lines
    for line in lines:
        type, rest = line

        if type == "code":
            source, no, codetext, title = rest
        
            effect_node = doc.createElement("effect")
            effect_node.setAttribute("source", source)
            effect_node.setAttribute("number", no)
            effect_node.setAttribute("title", title)

            codes_node = doc.createElement("codes")

            # Alternate game cartridges can have different codes
            alt_texts = codetext.split(" / ")
            for alt_index, alt_text in enumerate(alt_texts):
                # One effect can have multiple Game Genie codes you need to enter
                codes = alt_text.split(" + ")
                for code in codes:
                    code_node = doc.createElement("code")
                    code_node.setAttribute("genie", code)
                    if len(alt_texts) > 1:
                        if alt_index >= len(game2carts[game]):
                            # Galoob says there are two versions of Ultimate Stuntman, but http://bootgod.dyndns.org:7777/profile.php?id=354 says there is only one
                            # 'The Ultimate Stuntman (Aladdin Cart)' says the Aladdin version was never released, so I think that's it.
                            # Codemasters made both Game Genie and this game, so they could've made codes for a pre-release Aladdin game http://www.nesworld.com/codemasters.php
                            if game == "Ultimate Stuntman(tm) Game":
                                variant = "Aladdin Cart (unreleased)"
                            else:
                                variant = "unknown-%s" % (alt_index,)
                                assert False, "Warning: %s has multiple variants (%s), but only %s found\n" % (game, alt_index, game2carts[game])
                        else:
                            assert ALT2VARIANT.has_key(game), "Missing ALT2VARIANT for game %s -aka- %s" % (game, game2gn[game])
                                
                            variant = ALT2VARIANT[game][alt_index]
                            found = False
                            for that_hash, that_filename, that_variant in game2carts[game]:
                                if that_variant == variant:
                                    found = True
                                    break
                            assert found or "unreleased" in variant, "For game %s, ALT2VARIANT specifies invalid variant: %s, known %s" % (game, variant, game2carts[game])

                            if variant is None:
                                variant = "unknown"
                        code_node.setAttribute("applies", variant)

                    codes_node.appendChild(code_node) 

            effect_node.appendChild(codes_node)

            game_node.appendChild(effect_node)

        else:
            text = "\t".join(rest)

            info_node = doc.createElement("info")
            info_node.appendChild(doc.createTextNode(text))
            
            game_node.appendChild(info_node)

    root.appendChild(game_node)
doc.appendChild(root)

print doc.toprettyxml(indent=" ")
