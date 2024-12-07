import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { uniq, uniqBy } from "lodash";

type CoordStr = `${number},${number}`;

const parse = (coord: CoordStr): [number, number] =>
  coord.split(",").map((n) => Number.parseInt(n, 10)) as [number, number];

const getNeighbours = (coord: CoordStr): CoordStr[] => {
  const [x, y] = parse(coord);
  const neighbours: CoordStr[] = [];
  for (let i = x - 1; i <= x + 1; i++) {
    for (let j = y - 1; j <= y + 1; j++) {
      if (i === x && j === y) {
        continue; // Don't include ourselves.
      }
      neighbours.push(`${i},${j}`);
    }
  }
  return neighbours;
};

const input = readFileSync(resolve(__dirname, "input.txt"), "utf8");

const numbers: { value: number; coords: CoordStr[] }[] = [];
const symbols = new Map<CoordStr, string>();

// Build up the numbers, one digit at a time.
let buildingNumber: (readonly [string, CoordStr])[] | undefined;
const buildNumber = (char: string, coord: CoordStr) => {
  if (!buildingNumber) {
    buildingNumber = [[char, coord]];
  } else {
    buildingNumber.push([char, coord]);
  }
};
const tryFinishNumber = () => {
  if (buildingNumber) {
    const value = Number.parseInt(buildingNumber.map(([n]) => n).join(""), 10);
    const coords = buildingNumber.map(([, c]) => c);
    numbers.push({ value, coords });
    buildingNumber = undefined;
  }
};

input.split("\n").forEach((line, y) => {
  Array.from(line).forEach((char, x) => {
    if (char.match(/[0-9]/)) {
      // It's a digit, build up the number.
      buildNumber(char, `${x},${y}`);
    } else {
      // Not a digit.
      tryFinishNumber();

      if (char !== ".") {
        // It must be a symbol.
        symbols.set(`${x},${y}`, char);
      }
    }
  });
  tryFinishNumber();
});

// Sum all numbers that have a symbol next to them.
let part1 = 0;
for (const { value, coords } of numbers) {
  const neighbours = uniq(coords.flatMap(getNeighbours));
  if (neighbours.some((coord) => symbols.has(coord))) {
    part1 += value;
  }
}

const numbersByCoord = new Map(
  numbers.flatMap(({ value, coords }, id) =>
    coords.map((coord) => [coord, { id, value }] as const),
  ),
);

const gears = Array.from(symbols.entries())
  .filter(([, char]) => char === "*")
  .flatMap(([coord]) => {
    const neighbouringNumbers = getNeighbours(coord)
      .map((n) => numbersByCoord.get(n))
      .flatMap((n) => n ?? []); // Using as a type-aware filter
    const uniqueNeighbourNumbers = uniqBy(neighbouringNumbers, "id");
    if (uniqueNeighbourNumbers.length !== 2) {
      return [];
    }
    return {
      coord,
      // Gear ratio is the two numbers multiplied together.
      ratio: uniqueNeighbourNumbers.reduce((acc, { value }) => acc * value, 1),
    };
  });

// Sum all gear ratios.
const part2 = gears.reduce((acc, { ratio }) => acc + ratio, 0);

console.log(`The answer the Part 1 is: ${part1}`);
console.log(`The answer the Part 2 is: ${part2}`);
