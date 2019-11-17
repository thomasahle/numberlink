import sys
import random
import collections
import itertools
import argparse

from mitm import Mitm
from grid import Grid
import draw


# Number of tries at adding loops to the grid before redrawing the side paths.
LOOP_TRIES = 1000


parser = argparse.ArgumentParser(description='Generate Numberlink Puzzles')
parser.add_argument('width', type=int, default=10,
                    help='Width of the puzzle')
parser.add_argument('height', type=int, default=10,
                    help='Height of the puzzle')
parser.add_argument('n', type=int, default=1,
                    help='Number of puzzles to generate')
parser.add_argument('--min', type=int, default=-1,
                    help='Minimum number of pairs')
parser.add_argument('--max', type=int, default=-1,
                    help='Minimum number of pairs')
parser.add_argument('--verbose', action='store_true',
                    help='Print progress information')
parser.add_argument('--solve', action='store_true',
                    help='Print solution as well as puzzle')
parser.add_argument('--zero', action='store_true',
                    help='Print puzzle in zero format')
parser.add_argument('--no-colors', action='store_true',
                    help='Print puzzles without colors')
parser.add_argument('--no-pipes', action='store_true',
                    help='When printing solutions, don\'t use pipes')
parser.add_argument('--terminal-only', action='store_true',
                    help='Don\'t show the puzzle in matplotlib')


def has_loops(grid, uf):
    """ Check whether the puzzle has loops not attached to an endpoint. """
    groups = len({uf.find((x, y)) for y in range(grid.h) for x in range(grid.w)})
    ends = sum(bool(grid[x, y] in 'v^<>') for y in range(grid.h) for x in range(grid.w))
    return ends != 2 * groups


def has_pair(tg, uf):
    """ Check for a pair of endpoints next to each other. """
    for y in range(tg.h):
        for x in range(tg.w):
            for dx, dy in ((1, 0), (0, 1)):
                x1, y1 = x + dx, y + dy
                if x1 < tg.w and y1 < tg.h:
                    if tg[x, y] == tg[x1, y1] == 'x' \
                            and uf.find( (x, y)) == uf.find( (x1, y1)):
                        return True
    return False


def has_tripple(tg, uf):
    """ Check whether a path has a point with three same-colored neighbours.
        This would mean a path is touching itself, which is generally not
        allowed in pseudo-unique puzzles.
        (Note, this also captures squares.) """
    for y in range(tg.h):
        for x in range(tg.w):
            r = uf.find( (x, y))
            nbs = 0
            for dx, dy in ((1, 0), (0, 1), (-1, 0), (0, -1)):
                x1, y1 = x + dx, y + dy
                if 0 <= x1 < tg.w and 0 <= y1 < tg.h and uf.find( (x1, y1)) == r:
                    nbs += 1
            if nbs >= 3:
                return True
    return False


