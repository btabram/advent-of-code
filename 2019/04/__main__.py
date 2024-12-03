MIN = 125730
MAX = 579381


def is_valid_part1_pw(password):
    password_str = str(password) 
    if len(password_str) != 6:
        return False

    prev = 0 # Ensure first digit always passes the > prev check
    had_adjacent_match = False
    for digit in map(int, password_str):
        if digit < prev:
            return False
        if not had_adjacent_match and digit == prev:
            had_adjacent_match = True
        prev = digit

    return had_adjacent_match


print("The anser to Part 1 is {}.".format(
    sum(1 for _ in filter(is_valid_part1_pw, range(MIN, MAX + 1)))))


def is_valid_part2_pw(password):
    password_str = str(password) 
    if len(password_str) != 6:
        return False

    prev = 0 # Ensure first digit always passes the > prev check
    had_double = False
    just_had_interesting_double = False

    for digit in map(int, password_str):
        if digit < prev:
            return False

        if digit == prev:
            if just_had_interesting_double:
                # It's a triple not a double
                had_double = False
            elif not had_double:
                # We've found an interesting double
                had_double = True
                just_had_interesting_double = True
            else:
                # We've already had a double or it's longer than a triple
                just_had_interesting_double = False
        else:
            just_had_interesting_double = False

        prev = digit

    return had_double



print("The anser to Part 2 is {}.".format(
    sum(1 for _ in filter(is_valid_part2_pw, range(MIN, MAX + 1)))))
