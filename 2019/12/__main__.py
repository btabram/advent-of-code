from math import copysign, lcm
import re


POSITION_REGEX = re.compile("<x=(?P<x>-?\d+), y=(?P<y>-?\d+), z=(?P<z>-?\d+)>")


class Moon:
    def __init__(self, x, y, z):
        self.position = [x, y, z]
        self.velocity = [0, 0, 0] # Start at rest

    def __repr__(self):
        return f"Moon({self.position}, {self.velocity})"


    @property
    def energy(self):
        x, y, z = self.position
        v_x, v_y, v_z = self.velocity
        potential_energy = abs(x) + abs(y) + abs(z)
        kinetic_energy = abs(v_x) + abs(v_y) + abs(v_z)
        return potential_energy * kinetic_energy


def apply_gravity(moons):
    # Consider every moon pair, ensuring we don't consider A -> B if we have
    # already considered B -> A.
    for i, moon in enumerate(moons):
        for j in range(i + 1, len(moons)):
            other = moons[j]
            for axis in range(3):
                moon_pos = moon.position[axis]
                other_pos = other.position[axis]
                if moon_pos == other_pos:
                    continue
                # Velocity changes by +/-1 to bring the moons together.
                velocity_change = int(copysign(1, moon_pos - other_pos))
                moon.velocity[axis] -= velocity_change
                other.velocity[axis] += velocity_change


def update_positions(moons):
    for moon in moons:
        x, y, z = moon.position
        v_x, v_y, v_z = moon.velocity
        moon.position = [x + v_x, y + v_y, z + v_z]


with open("12/input.txt") as f:
    position_strings = f.readlines()

moons = []
for line in position_strings:
    m = POSITION_REGEX.match(line)
    assert m is not None
    moons.append(Moon(int(m.group("x")), int(m.group("y")), int(m.group("z"))))

initial_positions = [[], [], []]
initial_velocities = [[], [], []]
for moon in moons:
    for axis in range(3):
        initial_positions[axis].append(moon.position[axis])
        initial_velocities[axis].append(moon.velocity[axis])

repeat_periods = [None, None, None]
t = 0
while True:
    apply_gravity(moons)
    update_positions(moons)
    t += 1

    if t == 1000:
        print(f"The answer to Part 1 is {sum([m.energy for m in moons])}.")

    positions = [[], [], []]
    velocities = [[], [], []]
    for moon in moons:
        for axis in range(3):
            positions[axis].append(moon.position[axis])
            velocities[axis].append(moon.velocity[axis])

    # Part 2 is ticky and I had to look up the solution. The key is that we can
    # consider each dimension independently. The motion in each dimension has
    # an independent repeat period and the overall repeat period of the whole
    # system is the lowest common multiple of the individual repeat periods.
    for axis in range(3):
        if repeat_periods[axis] is None \
                and positions[axis] == initial_positions[axis] \
                and velocities[axis] == initial_velocities[axis]:
            # We've found the first repeated state in this dimension.
            repeat_periods[axis] = t

    # We've got all the information we need for Part 1 and Part 2.
    if t >= 1000 and None not in repeat_periods:
        break

overall_repeat_period = \
    lcm(repeat_periods[0], repeat_periods[1], repeat_periods[2])
print(f"The answer to Part 2 is {overall_repeat_period}.")
