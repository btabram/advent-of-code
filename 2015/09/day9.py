import re


INPUT_REGEX = "^(\w+) to (\w+) = ([0-9]+)$"


with open("input.txt") as f:
    input_lines = [line.strip() for line in f.readlines()]

places: set[str] = set()
costs_map: dict[tuple[str, str], int] = dict()
for line in input_lines:
    [a, b, cost] = re.match(INPUT_REGEX, line).groups()
    places.update([a, b])
    costs_map[(a, b)] = int(cost)
    costs_map[(b, a)] = int(cost)


def visit_all(current: str, cost_so_far: int, visited: set[str]) -> list[int]:
    visited.add(current)
    # Finished! Return the cost of this route.
    if len(visited) == len(places):
        return [cost_so_far]
    results = []
    for next in filter(lambda p: p not in visited, places):
        cost = costs_map[(current, next)]
        # Recurse, finding out the cost of all possible routes. Creating a new |visited|
        # set is important so that the different branches do not affect one another via
        # a shared reference.
        results += visit_all(next, cost_so_far + cost, set(visited))
    return results


possible_costs = []
for possible_start_place in places:
    possible_costs += visit_all(possible_start_place, 0, set())

print(f"The answer to Part 1 is {min(possible_costs)}.")
print(f"The answer to Part 2 is {max(possible_costs)}.")