def make(w, h, mitm, min_numbers=0, max_numbers=1000):
    """ Creates a grid of size  w x h  without any loops or squares.
        The mitm table should be genearted outside of make() to give
        the best performance.
        """

    def test_ready(grid):
        # Test if grid is ready to be returned.
        sg = grid.shrink()
        stg, uf = sg.make_tubes()
        numbers = list(stg.values()).count('x') // 2
        return min_numbers <= numbers <= max_numbers \
                and not has_loops(sg, uf) \
                and not has_pair(stg, uf) \
                and not has_tripple(stg, uf) \

    # Internally we work on a double size grid to handle crossings
    grid = Grid(2 * w + 1, 2 * h + 1)

    gtries = 0
    while True:
        # Previous tries may have drawn stuff on the grid
        grid.grid.clear()

        # Add left side path
        path = mitm.rand_path2(h, h, 0, -1)
        if not grid.test_path(path, 0, 0):
            continue
        grid.draw_path(path, 0, 0)
        # Draw_path doesn't know what to put in the first and last squares
        grid[0, 0], grid[0, 2 * h] = '\\', '/'

        # Add right side path
        path2 = mitm.rand_path2(h, h, 0, -1)
        if not grid.test_path(path2, 2 * w, 2 * h, 0, -1):
            continue
        grid.draw_path(path2, 2 * w, 2 * h, 0, -1)
        grid[2 * w, 0], grid[2 * w, 2 * h] = '/', '\\'

        # The puzzle might already be ready to return
        if test_ready(grid):
            return grid.shrink()

        # Add loops in the middle
        # Tube version of full grid, using for tracking orientations.
        # This doesn't make so much sense in terms of normal numberlink tubes.
        tg, _ = grid.make_tubes()
        # Maximum number of tries before retrying main loop
        for tries in range(LOOP_TRIES):
            x, y = 2 * random.randrange(w), 2 * random.randrange(h)

            # If the square square doen't have an orientation, it's a corner
            # or endpoint, so there's no point trying to add a loop there.
            if tg[x, y] not in '-|':
                continue

            path = mitm.rand_loop(clock=1 if tg[x, y] == '-' else -1)
            if grid.test_path(path, x, y):
                # A loop may not overlap with anything, and may even have
                # the right orientation, but if it 'traps' something inside it, that
                # might now have the wrong orientation.
                # Hence we clear the insides.
                grid.clear_path(path, x, y)

                # Add path and recompute orientations
                grid.draw_path(path, x, y, loop=True)
                tg, _ = grid.make_tubes()

                # Run tests to see if the puzzle is nice
                sg = grid.shrink()
                stg, uf = sg.make_tubes()
                numbers = list(stg.values()).count('x') // 2
                if numbers > max_numbers:
                    debug('Exceeded maximum number of number pairs.')
                    break
                if test_ready(grid):
                    debug(f'Finished in {tries} tries.')
                    debug(f'{numbers} numbers')
                    return sg

        debug(grid)
        debug(f'Gave up after {tries} tries')


def debug(s):
    try:
        if args.verbose:
            print(s, file=sys.stderr)
    except NameError:
        pass


def main():
    global args
    args = parser.parse_args()

    w, h = args.width, args.height
    if w < 4 or h < 4:
        print('Please choose width and height at least 4.')
        return

    n = int((w * h)**.5)
    min_numbers = n * 2 // 3 if args.min < 0 else args.min
    max_numbers = n * 3 // 2 if args.max < 0 else args.max

    debug('Preprocessing...')
    mitm = Mitm(lr_price=2, t_price=1)
    # Using a larger path length in mitm might increase puzzle complexity, but
    # 8 or 10 appears to be the sweet spot if we want small sizes like 4x4 to
    # work.
    mitm.prepare(min(20, max(h, 6)))
    debug('Generating puzzle...')

    for _ in range(args.n):
        grid = make(w, h, mitm, min_numbers, max_numbers)
        tube_grid, uf = grid.make_tubes()
        color_grid, mapping = draw.color_tubes(grid, no_colors=args.no_colors)

        # Print stuff
        debug(grid)

        print(w, h)
        if args.zero:
            # Print puzzle in 0 format
            for y in range(color_grid.h):
                for x in range(color_grid.w):
                    if grid[x, y] in 'v^<>':
                        print(color_grid[x, y], end=' ')
                    else:
                        print('0', end=' ')
                print()
        else:
            for y in range(color_grid.h):
                for x in range(color_grid.w):
                    if grid[x, y] in 'v^<>':
                        print(color_grid[x, y], end='')
                    else:
                        print('.', end='')
                print()
        print()

        if args.solve:
            print('Solution:')
            if not args.no_pipes:
                # Translate to proper pipe characters
                print(repr(color_grid).replace('-', '─').replace('|', '│'))
            else:
                for y in range(grid.h):
                    for x in range(grid.w):
                        print(mapping[uf.find( (x, y))], end='')
                    print()
            print()

        # Draw with pyplot
        if not args.terminal_only:
            draw.plot_puzzle(tube_grid, uf, include_solution=args.solve)


if __name__ == '__main__':
    main()
