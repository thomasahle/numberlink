import sys
import random
import string
import collections
import itertools
import argparse

from mitm import Mitm


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
parser.add_argument('--terminal-only', action='store_true',
                    help='Don\'t show the puzzle in matplotlib')


def sign(x):
    if x == 0: return x
    return -1 if x < 0 else 1


def union(uf, a, b):
    a_par, b_par = find(uf, a), find(uf, b)
    uf[a_par] = b_par


def find(uf, a):
    if uf.get(a, a) == a:
        return a
    par = find(uf, uf.get(a,a))
    uf[a] = par # Path compression
    return par


class Grid:
    def __init__(self, w, h):
        self.w, self.h = w, h
        self.grid = {}
    def __setitem__(self, key, val):
        self.grid[key] = val
    def __getitem__(self, key):
        return self.grid.get(key, ' ')
    def __repr__(self):
        res = []
        for y in range(self.h):
            res.append(''.join(self[x,y] for x in range(self.w)))
        return '\n'.join(res)
    def __iter__(self):
        return iter(self.grid.items())
    def __contains__(self, key):
        return key in self.grid
    def __delitem__(self, key):
        del self.grid[key]
    def add(self, other):
        if not other:
            return False
        assert self.w == other.w and self.h == other.h
        self.grid = {**self.grid, **other.grid}
        return True


def test_path(grid, path, x0, y0, dx0=0, dy0=1):
    """ Test whether the path is safe to draw on the grid, starting at x0, y0 """
    return all(0 <= x0-x+y < grid.w and 0 <= y0+x+y < grid.h
               and (x0-x+y, y0+x+y) not in grid for x, y in path.xys(dx0, dy0))


def draw_path(grid, path, x0, y0, dx0=0, dy0=1, loop=False):
    """ Draws path on the grid. Asserts this is safe (no overlaps).
        For non-loops, the first and the last character is not drawn,
        as we don't know what shape they should have. """
    ps = list(path.xys(dx0,dy0))
    # For loops, add the second character, so we get all rotational tripples:
    # abcda  ->  abcdab  ->  abc, bcd, cda, dab
    if loop:
        assert ps[0] == ps[-1], (path, ps)
        ps.append(ps[1])
    for i in range(1, len(ps)-1):
        xp, yp = ps[i-1]
        x, y = ps[i]
        xn, yn = ps[i+1]
        grid[x0-x+y,y0+x+y] = {
            (1,1,1): '<', (-1,-1,-1): '<',
            (1,1,-1): '>', (-1,-1,1): '>',
            (-1,1,1): 'v', (1,-1,-1): 'v',
            (-1,1,-1): '^', (1,-1,1): '^',
            (0,2,0): '\\', (0,-2,0): '\\',
            (2,0,0): '/', (-2,0,0): '/'
        }[xn-xp, yn-yp, sign((x-xp)*(yn-y)-(xn-x)*(y-yp))]


