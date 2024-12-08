import { readFileSync } from "node:fs";
import { resolve } from "node:path";

type Vec = { x: number; y: number };

const vecToStr = ({ x, y }: Vec): string => `${x}_${y}`;

const lines = readFileSync(resolve(__dirname, "input.txt"), "utf8").split("\n");

const getNullable = ({ x, y }: Vec): string | null => lines[y]?.[x] ?? null;
const get = (pos: Vec): string => {
  const value = getNullable(pos);
  if (!value) {
    throw new Error(
      `Expected value at (${pos.x}, ${pos.y}) but didn't find one`
    );
  }
  return value;
};

const getConnectingNeighbours = (pos: Vec): Vec[] => {
  const { x, y } = pos;

  const pipe = get(pos);
  switch (pipe) {
    case "|": // vertical pipe
      return [
        { x, y: y + 1 },
        { x, y: y - 1 },
      ];

    case "-": // horizontal pipe
      return [
        { x: x + 1, y },
        { x: x - 1, y },
      ];

    case "L": // north <-> east corner pipe
      return [
        { x, y: y - 1 }, // Note that positive y is downwards (south)
        { x: x + 1, y },
      ];

    case "J": // north <-> west corner pipe
      return [
        { x, y: y - 1 },
        { x: x - 1, y },
      ];

    case "7": // south <-> west corner pipe
      return [
        { x, y: y + 1 },
        { x: x - 1, y },
      ];

    case "F": // south <-> east corner pipe
      return [
        { x, y: y + 1 },
        { x: x + 1, y },
      ];

    default:
      throw new Error(`Unexpected pipe: ${pipe}`);
  }
};

const start = (() => {
  let s: Vec | undefined;
  lines.forEach((line, y) => {
    const startX = line.indexOf("S");
    if (startX !== -1) {
      s = { x: startX, y };
    }
  });
  if (!s) throw new Error("Failed to find start position");
  return s;
})();

/** Find the two pipes which connect to the start position (we know the start is part of a loop) */
const firstPipes = [
  { x: start.x + 1, y: start.y },
  { x: start.x - 1, y: start.y },
  { x: start.x, y: start.y + 1 },
  { x: start.x, y: start.y - 1 },
].filter((position) =>
  getConnectingNeighbours(position).some(
    (p) => p.x === start.x && p.y === start.y
  )
);
if (firstPipes.length !== 2) throw new Error("Failed to find first pipes");

const loopPipes = new Set<string>([start, ...firstPipes].map(vecToStr));

let prevPipes = firstPipes;
// biome-ignore lint/correctness/noConstantCondition: it's intentional here
for (let i = 1; true; i++) {
  const next = prevPipes
    .flatMap(getConnectingNeighbours)
    .filter((position) => !loopPipes.has(vecToStr(position)));

  if (!next.length) {
    // Part 1 wants the number of steps from the start to farthest away point on the loop
    console.log(`The answer the Part 1 is: ${i}`);
    break;
  }

  for (const position of next) {
    loopPipes.add(vecToStr(position));
  }

  prevPipes = next;
}

let tilesWithinLoop = 0;

lines.forEach((line, y) => {
  let loopCrossingsCount = 0;
  let followingLoopStart: string | null = null;
  line.split("").forEach((char, x) => {
    if (loopPipes.has(vecToStr({ x, y }))) {
      switch (char) {
        case "|":
          // Straightforward loop crossing
          loopCrossingsCount++;
          return;

        case "-":
          // Nothing to do, keep following the loop
          return;

        case "S": // Laziness, hardcode the type of pipe for "S" for my input instead of calculating
        case "L":
          // We're starting to follow the loop, don't know if we're crossing it yet
          followingLoopStart = "L";
          return;

        case "J":
          // We've finished following the loop, check if we crossed it (or it doubled back)
          if (followingLoopStart === "F") {
            loopCrossingsCount++;
          }
          return;

        case "7":
          // We've finished following the loop, check if we crossed it (or it doubled back)
          if (followingLoopStart === "L") {
            loopCrossingsCount++;
          }
          return;

        case "F":
          // We're starting to follow the loop, don't know if we're crossing it yet
          followingLoopStart = "F";
          return;

        default:
          throw new Error(`Unexpected pipe: ${char}`);
      }
    }

    // Credit to reddit for this clever solution! We count the number of times that we've crossed
    // the loop (starting from outside) and if we've crossed an odd number of times we must be in a
    // region contained within the loop. An even number of crossings means we're outside the loop.
    if (loopCrossingsCount % 2) {
      tilesWithinLoop++;
    }
  });
});

console.log(`The answer the Part 2 is: ${tilesWithinLoop}`);
