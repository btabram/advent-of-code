from intcode_computer import CompleteIntcodeComputer, Status


with open("19/input.txt") as f:
    program = list(map(int, f.readline().split(",")))

computer = CompleteIntcodeComputer(program)


count = 0
for y in range(50):
    for x in range(50):
        _, [output] = computer.start([x, y])
        count += output

print("The answer to Part 1 is {}.".format(count))


def is_in_beam(x, y):
    _, [output] = computer.start([x, y])
    return output == 1

x = 0
y = 100
done = False

print("This will take a little while...")
while True:
    # Go right till we hit the beam.
    while not is_in_beam(x, y):
        x += 1

    # Check if the beam is wide enough.
    if is_in_beam(x + 99, y):
        i = 0
        # Check if the beam is tall enough. If not, then move to the right
        # (if the beam is wide enough) to try to find a 100x100 square.
        while is_in_beam(x + i + 99, y):
            if is_in_beam(x + i, y + 99) and is_in_beam(x + i + 99, y + 99):
                done = True
                break
            i += 1

        if done:
            x += i
            break

    # Move down and try again.
    y += 1

print("The answer to Part 2 is {}.".format(x * 1000 + y))
