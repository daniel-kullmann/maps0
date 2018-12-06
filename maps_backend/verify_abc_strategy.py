#!/usr/bin/env/python3

# For some reason (load, probably), the OSM tile cache uses three different
# servers for servng the tiles. I wanted to know what strategy is used to
# decide which server is used.
#
# My guess was that for a request for
#
# https://c.tile.openstreetmap.org/13/3895/3178.png
#
# server c is used because (3895+3178) modulo 3 == 2. Those two last parts are
# the x and y positions of the tile, so:
# (x+y) % 3 == 0  => Use server a
# (x+y) % 3 == 1  => Use server b
# (x+y) % 3 == 2  => Use server c
#
# This piece of code verifies this guess. You need a cache with some tiles
# under tile_cache/ to make this work.
#

import os

for (dirpath, dirnames, filenames) in os.walk("./tile_cache/"):
    parts = dirpath.split("/")[2:]
    for filename in filenames:
        if filename.endswith(".png"):
            filename = filename[0:len(filename)-4]
        else:
            raise Exception("expected a .png file: " + dirpath + " " + filename)
        parts2 = parts + [filename]
        parts2 =  [parts2[0]] + [int(e) for e in parts2[1:]]
        parts2.append(sum(parts2[2:]) % 3)
        if parts2[0] == 'a':
            if parts2[4] != 0: print(parts2)
        elif parts2[0] == 'b':
            if parts2[4] != 1: print(parts2)
        elif parts2[0] == 'c':
            if parts2[4] != 2: print(parts2)
        else:
            print(parts2)


