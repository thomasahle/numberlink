import sys
import collections
import string

from grid import Grid, UnionFind


def color_tubes(grid, no_colors=False):
    """ Add colors and numbers for drawing the grid to the terminal. """
    if not no_colors:
        from colorama import Fore, Style, init
        init()
        colors = [
            Fore.BLUE,
            Fore.RED,
            Fore.WHITE,
            Fore.GREEN,
            Fore.YELLOW,
            Fore.MAGENTA,
            Fore.CYAN,
            Fore.BLACK]
        colors = colors + [c + Style.BRIGHT for c in colors]
        reset = Style.RESET_ALL + Fore.RESET
    else:
        colors = ['']
        reset = ''
    tube_grid, uf = grid.make_tubes()
    letters = string.digits[1:] + string.ascii_letters
    char = collections.defaultdict(lambda: letters[len(char)])
    col = collections.defaultdict(lambda: colors[len(col) % len(colors)])
    for x in range(tube_grid.w):
        for y in range(tube_grid.h):
            if tube_grid[x, y] == 'x':
                tube_grid[x, y] = char[uf.find( (x, y))]
            tube_grid[x, y] = col[uf.find( (x, y))] + tube_grid[x, y] + reset
    return tube_grid, char


def solution_lines(tg, uf):
    # Assumes that no two end points are adjacent
    done = {}
    for y in range(tg.h):
        for x in range(tg.w):
            if tg[x, y] == 'x':
                r = uf.find( (x, y))
                if r in done:
                    continue
                # Trace
                line = [(x, y)]
                while not (len(line) >= 2 and tg[line[-1]] == 'x'):
                    for dx, dy in ((-1, 0), (1, 0), (0, 1), (0, -1)):
                        x1, y1 = line[-1][0] + dx, line[-1][1] + dy
                        if 0 <= x1 < tg.w and 0 <= y1 < tg.h \
                                and uf.find( (x1, y1)) == r \
                                and (len(line) == 1 or (x1, y1) != line[-2]):
                            line.append((x1, y1))
                            break
                done[r] = line
    return done.values()


def plot_puzzle(tg, uf, include_solution=False, save_to=None):
    """ If include_solutions=False the tube_grid can simple contain x's at the
        endpoints. No actual tubes needed. """
    import matplotlib.pyplot as plt

    plt.figure(num=None, figsize=(tg.w*.6, tg.h*.6), dpi=80, facecolor='w', edgecolor='k')

    # Draw a grid
    middle_linestyle = dict(lw=1, color='grey', linestyle='--')
    end_linestyle = dict(lw=2, color='black',)
    for y in range(tg.h + 1):
        plt.gca().add_line(plt.Line2D((-1 / 2, tg.w - 1 / 2), (y - 1 / 2, y - 1 / 2),
                                      **(middle_linestyle if y not in (0, tg.h) else end_linestyle)))
    for x in range(tg.w + 1):
        plt.gca().add_line(plt.Line2D((x - 1 / 2, x - 1 / 2), (-1 / 2, tg.h - 1 / 2),
                                      **(middle_linestyle if x not in (0, tg.w) else end_linestyle)))

    # Draw the solution
    if include_solution:
        linestyle = dict(lw=3, color='cornflowerblue', solid_joinstyle='round')
        for line in solution_lines(tg, uf):
            plt.gca().add_line(plt.Line2D(*zip(*line), **linestyle))

    # Compute maximum number, used for font-size
    number = collections.defaultdict(lambda: len(number))
    for y in range(tg.h):
        for x in range(tg.w):
            if tg[x, y] == 'x':
                number[uf.find( (x, y))]
    max_len = len(str(len(number)))

    # Draw the numbers
    for y in range(tg.h):
        for x in range(tg.w):
            if tg[x, y] == 'x':
                n = number[uf.find( (x, y))] + 1
                plt.gca().add_artist(plt.Circle((x, y), 1 / 3, color='white', zorder=2))
                plt.text(x, y - .1, str(n),
                         ha='center', va='center',
                         fontsize=16 if max_len > 1 else 22,
                         family='fantasy'
                         )

    plt.axis('scaled')
    plt.axis('off')
    if save_to:
        plt.savefig(save_to, bbox_inches='tight')
    else:
        plt.show()


def main():
    n = 0
    for line in sys.stdin:
        line = line.strip()
        if not line or line.startswith('#'):
            continue
        w, h = map(int, line.split())
        grid = Grid(w, h)
        uf = {}
        for y in range(h):
            row = sys.stdin.readline().strip()
            for x, c in enumerate(row):
                if c not in '0.':
                    grid[x,y] = 'x'
                    uf[x,y] = c
        file_name = f'out_{n}.png'
        plot_puzzle(grid, UnionFind(uf), include_solution=False, save_to=file_name)
        n += 1


if __name__ == '__main__':
    main()

