import { readFileSync } from "node:fs";
import { resolve } from "node:path";

type Colour = "red" | "green" | "blue";

const lines = readFileSync(resolve(__dirname, "input.txt"), "utf8").split("\n");

const games = lines.map((line, i) => {
  // Lines are like "Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green".
  const rounds = line.split(";").map((round) => {
    const parsed = new Map<Colour, number>();
    for (const [, rawCount, colour] of round.matchAll(
      /([0-9]+) (red|green|blue)/g,
    )) {
      if (!rawCount || !colour) {
        throw new Error(`Failed to parse round: ${round}`);
      }
      parsed.set(colour as Colour, Number.parseInt(rawCount, 10));
    }
    return parsed;
  });
  return { id: i + 1, rounds }; // Assume games in input are in the right order
});

const LIMITS = {
  red: 12,
  green: 13,
  blue: 14,
} as const satisfies Record<Colour, number>;

const isValidRound = (round: Map<Colour, number>) => {
  for (const [colour, count] of round.entries()) {
    if (count > LIMITS[colour]) {
      return false;
    }
  }
  return true;
};

const part1 = games
  .filter(({ rounds }) => rounds.every(isValidRound))
  .reduce((acc, { id }) => acc + id, 0);

const getGamePower = ({ rounds }: (typeof games)[number]) => {
  const maxes = new Map<Colour, number>([
    ["red", 0],
    ["green", 0],
    ["blue", 0],
  ]);
  for (const round of rounds) {
    for (const [colour, count] of round.entries()) {
      if (count > maxes.get(colour)!) {
        maxes.set(colour, count);
      }
    }
  }
  return Array.from(maxes.values()).reduce((acc, v) => acc * v, 1);
};

const part2 = games.map(getGamePower).reduce((acc, v) => acc + v, 0);

console.log(`The answer the Part 1 is: ${part1}`);
console.log(`The answer the Part 2 is: ${part2}`);
