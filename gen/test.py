import gen
from mitm import Mitm
import time


def test_all_sizes():
    for low in range(4, 15):
        t = time.time()
        for _ in range(10):
            for w in range(4, 15):
                for h in range(4, 15):
                    if abs(w - h) > 5:
                        continue
                    mitm = Mitm(lr_price=2, t_price=1)
                    mitm.prepare(max(h, low))
                    n = int((w * h)**.5)
                    gen.make(w, h, mitm, min_numbers=n * 2 // 3, max_numbers=n * 3 // 2)
        print('low', low, 'time', time.time() - t)


test_all_sizes()
