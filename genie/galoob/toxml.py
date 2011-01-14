#!/usr/bin/python
# Created:20110113
# By Jeff Connelly

import csv
import xml

# Mapping of what Galoob calls a game, to its abbreviation (NOT unique), and GoodNES name
game2id = {}
game2gn = {}
for row in csv.reader(file("gamelist-galoob.csv", "rb"), delimiter="\t"):
    galoob, id, goodnes = row
    game2id[galoob] = id
    game2gn[galoob] = goodnes

for row in csv.reader(file("all-nev.csv", "rb"), delimiter="\t"):
    game, id, type = row[0:3]
    rest = row[3:]

    if type == "code":
        source, no, code, title = rest
    else:
        if len(rest) != 1:
            print rest


