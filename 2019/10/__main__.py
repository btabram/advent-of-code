from math import copysign


class Asteroid:
    def __init__(self, x, y):
        self.x = x
        self.y = y
        self.in_LOS = None

    def __repr__(self):
        return f"Asteroid(({self.x}, {self.y}), {len(self.in_LOS)})"

    def _get_highest_common_factor(self, a, b):
        if a == 0 or b == 0:
            return 0

        a = abs(a)
        b = abs(b)

        if a == b:
            return a
        lower = min(a, b)

        for i in range(lower, 1, -1):
            if a % i == 0 and b % i == 0:
                return i
        return 1

    def calculate_LOS(self, asteroid_positions):
        self.in_LOS = []
        for (x, y) in asteroid_positions:
            if x == self.x and y == self.y:
                continue

            x_step = self.x - x
            y_step = self.y - y

            hcf = self._get_highest_common_factor(x_step, y_step)
            if hcf == 1:
                # In this case the line of sight between self and the asteroid
                # at (x, y) does not go through any other points on the grid
                # and hence the line of sight cannot be obstructed.
                self.in_LOS.append((x, y))
                continue
            elif hcf == 0:
                # Line of sight is parallel to a grid direction.
                if x_step == 0:
                    potential_obstructions = \
                        [(x, y + int(copysign(i, y_step))) for i in range(1, abs(y_step))]
                else: # y_step == 0
                    potential_obstructions = \
                        [(x + int(copysign(i, x_step)), y) for i in range(1, abs(x_step))]
            else:
                # Positions on the grid which lie exactly on the line of sight
                # between self and the asteroid at (x, y). Only an asteroid in
                # one of these positions could block line of sight.
                potential_obstructions = \
                    [(x + i*(x_step // hcf), y + i*(y_step // hcf)) for i in range(1, hcf)]

            obstructions = \
                [p for p in potential_obstructions if p in asteroid_positions]

            if len(obstructions) == 0:
                self.in_LOS.append((x, y))


with open("10/input.txt") as f:
    grid = f.readlines()
grid = [line.strip() for line in grid]

width = len(grid[0])
height = len(grid)

asteroids = dict()
for y in range(height):
    for x in range(width):
        if grid[y][x] == "#":
            asteroids[(x, y)] = Asteroid(x, y)

asteroid_positions = set(asteroids.keys())
for asteroid in asteroids.values():
    asteroid.calculate_LOS(asteroid_positions)

best_monitoring_station = \
    max(asteroids.values(), key = lambda asteroid: len(asteroid.in_LOS))
print(f"The answer to Part 1 is {len(best_monitoring_station.in_LOS)}.")

monitoring_x = best_monitoring_station.x
monitoring_y = best_monitoring_station.y
relative_asteroid_positions = \
    set([(x - monitoring_x, y - monitoring_y) for (x, y) in asteroid_positions])
# Remove the monitoring station.
relative_asteroid_positions.remove((0, 0))
asteroid_count = len(relative_asteroid_positions)

# For Part 2 we need to work out the angle between the monitoring station and
# the other asteroids. Start my splitting the other asteroids into the left
# and right sides. For each half we can do y_step/x_step to calculate tanθ.
# Asteroids directly above and below the monitoring station need to be dealt
# with separately to avoid dividing by zero. Note that positive y is downwards
# because the origin of our grid is the top left corner.
above = [(x, y) for (x, y) in relative_asteroid_positions if x == 0 and y < 0]
below = [(x, y) for (x, y) in relative_asteroid_positions if x == 0 and y > 0]
lhs = [(x, y) for (x, y) in relative_asteroid_positions if x > 0]
rhs = [(x, y) for (x, y) in relative_asteroid_positions if x < 0]

lhs.sort(key = lambda pos: pos[1] / pos[0])
rhs.sort(key = lambda pos: pos[1] / pos[0])

sectors = [above, lhs, below, rhs]

laser = Asteroid(0, 0)
relative_vaporised = []
while len(relative_vaporised) < asteroid_count:
    for sector in sectors:
        laser.calculate_LOS(sector)
        for asteroid in laser.in_LOS:
            relative_vaporised.append(asteroid)
            sector.remove(asteroid)

vaporised = \
    [(x + monitoring_x, y + monitoring_y) for (x, y) in relative_vaporised]
(x_200, y_200) = vaporised[199]
print(f"The answer to Part 2 is {x_200 * 100 + y_200}.")