def shrink_grid(grid):
    small_grid = Grid(grid.w//2, grid.h//2)
    for y in range(grid.h//2):
        for x in range(grid.w//2):
            small_grid[x,y] = grid[2*x+1, 2*y+1]
    return small_grid


def make_tubes(grid):
    uf = {}
    tube_grid = Grid(grid.w, grid.h)
    for x in range(grid.w):
        d = '-'
        for y in range(grid.h):
            # We union things down and to the right.
            # This means ┌ gets to union twice.
            for dx, dy in {
                    '/-':[(0,1)], '\\-':[(1,0),(0,1)],
                    '/|':[(1,0)],
                    ' -':[(1,0)], ' |':[(0,1)],
                    'v|':[(0,1)], '>|':[(1,0)],
                    'v-':[(0,1)], '>-':[(1,0)],
                    }.get(grid[x,y]+d, []):
                union(uf, (x,y), (x+dx,y+dy))
            # We change alll <>v^ to x.
            tube_grid[x,y] = {
                    '/-':'┐', '\\-':'┌',
                    '/|':'└', '\\|':'┘',
                    ' -':'-', ' |':'|',
            }.get(grid[x,y]+d, 'x')
            # We change direction on v and ^, but not on < and >.
            if grid[x,y] in '\\/v^':
                d = '|' if d == '-' else '-'
    return tube_grid, uf


def color_tubes(grid, no_colors=False):
    """ Add colors and numbers for drawing the grid to the terminal. """
    if not no_colors:
        from colorama import Fore, Style, init
        init()
        colors = [Fore.BLUE, Fore.RED, Fore.WHITE, Fore.GREEN, Fore.YELLOW, Fore.MAGENTA, Fore.CYAN, Fore.BLACK]
        colors = colors + [c+Style.BRIGHT for c in colors]
        reset = Style.RESET_ALL + Fore.RESET
    else:
        colors = ['']
        reset = ''
    tube_grid, uf = make_tubes(grid)
    letters = string.digits[1:] + string.ascii_letters
    char = collections.defaultdict(lambda: letters[len(char)])
    col = collections.defaultdict(lambda: colors[len(col)%len(colors)])
    for x in range(tube_grid.w):
        for y in range(tube_grid.h):
            if tube_grid[x,y] == 'x':
                tube_grid[x,y] = char[find(uf,(x,y))]
            tube_grid[x,y] = col[find(uf,(x,y))] + tube_grid[x,y] + reset
    return tube_grid


def add_path(x0, y0, xn, yn, grid, tries=1):
    for i in range(tries):
        if grid.add(make_path(x0, y0, xn, yn, grid)):
            return 1+i
    return False


def has_loops(grid, uf):
    groups = len({find(uf, (x,y)) for y in range(grid.h) for x in range(grid.w)})
    ends = sum(bool(grid[x,y] in 'v^<>') for y in range(grid.h) for x in range(grid.w))
    return ends != 2*groups


def has_square(grid, uf):
    for y in range(grid.h-1):
        for x in range(grid.w-1):
            if find(uf, (x,y)) == find(uf, (x+1,y)) == find(uf, (x,y+1)) == find(uf, (x+1,y+1)):
                return True
    return False


def has_pair(tg, uf):
    for y in range(tg.h):
        for x in range(tg.w):
            for dx,dy in ((1,0), (0,1)):
                x1, y1 = x + dx, y + dy
                if x1 < tg.w and y1 < tg.h:
                    if tg[x,y] == tg[x1,y1] == 'x' and find(uf, (x,y)) == find(uf, (x1,y1)):
                        return True
    return False


def has_tripple(tg, uf):
    for y in range(tg.h):
        for x in range(tg.w):
            if tg[x,y] != 'x':
                continue
            r = find(uf, (x,y))
            nbs = 0
            for dx, dy in ((1,0), (0,1), (-1,0), (0,-1)):
                x1, y1 = x + dx, y + dy
                if 0 <= x1 < tg.w and 0 <= y1 < tg.h and find(uf, (x1,y1)) == r:
                    nbs += 1
            if nbs >= 2:
                return True
    return False


def make(w, h, min_numbers=0, max_numbers=100):
    """ Creates a grid of size  w x h  without any loops or squares. """

    # Internally we work on a double size grid to handle crossings
    grid = Grid(2*w+1, 2*h+1)

    debug('Preprocessing...')
    mitm = Mitm(lr_price=2, t_price=1)
    mitm.prepare(max(h,10))
    debug('Generating puzzle...')

    gtries = 0
    while True:
        # Previous tries may have drawn stuff on the grid
        grid.grid.clear()

        # Add left side path
        path = mitm.rand_path(h, h, 0, -1)
        if not test_path(grid, path, 0, 0):
            continue
        draw_path(grid, path, 0, 0)
        # Draw_path doesn't know what to put in the first and last squares
        grid[0,0], grid[0,2*h] = '\\', '/'

        # Add right side path
        path2 = mitm.rand_path(h, h, 0, -1)
        if not test_path(grid, path2, 2*w, 2*h, 0, -1):
            continue
        draw_path(grid, path2, 2*w, 2*h, 0, -1)
        grid[2*w,0], grid[2*w,2*h] = '/', '\\'

        # Add loops in the middle
        # Tube version of full grid, using for tracking orientations.
        # This doesn't make so much sense in terms of normal numberlink tubes.
        tg, _ = make_tubes(grid)
        # Maximum number of tries before retrying main loop
        for tries in range(1000):
            x, y = 2*random.randrange(w), 2*random.randrange(h)

            # If the square square doen't have an orientation, it's a corner
            # or endpoint, so there's no point trying to add a loop there.
            if tg[x,y] not in '-|':
                continue

            path = mitm.rand_loop(clock=1 if tg[x,y] == '-' else -1)
            if test_path(grid, path, x, y):
                new_path = Grid(2*w+1, 2*h+1)
                draw_path(new_path, path, x, y, loop=True)

                # A loop may not overlap with anything, and may even have
                # the right orientation, but if it 'traps' something inside it, that
                # might now have the wrong orientation.
                # Hence we clear the insides.
                inner = (key for key, val in make_tubes(new_path)[0] if val == '|')
                for key in inner:
                    grid.grid.pop(key, None)

                # Add path and recompute orientations
                grid.add(new_path)
                tg, _ = make_tubes(grid)
                sg = shrink_grid(grid)
                stg, uf = make_tubes(sg)

                numbers = list(stg.grid.values()).count('x')//2
                if numbers > max_numbers:
                    break

                # Run tests to see if the puzzle is nice
                if not has_square(sg, uf) \
                        and not has_loops(sg, uf) \
                        and not has_pair(stg, uf) \
                        and not has_tripple(stg, uf) \
                        and numbers >= min_numbers:
                    debug(f'Finished in {tries} tries.')
                    debug(f'{numbers} numbers')
                    return shrink_grid(grid)

        debug(grid)
        debug(f'Gave up after {tries} tries')


def solution_lines(grid):
    # Assumes that no two end points are adjacent
    tg, uf = make_tubes(grid)
    done = {}
    for y in range(tg.h):
        for x in range(tg.w):
            if tg[x,y] == 'x':
                r = find(uf, (x,y))
                if r in done:
                    continue
                # Trace
                line = [(x,y)]
                while not (len(line) >= 2 and tg[line[-1]] == 'x'):
                    for dx, dy in ((-1,0),(1,0),(0,1),(0,-1)):
                        x1, y1 = line[-1][0]+dx, line[-1][1]+dy
                        if 0 <= x1 < tg.w and 0 <= y1 < tg.h \
                                and find(uf,(x1,y1)) == r \
                                and (len(line) == 1 or (x1,y1) != line[-2]):
                            line.append((x1,y1))
                            break
                done[r] = line
    return done.values()


def plot_puzzle(grid, include_solution=False):
    import matplotlib.pyplot as plt
    tg, uf = make_tubes(grid)

    # Draw a grid
    middle_linestyle = dict( lw=1, color='grey', linestyle='--')
    end_linestyle = dict( lw=2, color='black',)
    for y in range(tg.h+1):
        plt.gca().add_line(plt.Line2D((-1/2, tg.w-1/2), (y-1/2,y-1/2),
            **(middle_linestyle if y not in (0,tg.h) else end_linestyle)))
    for x in range(tg.w+1):
        plt.gca().add_line(plt.Line2D((x-1/2,x-1/2), (-1/2, tg.h-1/2),
            **(middle_linestyle if x not in (0,tg.w) else end_linestyle)))

    # Draw the solution
    if include_solution:
        linestyle = dict( lw=3, color='cornflowerblue', solid_joinstyle='round')
        for line in solution_lines(grid):
            plt.gca().add_line(plt.Line2D(*zip(*line), **linestyle))

    # Compute maximum number, used for font-size
    number = collections.defaultdict(lambda: len(number))
    for y in range(tg.h):
        for x in range(tg.w):
            if tg[x,y] == 'x':
                number[find(uf, (x,y))]
    max_len = len(str(len(number)))

    # Draw the numbers
    for y in range(tg.h):
        for x in range(tg.w):
            if tg[x,y] == 'x':
                n = number[find(uf, (x,y))] + 1
                plt.gca().add_artist(plt.Circle((x,y),1/3,color='white',zorder=2))
                plt.text(x, y-.1, str(n),
                        ha='center', va='center',
                        fontsize=16 if max_len > 1 else 22,
                        family='fantasy'
                        )

    plt.axis('scaled')
    plt.axis('off')
    plt.show()


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

    n = int((w*h)**.5)
    min_numbers = n*2//3 if args.min < 0 else args.min
    max_numbers = n*3//2 if args.max < 0 else args.max

    for _ in range(args.n):
        grid = make(w, h, min_numbers, max_numbers)
        color_grid = color_tubes(grid, no_colors=args.no_colors)

        # Print stuff
        debug(grid)
        if args.solve:
            print(color_grid)

        print(w, h)
        if args.zero:
            # Print puzzle in 0 format
            for y in range(color_grid.h):
                for x in range(color_grid.w):
                    if grid[x,y] in 'v^<>':
                        print(color_grid[x,y], end=' ')
                    else: print('0', end=' ')
                print()
        else:
            for y in range(color_grid.h):
                for x in range(color_grid.w):
                    if grid[x,y] in 'v^<>':
                        print(color_grid[x,y], end='')
                    else: print('.', end='')
                print()
        print()

        # Draw with pyplot
        if not args.terminal_only:
            plot_puzzle(grid, include_solution=args.solve)


if __name__ == '__main__':
    main()

