with open("input.txt") as f:
    input_lines = [line.strip() for line in f.readlines()]

# Count the number of code chars that are lost in the in-memory representation.
part1 = 0
for line in input_lines:
    # Count and strip the outer double quotes.
    line_content = line[1 : len(line) - 1]
    part1 += 2

    # We need to be careful not to double-count any escapes so we iterate through the
    # sting, counting as we go.
    i = 0
    while i < len(line_content):
        if line_content[i] == "\\":
            match line_content[i + 1]:
                case "\\":
                    chars_to_escape = 1
                case '"':
                    chars_to_escape = 1
                case "x":
                    chars_to_escape = 3
        else:
            chars_to_escape = 0
        part1 += chars_to_escape
        i += 1 + chars_to_escape

# Count the number of extra code chars needed to escape the strings. We don't need to
# worry about double-counting here. We simply need to add a double quote to each end
# and to escape all backslashes and double quotes with a blackslash.
part2 = 0
for line in input_lines:
    part2 += 2 + line.count('"') + line.count("\\")

print(f"The answer to Part 1 is {part1}.")
print(f"The answer to Part 2 is {part2}.")
