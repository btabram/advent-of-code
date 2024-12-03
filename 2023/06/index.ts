import { readFileSync } from "fs";
import { zip } from "lodash";
import { resolve } from "path";

const input = readFileSync(resolve(__dirname, "input.txt"), "utf8");

const [times, distances] = input.split("\n").map((line) =>
  line
    .split(" ")
    .slice(1)
    .filter(Boolean)
    .map((n) => parseInt(n, 10))
);

const races = (zip(times, distances) as [number, number][]).map(
  ([raceDuration, distanceRecord]) => ({
    raceDuration,
    distanceRecord,
  })
);
type Race = (typeof races)[number];

const waysToBeatRace = ({ raceDuration, distanceRecord }: Race) => {
  let ways = 0;
  let timeHoldingButton = 1;
  while (timeHoldingButton < raceDuration) {
    const distance = (raceDuration - timeHoldingButton) * timeHoldingButton;
    if (distance > distanceRecord) {
      ways += 1;
    }
    timeHoldingButton += 1;
  }
  return ways;
};

const part1 = races.map(waysToBeatRace).reduce((acc, score) => acc * score, 1);

const part2Input = input
  .split("\n")
  .map((line) => parseInt(line.split(":")[1]!.replaceAll(" ", ""), 10));
const part2Race = {
  raceDuration: part2Input[0]!,
  distanceRecord: part2Input[1]!,
};

// In a race we have:
// distance = (raceDuration - timeHoldingButton) * timeHoldingButton
//
// We can write this as a quadratic equation:
// y = (rd - x) * x = -x^2 + rd*x
//
// We want to find the x values where y = distanceRecord so we can can work out
// the number of ways to win the race.
//
// dr = -x^2 + rd*x
// -x^2 + rd*x - dr = 0
//
// We can use the quadratic formula to solve this with a = -1, b = rd & c = -dr:
const getInterceptPoints = (race: Race) => {
  const { raceDuration: rd, distanceRecord: dr } = race;
  const sqrtPart = Math.sqrt(rd * rd - 4 * -1 * -dr);
  return [(-rd + sqrtPart) / -2, (-rd - sqrtPart) / -2] as const;
};

const [firstIntercept, secondIntercept] = getInterceptPoints(part2Race);
const part2 = Math.floor(secondIntercept) - Math.ceil(firstIntercept) + 1;

console.log(`The answer the Part 1 is: ${part1}`);
console.log(`The answer the Part 2 is: ${part2}`);
