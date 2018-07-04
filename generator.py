import random
import math
import argparse

parser = argparse.ArgumentParser(description='Generate a cave.')
parser.add_argument('--count', type=int, default=100,
                    help='The number of caverns to create')
parser.add_argument('--width', type=int, default=400,
                    help='The width of the cave')
parser.add_argument('--height', type=int, default=200,
                    help='The height of the cave')
parser.add_argument('--connectivity', type=int, default=50,
                    help='The chance for 2 caverns to be connected')
parser.add_argument('--radius', type=int, default=30,
                    help='The max radius for a cavern\'s connections')

args = parser.parse_args()


items = []

def exists(x,y):
    for i in items:
        if i["x"] == x and i["y"] == y:
            return True

    return False


def gen_item(xmin, xmax, ymin, ymax):
    x = random.randint(xmin, xmax)
    y = random.randint(ymin, ymax)

    if exists(x, y):
        return gen_item(xmin, xmax, ymin, ymax)

    return x, y

def distance(a, b):
    return math.sqrt(math.pow(a["x"]-b["x"], 2) + math.pow(a["y"] - b["y"], 2))


with open("generated.cav", "w") as f:
    f.write(str(args.count))
    for n in range(args.count):
        if n == 0:
            x, y = gen_item(0, args.width/3, 0, args.height/3)
        elif n == args.count-1:
            x, y = gen_item(args.width/3*2, args.width, args.height/3*2, args.height)
        else:
            x, y = gen_item(0, args.width, 0, args.height)
        items.append({
            "x": x,
            "y": y
        })
        f.write(",{},{}".format(x,y))

    for n in range(args.count):
        for i in range(args.count):
            if n == i:
                f.write(",0")
            elif distance(items[n], items[i]) > args.radius:
                f.write(",0")
            else:
                c = random.randint(0,100)
                f.write(",{}".format(int(c <= args.connectivity)))
