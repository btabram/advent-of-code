import { readFileSync } from "fs";
import { chunk, min } from "lodash";
import { resolve } from "path";

const input = readFileSync(resolve(__dirname, "input.txt"), "utf8");

const [seedStr, ...mapStrs] = input.split("\n\n");

const seeds = seedStr!
  .split(" ")
  .slice(1)
  .map((n) => parseInt(n, 10));

const maps = mapStrs.map((mapStr) => {
  const rawMap = mapStr.split("\n");
  const [, source, dest] = rawMap[0]?.match(/(\w+)-to-(\w+) map:/) ?? [];
  if (!source || !dest) {
    throw new Error(`Failed to parse map name: ${rawMap[0]}`);
  }
  const ranges = rawMap.slice(1).map((line) => {
    const [destRangeStart, sourceRangeStart, rangeLength] = line
      .split(" ")
      .map((n) => parseInt(n, 10)) as [number, number, number];
    return {
      destRangeStart,
      sourceRangeStart,
      sourceRangeEnd: sourceRangeStart + rangeLength,
    };
  });
  return { source, dest, ranges };
});

type Map = (typeof maps)[number];

const convertWithMap = ({ ranges }: Map, sourceNumber: number) => {
  for (const { destRangeStart, sourceRangeStart, sourceRangeEnd } of ranges) {
    // Ranges are inclusive at the start but not the end.
    if (sourceNumber >= sourceRangeStart && sourceNumber < sourceRangeEnd) {
      return destRangeStart + (sourceNumber - sourceRangeStart);
    }
  }
  return sourceNumber;
};

// Assuming that the maps are all in the right order here.
const fullyConvertedSeeds = seeds.map((seed) =>
  maps.reduce((acc, map) => convertWithMap(map, acc), seed)
);
const part1 = min(fullyConvertedSeeds);

const seedRanges = (chunk(seeds, 2) as [number, number][]).map(
  ([start, length]) => ({
    start,
    end: start + length,
  })
);
type Range = (typeof seedRanges)[number];

const convertRangeWithMap = ({ ranges }: Map, sourceRange: Range) => {
  const { start, end } = sourceRange;
  const destRages: Range[] = [];

  let current = start;
  while (current < end) {
    const [inMappingRange] = ranges.filter(
      ({ sourceRangeStart, sourceRangeEnd }) =>
        current >= sourceRangeStart && current < sourceRangeEnd
    );
    const nextMappingRangeStart = inMappingRange
      ? undefined
      : min(
          ranges
            .filter(({ destRangeStart }) => destRangeStart >= current)
            .map(({ destRangeStart }) => destRangeStart)
        );

    const destRangeStart = inMappingRange
      ? inMappingRange.destRangeStart +
        (current - inMappingRange.sourceRangeStart)
      : current;

    // Out input range is either at least as long as the current mapping range
    // (or gap between mapping ranges) or it fits inside the current range.
    const lengthInThisRange = Math.min(
      inMappingRange
        ? inMappingRange.sourceRangeEnd - current
        : nextMappingRangeStart! - current,
      end - current
    );

    destRages.push({
      start: destRangeStart,
      end: destRangeStart + lengthInThisRange,
    });

    current += lengthInThisRange;
  }
  return destRages;
};

const fullyConvertedRanges = maps.reduce(
  (acc, map) => acc.map((r) => convertRangeWithMap(map, r)).flat(),
  seedRanges
);
const part2 = min(fullyConvertedRanges.map(({ start }) => start));

console.log(`The answer the Part 1 is: ${part1}`);
console.log(`The answer the Part 2 is: ${part2}`);
