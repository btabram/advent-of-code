import { readFileSync } from "fs";
import { sum } from "lodash";
import { resolve } from "path";

const lines = readFileSync(resolve(__dirname, "input.txt"), "utf8").split("\n");

const valueHistories = lines.map((line) =>
  line.split(" ").map((v) => parseInt(v, 10))
);

function extrapolate(
  valueHistory: number[],
  direction: "forwards" | "backwards"
): number {
  const diffs = valueHistory.slice(1).map((v, i) => v - valueHistory[i]!);
  // If the diffs are all the same then we don't need another level of diffing
  // since it will be all zero, we can start working out the result value.
  // Otherwise, we need to keep diffing.
  const changeToApply = diffs.every((d) => d === diffs[0])
    ? diffs[0]!
    : extrapolate(diffs, direction);
  if (direction === "forwards") {
    return valueHistory[valueHistory.length - 1]! + changeToApply;
  } else {
    return valueHistory[0]! - changeToApply;
  }
}

const part1 = sum(valueHistories.map((vh) => extrapolate(vh, "forwards")));
const part2 = sum(valueHistories.map((vh) => extrapolate(vh, "backwards")));

console.log(`The answer the Part 1 is: ${part1}`);
console.log(`The answer the Part 2 is: ${part2}`);
