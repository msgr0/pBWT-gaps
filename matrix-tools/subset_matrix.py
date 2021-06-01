import sys

x = int(sys.argv[1])
y = int(sys.argv[2])
dx = int(sys.argv[3])
dy = int(sys.argv[4])

for i, l in enumerate(sys.stdin):
    if i < y:
        continue
    if i >= y + dy:
        break

    print(l[x:x+dx].rstrip())
