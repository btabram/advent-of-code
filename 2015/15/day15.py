import re
from typing import Iterator


INPUT_REGEX = (
    "^(\w+): capacity (-?[0-9]+), durability (-?[0-9]+),"
    + " flavor (-?[0-9]+), texture (-?[0-9]+), calories ([0-9]+)$"
)


class Ingredient:
    def __init__(self, s: str) -> None:
        [name, cap, dur, fla, tex, cal] = re.match(INPUT_REGEX, s).groups()
        self.name = name
        self.cap = int(cap)
        self.dur = int(dur)
        self.fla = int(fla)
        self.tex = int(tex)
        self.cal = int(cal)


def ingredient_combinations(
    partial_combinations: list[int], n: int
) -> Iterator[list[int]]:
    if n == 1:
        yield partial_combinations + [100 - sum(partial_combinations)]
        return

    for i in range(max(101 - sum(partial_combinations), 1)):
        for combination in ingredient_combinations(partial_combinations + [i], n - 1):
            yield combination


with open("input.txt") as f:
    ingredients = [Ingredient(line.strip()) for line in f.readlines()]


part1 = 0
part2 = 0
for combination in ingredient_combinations([], len(ingredients)):
    capacity = 0
    durability = 0
    flavour = 0
    texture = 0
    calories = 0
    for ingredient, amount in zip(ingredients, combination):
        capacity += ingredient.cap * amount
        durability += ingredient.dur * amount
        flavour += ingredient.fla * amount
        texture += ingredient.tex * amount
        calories += ingredient.cal * amount

    total_score = (
        max(capacity, 0) * max(durability, 0) * max(flavour, 0) * max(texture, 0)
    )

    if total_score > part1:
        part1 = total_score

    if calories == 500 and total_score > part2:
        part2 = total_score

print(f"The answer to Part 1 is {part1}.")
print(f"The answer to Part 2 is {part2}.")
