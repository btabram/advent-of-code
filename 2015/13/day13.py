import re


INPUT_REGEX = (
    "^(\w+) would (gain|lose) ([0-9]+) happiness units by sitting next to (\w+).$"
)


with open("input.txt") as f:
    input_lines = [line.strip() for line in f.readlines()]

happiness_effects_map: dict[str, dict[str, int]] = dict()
for line in input_lines:
    [a, sign, amount, b] = re.match(INPUT_REGEX, line).groups()
    if a not in happiness_effects_map:
        happiness_effects_map[a] = {}
    happiness_effects_map[a][b] = int(amount) if sign == "gain" else -int(amount)


def do_seating_arragement(seated: list[str], unseated: set[str], score: int) -> int:
    if len(unseated) == 0:
        first_person = seated[0]
        last_person = seated[-1]
        score += happiness_effects_map[first_person][last_person]
        score += happiness_effects_map[last_person][first_person]
        return score

    final_scores = []
    for person_to_seat in unseated:
        new_score = score
        if len(seated) > 0:
            end_person = seated[-1]
            new_score += happiness_effects_map[end_person][person_to_seat]
            new_score += happiness_effects_map[person_to_seat][end_person]
        new_seated = seated + [person_to_seat]
        new_unseated = unseated - set([person_to_seat])
        final_scores.append(do_seating_arragement(new_seated, new_unseated, new_score))
    return max(final_scores)


part1 = do_seating_arragement([], set(happiness_effects_map.keys()), 0)

other_people = list(happiness_effects_map.keys())
happiness_effects_map["me"] = {}
for person in other_people:
    happiness_effects_map["me"][person] = 0
    happiness_effects_map[person]["me"] = 0

part2 = do_seating_arragement([], set(happiness_effects_map.keys()), 0)


print(f"The answer to Part 1 is {part1}.")
print(f"The answer to Part 2 is {part2}.")
