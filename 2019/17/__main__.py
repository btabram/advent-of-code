import re

from intcode_computer import CompleteIntcodeComputer, Status


with open("17/input.txt") as f:
    program = list(map(int, f.readline().split(",")))

computer = CompleteIntcodeComputer(program)
status, output = computer.start()

row = 0
grid = []
for val in output:
    if val == 10:
        row += 1
    else:
        if len(grid) == row:
            grid.append([])
        grid[row].append(val)

height = len(grid)
width = len(grid[0])


# Find the robot starting position and print the grid.
start = None
for y in range(height):
    line = ""
    for x in range(width):
        line += chr(grid[y][x])
        if line[-1] == "^":
            start = (x, y)
    print(line)

# Work out the sum of the alignment parameters.
total_alignment = 0
for y in range(1, height - 1):
    for x in range(1, width - 1):
        if grid[y][x] != 35:
            continue

        neighbours = [
            grid[y + 1][x],
            grid[y - 1][x],
            grid[y][x + 1],
            grid[y][x - 1],
        ]

        if len([n for n in neighbours if n == 35]) == 4:
            total_alignment += x * y
print("The answer to Part 1 is {}.".format(total_alignment))


# For Part 2 we need the robot to navigate the whole scaffold. We go in a
# straight line for as long as we can and then turn to follow the scaffold.
directions = [(1, 0), (0, 1), (-1, 0), (0, -1)]
pos = start
current_direction = (0, -1) # Start facing up
current_straight = 0

def try_direction(direction):
    x, y = pos
    x_step, y_step = direction
    try:
        return grid[y + y_step][x + x_step] == 35
    except IndexError:
        return False

# Build up the path the robot takes to traverse the whole scaffold using the
# syntax which the vacuum robot understands.
path = ""
while True:
    if try_direction(current_direction):
        current_straight += 1
        x, y = pos
        x_step, y_step = current_direction
        pos =  (x + x_step, y + y_step)
        continue

    current_direction_index = directions.index(current_direction)
    left_direction = directions[current_direction_index - 1]
    right_direction = directions[(current_direction_index + 1) % 4]

    if try_direction(left_direction):
        path += str(current_straight) + ",L,"
        current_straight = 0
        current_direction = left_direction
    elif try_direction(right_direction):
        path += str(current_straight) + ",R,"
        current_straight = 0
        current_direction = right_direction
    else:
        path += str(current_straight) + "," # Required to make regex work
        break

if path.startswith("0,"):
    path = path[2:]

# We need to factorise the path into 3 repeated units which are each no more
# that 20 characters to fit the instructions in the robot's memory. Someone on
# the internet had the idea to use a regex which is a very elegant solution.
m = re.match(r"^(.{2,20})\1*(.{2,20})(?:\1|\2)*(.{2,20})(?:\1|\2|\3)*$", path)
functions = [("A", m.group(1)), ("B", m.group(2)), ("C", m.group(3))]

main_routine = path
for letter, function in functions:
    main_routine = main_routine.replace(function, letter + ",")
assert len(main_routine) <= 20

# Strip trailing commas.
main_routine = re.sub(",$", "", main_routine)
functions = [(l, re.sub(",$", "", f)) for l, f in functions]

# Translate our string instructions into ASCII values for the robot.
robot_instructions = list(map(ord, main_routine))
robot_instructions += [10] # new line
for _, function in functions:
    robot_instructions += list(map(ord, function))
    robot_instructions += [10] # new line
robot_instructions += [ord("n"), 10] # Decline continuous video feed
robot_instructions.reverse()

part2_program = program
part2_program[0] = 2
status, output = computer.start()
while status == Status.INPUT_REQUIRED:
    status, output = computer.resume(robot_instructions.pop())
print("The answer to Part 2 is {}.".format(output[-1]))
