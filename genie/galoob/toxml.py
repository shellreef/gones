#!/usr/bin/python
# Created:20110113
# By Jeff Connelly

import csv
import xml.dom.minidom

# Mapping of what Galoob calls a game, to its abbreviation (NOT unique), and GoodNES name
game2id = {}
game2gn = {}
for row in csv.reader(file("gamelist-galoob.csv", "rb"), delimiter="\t"):
    galoob, id, goodnes = row
    game2id[galoob] = id
    game2gn[galoob] = goodnes

# Read comprehensive code file, parsed
game_lines = {}
for row in csv.reader(file("all-nev.csv", "rb"), delimiter="\t"):
    game, id, type = row[0:3]
    rest = row[3:]

    if not game_lines.has_key(game):
        game_lines[game] = []

    game_lines[game].append((type, rest))

# Write
doc = xml.dom.minidom.Document()
root = doc.createElement("cheats")
for game, lines in game_lines.iteritems():
    game_node = doc.createElement("game")
    game_node.setAttribute("galoob-name", game)
    game_node.setAttribute("galoob-id", game2id[game])
    game_node.setAttribute("fullname", game2gn[game])
    for line in lines:
        type, rest = line

        if type == "code":
            source, no, code, title = rest
        
            code_node = doc.createElement("codeline")
            code_node.setAttribute("source", source)
            code_node.setAttribute("number", no)
            code_node.setAttribute("code", code)
            code_node.setAttribute("title", title)

            game_node.appendChild(code_node)

        else:
            text = "\t".join(rest)

            info_node = doc.createElement("info")
            info_node.appendChild(doc.createTextNode(text))
            
            game_node.appendChild(info_node)

    root.appendChild(game_node)
doc.appendChild(root)

print doc.toprettyxml(indent=" ")
