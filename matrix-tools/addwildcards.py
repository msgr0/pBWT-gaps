import collections
import sys
import random

wildcards_rate = float(sys.argv[2])

linelen = 0
nlines = 0

with open(sys.argv[1], 'r') as f:
    linelen = len(f.readline().rstrip())
    nlines = 1
    for line in f:
        nlines += 1

print("no. of lines = ", nlines, file=sys.stderr)
print("line length = ", linelen, file=sys.stderr)
print("wildcard rate = ", wildcards_rate, file=sys.stderr)

n_wildcards = int(nlines * wildcards_rate * linelen)
print("no. of wildcards = ", n_wildcards, file=sys.stderr)

positions = set()
while len(positions) < n_wildcards:
    x = random.randrange(linelen)
    y = random.randrange(nlines)
    positions.add((y, x))

positions = sorted(list(positions))

hist = collections.defaultdict(int)
for (line, x) in positions:
    hist[line] += 1


print("distribution of no. of wilcards per line", file=sys.stderr)
print("# no. of wildcards per line,  no. of lines", file=sys.stderr)
hist2 = collections.defaultdict(int)
for (line, nwild) in hist.items():
    hist2[nwild] += 1
for nwild in sorted(hist2):
    print("#", nwild, hist2[nwild], file=sys.stderr)

# Add a sentinel
positions.append((nlines + 1, 0))

next_position = 0
with open(sys.argv[1], 'r') as f:
    for i, line in enumerate(f):
        while positions[next_position][0] == i:
            position = positions[next_position][1]
            line = line[:position] + '*' + line[position+1:]
            next_position += 1
        print(line, end='')
