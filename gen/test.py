import gen

def test_all_sizes():
    for w in range(4,10):
        for h in range(4,10):
            print(w,h)
            gen.make(w,h)

test_all_sizes()
