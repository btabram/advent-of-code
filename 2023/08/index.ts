import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { lcm } from "mathjs";

const lines = readFileSync(resolve(__dirname, "input.txt"), "utf8").split("\n");

const instructions = lines[0]!;

const nodes = new Map(
  lines.slice(2).map((nodeStr) => {
    // Network definitions are like "AAA = (BBB, CCC)".
    const nodeRe = /^([A-Z]+) = \(([A-Z]+), ([A-Z]+)\)$/;
    const [, name, left, right] = nodeStr.match(nodeRe) ?? [];
    if (!name || !left || !right) {
      throw new Error(`Failed to parse node: ${nodeStr}`);
    }
    return [name, { name, left, right }] as const;
  }),
);

let current = "AAA";
let step = 0;
while (current !== "ZZZ") {
  const directionToTake = instructions[step % instructions.length]!;
  step += 1;

  const currentNode = nodes.get(current)!;
  const nextNode =
    directionToTake === "L" ? currentNode.left : currentNode.right;

  current = nextNode;
}
const part1 = step;

const lastLetterIsA = (s: string) => s[s.length - 1] === "A";
const lastLetterIsZ = (s: string) => s[s.length - 1] === "Z";

let ghostPositions = Array.from(nodes.keys()).filter(lastLetterIsA);
step = 0;
// Work out when each independent ghost first finishes. It's a bit of an
// assumption (which does work because of carefully designed input!) but say
// that each path is cyclical and aligned with the instructions period, then
// the first time that all ghosts are finished together is the lowest common
// multiple of the individual cycle lengths. Then also assume that the cycle
// lengths are the same as the start -> first finish lengths.
const firstFinishes = new Array<number>(ghostPositions.length).fill(0);
while (!firstFinishes.every(Boolean)) {
  const directionToTake = instructions[step % instructions.length]!;
  step += 1;

  ghostPositions = ghostPositions.map((cur, i) => {
    const currentNode = nodes.get(cur)!;
    const nextNode =
      directionToTake === "L" ? currentNode.left : currentNode.right;

    if (!firstFinishes[i] && lastLetterIsZ(nextNode)) {
      firstFinishes[i] = step;
    }

    return nextNode;
  });
}
const part2 = firstFinishes.reduce((acc, value) => lcm(acc, value));

console.log(`The answer the Part 1 is: ${part1}`);
console.log(`The answer the Part 2 is: ${part2}`);
