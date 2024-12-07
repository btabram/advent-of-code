import { readFileSync } from "node:fs";
import { resolve } from "node:path";

const digitWords: Record<string, string> = {
  one: "1",
  two: "2",
  three: "3",
  four: "4",
  five: "5",
  six: "6",
  seven: "7",
  eight: "8",
  nine: "9",
};
const matchToDigit = (match: string) => digitWords[match] ?? match;

const part1DigitRegex = "([0-9])";
const part2DigitRegex = `([0-9]|${Object.keys(digitWords).join("|")})`;

const lines = readFileSync(resolve(__dirname, "input.txt"), "utf8").split("\n");

const sumCalibrationValues = (baseDigitRegex: string) =>
  lines
    .map((line) => {
      const [, first] = line.match(new RegExp(baseDigitRegex)) ?? [];
      // Regex wildcards are greedy and will find the longest possible match.
      const [, last] = line.match(new RegExp(`.*${baseDigitRegex}`)) ?? [];
      if (!first || !last) {
        throw new Error(`Failed to get calibration value for: ${line}`);
      }
      return Number(`${matchToDigit(first)}${matchToDigit(last)}`);
    })
    .reduce((a, b) => a + b, 0);

console.log(
  `The answer the Part 1 is: ${sumCalibrationValues(part1DigitRegex)}`,
);
console.log(
  `The answer the Part 2 is: ${sumCalibrationValues(part2DigitRegex)}`,
);
