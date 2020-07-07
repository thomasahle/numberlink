Numberlink
==========

Numberlink is a very fast program for generating and solving large instances of the puzzle [Numberlink / Arukone / Nanbarinku / Flow Free](https://en.wikipedia.org/wiki/Numberlink). Instances of size 40x40 are casually solved and generated, and the complete set of [270 puzzles at Janko.at](https://www.janko.at/Raetsel/Arukone-3/index.htm)
is solved in 0.14 seconds(!).

Example generated (and solved) puzzle:
![Screenshot](https://raw.githubusercontent.com/thomasahle/numberlink/master/images/example_0.png)

Generating Puzzles
------------------

The python code in `gen/gen.py` can be used to generate very large puzzles. To create such puzzles, simply run the script with your intended width and height.
(Notice that puzzles with heights above 20 may take a while to generate.)

In the example below, we generate a puzzle with width=40 and height=10:

    $ python gen/gen.py 40 10
    40 10
    1.......................................
    ......................6.................
    ....24.4................8...............
    .......5....................5..93.......
    ..............................7....9....
    ....3...................................
    ....................................8...
    ...................627..................
    ........................................
    ......................1.................

By default it tries to put around `sqrt(width * height)` number pairs in a puzzle, but this can be controlled with the `--min` and `--max flags`.
See `gen.py --help` for more options.

Numberlink Solver
----------

The Numberlink solver is written in the Go Programming Language and is compiled using
`$ go install numberlink`. This won't install anything on your system. For more
information on compiling, see the INSTALL file.

When you have created the binary, you can run `$ bin/numberlink [options]`.
Numberlink will then read puzzles from standard input in the following format:

    5 4
    C...B
    A.BA.
    ...C.
    .....

The first line consists of the width and height of the puzzle.
The following lines contain the puzzle where `.` represents an empty square and
letters or digits are the sources that must be connected.

Numberlink then prints the solved puzzle to standard output, either in the
format below, or as specified by command-line flags:

    5 4
    CCBBB
    ACBAA
    ACCCA
    AAAAA

If the puzzle wasn't solvable, `IMPOSSIBLE` is printed.

To learn about the available command-line flags, see `$ bin/numberlink --help`. 

Old Generator
-------------

Examples of generating a puzzle and immediately solving it:

    $ bin/numberlink -generate=10x10 | bin/numberlink -tubes
    0─────┐120
    3────┐│12│
    ┌──45│└──┘
    │┌──┘└───3
    45┌──────6
    ┌─┘┌─┐┌──8
    │7─┘9││┌─a
    └──6││││8┐
    ┌───┘││└a│
    97───┘└──┘

and larger:

    $ bin/numberlink -generate=20x20 | bin/numberlink -tubes
    0────┐┌─┐223┌─4┌───┐
    6───┐│17│┌─┘48┌┘9ab│
    c──c││┌┘││de┌┘│┌┘a││
    ┌──f││71┘3de8┌┘│5b┘│
    │g─┐6└───┐┌──┘9┘└──┘
    └─┐└───gh││ij┐┌────┐
    kk└f┌─┐┌┘││└┐││┌──5i
    ┌─┐l┘nlho│└┐│j││┌──┐
    │q│┌r└┐┌┘│s│└─┘││┌t│
    m│││uu│o┌┘│└─┐v│││p┘
    ┌┘│└─r│┌┘┌┘ww│v││└─t
    │x│┌─┐n│y│zA┐└─┘└──┐
    qxm│B└─┘││└┐└─┐CD─┐│
    ┌─┐│└┐┌─┘└┐└─┐│└─C│p
    │0│└┐B│E─┐└─┐││┌─┐└D
    ││└y│┌┘┌s└─E│z││G│HH
    │└──┘│I│JJ┌─┘A┘│G└─F
    │K──K│││LL│M──┐└┐┌─┐
    └────┘│└──┘┌─O│P│NF│
    I─────┘O───┘M─┘P└─┘N

See https://stackoverflow.com/a/14007585/205521 for an explanation of the algorithm.

What Numberlink is not
----------------------

You can't use numberlink for checking if a puzzle is unique. Indeed numberlink
will assume the puzzle has just one solution and only it such that the solution
uses 100% of the paper and no link touches itself. Hence some non-unique puzzles
will be solved, while others will be `IMPOSSIBLE`.

If you want to find the number of solution to a general numberlink puzzle, I
suggest using this solver by ~imos: https://github.com/imos/Puzzle/tree/master/NumberLink

How it works
------------

Numberlink solves puzzles using a heavily pruned backtracking search. In
particular the following pruning heuristics are used:

* Partial links
* A dual representation based on link corners
* Optimistic validation

There are multiple ways to do backtracking on numberlink puzzles. The most
obvious is to start at a source, choose a link to its other end and recurse.
Alternatively one can start at all sources at the same time, or you can ignore
the sources and systematically fill out the squares on the paper in some order.

Numberlink uses the later approach: It fills out the paper starting in the upper
left corner and continuing along the SW-diagonals. For a 4x4 paper the order in
which squares are visited is (in base 16):

    0136
    247a
    58bd
    9cef

Backtracking in this systematic give us a lot of advantages compared to starting
at the sources:

* We never get unconnected squares. We simply always connect a square as we go
  over it.
* We never block a source from its other end. To see this notice that blocking
  a source requires us to have passed it. In passing it we must have connected
  it to a partial link. The other end of this partial link can not be connected
  to any other sources or to the side, so it must be in the 'active diagonal'.
* We always know exactly what squares around us have already been connected.
  That's the one above us and the one to the left. The directions we need to
  care about is down and right.

The challenge with this approach is that we need to manage 'partial links' that
aren't yet connected to anything. We don't want to accidentally connect a link
to itself, or to connect the ends to different labelled sources.
This problem can be solved efficiently by the disjoint-set data structure, but
it is simpler for us to just keep an array such that if `pos` is a 'link head'
then `end[pos]` is the current position of the other end. Initially
`end[pos]=pos` and if `pos` has degree two, `end[pos]=-1`. (Actually the last
part is unnecessary due to the systematic order of connection). This array is
easily updated when two links are merged by no more than two array assignments.

The corner dual heuristic is the most important part of what makes Numberlink
fast. It relies on the observation that if a square is filled out with a ┐ (a
south turn of a link, we'll call it a SW 'corner') the lower left square will
either have to be a source or to be a SW corner as well. Anything else will
force a self-touching link.

Taking the inductive closure of the above observation, we see that all
corners must be found in 'spikes' rooted at the sources. Indeed a source can't
even have such a spike in two opposite directions, as it would create a link
surrounding the source. All in all, we conclude that any solution to a numberlink
puzzle can be represented uniquely as a set of signed integer pairs, one pair
for each source, describing the length of its two spikes.

We don't directly use the above representation however, as it doesn't seem to
suggest an easy way to backtrack. Instead, we backtrack on the partial link
representation, but make sure that no connections are made, which would create
an illegal situation in the dual representation. It is worth noticing that the
dual representation means especially very sparse puzzles can be efficiently
solved.

The corner heuristic also protects us from a lot of illegal states in the
primary representation, for example self touching links are very rarely
explored. It isn't however totally safe to rely on, as this example shows:

    4 4
    ....
    .ab.
    ..b.
    a...

Numberlinks approach to this kind of situation is to assume that they won't
happen very often, and hence we don't need to check for them during solving.
This makes searching more efficient and only once the entire paper is filled out
do we check if we have done something illegal.

It's still unexplored however exactly how much extra pruning we could make, if
we detected self-touch early on. Another kind of self touch we might be able to
prune early is this one:

    ─y┌y
    ──┘z

The last question one may ask is 'why search diagonally?' Instead, one could have
walked row by row, or with an expanding boundary like a bfs search. While the
later approach may allow us to fill out some obvious squares higher up in the
tree, it doesn't give us much predictability in the structure of the filled out
squares, something that simplifies the search code greatly. Filling by rows is
very similar to diagonals, but with diagonals the tree is often twice as high.

History
-------

Numberlink was written by Thomas Dybdahl Ahle for [a competition at Oxford
University](https://www.cs.ox.ac.uk/news/622-full.html) arranged by [Michael Spivey](http://spivey.oriel.ox.ac.uk). The
description of the competition is available at
[spivey.oriel.ox.ac.uk/wiki/index.php/Programming_competition_2012](https://web.archive.org/web/20160416150900/http://spivey.oriel.ox.ac.uk/corner/Programming_competition_2012)

Legal
-----

Numberlink is released under the Affero GPL3.

Read LICENSE for more details.

If you generate puzzles using the scripts and post them in books or use them on your website, put clear reference attrtibuting them to this author.
