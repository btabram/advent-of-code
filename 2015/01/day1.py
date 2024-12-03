with open("input.txt") as f:
    input = f.readline().strip()

floor = 0
enter_basement_pos = None
for (i, c) in enumerate(input):
    floor = floor + 1 if c == "(" else floor - 1
    if enter_basement_pos is None and floor == -1:
        enter_basement_pos = i + 1

print(f"The answer to Part 1 is {floor}.")
print(f"The answer to Part 2 is {enter_basement_pos}.")
