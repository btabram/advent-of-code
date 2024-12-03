import re


INPUT_REGEX = (
    "^(\w+) can fly ([0-9]+) km/s for ([0-9]+) seconds,"
    + " but then must rest for ([0-9]+) seconds.$"
)


class Reindeer:
    def __init__(self, s: str) -> None:
        [name, speed, flying_time, resting_time] = re.match(INPUT_REGEX, s).groups()
        self.name = name
        self.speed = int(speed)
        self.flying_time = int(flying_time)
        self.resting_time = int(resting_time)

    def get_distance(self, time: int) -> int:
        repeats = time // (self.flying_time + self.resting_time)
        remainder = time % (self.flying_time + self.resting_time)

        distance_per_repeat = self.speed * self.flying_time

        return (repeats * distance_per_repeat) + min(
            remainder * self.speed, distance_per_repeat
        )


with open("input.txt") as f:
    reindeer = [Reindeer(line.strip()) for line in f.readlines()]

part1 = max([r.get_distance(2503) for r in reindeer])

reindeer_scores = dict([(r.name, 0) for r in reindeer])
for t in range(1, 2504):
    winning_distance = max([r.get_distance(t) for r in reindeer])
    winners = [r.name for r in reindeer if r.get_distance(t) == winning_distance]
    for winner in winners:
        reindeer_scores[winner] += 1
part2 = max(reindeer_scores.values())

print(f"The answer to Part 1 is {part1}.")
print(f"The answer to Part 2 is {part2}.")
