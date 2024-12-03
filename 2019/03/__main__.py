from dataclasses import dataclass


@dataclass(frozen=True)
class Coord(object):
    x: int
    y: int

    def distance(self, other):
        return abs(self.x - other.x) + abs(self.y - other.y)

    def __add__(self, other):
        return Coord(self.x + other.x, self.y + other.y)


@dataclass
class Direction(object):
    step: Coord
    distance: int

    def __init__(self, direction_string):
        cardinal = direction_string[0]
        if cardinal == "U":
            self.step = Coord(0, 1)
        elif cardinal == "D":
            self.step = Coord(0, -1)
        elif cardinal == "L":
            self.step = Coord(-1, 0)
        elif cardinal == "R":
            self.step = Coord(1, 0)
        else:
            raise Exception("Invalid cardinal: {}".format(cardinal))
        self.distance = int(direction_string[1:])


class Wire(object):
    def __init__(self, directions):
        self.current_position = Coord(0, 0)
        self.distance_travelled = 0
        self.distances = dict()
        self.points_crossed = set()

        for d in directions:
            self.do_direction(d)

    def do_direction(self, direction):
        """
        Extend a wire from |start| using |direction|.
        Return the coordinates of all points the wire has now crossed.
        """
        for _ in range(direction.distance):
            self.current_position += direction.step
            self.distance_travelled += 1
            self.points_crossed.add(self.current_position)
            # Keep the initial distance value if we revisit a coordinate.
            if self.current_position not in self.distances:
                self.distances[self.current_position] = self.distance_travelled


with open("03/input.txt") as f:
    wire_1 = Wire(map(Direction, f.readline().split(",")))
    wire_2 = Wire(map(Direction, f.readline().split(",")))

intersections = wire_1.points_crossed.intersection(wire_2.points_crossed)

min_dist = None
origin = Coord(0, 0)
for inter in intersections:
    dist = origin.distance(inter)
    if min_dist is None or dist < min_dist:
        min_dist = dist
print("The answer to Part 1 is {}.".format(min_dist))

min_dist = None
for inter in intersections:
    dist = wire_1.distances[inter] + wire_2.distances[inter]
    if min_dist is None or dist < min_dist:
        min_dist = dist
print("The answer to Part 2 is {}.".format(min_dist))
