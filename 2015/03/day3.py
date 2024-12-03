def get_next_pos(pos: tuple[int, int], instruction: str) -> tuple[int, int]:
    (x, y) = pos
    match instruction:
        case "^":
            return (x, y + 1)
        case ">":
            return (x + 1, y)
        case "v":
            return (x, y - 1)
        case "<":
            return (x - 1, y)


with open("input.txt") as f:
    input = f.readline().strip()

pos = (0, 0)
visited = set([pos])  # Deliver a present to the initial house
for instruction in input:
    pos = get_next_pos(pos, instruction)
    visited.add(pos)
print(f"The answer to Part 1 is {len(visited)}.")

santa_pos = (0, 0)
robo_santa_pos = (0, 0)
visited = set([santa_pos, robo_santa_pos])  # They both deliver to the initial house
for i in range(0, len(input), 2):
    santa_pos = get_next_pos(santa_pos, input[i])
    robo_santa_pos = get_next_pos(robo_santa_pos, input[i + 1])
    visited.update([santa_pos, robo_santa_pos])
print(f"The answer to Part 2 is {len(visited)}.")
