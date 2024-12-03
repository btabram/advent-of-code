import re


INPUT_REGEX = "^Sue ([0-9]+): (.+)$"


PRESENT_GIVER_ATTRIBUTES = {
    "children": 3,
    "cats": 7,
    "samoyeds": 2,
    "pomeranians": 3,
    "akitas": 0,
    "vizslas": 0,
    "goldfish": 5,
    "trees": 3,
    "cars": 2,
    "perfumes": 1,
}


class AuntSue:
    def __init__(self, s: str) -> None:
        [number, attributes] = re.match(INPUT_REGEX, s).groups()
        self.number = int(number)
        self.attibutes: dict[str, int] = {}
        for attibute in attributes.split(", "):
            [name, count] = attibute.split(": ")
            self.attibutes[name] = int(count)


def could_be_part1_present_giver(sue: AuntSue) -> bool:
    for name, pg_count in PRESENT_GIVER_ATTRIBUTES.items():
        sue_count = sue.attibutes.get(name)
        if sue_count is not None and sue_count != pg_count:
            return False  # This Sue does not match the present giver
    return True


def could_be_part2_present_giver(sue: AuntSue) -> bool:
    for name, pg_count in PRESENT_GIVER_ATTRIBUTES.items():
        sue_count = sue.attibutes.get(name)
        if sue_count is not None:
            if name in ["cats", "trees"]:  # These attribute counts are a lower bound
                if sue_count <= pg_count:
                    return False
            elif name in ["pomeranians", "goldfish"]:  # These counts are an upper bound
                if sue_count >= pg_count:
                    return False
            else:
                if sue_count != pg_count:
                    return False
    return True


with open("input.txt") as f:
    sues = [AuntSue(line.strip()) for line in f.readlines()]

part1 = [sue for sue in sues if could_be_part1_present_giver(sue)]
part2 = [sue for sue in sues if could_be_part2_present_giver(sue)]

assert len(part1) == 1
assert len(part2) == 1

print(f"The answer to Part 1 is {part1[0].number}.")
print(f"The answer to Part 2 is {part2[0].number}.")
