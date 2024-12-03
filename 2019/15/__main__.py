from intcode_computer import CompleteIntcodeComputer, Status


NORTH = 1
SOUTH = 2
WEST = 3
EAST = 4
DIRECTIONS = [NORTH, EAST, SOUTH, WEST]


class RepairDroid:
    def __init__(self, program):
        self.computer = CompleteIntcodeComputer(program)
        self.position = (0, 0)
        self.direction = NORTH
        self.grid = dict()
        self.first_wall = None

    def rotate_left(self):
        index = DIRECTIONS.index(self.direction)
        self.direction = DIRECTIONS[index - 1]

    def rotate_right(self):
        index = DIRECTIONS.index(self.direction)
        self.direction = DIRECTIONS[(index + 1) % 4]

    @property
    def next_position(self):
        x, y = self.position
        assert self.direction in DIRECTIONS
        if self.direction == NORTH:
            y -= 1
        elif self.direction == SOUTH:
            y += 1
        elif self.direction == EAST:
            x += 1
        else: # WEST
            x -= 1
        return (x, y)

    def build_grid(self):
        status, _ = self.computer.start()

        # Follow the left wall until we're back where we started. Turn left
        # after each move and then try to move forwards, turning right if we
        # hit a wall.
        while status == Status.INPUT_REQUIRED:
            self.rotate_left()
            while True:
                status, [output] = self.computer.resume(self.direction)
                if output != 0:
                    break

                # Stop if we're back where we started.
                if self.first_wall is None:
                    self.first_wall = self.next_position
                elif self.next_position == self.first_wall \
                        and self.position == (0, 0):
                    return

                self.grid[self.next_position] = 0
                self.rotate_right()

            assert output in [1, 2]
            self.grid[self.next_position] = output
            self.position = self.next_position

    def print_grid(self):
        max_x = max(self.grid.keys(), key = lambda pos: pos[0])[0]
        min_x = min(self.grid.keys(), key = lambda pos: pos[0])[0]
        max_y = max(self.grid.keys(), key = lambda pos: pos[1])[1]
        min_y = min(self.grid.keys(), key = lambda pos: pos[1])[1]

        char_conversions = {
            -1: " ", # Not visited
            0 : "█", # Wall
            1 : ".", # Visited corridor
            2 : "O", # Oxygen system
            3 : "D", # Doid starting location
        }

        for y in range(min_y, max_y + 1):
            line = ""
            for x in range(min_x, max_x + 1):
                if (x, y) == (0, 0):
                    val = 3 # Droid starting position
                elif (x, y) in self.grid:
                    val = self.grid[(x, y)]
                else:
                    val = -1
                line += char_conversions[val]
            print(line)


with open("15/input.txt") as f:
    program = list(map(int, f.readline().split(",")))

rp = RepairDroid(program)
rp.build_grid()
print("The map of the ship section is:")
rp.print_grid()

# Track the edge of the oxygenated region.
oxygen_border = set()
oxygen_border.add([pos for pos, val in rp.grid.items() if val == 2][0])

t = 0
while True:
    new_border = set()
    # For every square on the oxygen border we identify neighbours which the
    # oxygen will spread to in the next minute.
    for x, y in oxygen_border:
        neighbours = [(x + 1, y), (x - 1, y), (x, y + 1), (x, y - 1)]
        [new_border.add(n) for n in neighbours if rp.grid[n] == 1]

    if len(new_border) == 0:
        # The oxygen has spread to every possible square.
        break

    if (0, 0) in oxygen_border:
        # |t| at this point is equal to the number of steps from the starting
        # position to the oxygen system.
        print(f"The answer to Part 1 is {t}.")

    t += 1
    oxygen_border = new_border

    # Update the squares which oxygen has spread to.
    for pos in oxygen_border:
        rp.grid[pos] = 2
print(f"The answer to Part 2 is {t}.")
