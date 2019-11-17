def sign(x):
    if x == 0:
        return x
    return -1 if x < 0 else 1


class UnionFind:
    def __init__(self, initial=None):
        self.uf = initial or {}

    def union(self, a, b):
        a_par, b_par = self.find(a), self.find(b)
        self.uf[a_par] = b_par

    def find(self, a):
        if self.uf.get(a, a) == a:
            return a
        par = self.find(self.uf.get(a, a))
        # Path compression
        self.uf[a] = par
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
            res.append(''.join(self[x, y] for x in range(self.w)))
        return '\n'.join(res)

    def __iter__(self):
        return iter(self.grid.items())

    def __contains__(self, key):
        return key in self.grid

    def __delitem__(self, key):
        del self.grid[key]

    def clear(self):
        self.grid.clear()

    def values(self):
        return self.grid.values()

    def shrink(self):
        """ Returns a new grid of half the height and width """
        small_grid = Grid(self.w // 2, self.h // 2)
        for y in range(self.h // 2):
            for x in range(self.w // 2):
                small_grid[x, y] = self[2 * x + 1, 2 * y + 1]
        return small_grid

    def test_path(self, path, x0, y0, dx0=0, dy0=1):
        """ Test whether the path is safe to draw on the grid, starting at x0, y0 """
        return all(0 <= x0 - x + y < self.w and 0 <= y0 + x + y < self.h
                   and (x0 - x + y, y0 + x + y) not in self for x, y in path.xys(dx0, dy0))

    def draw_path(self, path, x0, y0, dx0=0, dy0=1, loop=False):
        """ Draws path on the grid. Asserts this is safe (no overlaps).
            For non-loops, the first and the last character is not drawn,
            as we don't know what shape they should have. """
        ps = list(path.xys(dx0, dy0))
        # For loops, add the second character, so we get all rotational tripples:
        # abcda  ->  abcdab  ->  abc, bcd, cda, dab
        if loop:
            assert ps[0] == ps[-1], (path, ps)
            ps.append(ps[1])
        for i in range(1, len(ps) - 1):
            xp, yp = ps[i - 1]
            x, y = ps[i]
            xn, yn = ps[i + 1]
            self[x0 - x + y, y0 + x + y] = {
                (1, 1, 1): '<', (-1, -1, -1): '<',
                (1, 1, -1): '>', (-1, -1, 1): '>',
                (-1, 1, 1): 'v', (1, -1, -1): 'v',
                (-1, 1, -1): '^', (1, -1, 1): '^',
                (0, 2, 0): '\\', (0, -2, 0): '\\',
                (2, 0, 0): '/', (-2, 0, 0): '/'
            }[xn - xp, yn - yp, sign((x - xp) * (yn - y) - (xn - x) * (y - yp))]

    def make_tubes(self):
        uf = UnionFind()
        tube_grid = Grid(self.w, self.h)
        for x in range(self.w):
            d = '-'
            for y in range(self.h):
                # We union things down and to the right.
                # This means ┌ gets to union twice.
                for dx, dy in {
                        '/-': [(0, 1)], '\\-': [(1, 0), (0, 1)],
                        '/|': [(1, 0)],
                        ' -': [(1, 0)], ' |': [(0, 1)],
                        'v|': [(0, 1)], '>|': [(1, 0)],
                        'v-': [(0, 1)], '>-': [(1, 0)],
                }.get(self[x, y] + d, []):
                    uf.union((x, y), (x + dx, y + dy))
                # We change alll <>v^ to x.
                tube_grid[x, y] = {
                    '/-': '┐', '\\-': '┌',
                    '/|': '└', '\\|': '┘',
                    ' -': '-', ' |': '|',
                }.get(self[x, y] + d, 'x')
                # We change direction on v and ^, but not on < and >.
                if self[x, y] in '\\/v^':
                    d = '|' if d == '-' else '-'
        return tube_grid, uf

    def clear_path(self, path, x, y):
        """ Removes everything contained in the path (loop) placed at x, y. """
        path_grid = Grid(self.w, self.h)
        path_grid.draw_path(path, x, y, loop=True)
        for key, val in path_grid.make_tubes()[0]:
            if val == '|':
                self.grid.pop(key, None)


