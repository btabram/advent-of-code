import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { groupBy, sortBy } from "lodash";

const cardToStrength: Record<string, string> = {
  2: "b",
  3: "c",
  4: "d",
  5: "e",
  6: "f",
  7: "g",
  8: "h",
  9: "i",
  T: "j",
  J: "k",
  Q: "l",
  K: "m",
  A: "n", // Ace is high
};

const handTypesToStrength = {
  highCard: "o",
  onePair: "p",
  twoPairs: "q",
  threeOfAKind: "r",
  fullHouse: "s",
  fourOfAKind: "t",
  fiveOfAKind: "u",
};

const lines = readFileSync(resolve(__dirname, "input.txt"), "utf8").split("\n");

const hands = lines.map((line) => {
  // A line is like "32T3K 765".
  const [cards, bidStr] = line.split(" ") as [string, string];

  const bid = Number.parseInt(bidStr, 10);
  const cardStrengths = cards
    .split("")
    .map((card) => cardToStrength[card]!)
    .join("");

  // Take a sorted (biggest to smallest) list of repeats in the hand and work
  // out the hand type strength, e.g. [3, 2] => "fullHouse" => "s".
  const repeatsToTypeStrength = (repeats: number[]) => {
    const biggestRepeat = repeats[0];
    const type = (() => {
      switch (biggestRepeat) {
        case 5:
          return "fiveOfAKind";
        case 4:
          return "fourOfAKind";
        case 3:
          return repeats[1] === 2 ? "fullHouse" : "threeOfAKind";
        case 2:
          return repeats[1] === 2 ? "twoPairs" : "onePair";
        default:
          return "highCard";
      }
    })();
    return handTypesToStrength[type];
  };

  // For part 1.
  const normalRepeats = Object.values(groupBy(cards))
    .map((group) => group.length)
    .sort()
    .reverse();
  const normalTypeStrength = repeatsToTypeStrength(normalRepeats);

  // For part 2, where "J" is a joker wildcard.
  const cardsWithoutJokers = cards.replaceAll("J", "");
  const jokerCount = 5 - cardsWithoutJokers.length;

  const repeatsWithoutJokers = Object.values(groupBy(cardsWithoutJokers))
    .map((group) => group.length)
    .sort()
    .reverse();

  // The best use of joker wildcards is to make the biggest repeat bigger.
  const wildcardRepeats = [
    jokerCount + (repeatsWithoutJokers[0] ?? 0),
    ...repeatsWithoutJokers.slice(1),
  ];
  const wildTypeStrength = repeatsToTypeStrength(wildcardRepeats);

  return {
    normalTypeStrength,
    wildTypeStrength,
    cardStrengths,
    bid,
  };
});

const scoreRankedHands = (rankedHands: typeof hands) =>
  rankedHands.reduce((acc, hand, i) => acc + hand.bid * (i + 1), 0);

const part1 = scoreRankedHands(
  sortBy(hands, "normalTypeStrength", "cardStrengths"),
);

const part2 = scoreRankedHands(
  sortBy(
    // Jokers (strength "k") are the lowest strength under part 2 rules.
    hands.map(({ cardStrengths, ...rest }) => ({
      cardStrengths: cardStrengths.replaceAll("k", "a"),
      ...rest,
    })),
    "wildTypeStrength",
    "cardStrengths",
  ),
);

console.log(`The answer the Part 1 is: ${part1}`);
console.log(`The answer the Part 2 is: ${part2}`);
