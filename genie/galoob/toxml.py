#!/usr/bin/python
# Created:20110113
# By Jeff Connelly

import csv
import xml.dom.minidom

import os
import hashlib

ROM_ROOT = "../../roms/best"

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
    for filename in filenames:
        # Calculate identifying hash. This is technically defined as hash(PRG+CHR),
        # but to avoid parsing the whole header, we just skip it and assume PRG+CHR
        # directly follows. Though the format does allow a "trainer" header, and a
        # hintscreen trailer, almost no files have either, so this usually works.
        f = file(os.path.join(dir, filename), "rb")
        if not filename.endswith(".nes"):
            # UNIF could be supported, but would require parsing out PRGx+CHRx
            print "can only analyze iNes, sorry", filename
            raise SystemExit
        f.read(0x10)
        sha1 = hashlib.sha1(f.read()).hexdigest().upper()
        print sha1, filename
        # Quick and dirty check the hash is valid
        ret = os.system("grep %s ~/games/nese/gones/data/cartdb.xml" % (sha1,))
        if ret != 0:
            print "Failed to match %s in %s" % (sha1, filename)



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
    game_node.setAttribute("fullname", game2gn[game])

    if game_intro.has_key(game):
        intro_node = doc.createElement("intro")
        intro_node.appendChild(doc.createTextNode("\n".join(game_intro[game])))
        game_node.appendChild(intro_node)

    for line in lines:
        type, rest = line

        if type == "code":
            source, no, codetext, title = rest
        
            effect_node = doc.createElement("effect")
            effect_node.setAttribute("source", source)
            effect_node.setAttribute("number", no)
            effect_node.setAttribute("title", title)

            codes_node = doc.createElement("codes")
            # One effect can have multiple Game Genie codes you need to enter
            # TODO: need to do something about "alternate" codes (" / ")
            codes = codetext.split(" + ")
            for code in codes:
                code_node = doc.createElement("code")
                code_node.setAttribute("genie", code)

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
