from itertools import count


class MazeElement:
    def __init__(self, value):
        self.value = value
        self.portal = None
        self.portal_is_outer = None

    def __repr__(self):
        return "*" if self.portal else self.value

    def add_portal(self, portal_name, outer):
        assert self.value == "."
        self.portal = portal_name
        self.portal_is_outer = outer

    # We use isupper to detect the portals in the grid (it returns false for
    # " ", "." and "#") so this class needs to implement it.
    def isupper(self):
        return self.value.isupper()


with open("20/input.txt") as f:
    input = list(map(lambda line: [c for c in line], f.read().splitlines()))

input_height = len(input)
input_width = len(input[0])

# Replace maze elements with MazeElement instances.
for y in range(input_height):
    for x in range(input_width):
        if (value:= input[y][x]) in ["#", "."]:
            input[y][x] = MazeElement(value)

# Parse vertical portals.
# Ignore space either side of the maze.
for x in range(2, input_width - 2):
    # No need to check the last column.
    for y in range(input_height - 1):
        value = input[y][x]
        next = input[y + 1][x]
        if value.isupper() and next.isupper():
            portal_name = value + next
            outer = y in [0, input_height - 2]
            try:
                input[y - 1][x].add_portal(portal_name, outer)
            except AttributeError:
                input[y + 2][x].add_portal(portal_name, outer)

# Parse horizontal portals.
# Ignore space above and below the maze.
for y in range(2, input_height - 2):
    # No need to check the last column.
    for x in range(input_width - 1):
        value = input[y][x]
        next = input[y][x + 1]
        if value.isupper() and next.isupper():
            portal_name = value + next
            outer = x in [0, input_width - 2]
            try:
                input[y][x - 1].add_portal(portal_name, outer)
            except AttributeError:
                input[y][x + 2].add_portal(portal_name, outer)

# Construct a maze dict and collect the portal locations.
maze = dict()
portal_locations = dict()
for y in range(2, input_height - 2):
    for x in range(2, input_width - 2):
        if isinstance(element:= input[y][x], MazeElement):
            maze[(x, y)] = element.value
            if (portal := element.portal) is not None:
                outer = element.portal_is_outer
                try:
                    portal_locations[element.portal].append(((x, y), outer))
                except KeyError:
                    portal_locations[element.portal] = [((x, y), outer)]

# Construct dicts of portal locations so we can easily travel through them.
portal_map = dict()
inwards_portal_map = dict()
outwards_portal_map = dict()
for locations in portal_locations.values():
    # The start and end portals only have a single location.
    if len(locations) != 2:
        continue
    (a_pos, a_outer), (b_pos, b_outer) = locations
    portal_map[a_pos] = b_pos
    portal_map[b_pos] = a_pos
    assert a_outer != b_outer
    if a_outer:
        outwards_portal_map[a_pos] = b_pos
        inwards_portal_map[b_pos] = a_pos
    else:
        outwards_portal_map[b_pos] = a_pos
        inwards_portal_map[a_pos] = b_pos


def is_valid_neighbour(neighbour):
    try:
        return maze[neighbour] != "#" # Can't walk through walls
    except KeyError:
        return False


def part1():
    def get_neighbours(position):
        x, y = position
        neighbours = [(x + 1, y), (x, y + 1), (x - 1, y), (x, y - 1)]
        try:
            neighbours.append(portal_map[position])
        except KeyError:
            pass
        return set(filter(is_valid_neighbour, neighbours))

    start = portal_locations["AA"][0][0]
    finish = portal_locations["ZZ"][0][0]

    visited = set([start])
    boundary = set([start])
    for i in count(1):
        new_boundary = set()
        for position in boundary:
            new_boundary |= get_neighbours(position)
        new_boundary -= visited

        if finish in new_boundary:
            print("The answer to Part 1 is {}.".format(i))
            return

        visited |= new_boundary
        boundary = new_boundary


def part2():
    def get_neighbours(positionAndLevel):
        x, y, l = positionAndLevel
        position = (x, y)
        neighbours = \
            [(x + 1, y, l), (x, y + 1, l), (x - 1, y, l), (x, y - 1, l)]
        # Always consider inwards portals.
        try:
            x, y = inwards_portal_map[position]
            neighbours.append((x, y, l + 1))
        except KeyError:
            pass
        # Consider outwards portals if we're not on the top level.
        if l > 0:
            try:
                x, y = outwards_portal_map[position]
                neighbours.append((x, y, l - 1))
            except KeyError:
                pass
        def is_valid(neighbourAndLevel):
            x, y, _ = neighbourAndLevel
            return is_valid_neighbour((x, y))
        return set(filter(is_valid, neighbours))

    s_x, s_y = portal_locations["AA"][0][0]
    start = (s_x, s_y, 0)
    f_x, f_y = portal_locations["ZZ"][0][0]
    finish = (f_x, f_y, 0)

    visited = set([start])
    boundary = set([start])
    for i in count(1):
        new_boundary = set()
        for position in boundary:
            new_boundary |= get_neighbours(position)
        new_boundary -= visited

        if finish in new_boundary:
            print("The answer to Part 2 is {}.".format(i))
            return

        visited |= new_boundary
        boundary = new_boundary


part1()
part2()
