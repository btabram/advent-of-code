from collections import deque


def increment(s: str) -> str:
    new_s = deque()
    done_increment = False
    for c in reversed(s):
        if not done_increment:
            if c == "z":
                c = "a"
            else:
                c = chr(ord(c) + 1)
                done_increment = True
        new_s.appendleft(c)
    return "".join(new_s)


def is_valid_password(s: str) -> bool:
    if any([confusing_c in s for confusing_c in "iol"]):
        return False

    has_straight = False
    for i in range(len(s) - 2):
        curr = ord(s[i])
        if curr + 1 == ord(s[i + 1]) and curr + 2 == ord(s[i + 2]):
            has_straight = True
            break
    if not has_straight:
        return False

    double_count = 0
    i = 0
    while i < len(s) - 1:
        if s[i] == s[i + 1]:
            double_count += 1
            i += 1  # Skip ahead, we only want non-overlapping doubles
        i += 1
    return double_count >= 2


def find_next_valid_password(s: str) -> str:
    while True:
        s = increment(s)
        if is_valid_password(s):
            return s


with open("input.txt") as f:
    input = f.readline().strip()

part1 = find_next_valid_password(input)
part2 = find_next_valid_password(part1)

print(f"The answer to Part 1 is {part1}.")
print(f"The answer to Part 2 is {part2}.")
