with open("input.txt") as f:
    containers = [int(line.strip()) for line in f.readlines()]
container_count = len(containers)

possible_combinations = []
for x in range(1 << container_count):  # Use a bitmask to list all possible combinations
    combination = []
    for i in range(container_count):
        if x & (1 << i) != 0:
            combination.append(containers[i])
    if sum(combination) == 150:
        possible_combinations.append(combination)

part1 = len(possible_combinations)

fewest_required_containers = min(map(len, possible_combinations))
part2 = len([c for c in possible_combinations if len(c) == fewest_required_containers])

print(f"The answer to Part 1 is {part1}.")
print(f"The answer to Part 2 is {part2}.")
