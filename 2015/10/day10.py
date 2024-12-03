def look_and_say(value: list[int]) -> list[int]:
    result = []

    i = 0
    while i < len(value):
        current = value[i]

        run_length = 1
        while i + run_length < len(value) and value[i + run_length] == current:
            run_length += 1

        # [1, 1, 1] becomes "three ones" AKA [3, 1] etc.
        result += [run_length, current]

        i += run_length

    return result


with open("input.txt") as f:
    input = [int(c) for c in f.readline().strip()]

value = input
for n in range(1, 51):
    value = look_and_say(value)

    if n == 40:
        print(f"The answer to Part 1 is {len(value)}.")
    if n == 50:
        print(f"The answer to Part 2 is {len(value)}.")
