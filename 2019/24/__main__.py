from collections import deque
from copy import deepcopy
from typing import Deque, List, Set, Tuple


BUG = "#"
EMPTY = "."
RECURSION = "?"


class ErisSurface(object):
    def __init__(self, input: str) -> None:
        self.data: List[List[int]] = []
        for line in input.splitlines():
            row = []
            for char in line:
                row.append(char)
            self.data.append(row)

    def advance(self) -> None:
        # Don't update |self.data| during the loop because bug updates happen simultaneously.
        new_data = deepcopy(self.data)
        for y, row in enumerate(self.data):
            for x, tile in enumerate(row):
                neighbouring_bugs = self.getNeighbouringBugs(x, y)
                if tile == BUG:
                    if neighbouring_bugs != 1:
                        new_data[y][x] = EMPTY
                elif tile == EMPTY:
                    if neighbouring_bugs in [1, 2]:
                        new_data[y][x] = BUG
                else:
                    raise Exception("Unexpected item found on Eris' surface!")
        self.data = new_data

    def getNeighbouringBugs(self, x: int, y: int) -> int:
        all_neighbours = [(x + 1, y), (x, y + 1), (x - 1, y), (x, y - 1)]
        valid_neighbours = [(x, y) for (x, y) in all_neighbours if 0 <= x <= 4 and 0 <= y <= 4]
        bug_count = 0
        for n_x, n_y in valid_neighbours:
            if self.data[n_y][n_x] == BUG:
                bug_count += 1
        return bug_count

    def getBiodiversityRating(self) -> int:
        i = 0
        rating = 0
        for row in self.data:
            for tile in row:
                if tile == BUG:
                    rating += 1 << i
                i += 1
        return rating

    def __str__(self) -> str:
        s = ""
        for row in self.data:
            for tile in row:
                s += tile
            s += "\n"
        return s


class RecursiveErisSurface(object):
    def __init__(self, input: str) -> None:
        self.data: Deque[List[List[str]]] = deque()
        self.data.append([])
        for line in input.splitlines():
            row = []
            for char in line:
                row.append(char)
            self.data[0].append(row)
        self.data[0][2][2] = RECURSION
        self.depth_0_index = 0

    def newLayer(self) -> List[List[str]]:
        layer = [[EMPTY for _ in range(5)] for _ in range(5)]
        layer[2][2] = RECURSION
        return layer

    def isEmptyLayer(self, z: int) -> bool:
        for row in self.data[z]:
            if BUG in row:
                return False
        return True

    def advance(self) -> None:
        # We need to have consider above and below our currently populated layers since the bugs # could spread there.
        self.data.appendleft(self.newLayer())
        self.data.append(self.newLayer())
        self.depth_0_index += 1

        # Don't update |self.data| during the loop because bug updates happen simultaneously.
        new_data = deepcopy(self.data)
        for z, layer in enumerate(self.data):
            for y, row in enumerate(layer):
                for x, tile in enumerate(row):
                    # The middle tile doens't really exist, it's a whole new recursive layer instead.
                    if x == 2 and y == 2:
                        continue
                    neighbouring_bugs = self.getNeighbouringBugs(x, y, z)
                    if tile == BUG:
                        if neighbouring_bugs != 1:
                            new_data[z][y][x] = EMPTY
                    elif tile == EMPTY:
                        if neighbouring_bugs in [1, 2]:
                            new_data[z][y][x] = BUG
                    else:
                        print(x, y)
                        raise Exception("Unexpected item found on Eris' surface!")
        self.data = new_data

        # Prune empty layers.
        while self.isEmptyLayer(0):
            self.data.popleft()
            self.depth_0_index -= 1
        while self.isEmptyLayer(len(self.data) - 1):
            self.data.pop()

    def getNeighbouringBugs(self, x: int, y: int, z: int) -> int:
        neighbours = [(x + 1, y, z), (x, y + 1, z), (x - 1, y, z), (x, y - 1, z)]

        def handleLayerAbove(neighbour: Tuple[int, int, int]) -> Tuple[int, int, int]:
            n_x, n_y, n_z = neighbour

            # We add an extra empty layer above at the start of every timestep so we know there's no bugs above the top
            # layer and so don't need to consider it.
            if n_z == 0:
                return neighbour

            if n_x == -1:
                return (1, 2, n_z - 1)
            elif n_x == 5:
                return (3, 2, n_z - 1)
            elif n_y == -1:
                return (2, 1, n_z - 1)
            elif n_y == 5:
                return (2, 3, n_z - 1)
            else:
                return neighbour

        neighbours = [handleLayerAbove(n) for n in neighbours]

        # Handle neighbours on the layer below. We add an extra empty layer below at the start of every timestep so we
        # know there's no bugs below the bottomt layer and so don't need to consider it.
        if z != len(self.data) - 1 and (2, 2, z) in neighbours:
            neighbours.remove((2, 2, z))
            if x == 1:
                neighbours += [(0, n_y, z + 1) for n_y in range(5)]
            elif x == 3:
                neighbours += [(4, n_y, z + 1) for n_y in range(5)]
            elif y == 1:
                neighbours += [(n_x, 0, z + 1) for n_x in range(5)]
            elif y == 3:
                neighbours += [(n_x, 4, z + 1) for n_x in range(5)]
            else:
                raise Exception("Unexpected ({},{},{}) neighbours".format(x, y, z))

        # We may still have some invalid neighbours to filter out.
        valid_neighbours = [(x, y, z) for (x, y, z) in neighbours if 0 <= x <= 4 and 0 <= y <= 4]

        bug_count = 0
        for n_x, n_y, n_z in valid_neighbours:
            if self.data[n_z][n_y][n_x] == BUG:
                bug_count += 1
        return bug_count

    def countBugs(self) -> int:
        count = 0
        for layer in self.data:
            for row in layer:
                for tile in row:
                    if tile == BUG:
                        count += 1
        return count

    def __str__(self) -> str:
        s = ""
        for i, layer in enumerate(self.data):
            s += "Depth {}:\n".format(i - self.depth_0_index)
            for row in layer:
                for tile in row:
                    s += tile
                s += "\n"
            s += "\n"
        return s


with open("24/input.txt") as f:
    input = f.read()

# Part 1.
eris_surface = ErisSurface(input)
previous_layouts: Set[str] = set()
previous_layouts.add(str(eris_surface))
while True:
    eris_surface.advance()
    layout = str(eris_surface)
    if layout in previous_layouts:
        print("The answer to Part 1 is {}.".format(eris_surface.getBiodiversityRating()))
        break
    previous_layouts.add(layout)

# Part 2.
recursive_eris_surface = RecursiveErisSurface(input)
for _ in range(200):
    recursive_eris_surface.advance()
print("The answer to Part 2 is {}.".format(recursive_eris_surface.countBugs()))