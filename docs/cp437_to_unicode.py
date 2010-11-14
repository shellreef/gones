#!/usr/bin/python
# Created:20101113
# By Jeff Connelly
#
# Convert CP437 (DOS style) files to UTF-8 Unicode
# This is useful for older documentation that uses line drawing characters,
# which can be viewed in DOS, but Unicode is supported in more places (like on a Mac)
# For example: http://nesdev.parodius.com/2A03%20technical%20reference.txt

import sys

if len(sys.argv) < 2:
    print "Usage: %s filename"  % (sys.argv[0],)
    raise SystemExit

bytes = file(sys.argv[1], "rb").read()
file("utf8-" + sys.argv[1], "wb").write(unicode(bytes, "cp437").encode("utf8"))
