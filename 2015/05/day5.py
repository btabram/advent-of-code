def is_nice_part1(string: str) -> bool:
    if any([bad_substring in string for bad_substring in ["ab", "cd", "pq", "xy"]]):
        return False

    has_repeat = False
    prev_char = None
    vowel_count = 0
    for char in string:
        if not has_repeat:
            has_repeat = char == prev_char
        if char in "aeiou":
            vowel_count += 1
        prev_char = char

    return has_repeat and vowel_count >= 3


def is_nice_part2(string: str) -> bool:
    has_repeated_pair = False
    has_separated_repeat = False
    for i in range(len(string) - 2):
        if not has_repeated_pair:
            pair = string[i : i + 2]
            has_repeated_pair = string.find(pair, i + 2) != -1
        if not has_separated_repeat:
            has_separated_repeat = string[i] == string[i + 2]

    return has_repeated_pair and has_separated_repeat


with open("input.txt") as f:
    input_strings = [line.strip() for line in f.readlines()]

part1 = 0
part2 = 0
for string in input_strings:
    if is_nice_part1(string):
        part1 += 1
    if is_nice_part2(string):
        part2 += 1

print(f"The answer to Part 1 is {part1}.")
print(f"The answer to Part 2 is {part2}.")
