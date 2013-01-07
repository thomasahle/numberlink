Numberlink
==========

Numberlink is a small, but very fast, program to solve puzzles of
Numberlink/Arukone/Nanbarinku. The puzzle involves finding paths to connect
numbers or letters in a grid.

See http://en.wikipedia.org/wiki/Numberlink for a detailed description.

Running it
----------

Download Numberlink from https://github.com/thomasahle/numberlink

Numberlink is written in the Go Programming Language and the binary can be
created by running

    go install numberlink

For more information on compiling, read the INSTALL file.

When you have created the binary, you can run `bin/numberlink [options]`.
Numberlink will then read puzzles from standard input in the following format:

    5 4
    C...B
    A.BA.
    ...C.
    .....

The first line consists of the width and height of the puzzle.
The following lines contains the puzzle where '.' represents an empty square and
[a-zA-Z0-9] are sources that must be connected.

Numberlink then prints the solved puzzle to standard input, either in the format
below, or as specified by command line flags:

    5 4
    CCBBB
    ACBAA
    ACCCA
    AAAAA

Also read the INSTALL file
See bin/numberlink --help

What Numberlink is not
----------------------

You can't use numberlink for checking if a puzzle is unique. Indeed numberlink
will only solve puzzles where the solution use 100% of the paper and no path
touches itself.

If you want to find the number of solution to a general numberlink puzzle I
suggest using this solver by ~imos: https://github.com/imos/Puzzle/tree/master/NumberLink

How it works
------------

Numberlink solves puzzles using a heavily pruned backtracking search on an
optimied datastructure.

In particular the following heuristics are used:

* Partial paths
* Corner heuristic
* Optimistic (late validation)

There are multiple ways to do backtracking on numberlink puzzles.
The most obvious is to start at a source, choose a path to its other end and
recurse. Alternatively one can start at all sources at the same time, or
systematically fill out the squares on the paper in some predefined order.

Numberlink uses the later aproach: It fills out the paper starting in the upper
left corner, and continues to the SW facing diagonals. For a 4x4 paper it will
look like this (base 16):

    0136
    247a
    58bd
    9cef

Backtracking in this systematic give us a lot of advantages compared to starting
at the sources:

* We never get empty squares
* We never block a source from its other end
* We always know exactly what squares around us have already been connected

The challange with this approach is that we need to manage 'partial paths' that
aren't yet connected to anything. We could do this by disjoint-set, but it is
simpler to just keep an array such that if pos is a start of a path then
end[pos] is the position of the other end. This is easily updated when two paths
are merged, and can be used to ensure different sources aren't connected.

The last question one may ask is 'why diagonally?' Instead one could have done
row by row, or with an expanding boundary like a bfs search. While the later
approach may allow us to fill out some obvious squares higher up in the tree, it
doesn't give us much predictability in the structure of the filled out squares,
something we'll need for the 'corner heuristic'. Filling by rows is very similar
to diagonals, but with diagonals the tree is heigher.

The corner heuristic is the most important part of what makes Numberlink fast.
It relies on the obvervation that if a square is filled out with a ┐ the only
option for the lower left square is another ┐ or a source. Anything else will
either force a self touching path or connect into the side of the paper.

Taking the inductive closure of the above obvervation we see that all
path-turns, or 'corners', must be found in 'spikes' 'coming out' of the sources.
Indeed a source can't even have such a spike in two opposite direcitons, as it
would create a flow surounding the source. In conclusion we see that a solution
to a numberlink puzzle can be represented uniquely as a set of signed integer
pairs for each source, describing the length of its two spikes.

We don't directly use the above representation, as it doesn't seem to suggest an
easy way to backtrack. Instead we make sure that no connections are made, which
would create an illegal situation in the dual representation. It is worth
noticing that the dual representation means especially very sparse puzzles can
be efficiently solved.

The corner heuristic also protects us from a lot of illegal solutions, like
connecting a path head to the middle of another path. It doesn't however quite
save us from self touching paths, as this example shows:

    4 4
    ....
    .ab.
    ..b.
    a...

Numberlinks approach to such situations is to assume that these situations won't
happen very often, and hence we don't need to check for the case during solving.
Only once the whole paper is filled out do we check if we have done something
illegal.

History
-------

Numberlink was written by Thomas Dybdahl Ahle for a competition at Oxford
University arranged by Michael Spivey (spivey.oriel.ox.ac.uk). The description
of the competition is still available at
http://spivey.oriel.ox.ac.uk/wiki/index.php/Programming_competition_2012

Legal
-----

Numberlink is released under the GPL3
Read LICENSE for more details
