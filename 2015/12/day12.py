from typing import Any
import json


with open("input.txt") as f:
    input = json.load(f)

part2 = False


def count(x: Any) -> int:
    if isinstance(x, int):
        return x
    elif isinstance(x, str):
        return 0  # Ignore strings
    elif isinstance(x, list):
        return sum(map(count, x))
    elif isinstance(x, dict):
        if part2 and "red" in x.values():
            return 0  # Ignore 'red' objects in part 2
        return sum(map(count, x.values()))
    else:
        raise Exception(f"Unexpected value: {x}")


print(f"The answer to Part 1 is {count(input)}.")
part2 = True
print(f"The answer to Part 2 is {count(input)}.")
