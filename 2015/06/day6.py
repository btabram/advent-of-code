import re


INPUT_REGEX = "^(turn on|turn off|toggle) ([0-9]+),([0-9]+) through ([0-9]+),([0-9]+)$"


class Instruction:
    def __init__(self, instruction_string: str) -> None:
        [action, x1, y1, x2, y2] = re.match(INPUT_REGEX, instruction_string).groups()
        self.action = action
        self.x1 = int(x1)
        self.y1 = int(y1)
        self.x2 = int(x2)
        self.y2 = int(y2)


def follow_instructions(instructions: list[Instruction], is_part1: bool) -> int:
    lights = [[0 for _ in range(1000)] for _ in range(1000)]

    for inst in instructions:
        for x in range(inst.x1, inst.x2 + 1):
            for y in range(inst.y1, inst.y2 + 1):
                val = lights[x][y]
                match inst.action:
                    case "turn on":
                        if is_part1:
                            val = 1
                        else:
                            val += 1
                    case "turn off":
                        if is_part1:
                            val = 0
                        else:
                            val = max(val - 1, 0)
                    case "toggle":
                        if is_part1:
                            val = 0 if val == 1 else 1
                        else:
                            val += 2
                lights[x][y] = val

    return sum([sum(line) for line in lights])


with open("input.txt") as f:
    instructions = [Instruction(line.strip()) for line in f.readlines()]

print(f"The answer to Part 1 is {follow_instructions(instructions, is_part1=True)}.")
print(f"The answer to Part 2 is {follow_instructions(instructions, is_part1=False)}.")
