from intcode_computer import CompleteIntcodeComputer, Status


UP = 0
RIGHT = 1
DOWN = 2
LEFT = 3


class PaintingRobot:
    DIRECTIONS = [UP, RIGHT, DOWN, LEFT]

    def __init__(self, program):
        self.computer = CompleteIntcodeComputer(program)
        self.direction = UP
        self.panels = dict()
        self.position = (0, 0)

    def paint(self, value):
        assert value in [0, 1]
        self.panels[self.position] = value

    def rotate(self, value):
        assert value in [0, 1]
        index = self.DIRECTIONS.index(self.direction)
        if value == 0: # Turn left
            index -= 1
        else: # Turn right
            index = (index + 1) % 4
        self.direction = self.DIRECTIONS[index]

    # Note that positive y is downwards
    def move(self):
        assert self.direction in self.DIRECTIONS
        x, y = self.position
        if self.direction == UP:
            y -= 1
        elif self.direction == RIGHT:
            x += 1
        elif self.direction == DOWN:
            y += 1
        else: # LEFT
            x -= 1
        self.position = (x, y)

    def detect(self):
        if self.position in self.panels:
            return self.panels[self.position]
        else:
            return 0 # All panels are initially black

    def run(self, starting_panel):
        self.direction = UP
        self.panels = dict()
        self.position = (0, 0)
        status, output = self.computer.start(starting_panel)
        while status == Status.INPUT_REQUIRED:
            assert len(output) == 2
            self.paint(output[0])
            self.rotate(output[1])
            self.move()
            current_panel_colour = self.detect()
            status, output = self.computer.resume(current_panel_colour)


with open("11/input.txt") as f:
    program = list(map(int, f.readline().split(",")))

robot = PaintingRobot(program)
# Start on a black panel for Part 1.
robot.run(0)
print(f"The answer to Part 1 is {len(robot.panels)}.")

# Start on a white panel for Part 2.
robot.run(1)
max_x = max(robot.panels.keys(), key = lambda pos: pos[0])[0]
min_x = min(robot.panels.keys(), key = lambda pos: pos[0])[0]
max_y = max(robot.panels.keys(), key = lambda pos: pos[1])[1]
min_y = min(robot.panels.keys(), key = lambda pos: pos[1])[1]

print("The painted pattern in Part 2 is:")
for y in range(min_y, max_y + 1):
    line = ""
    for x in range(min_x, max_x + 1):
        colour = 0 # All panels were initially black
        if (x, y) in robot.panels:
            colour = robot.panels[(x, y)]
        line += "#" if colour == 1 else " "
    print(line)
