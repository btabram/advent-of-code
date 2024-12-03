from itertools import combinations

from intcode_computer import CompleteIntcodeComputer, Status


with open("25/input.txt") as f:
    program = list(map(int, f.readline().split(",")))

computer = CompleteIntcodeComputer(program, ascii=True)
status, output = computer.start()
assert status == Status.INPUT_REQUIRED
print(output)

def run_command(command: str):
    status, output = computer.resume(command + "\n")
    print(output)
    if status == Status.SUCCESS:
        exit(0)

# I manually worked out these commands.
collect_all_items_commands = [
    "south",
    "east",
    "take space heater",
    "west",
    "west",
    "take shell",
    "east",
    "north",
    "west",
    "north",
    "take jam",
    "north",
    "take astronaut ice cream",
    "north",
    "east",
    "south",
    "take space law space brochure",
    "north",
    "west",
    "south",
    "south",
    "east",
    "south",
    "take asterisk",
    "south",
    "take klein bottle",
    "east",
    "take spool of cat6",
    "west",
    "north",
    "north",
    "west",
    "south",
    "east",
    "west",
    "south",
    "west",
    "south"
]

# Not actually all the itmes but it's all the items you can pick up without ending the game.
all_items = [
    "spool of cat6",
    "space law space brochure",
    "asterisk",
    "jam",
    "shell",
    "astronaut ice cream",
    "space heater",
    "klein bottle"
]

# Collect all items and move to the room before the pressure-sensitive floor.
for command in collect_all_items_commands:
    run_command(command)

# Drop all items in the room before the pressure-sensitive floor.
for item in all_items:
    run_command("drop " + item)

# Try all possible item combinations get the correct weight for the pressure-sensitive floor.
for n in range(1, len(all_items)):
    for item_combination in combinations(all_items, n):
        # Pick up all the items in this attempt.
        for item in item_combination:
            run_command("take " + item)
        # Try to pass the pressure-sensitive floor. We don't move anywhere if we fail.
        run_command("south")
        # Drop the items from this attempt ready for the next one.
        for item in item_combination:
            run_command("drop " + item)

# Interative mode, if necessary.
interactive_commands = []
while True:
    command = input()
    interactive_commands.append(command)
    if command == "exit":
        break
    status, output = computer.resume(command + "\n")
    print(output)
    if status == Status.SUCCESS:
        break
print(interactive_commands)